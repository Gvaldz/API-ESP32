package domain

type Food struct {
	IDalimento int `json:"idalimento"`
	IDHamster int32 `json:"idhamster"`
	Alimento bool `json:"alimento"`
	Porcentaje float32 `json:"porcentaje"`
	HoraRegistro string `json:"hora_registro"`
}