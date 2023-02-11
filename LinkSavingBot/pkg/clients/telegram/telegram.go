package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const ( // method names for requests (taken from tg)
	sendMessageMethod = "sendMessage"
	getUpdatesMethod  = "getUpdates"
)

type Client struct { // basic client struct
	host     string
	basePath string
	client   http.Client
}

func NewClient(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}                      // configuration for every sent message
	q.Add("chat_id", strconv.Itoa(chatID)) // where to send msg
	q.Add("text", text)                    // what to send in msg

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return fmt.Errorf("send message error : %s", err)
	}

	return nil
}

func (c *Client) Updates(offset int, limit int) ([]Update, error) {
	q := url.Values{} // getting updates from the chats
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, fmt.Errorf("get updates error : %s", err)
	}

	var res UpdateResponse
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal data : %s", err)
	}

	return res.Result, nil
}

func (c *Client) doRequest(method string, query url.Values) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("can't creat a get request : %s", err.Error())
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't do a get request : %s", err.Error())
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read the responce body : %s", err.Error())
	}

	return body, nil
}
