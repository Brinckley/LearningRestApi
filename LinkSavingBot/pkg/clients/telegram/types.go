package telegram

type Update struct {
	UpdateID int              `json:"update_id"`
	Message  *IncomingMessage `json:"message"`
}

type UpdateResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type IncomingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

type From struct {
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"ID"`
}
