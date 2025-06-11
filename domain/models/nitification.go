package models


type Notification struct {
	Payload   KafkaEventPayload `json:"payload"`
}

type KafkaEventPayload struct {
	UserID       int    `json:"userId"`
	Amount       int    `json:"amount"`
	FromCurrency string `json:"FromCurrency"`
	ToCurrency   string `json:"ToCurrency"`
}
