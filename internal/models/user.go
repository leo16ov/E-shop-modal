package models

type User struct {
	ID         int    `json:"id"`
	Nombre     string `json:"nombre"`
	Apellido   string `json:"apellido"`
	Email      string `json:"email"`
	Contrasena string `json:"contrasena"`
	Rol        string `json:"rol"`
	DNI        int    `json:"dni"`
	Telefono   string `json:"telefono"`
	Provider   string `json:"provider"`
}
