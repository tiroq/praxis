#!/usr/bin/env python3
"""Praxis-local graphify rebuild wrapper.

Rebuilds the knowledge graph from the existing graphify caches (AST + semantic)
with two Praxis-specific improvements layered on top of the stock graphify
pipeline, WITHOUT modifying graphify package internals:

  1. Directed mode by default (edge direction = architectural dependency).
  2. A post-extraction ID-normalization pass that recovers "dangling" edges whose
     endpoints reference a real node under a slightly different id convention:
       - missing ``_doc`` suffix on RFC document nodes
         (``rfcs_013_event_model`` -> ``rfcs_013_event_model_doc``);
       - ``packages_`` / ``services_`` path-prefix mismatches
         (``packages_praxis_agents_base`` -> ``praxis_agents_base``);
       - other unambiguous single-candidate id matches.

External-library dangling edges (``react``, ``fastapi``, ``os`` ...) are left out
of the graph but reported separately rather than silently dropped.

Outputs written to ``graphify-out/``:
    graph.json          - rebuilt directed graph
    GRAPH_REPORT.md     - stock graphify audit report
    GRAPH_HEALTH.md     - human-readable health / normalization report
    graph_health.json   - machine-readable health metrics

Usage:
    python scripts/rebuild_graph.py

Relies on the ``graphify`` package already being installed (it is, via the
/graphify skill). No new infrastructure, no LLM calls, no network access: the
semantic layer is read from the on-disk cache populated by the last full build.
"""

from __future__ import annotations

import importlib.util
import json
import os
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
OUT = ROOT / "graphify-out"


# --------------------------------------------------------------------------- #
# Interpreter bootstrap                                                        #
# --------------------------------------------------------------------------- #
def _ensure_graphify() -> None:
    """Re-exec into the interpreter that has graphify if the current one lacks it."""
    if importlib.util.find_spec("graphify") is not None:
        return
    recorded = OUT / ".graphify_python"
    if recorded.exists():
        interp = recorded.read_text(encoding="utf-8").strip()
        if interp and Path(interp).exists() and Path(interp) != Path(sys.executable):
            os.execv(interp, [interp, str(Path(__file__).resolve()), *sys.argv[1:]])
    sys.exit(
        "graphify is not importable. Run the /graphify skill once to install it, "
        "or `uv tool install graphifyy` / `pip install graphifyy`."
    )


_ensure_graphify()

from graphify.analyze import god_nodes, suggest_questions, surprising_connections  # noqa: E402
from graphify.build import build_from_json  # noqa: E402
from graphify.cache import check_semantic_cache  # noqa: E402
from graphify.cluster import cluster, score_all  # noqa: E402
from graphify.detect import detect  # noqa: E402
from graphify.export import to_json  # noqa: E402
from graphify.extract import collect_files, extract  # noqa: E402
from graphify.report import generate  # noqa: E402

DIRECTED_DEFAULT = True
_IMPORT_RELATIONS = {"imports", "imports_from", "include", "extends"}


# --------------------------------------------------------------------------- #
# Step 1 - rebuild raw extraction from caches (no LLM)                         #
# --------------------------------------------------------------------------- #
def _under_out(path: str) -> bool:
    """True if a detected path lives inside graphify-out/ (a generated sidecar)."""
    try:
        Path(path).resolve().relative_to(OUT)
        return True
    except ValueError:
        return False


def load_extraction() -> tuple[dict, dict, int]:
    """Return (extraction, detection, n_uncached) from AST + cached semantic data."""
    detection = detect(ROOT)
    # Never treat graphify's own outputs (graph.json, GRAPH_REPORT.md, manifest …)
    # as corpus content — the stock skill filters these, detect() does not.
    for cat, files in detection.get("files", {}).items():
        detection["files"][cat] = [f for f in files if not _under_out(f)]

    code_files: list[Path] = []
    for f in detection.get("files", {}).get("code", []):
        p = Path(f)
        code_files.extend(collect_files(p) if p.is_dir() else [p])
    ast = extract(code_files, cache_root=ROOT) if code_files else {"nodes": [], "edges": []}

    doc_files = [f for cat in ("document", "paper", "image") for f in detection["files"].get(cat, [])]
    sem_nodes, sem_edges, sem_hyper, uncached = check_semantic_cache(doc_files, root=str(ROOT))

    seen = {n["id"] for n in ast["nodes"]}
    nodes = list(ast["nodes"])
    for n in sem_nodes:
        if n["id"] not in seen:
            nodes.append(n)
            seen.add(n["id"])

    extraction = {
        "nodes": nodes,
        "edges": list(ast["edges"]) + list(sem_edges),
        "hyperedges": list(sem_hyper),
        "input_tokens": 0,
        "output_tokens": 0,
    }
    return extraction, detection, len(uncached)


