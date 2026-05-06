from abc import ABC, abstractmethod
from app.domain.models import SubmissionResult


class APIClient(ABC):

    @abstractmethod
    def update_submission(self, result: SubmissionResult) -> None:
        pass