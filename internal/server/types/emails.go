package server

type Email struct {
	Id        string `json:"id"`
	MessageID string `json:"messageID"`
	Subject   string `json:"subject"`
	XFrom     string `json:"xFrom"`
	From      string `json:"from"`
	To        string `json:"to"`
	Date      string `json:"date"`
	Body      string `json:"body"`
}

type EmailsResponse struct {
	Emails []Email `json:"emails"`
	Total  int     `json:"total"`
}
