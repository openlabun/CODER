from dataclasses import dataclass
from typing import Optional


@dataclass
class SubmissionResult:
    submission_id: str
    code: str
    result_id: str
    status: str
    var_type: type
    time_limit_ms: int
    memory_limit_mb: int
    output: Optional[str]
    error: Optional[str]
    execution_time_ms: int
