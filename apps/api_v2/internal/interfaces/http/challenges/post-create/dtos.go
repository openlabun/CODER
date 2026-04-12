package postcreate

import examDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"

type RequestDTO struct {
	Title             string                   `json:"title"`
	Description       string                   `json:"description"`
	Tags              []string                 `json:"tags"`
	Status            string                   `json:"status"`
	Difficulty        string                   `json:"difficulty"`
	WorkerTimeLimit   int                      `json:"worker_time_limit"`
	WorkerMemoryLimit int                      `json:"worker_memory_limit"`
	CodeTemplates     map[string]string        `json:"code_templates"`
	InputVariables    []examDtos.IOVariableDTO `json:"input_variables"`
	OutputVariable    examDtos.IOVariableDTO   `json:"output_variable"`
	Constraints       string                   `json:"constraints"`
}
