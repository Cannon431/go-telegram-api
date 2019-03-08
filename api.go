package telegram_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type API struct {
	token string
}

type Params map[string]string

func (api *API) getUrl(method string, params Params) (string, error) {
	s := "https://api.telegram.org/bot" + api.token + "/" + method
	u, err := url.Parse(s)

	if err != nil {
		return "", err
	}

	query := u.Query()

	for key, value := range params {
		query.Add(key, value)
	}

	s = s + "?" + query.Encode()

	return s, nil
}

func (api *API) Request(method string, params ...Params) (*Response, error) {
	var r Response
	var p Params

	if len(params) == 0 {
		p = Params{}
	} else {
		p = params[0]
	}

	u, err := api.getUrl(method, p)
	r.Url = u

	if err != nil {
		return &r, err
	}

	res, err := http.Get(u)

	if err != nil {
		return &r, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return &r, err
	}

	err = json.Unmarshal(body, &r)

	if err != nil {
		return &r, err
	}

	if !r.Ok {
		return &r, errors.New(fmt.Sprintf("[%d] %s", r.ErrorCode, r.Description))
	}

	return &r, nil
}

func (api *API) MustRequest(method string, params ...Params) *Response {
	result, err := api.Request(method, params...)

	if err != nil {
		panic(err)
	}

	return result
}

func (api *API) GetUpdates(params ...Params) (*Response, error) {
	resp, err := api.Request("getUpdates", params...)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (api *API) ClearUpdates() error {
	updates, err := api.GetUpdates()

	if err != nil {
		return err
	}

	result := updates.Result

	if len(result) > 0 {
		newOffset := result[len(result) - 1].UpdateID + 1
		_, err := api.GetUpdates(Params{"offset": strconv.Itoa(newOffset)})

		if err != nil {
			return err
		}
	}

	return nil
}

func (api *API) SendMessage(chatId int, text string, params ...Params) (*Response, error) {
	var p Params

	if len(params) == 0 {
		p = Params{}
	} else {
		p = params[0]
	}

	p["chat_id"] = strconv.Itoa(chatId)
	p["text"] = text

	return api.Request("sendMessage", p)
}
