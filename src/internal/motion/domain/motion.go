package domain

type Motion struct {
	IDMovimiento int `json:"idmovimiento"`
	IDHamster string `json:"idhamster"`
	Movimiento bool `json:"movimiento"`
	HoraRegistro string `json:"hora_registro"`
}