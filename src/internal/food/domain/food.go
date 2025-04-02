package domain

type Food struct {
	IDalimento int `json:"idalimento"`
	IDHamster string `json:"idhamster"`
	Alimento int `json:"alimento"`
	Porcentaje float32 `json:"porcentaje"`
	HoraRegistro string `json:"hora_registro"`
}