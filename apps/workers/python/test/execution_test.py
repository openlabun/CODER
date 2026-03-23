import os
import sys
import unittest


# Allow running this test directly from the worker root.
CURRENT_DIR = os.path.dirname(__file__)
WORKER_ROOT = os.path.abspath(os.path.join(CURRENT_DIR, ".."))
if WORKER_ROOT not in sys.path:
	sys.path.insert(0, WORKER_ROOT)

from app.adapters.driven.docker_executor import DockerExecutor
from app.application.process_submission import ProcessSubmission
from app.domain.models import SubmissionResult


class FakeAPIClient:
	def __init__(self):
		self.last_result = None

	def update_submission(self, result: SubmissionResult) -> None:
		self.last_result = result


class ExecutionTest(unittest.TestCase):
	def test_simple_sum_code_returns_expected_value(self):
		executor = DockerExecutor()
		executor.ensure_image_cached()
		api_client = FakeAPIClient()
		use_case = ProcessSubmission(executor, api_client)

		submission = SubmissionResult(
			submission_id="sub-test-1",
			result_id="res-test-1",
			code="print(2 + 3)",
			status="PENDING",
			time_limit_ms=5000,
            memory_limit_mb=128,
			var_type=int,
			output=None,
			error=None,
			execution_time_ms=0,
		)

		use_case.execute(submission)

		self.assertEqual(submission.status, "SUCCESS")
		self.assertEqual(submission.output, 5)
		self.assertIsNone(submission.error)
		self.assertGreaterEqual(submission.execution_time_ms, 0)
		self.assertIsNotNone(api_client.last_result)
		self.assertEqual(api_client.last_result.output, 5)


if __name__ == "__main__":
	unittest.main(verbosity=2)
