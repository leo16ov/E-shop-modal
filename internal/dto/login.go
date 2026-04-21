package dto

type UserLogin struct {
	ID       int    `json:"id"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Email    string `json:"email"`
	Rol      string `json:"rol"`
}

type LoginResponse struct {
	User  *UserLogin `json:"user"`
	Token string     `json:"token"`
}

type UserOAuth struct {
	ID       int    `json:"id"`
	Nombre   string `json:"nombre"`
	Email    string `json:"email"`
	Provider string `json:"provider"`
	Rol      string `json:"rol"`
}

type UserGoogle struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}
