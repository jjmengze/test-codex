package auth

type TokenPayload struct {
	ProductId  string `json:"prod_id"`
	CompanyId  string `json:"company_id"`
	ComputerId string `json:"computer_id"`
}