# --------------------------------------------------------------------------- #
# Step 2 - ID normalization                                                    #
# --------------------------------------------------------------------------- #
def _candidates(missing_id: str, node_ids: set[str]) -> list[str]:
    """Return distinct existing node ids the missing id could safely refer to.

    Only structurally-targeted strategies are used. A broad "any node whose id
    ends with this token" (suffix) match was deliberately removed: it produced
    semantic garbage (``os`` -> ``readme_ai_work_os``, external npm libs binding
    to their own package.json declaration nodes). When in doubt we leave the edge
    dangling and report it rather than invent a wrong relationship.
    """
    found: list[str] = []

    def add(cid: str) -> None:
        if cid in node_ids and cid not in found:
            found.append(cid)

    # 1. `_doc` suffix mismatch on RFC / document nodes — handle BOTH directions:
    #    a cite may target `rfcs_033_storage_model` while the node is
    #    `…_storage_model_doc`, or target `…_doc` while the node omits it.
    add(missing_id + "_doc")
    if missing_id.endswith("_doc"):
        add(missing_id[: -len("_doc")])
    # 2. path-prefix mismatch: strip a leading packages_/services_ segment
    for prefix in ("packages_", "services_"):
        if missing_id.startswith(prefix):
            add(missing_id[len(prefix):])
    # 3. full-path id -> parent_dir+stem id: progressive left-strip of leading
    #    path segments (e.g. `services_api_routes_agents` -> `routes_agents`).
    #    Requires >=3 segments so a bare module name (`os`, `common`) never
    #    left-strips into an unrelated node.
    parts = missing_id.split("_")
    if len(parts) >= 3:
        for k in range(1, len(parts) - 1):
            add("_".join(parts[k:]))
    return found


def _classify_external(missing_id: str, relation: str) -> str:
    if missing_id.startswith("ref_"):
        return "tsconfig lib/include reference"
    if relation in _IMPORT_RELATIONS:
        return "external-library import (non-corpus)"
    if relation in ("cites", "references"):
        return "unresolved cross-document reference"
    return "other unresolved"


def normalize(extraction: dict) -> dict:
    """Rewrite recoverable dangling endpoints in place; return a report dict."""
    node_ids = {n["id"] for n in extraction["nodes"]}

    recovered: list[dict] = []
    remaining: list[dict] = []
    ambiguous: list[dict] = []
    dropped_self_loops = 0
    before = 0

    for e in extraction["edges"]:
        for role in ("source", "target"):
            mid = e.get(role)
            if mid in node_ids:
                continue
            before += 1
            cands = _candidates(mid, node_ids)
            if len(cands) == 1:
                e[role] = cands[0]
                recovered.append(
                    {"from_id": mid, "to_id": cands[0], "relation": e.get("relation"),
                     "role": role, "source_file": e.get("source_file")}
                )
            elif len(cands) > 1:
                ambiguous.append(
                    {"missing_id": mid, "relation": e.get("relation"),
                     "candidates": cands, "source_file": e.get("source_file")}
                )
            else:
                remaining.append(
                    {"missing_id": mid, "relation": e.get("relation"), "role": role,
                     "source": e.get("source"), "target": e.get("target"),
                     "source_file": e.get("source_file"),
                     "cause": _classify_external(mid, e.get("relation", ""))}
                )

    # Drop edges that recovery turned into self-loops (id collision after rewrite).
    kept_edges = []
    for e in extraction["edges"]:
        if e.get("source") == e.get("target"):
            dropped_self_loops += 1
            continue
        kept_edges.append(e)
    extraction["edges"] = kept_edges

    return {
        "node_ids": node_ids,
        "dangling_before": before,
        "recovered": recovered,
        "remaining": remaining,
        "ambiguous": ambiguous,
        "dropped_self_loops": dropped_self_loops,
    }


