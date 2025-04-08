package domain

type Humidity struct {
	IDhumedad int `json:"idhumedad"`
	IDHamster string `json:"idhamster"`
	Humedad float64 `json:"humedad"`
	HoraRegistro string `json:"hora_registro"`
}