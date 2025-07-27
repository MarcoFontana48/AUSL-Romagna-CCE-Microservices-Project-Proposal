package response

type HealthCheck struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

type Error struct {
	Error string `json:"error"`
}
