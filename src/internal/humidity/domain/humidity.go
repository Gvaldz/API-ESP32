package domain

type Humidity struct {
	IDhumedad int `json:"idhumedad"`
	IDHamster int32 `json:"idhamster"`
	Humedad float64 `json:"humedad"`
	HoraRegistro string `json:"hora_registro"`
}