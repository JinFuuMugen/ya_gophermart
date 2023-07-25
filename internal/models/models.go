package models

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Order struct {
	Dateadd string  `json:"processed_at"`
	Number  string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}