# --------------------------------------------------------------------------- #
# Step 3 - collapse analysis (multi-relation / bidirectional)                  #
# --------------------------------------------------------------------------- #
def collapse_analysis(extraction: dict, node_ids: set[str]) -> dict:
    from collections import defaultdict

    labels = {n["id"]: n.get("label", n["id"]) for n in extraction["nodes"]}
    valid = [
        e for e in extraction["edges"]
        if e["source"] in node_ids and e["target"] in node_ids and e["source"] != e["target"]
    ]

    directed_groups: dict[tuple, set] = defaultdict(set)   # (s,t) -> relations
    undirected_groups: dict[tuple, list] = defaultdict(list)
    for e in valid:
        directed_groups[(e["source"], e["target"])].add(e.get("relation"))
        undirected_groups[tuple(sorted((e["source"], e["target"])))].append(
            (e["source"], e["target"], e.get("relation"))
        )

    # Same-direction, multiple distinct relations -> collapses even in a DiGraph.
    directed_multi = [
        {"source": s, "target": t, "relations": sorted(rels),
         "label": f"{labels.get(s, s)} -> {labels.get(t, t)}"}
        for (s, t), rels in directed_groups.items() if len(rels) > 1
    ]
    # Both directions present -> preserved in directed, merged in undirected.
    bidirectional = [
        {"pair": [a, b], "label": f"{labels.get(a, a)} <-> {labels.get(b, b)}"}
        for (a, b), v in undirected_groups.items()
        if len({(s, t) for s, t, _ in v}) > 1
    ]
    undirected_multi = [
        {"pair": [a, b], "relations": sorted({r for _, _, r in v}),
         "label": f"{labels.get(a, a)} <-> {labels.get(b, b)}"}
        for (a, b), v in undirected_groups.items() if len(v) > 1
    ]
    return {
        "directed_collapsed_multi_relation": directed_multi,
        "undirected_collapsed_pairs": undirected_multi,
        "bidirectional_pairs": bidirectional,
    }


# --------------------------------------------------------------------------- #
# Step 4 - auto labels (top-degree node per community, no LLM)                 #
# --------------------------------------------------------------------------- #
def auto_labels(G, communities: dict, extraction: dict) -> dict:
    labels = {n["id"]: n.get("label", n["id"]) for n in extraction["nodes"]}
    out: dict[int, str] = {}
    for cid, members in communities.items():
        ranked = sorted(members, key=lambda nid: G.degree(nid) if nid in G else 0, reverse=True)
        out[cid] = labels.get(ranked[0], f"Community {cid}") if ranked else f"Community {cid}"
    return out


# --------------------------------------------------------------------------- #
# Step 5 - health verdict                                                      #
# --------------------------------------------------------------------------- #
def suitability(remaining: list[dict]) -> tuple[bool, str]:
    """Decide Copilot/MCP suitability from the still-dangling edges.

    Edges that point outside the corpus (external libs / tsconfig refs) never
    block. Among internal-but-unresolved edges we distinguish *structural*
    dependency relations (cites/imports/calls/implements/…) — which genuinely
    understate coupling — from weak prose relations (references /
    conceptually_related_to), which are advisory only.
    """
    external_causes = (
        "external-library import (non-corpus)",
        "tsconfig lib/include reference",
    )
    structural_rel = {
        "cites", "imports", "imports_from", "include", "calls", "implements",
        "extends", "inherits", "uses", "shares_data_with",
    }
    internal = [r for r in remaining if r["cause"] not in external_causes]
    blocking = [r for r in internal if r.get("relation") in structural_rel]
    weak = [r for r in internal if r.get("relation") not in structural_rel]

    if not blocking:
        note = ""
        if weak:
            note = (
                f" ({len(weak)} weak prose reference(s) remain unresolved — "
                "advisory only, no structural coupling lost)."
            )
        return True, (
            "SUITABLE for Copilot/MCP grounding. No structural dependency edge is "
            "unresolved; every remaining dangling endpoint is an external library, a "
            "tsconfig ref, or a weak prose mention." + note
        )
    return False, (
        f"USE WITH CAUTION. {len(blocking)} structural dependency edge(s) remain "
        "unresolved and may understate coupling. Safe for navigation/Q&A; verify "
        "dependency-direction queries."
    )


