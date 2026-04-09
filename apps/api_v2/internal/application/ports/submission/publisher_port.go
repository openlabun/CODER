package submission_ports

import (
	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"
)

type SubmissionPublisherPort interface {
	PublishSubmission(dto dtos.SubmissionResultPublishedDTO) error
}