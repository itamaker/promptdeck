package app

type PromptCandidate struct {
	Index  string         `json:"index"`
	Vars   map[string]any `json:"vars"`
	Prompt string         `json:"prompt"`
}

type CandidateScore struct {
	Index string             `json:"index"`
	Score float64            `json:"score"`
	Notes string             `json:"notes,omitempty"`
	Extra map[string]float64 `json:"extra,omitempty"`
}

type RankedCandidate struct {
	Index  string         `json:"index"`
	Score  float64        `json:"score"`
	Vars   map[string]any `json:"vars"`
	Prompt string         `json:"prompt"`
	Notes  string         `json:"notes,omitempty"`
}

type FactorEffect struct {
	Variable string  `json:"variable"`
	Value    string  `json:"value"`
	Count    int     `json:"count"`
	AvgScore float64 `json:"avg_score"`
}

type OptimizationReport struct {
	Candidates    int               `json:"candidates"`
	Scored        int               `json:"scored"`
	Best          *RankedCandidate  `json:"best,omitempty"`
	Top           []RankedCandidate `json:"top,omitempty"`
	FactorEffects []FactorEffect    `json:"factor_effects,omitempty"`
}
