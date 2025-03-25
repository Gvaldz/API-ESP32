package domain

type Motion struct {
	IDMovimiento int `json:"idmovimiento"`
	IDHamster int32 `json:"idhamster"`
	Movimiento bool `json:"movimiento"`
	HoraRegistro string `json:"hora_registro"`
}