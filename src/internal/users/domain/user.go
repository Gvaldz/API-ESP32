package domain

type User struct {
	IdUsuario int32  `json:"iduser"`
	Nombre   string `json:"nombre"`
	Correo 	string `json:"correo"`
	Contrasena string `json:"contrasena"`
	Tipo string `json:"tipo"`
	FCMToken string `json:"fcm_token"`
}