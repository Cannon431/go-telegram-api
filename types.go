package telegram_api

type Response struct {
	Ok          bool      `json:"ok"`
	Result      []*Result `json:"result"`
	ErrorCode   int       `json:"error_code"`
	Description string    `json:"description"`
	Url         string
}

type Result struct {
	UpdateID int `json:"update_id"`
	Message  `json:"message"`
}

type Message struct {
	MessageID int `json:"message_id"`
	From      `json:"from"`
	Chat      `json:"chat"`
	Date      int    `json:"date"`
	Text      string `json:"text"`
}

type From struct {
	ID           int    `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

type Chat struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}

func (c *Chat) IsPrivate() bool {
	return c.Type == "private"
}

func (c *Chat) IsGroup() bool {
	return c.Type == "group"
}

func (c *Chat) IsSuperGroup() bool {
	return c.Type == "supergroup"
}

func (c *Chat) IsChannel() bool {
	return c.Type == "channel"
}

type ReplyKeyboardMarkup struct {
	Keyboard        [][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool               `json:"resize_keyboard"`
	OneTimeKeyboard bool               `json:"one_time_keyboard"`
	Selective       bool               `json:"selective"`
}

type KeyboardButton struct {
	Text            string `json:"text"`
	RequestContact  bool   `json:"request_contact"`
	RequestLocation bool   `json:"request_location"`
}

func (kb *ReplyKeyboardMarkup) Row(row []KeyboardButton) *ReplyKeyboardMarkup {
	kb.Keyboard = append(kb.Keyboard, row)

	return kb
}
