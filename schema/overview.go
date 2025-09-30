package schema

type OverviewResult struct {
	Id                           uint    `json:"id"`
	SloName                      string  `json:"slo_name"`
	TargetSlo                    string  `json:"target_slo"`
	CurrentSlo                   float32 `json:"current_slo"`
	RemErrBudget                 float32 `json:"remaining_err_budget"`
	Num_incidents                int     `json:"num_incidents"`
	Num_incidents_false_positive int     `json:"num_incidents_false_positive"`
}
