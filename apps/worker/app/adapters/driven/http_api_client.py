import requests
import time
from app.application.ports.api_client import APIClient
from app.domain.models import SubmissionResult
from app.domain.errors import RetryableSubmissionUpdateError, PermanentSubmissionUpdateError
from config import API_BASE_URL, API_TOKEN


class HTTPAPIClient(APIClient):

    def update_submission(self, result: SubmissionResult) -> None:
        # PATCH /submissions/results/{resultId} (internal worker endpoint)
        url = f"{API_BASE_URL}/submissions/results/{result.result_id}"
        print(f"Updating submission result to API at {url} with status {result.status}", flush=True)

        payload = {
            "status": result.status,
            "timeExecution": result.execution_time_ms,
            "output": result.output,
            "error": result.error
        }

        headers = {
            "WorkerKey": API_TOKEN,
            "Content-Type": "application/json"
        }

        for i in range(3):
            try:
                response = requests.patch(url, json=payload, headers=headers, timeout=15)

                if 200 <= response.status_code < 300:
                    print(f"Successfully updated submission of id {result.result_id} (status: {result.status})", flush=True)
                    return
                if 400 <= response.status_code < 500:
                    raise PermanentSubmissionUpdateError(
                        f"Permanent API error {response.status_code}: {response.text}"
                    )
                print(f"Failed to update submission of id {result.result_id} (status: {response.status_code}, body: {response.text})", flush=True)
            except Exception as e:
                print(f"Error updating submission of id {result.result_id}: {e}", flush=True)
                if isinstance(e, PermanentSubmissionUpdateError):
                    raise
                pass
            time.sleep(i * 2)
        raise RetryableSubmissionUpdateError("Failed to update submission after retries")