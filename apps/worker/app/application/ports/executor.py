from abc import ABC, abstractmethod
from app.domain.models import SubmissionResult


class Executor(ABC):

    @abstractmethod
    async def execute(self, submission: SubmissionResult) -> SubmissionResult:
        pass