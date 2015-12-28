// Package gohue provides an interface to Philips Hue devices
package gohue

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Client represents the interface to the network of Hue devices
type Client struct {
	username, ip string
	Lights       []Light
}

// NewClient returns a pointer to a new Client. This instance does not connect to your network, yet.
// To fetch all available lights, call Connect() on the returned instance
func NewClient(username, ip string) *Client {
	c := Client{}
	c.ip = ip
	c.username = username
	c.Lights = make([]Light, 0)

	return &c
}

// Connect fetches all available lights in the network and stores them inside the Clients Lights list
func (c *Client) Connect() error {
	url := getAPIBaseURL(c.ip, c.username) + "lights"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	lightsData, _ := ioutil.ReadAll(res.Body)
	var jsonResponse jsonLightMap
	err = json.Unmarshal(lightsData, &jsonResponse)
	if err != nil {
		return err
	}

	for id, jsonLight := range jsonResponse {
		light := newLight(id)
		light.state.Brightness = jsonLight.Bri
		light.state.Hue = jsonLight.Hue
		light.state.Saturation = jsonLight.Sat
		light.state.On = jsonLight.On
		light.name = jsonLight.Name
		light.state.Xy = jsonLight.Xy
		light.client = c

		c.Lights = append(c.Lights, *light)
	}

	return nil
}

func getAPIBaseURL(ip, username string) string {
	return "http://" + ip + "/api/" + username + "/"
}
