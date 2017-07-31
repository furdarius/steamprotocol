package cmlist

import (
	"net/http"

	"encoding/json"

	"math/rand"

	"github.com/pkg/errors"
)

type CMList struct {
	httpCl         *http.Client
	serversList    []string
	websocketsList []string
}

func NewCMList(httpCl *http.Client) *CMList {
	return &CMList{
		httpCl: httpCl,
	}
}

// RefreshList refresh servers and websockets ips list
func (c *CMList) RefreshList() error {
	// TODO: Get cellID as param
	// 7 used randomly now
	resp, err := c.httpCl.Get("https://api.steampowered.com/ISteamDirectory/GetCMList/v1/?cellId=7")
	if err != nil {
		return errors.Wrap(err, "failed to get cm list from steam api")
	}

	if resp.StatusCode != 200 {
		return errors.New(http.StatusText(resp.StatusCode))
	}

	type Response struct {
		Response struct {
			Servers    []string `json:"serverlist"`
			WebSockets []string `json:"serverlist_websockets"`
			Result     int      `json:"result"`
			Message    string   `json:"message"`
		} `json:"response"`
	}

	var response Response
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return errors.Wrap(err, "failed to decode response")
	}

	resp.Body.Close()

	c.serversList = response.Response.Servers
	c.websocketsList = response.Response.WebSockets

	return nil
}

// GetRandomServer refresh servers list, if empty and
// return random server ip from list
func (c *CMList) GetRandomServer() (string, error) {
	if c.serversList == nil {
		err := c.RefreshList()
		if err != nil {
			return "", errors.Wrap(err, "failed to refresh servers list")
		}
	}

	n := len(c.serversList)
	return c.serversList[rand.Intn(n)], nil
}
