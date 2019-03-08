Simple library for working with Telegram bot API.
## Usage example
```go
package main

import (
	"log"
	ta "telegram-api/telegram-api"
	"time"
)

const token = "YOUR TOKEN"

func main() {
	client := ta.New(token)

	client.Command("start", func(m ta.Message) {
		client.SendMessage(m.Chat.ID, "Welcome!")
	})

	client.Regexp("[0-9]+", func(m ta.Message) {
		client.SendMessage(m.Chat.ID, "Regexp")
	})

	err := client.Pooling(time.Second * 2)

	if err != nil {
		log.Fatal(err)
	}
}
```