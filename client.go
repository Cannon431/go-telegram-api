package telegram_api

import (
	"log"
	"regexp"
	"time"
)

type Client struct {
	API
	handlers
	updates   []Result
	delivered []int
}

type Handler struct {
	_type   string
	trigger string
	handler handlerFunc
}

type handlerFunc func(m Message)
type handlers []*Handler

func New(token string) *Client {
	return &Client{API: API{token: token}}
}

func (c *Client) AddHandler(_type string, trigger string, handler handlerFunc) {
	c.handlers = append(c.handlers, &Handler{
		_type:   _type,
		trigger: trigger,
		handler: handler,
	})
}

func (c *Client) Text(text string, handler handlerFunc) {
	c.AddHandler("text", text, handler)
}

func (c *Client) Command(command string, handler handlerFunc) {
	if string(command[0]) != "/" {
		command = "/" + command
	}

	c.Text(command, handler)
}

func (c *Client) Regexp(regexp string, handler handlerFunc) {
	c.AddHandler("regexp", regexp, handler)
}

func (c *Client) Pooling(timeout time.Duration, params ...Params) error {
	for {
		updates, err := c.GetUpdates(params...)
		var unhandledUpdates [] *Result

		if err != nil {
			return err
		}

		for _, update := range updates.Result {
			if !intInSlice(update.UpdateID, c.delivered) {
				unhandledUpdates = append(unhandledUpdates, update)
			}
		}

		mid := len(unhandledUpdates) / 2
		thirstPiece := unhandledUpdates[:mid]
		secondPiece := unhandledUpdates[mid:]

		go c.handle(thirstPiece)
		go c.handle(secondPiece)

		time.Sleep(timeout)
	}
}

func (c *Client) MustPooling(timeout time.Duration, params ...Params) {
	err := c.Pooling(timeout, params...)

	if err != nil {
		panic(err)
	}
}

func (c *Client) handle(updates []*Result) {
	for _, update := range updates {
		c.delivered = append(c.delivered, update.UpdateID)
		handled := false

		for _, handler := range c.handlers {
			switch handler._type {
			case "text":
				if update.Message.Text != "" && update.Message.Text == handler.trigger {
					go handler.handler(update.Message)
					handled = true
				}
			case "regexp":
				if update.Message.Text != "" && regexp.MustCompile(handler.trigger).MatchString(update.Message.Text) {
					go handler.handler(update.Message)
					handled = true
				}
			}
		}

		if handled {
			log.Printf("New handled message from id %d\n", update.From.ID)
		} else {
			log.Printf("New unhandled message from id %d\n", update.From.ID)
		}
	}
}
