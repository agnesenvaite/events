package status

type ResponseStatus struct {
	Status string `json:"status" example:"OK"`
}

type ResponseDetails struct {
	Status map[string]string `json:"status"  swaggertype:"object,string" example:"mysql: OK"`
}
