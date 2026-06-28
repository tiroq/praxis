from __future__ import annotations

from dataclasses import dataclass


@dataclass(slots=True)
class AgentResult:
    agent_name: str
    summary: str
    status: str = "scaffold"


class BaseAgent:
    name = "base"

    def run(self, prompt: str) -> AgentResult:
        return AgentResult(agent_name=self.name, summary=f"TODO: implement {self.name} agent for: {prompt}")
