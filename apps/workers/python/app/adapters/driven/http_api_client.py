import requests
import time
from app.application.ports.api_client import APIClient
from app.domain.models import SubmissionResult
from config import API_BASE_URL, API_TOKEN


class HTTPAPIClient(APIClient):

    def update_submission(self, result: SubmissionResult) -> None:
        url = f"{API_BASE_URL}/submissions/{result.submission_id}"

        payload = {
            "status": result.status,
            "output": result.output,
            "error": result.error,
            "execution_time_ms": result.execution_time_ms
        }

        headers = {
            "Authorization": f"Bearer {API_TOKEN}",
            "Content-Type": "application/json"
        }

        for i in range(5):
            try:
                response = requests.patch(url, json=payload, headers=headers, timeout=5)

                if response.status_code < 300:
                    return

            except Exception:
                pass

            time.sleep(i * 2)

        raise Exception("Failed to update submission after retries")