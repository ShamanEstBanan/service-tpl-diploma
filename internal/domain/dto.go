package domain

type NewUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

const (
	OrderStatusNEW        = "NEW"
	OrderStatusPROCESSING = "PROCESSING"
	OrderStatusINVALID    = "INVALID"
	OrderStatusPROCESSED  = "PROCESSING"
)
