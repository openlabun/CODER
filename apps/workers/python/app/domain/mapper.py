from app.domain.models import SubmissionResult

def MapSubmissionResult(data: dict) -> SubmissionResult:
    return SubmissionResult(
        submission_id=data["submission_id"],
        code=data["code"],
        input=data["input"],
        result_id=data["result_id"],
        time_limit_ms=data["time_limit_ms"],
        memory_limit_mb=data["memory_limit_mb"],
        var_type=standarizeTypes(data["type"]),
        status=data["status"],
        output=None,
        error=None,
        execution_time_ms=0
    )

def standarizeTypes(var_type: str) -> type:
    if var_type in ["int", "integer"]:
        return int
    elif var_type in ["str", "string"]:
        return str
    elif var_type in ["float", "double"]:
        return float
    elif var_type in ["bool", "boolean"]:
        return bool
    else:
        return str