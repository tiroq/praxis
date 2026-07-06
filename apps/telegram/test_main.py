from __future__ import annotations

import unittest
from types import SimpleNamespace
from unittest.mock import AsyncMock

from apps.telegram.main import RuntimeState, handle_output_message, render_reply_text


class ReplyPipelineTests(unittest.IsolatedAsyncioTestCase):
    async def test_handle_output_message_uses_output_metadata_chat_id(self) -> None:
        app = SimpleNamespace(bot=SimpleNamespace(send_message=AsyncMock()))
        state = RuntimeState()
        pending_chat_by_event_id: dict[str, str] = {}
        output = {
            "input_event_id": "evt_1",
            "status": "ok",
            "metadata": {"chat_id": "42"},
            "result": {"decision": {"outcome": "approve"}, "actions": [{"id": "act_1"}]},
        }

        await handle_output_message(app, state, pending_chat_by_event_id, output)

        app.bot.send_message.assert_awaited_once_with(
            chat_id=42,
            text="Decision: approve\nActions: 1",
        )
        self.assertEqual(state.snapshot()["metrics"]["output_messages_total"], 1)
        self.assertEqual(state.snapshot()["metrics"]["replies_sent_total"], 1)

    async def test_handle_output_message_falls_back_to_pending_chat_id(self) -> None:
        app = SimpleNamespace(bot=SimpleNamespace(send_message=AsyncMock()))
        state = RuntimeState()
        pending_chat_by_event_id = {"evt_2": "99"}
        output = {
            "input_event_id": "evt_2",
            "status": "error",
            "error": "kernel blew up",
        }

        await handle_output_message(app, state, pending_chat_by_event_id, output)

        app.bot.send_message.assert_awaited_once_with(
            chat_id=99,
            text="Praxis error: kernel blew up",
        )
        self.assertNotIn("evt_2", pending_chat_by_event_id)
        self.assertEqual(state.snapshot()["metrics"]["replies_sent_total"], 1)

    async def test_handle_output_message_ignores_outputs_without_chat_id(self) -> None:
        app = SimpleNamespace(bot=SimpleNamespace(send_message=AsyncMock()))
        state = RuntimeState()
        pending_chat_by_event_id: dict[str, str] = {}
        output = {"input_event_id": "evt_3", "status": "ok"}

        await handle_output_message(app, state, pending_chat_by_event_id, output)

        app.bot.send_message.assert_not_awaited()
        self.assertEqual(state.snapshot()["metrics"]["output_messages_total"], 1)
        self.assertEqual(state.snapshot()["metrics"]["replies_sent_total"], 0)

    def test_render_reply_text(self) -> None:
        self.assertEqual(
            render_reply_text({
                "status": "ok",
                "assistant_reply": "Hello from LLM",
                "result": {"decision": {"outcome": "approve"}, "actions": [{"id": "a1"}]},
            }),
            "Hello from LLM",
        )
        self.assertEqual(
            render_reply_text({
                "status": "ok",
                "result": {"decision": {"outcome": "approve"}, "actions": [{"id": "a1"}]}
            }),
            "Decision: approve\nActions: 1",
        )
        self.assertEqual(
            render_reply_text({
                "status": "ok",
                "result": {"decision": {"outcome": "reject"}, "actions": [{"id": "a1"}, {"id": "a2"}]}
            }),
            "Decision: reject\nActions: 2",
        )
        self.assertEqual(
            render_reply_text({"status": "error", "error": "boom"}),
            "Praxis error: boom",
        )
        self.assertEqual(
            render_reply_text({"status": "ok", "result": {"decision": {"outcome": "hold"}}}),
            "Decision: hold\nActions: 0",
        )


if __name__ == "__main__":
    unittest.main()
