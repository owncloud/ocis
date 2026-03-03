package droneextractor

type PipelineInfo struct {
	BuildURL        string          `json:"build_url"`
	Started         int64           `json:"started,omitempty"`         // Unix timestamp
	Finished        int64           `json:"finished,omitempty"`        // Unix timestamp
	DurationMinutes float64         `json:"duration_minutes,omitempty"` // Calculated: (finished - started) / 60
	PipelineStages  []PipelineStage `json:"pipeline_stages"`
}

type PipelineStage struct {
	StageNumber int           `json:"stage_number"`
	StageName   string        `json:"stage_name,omitempty"`
	Status      string        `json:"status,omitempty"`
	Steps       []PipelineStep `json:"steps"`
}

type PipelineStep struct {
	StepNumber int    `json:"step_number"`
	StepName   string `json:"step_name"`
	Status     string `json:"status"` // "success", "failure", "error", "running", etc.
	URL        string `json:"url,omitempty"`
	Logs       string `json:"logs,omitempty"` // Console logs content
}

type PipelineInfoSummary struct {
	BuildURL        string       `json:"build_url"`
	Started         int64        `json:"started,omitempty"`
	Finished        int64        `json:"finished,omitempty"`
	DurationMinutes float64      `json:"duration_minutes,omitempty"`
	Status          string       `json:"status,omitempty"` // Overall build status
	FailedStages    []FailedStage `json:"failed_stages,omitempty"`
}

type FailedStage struct {
	StageNumber int    `json:"stage_number"`
	StageName   string `json:"stage_name"`
	Status      string `json:"status"`
}
