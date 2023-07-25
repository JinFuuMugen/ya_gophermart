package models

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Order struct {
	Dateadd string `json:"processed_at"`
	Number  int    `json:"order"`
	Status  string `json:"status"`
	Accrual int    `json:"accrual"`
}
