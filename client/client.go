package client

type Client struct {
	Id        string `json:"id" db:"id"`
	FirstName string `json:"firstName" db:"first_name"`
	LastName  string `json:"lastName" db:"last_name"`
	Login     string `json:"login" db:"login"`
	ChatId    string `json:"chatId" db:"chat_id"`
	Active    bool   `json:"active" db:"active"`
}