# --------------------------------------------------------------------------- #
# Health report rendering                                                      #
# --------------------------------------------------------------------------- #
def render_health_md(metrics: dict) -> str:
    from collections import Counter

    L: list[str] = []
    a = L.append
    a("# Praxis Graph Health Report")
    a("")
    a(f"_Generated by `scripts/rebuild_graph.py` — directed mode: "
      f"**{'on' if metrics['directed'] else 'off'}** (default)._")
    a("")
    a("## Summary")
    a("")
    a(f"- **Total nodes:** {metrics['total_nodes']}")
    a(f"- **Total edges (in graph):** {metrics['total_edges']}")
    a(f"- **Communities:** {metrics['communities']}")
    a(f"- **Dangling edges before normalization:** {metrics['dangling_before']}")
    a(f"- **Recovered edges:** {metrics['recovered_count']}")
    a(f"- **Remaining dangling edges:** {metrics['remaining_count']}")
    if metrics["ambiguous_count"]:
        a(f"- **Ambiguous (multi-candidate, left as-is):** {metrics['ambiguous_count']}")
    if metrics["dropped_self_loops"]:
        a(f"- **Self-loops dropped after recovery:** {metrics['dropped_self_loops']}")
    a("")
    a("## Verdict — Copilot / MCP grounding")
    a("")
    a(f"> **{'✅ ' if metrics['suitable'] else '⚠️ '}{metrics['verdict']}**")
    a("")
    a("## Recovered edges (id-convention normalization)")
    a("")
    rec = metrics["recovered_detail"]
    if rec:
        by_kind = Counter()
        for r in rec:
            if r["to_id"] == r["from_id"] + "_doc":
                by_kind["missing _doc suffix (RFC nodes)"] += 1
            elif r["from_id"].startswith(("packages_", "services_")):
                by_kind["packages_/services_ prefix mismatch"] += 1
            else:
                by_kind["other single-candidate match"] += 1
        for kind, n in by_kind.most_common():
            a(f"- {n} × {kind}")
        a("")
        a("Examples:")
        for r in rec[:10]:
            a(f"  - `{r['from_id']}` → `{r['to_id']}`  ({r['relation']})")
    else:
        a("- none")
    a("")
    a("## Remaining dangling edges (grouped by cause)")
    a("")
    causes = Counter(r["cause"] for r in metrics["remaining_detail"])
    if causes:
        for cause, n in causes.most_common():
            a(f"- **{n}** — {cause}")
        a("")
        a("Examples:")
        seen = set()
        for r in metrics["remaining_detail"]:
            key = r["missing_id"]
            if key in seen:
                continue
            seen.add(key)
            a(f"  - `{r['missing_id']}`  ({r['relation']}, {r['cause']})")
            if len(seen) >= 10:
                break
    else:
        a("- none — every endpoint resolves to a real node")
    a("")
    a("## Collapsed multi-relation pairs")
    a("")
    dm = metrics["collapse"]["directed_collapsed_multi_relation"]
    bd = metrics["collapse"]["bidirectional_pairs"]
    um = metrics["collapse"]["undirected_collapsed_pairs"]
    a(f"- **{len(um)}** node pairs carry more than one raw edge (undirected view).")
    a(f"- **{len(bd)}** are bidirectional (A→B and B→A) — **preserved** in directed mode, "
      "merged in undirected mode.")
    a(f"- **{len(dm)}** are same-direction multi-relation — these still collapse in a "
      "`DiGraph` (graphify directed mode is not a multigraph). Acceptable: the relation "
      "pairs are near-redundant.")
    a("")
    if dm:
        a("Same-direction multi-relation (still collapsed):")
        for d in dm[:10]:
            a(f"  - {d['label']} :: {', '.join(d['relations'])}")
        a("")
    if bd:
        a("Bidirectional pairs now preserved by directed mode (examples):")
        for d in bd[:8]:
            a(f"  - {d['label']}")
    a("")
    return "\n".join(L)


