package models

type Product struct {
	ID       int     `json:"id"`
	Tipo     string  `json:"tipo"` //Se aclara si es remera, pantalon, vestido y que tipo
	Precio   float32 `json:"precio"`
	Cantidad int     `json:"cantidad"`
	Talles   string  `json:"talles"`  //[]int
	Colores  string  `json:"colores"` //[]string
}
