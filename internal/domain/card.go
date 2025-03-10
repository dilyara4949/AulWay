package domain

type CardDetails struct {
	Number   string `json:"number"`
	ExpMonth int64  `json:"exp_month"`
	ExpYear  int64  `json:"exp_year"`
	CVC      string `json:"cvc"`
}
