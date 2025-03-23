package domain

type Temperature struct {
	IDtemperatura int `json:"idtemperatura"`
	IDHamster int32 `json:"idhamster"`
	Temperatura float64 `json:"temperatura"`
	HoraRegistro string `json:"hora_registro"`
}