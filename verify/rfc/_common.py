"""Shared, read-only helpers for the Phase 0 RFC hygiene tools.

Standard library only. These helpers never write to or mutate RFC files.

This module is internal scaffolding for `verify/rfc/`. It exists to avoid duplicating
RFC-parsing logic across the individual checkers. It does not encode any architectural
decision and does not interpret Canonical Object / Artifact semantics.
"""

from __future__ import annotations

import re
from pathlib import Path

# --- Patterns ---------------------------------------------------------------

# RFC source file name, e.g. "052-work-space.md" -> number "052".
RFC_FILENAME_RE = re.compile(r"^(\d{3})-.*\.md$")

# Top-level RFC heading, e.g. "# RFC-052 Work Space" or "# RFC-052: Work Space".
RFC_H1_RE = re.compile(r"^#\s+RFC-(\d{3})\s*[:\-\u2014]?\s*(.*?)\s*$")

# Inline reference to an RFC, with an optional parenthetical title claim, e.g.
#   "RFC-050"                      -> ("050", None)
#   "RFC-050 (Business Space)"     -> ("050", "Business Space")
# A parenthetical only counts when "(" immediately follows the number (allowing spaces).
RFC_REF_RE = re.compile(r"RFC-(\d{3})(?:\s*\(([^)]*)\))?")

# README numbering-table row, e.g. "| 052  | Work Space | Draft | ... |".
README_ROW_RE = re.compile(r"^\|\s*(\d{3})\s*\|\s*([^|]+?)\s*\|")

# Markdown heading of any level.
HEADING_RE = re.compile(r"^(#{1,6})\s+(.*\S)\s*$")

# Parentheticals that are clearly NOT title claims.
_TITLE_STOPWORDS = {
    "this rfc",
    "this document",
    "current",
    "draft",
    "accepted",
    "tbd",
    "wip",
    "see above",
    "note",
    "initial draft",
}


# --- Repo / file discovery --------------------------------------------------

def find_repo_root(start: Path | None = None) -> Path:
    """Walk upward from ``start`` until a directory containing ``rfcs/`` is found."""
    here = (start or Path(__file__)).resolve()
    for candidate in [here, *here.parents]:
        if (candidate / "rfcs").is_dir():
            return candidate
    raise FileNotFoundError("Could not locate repo root containing an 'rfcs/' directory.")


def rfcs_dir(root: Path) -> Path:
    return root / "rfcs"


def list_rfc_files(rfcs_path: Path) -> dict[str, Path]:
    """Map zero-padded RFC number -> file path for every ``NNN-*.md`` in ``rfcs_path``."""
    out: dict[str, Path] = {}
    for path in sorted(rfcs_path.glob("*.md")):
        match = RFC_FILENAME_RE.match(path.name)
        if match:
            out[match.group(1)] = path
    return out


def read_text(path: Path) -> str:
    return path.read_text(encoding="utf-8")


# --- Title / heading parsing -----------------------------------------------

def find_h1_rfc_headings(text: str) -> list[tuple[int, str, str]]:
    """Return ``(line_number, rfc_number, title)`` for each top-level RFC heading."""
    results: list[tuple[int, str, str]] = []
    for lineno, line in enumerate(text.splitlines(), start=1):
        match = RFC_H1_RE.match(line)
        if match:
            results.append((lineno, match.group(1), match.group(2).strip()))
    return results


def first_h1_title_map(files_map: dict[str, Path]) -> dict[str, str]:
    """Map RFC number -> the title of its first top-level RFC heading (or "")."""
    out: dict[str, str] = {}
    for number, path in files_map.items():
        headings = find_h1_rfc_headings(read_text(path))
        out[number] = headings[0][2] if headings else ""
    return out


def parse_readme_titles(readme_path: Path) -> dict[str, str]:
    """Parse the README RFC numbering table into ``{number: title}``."""
    out: dict[str, str] = {}
    if not readme_path.exists():
        return out
    for line in read_text(readme_path).splitlines():
        match = README_ROW_RE.match(line)
        if match:
            out[match.group(1)] = normalize_title(match.group(2))
    return out


# --- Title normalization ----------------------------------------------------

def normalize_title(s: str) -> str:
    return re.sub(r"\s+", " ", s).strip()


def title_key(s: str) -> str:
    """Case-insensitive comparison key for titles."""
    return normalize_title(s).lower().rstrip(".")


def is_probable_title(parenthetical: str) -> bool:
    """Heuristic: does a parenthetical after an RFC reference look like a title claim?

    Rejects descriptive/enumerated parentheticals (commas, "etc.") since RFC titles in this
    corpus contain neither; this avoids flagging range annotations like
    "(Core, Data, Product, etc.)" as title claims.
    """
    s = parenthetical.strip()
    if not s:
        return False
    if not (s[0].isalpha() and s[0].isupper()):
        return False
    if s.lower() in _TITLE_STOPWORDS:
        return False
    if "," in s or re.search(r"\betc\b", s, re.IGNORECASE):
        return False
    return bool(re.fullmatch(r"[A-Za-z0-9 &/()'\-]+", s))