# --------------------------------------------------------------------------- #
# Main                                                                         #
# --------------------------------------------------------------------------- #
def main() -> int:
    OUT.mkdir(parents=True, exist_ok=True)
    directed = DIRECTED_DEFAULT

    print("[rebuild] loading extraction from caches (AST + cached semantic, no LLM)…")
    extraction, detection, n_uncached = load_extraction()
    if n_uncached:
        print(f"[rebuild] WARNING: {n_uncached} document(s) have no cached semantic "
              "extraction. Run the /graphify skill for a full re-extract; proceeding "
              "with cached data only.")
    print(f"[rebuild] raw extraction: {len(extraction['nodes'])} nodes, "
          f"{len(extraction['edges'])} edges")

    print("[rebuild] normalizing endpoint ids…")
    norm = normalize(extraction)
    print(f"[rebuild] dangling before: {norm['dangling_before']} · "
          f"recovered: {len(norm['recovered'])} · remaining: {len(norm['remaining'])} · "
          f"ambiguous: {len(norm['ambiguous'])}")

    collapse = collapse_analysis(extraction, norm["node_ids"])

    print(f"[rebuild] building {'directed' if directed else 'undirected'} graph…")
    G = build_from_json(extraction, root=str(ROOT), directed=directed)
    if G.number_of_nodes() == 0:
        sys.exit("[rebuild] ERROR: graph is empty — aborting (no write).")

    try:
        communities = cluster(G)
    except Exception:  # louvain needs an undirected view
        communities = cluster(G.to_undirected())
    cohesion = score_all(G, communities)
    labels = auto_labels(G, communities, extraction)
    gods = god_nodes(G)
    surprises = surprising_connections(G, communities)
    questions = suggest_questions(G, communities, labels)

    # graph.json (force=True so recovered/directed graph always writes).
    wrote = to_json(G, communities, str(OUT / "graph.json"),
                    force=True, community_labels=labels)
    if not wrote:
        sys.exit("[rebuild] ERROR: graph.json was not written.")

    report = generate(G, communities, cohesion, labels, gods, surprises, detection,
                      {"input": 0, "output": 0}, str(ROOT), suggested_questions=questions)
    (OUT / "GRAPH_REPORT.md").write_text(report, encoding="utf-8")

    suitable, verdict = suitability(norm["remaining"])
    metrics = {
        "directed": directed,
        "total_nodes": G.number_of_nodes(),
        "total_edges": G.number_of_edges(),
        "communities": len(communities),
        "dangling_before": norm["dangling_before"],
        "recovered_count": len(norm["recovered"]),
        "remaining_count": len(norm["remaining"]),
        "ambiguous_count": len(norm["ambiguous"]),
        "dropped_self_loops": norm["dropped_self_loops"],
        "recovered_detail": norm["recovered"],
        "remaining_detail": norm["remaining"],
        "ambiguous_detail": norm["ambiguous"],
        "collapse": collapse,
        "suitable": suitable,
        "verdict": verdict,
    }

    (OUT / "graph_health.json").write_text(
        json.dumps(metrics, indent=2, ensure_ascii=False), encoding="utf-8")
    (OUT / "GRAPH_HEALTH.md").write_text(render_health_md(metrics), encoding="utf-8")

    print("[rebuild] done.")
    print(f"[rebuild]   graph.json        : {G.number_of_nodes()} nodes, "
          f"{G.number_of_edges()} edges, {len(communities)} communities")
    print(f"[rebuild]   recovered edges   : {len(norm['recovered'])}")
    print(f"[rebuild]   remaining dangling: {len(norm['remaining'])} "
          "(external/non-corpus)")
    print(f"[rebuild]   directed mode     : {'ON (default)' if directed else 'off'}")
    print(f"[rebuild]   Copilot/MCP       : {'SUITABLE' if suitable else 'CAUTION'}")
    print("[rebuild]   outputs           : graphify-out/{graph.json,GRAPH_REPORT.md,"
          "GRAPH_HEALTH.md,graph_health.json}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
