PYTHON ?= python3

.PHONY: verify run-api run-worker run-chiefly run-telegram

verify:
	$(PYTHON) -m compileall apps services packages scripts
	$(PYTHON) scripts/verify/verify_db.py
	$(PYTHON) scripts/verify/verify_queue.py
	$(PYTHON) scripts/verify/verify_llm_router.py
	$(PYTHON) scripts/verify/verify_json_sanitizer.py
	$(PYTHON) scripts/verify/verify_work_item_pipeline.py
	$(PYTHON) scripts/verify/verify_upwork_pipeline.py
	$(PYTHON) scripts/verify/verify_e2e.py

run-api:
	$(PYTHON) services/api/main.py

run-worker:
	$(PYTHON) services/worker/main.py

run-chiefly:
	$(PYTHON) apps/chiefly/main.py

run-telegram:
	$(PYTHON) apps/telegram/main.py
