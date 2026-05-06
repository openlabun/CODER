import asyncio
import aiohttp
from app.application.ports.api_client import APIClient
from app.domain.models import SubmissionResult
from app.domain.errors import RetryableSubmissionUpdateError, PermanentSubmissionUpdateError
from config import API_BASE_URL, API_TOKEN


class HTTPAPIClient(APIClient):

    async def update_submission(self, result: SubmissionResult) -> None:
        url = f"{API_BASE_URL}/submissions/results/{result.result_id}"
        print(f"Updating submission result to API at {url} with status {result.status}", flush=True)

        payload = {
            "status": result.status,
            "timeExecution": result.execution_time_ms,
            "output": result.output,
            "error": result.error
        }

        headers = {"WorkerKey": API_TOKEN}

        for i in range(3):
            try:
                async with aiohttp.ClientSession() as session:
                    async with session.patch(
                        url,
                        json=payload,
                        headers=headers,
                        timeout=aiohttp.ClientTimeout(total=15),
                    ) as response:
                        if 200 <= response.status < 300:
                            print(f"Successfully updated submission of id {result.result_id} (status: {result.status})", flush=True)
                            return
                        if 400 <= response.status < 500:
                            text = await response.text()
                            raise PermanentSubmissionUpdateError(
                                f"Permanent API error {response.status}: {text}"
                            )
                        text = await response.text()
                        print(f"Failed to update submission of id {result.result_id} (status: {response.status}, body: {text})", flush=True)
            except Exception as e:
                print(f"Error updating submission of id {result.result_id}: {e}", flush=True)
                if isinstance(e, PermanentSubmissionUpdateError):
                    raise
            await asyncio.sleep(i * 2)
        raise RetryableSubmissionUpdateError("Failed to update submission after retries")
