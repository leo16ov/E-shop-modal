package dto

type Payload struct {
	Action string `json:"action"`
	Type   string `json:"type"`
	Data   struct {
		ID int64 `json:"id"`
	} `json:"data"`
}
