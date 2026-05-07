import time
from app.application.ports.executor import Executor
from app.application.ports.api_client import APIClient
from app.domain.models import SubmissionResult


class ProcessSubmission:

    def __init__(self, executor: Executor, api_client: APIClient):
        self.executor = executor
        self.api_client = api_client

    async def execute(self, submission: SubmissionResult):
        submission.status = "running"
        await self.api_client.update_submission(submission)

        start = time.time()

        result = await self.executor.execute(submission)

        result.execution_time_ms = int((time.time() - start) * 1000)

        await self.api_client.update_submission(result)