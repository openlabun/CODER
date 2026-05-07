class RetryableSubmissionUpdateError(Exception):
    """Raised when submission update can be retried later."""


class PermanentSubmissionUpdateError(Exception):
    """Raised when submission update should not be retried."""
