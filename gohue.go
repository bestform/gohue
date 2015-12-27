// Package gohue provides an interface to Philips Hue devices
package gohue

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// Light represents one light in the network of a Client.
type Light struct {
	id     string
	name   string
	State  State
	client *Client
}

// Client represents the interface to the network of Hue devices
type Client struct {
	username, ip string
	Lights       []Light
}

// State represents the current or desired state of a light
type State struct {
	Hue, Brightness, Saturation int
	On                          bool
	Xy                          []float64
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
		light.State.Brightness = jsonLight.Bri
		light.State.Hue = jsonLight.Hue
		light.State.Saturation = jsonLight.Sat
		light.State.On = jsonLight.On
		light.name = jsonLight.Name
		light.State.Xy = jsonLight.Xy
		light.client = c

		c.Lights = append(c.Lights, *light)
	}

	return nil
}

func newLight(id string) *Light {
	l := Light{}
	l.id = id
	l.State = State{10000, 254, 254, true, []float64{0.0, 0.0}}

	return &l
}

// UpdateState sends the current state of the light to its hardware counterpart.
func (light *Light) UpdateState(useXY bool) error {
	url := getAPIBaseURL(light.client.ip, light.client.username) + "lights/" + light.id + "/state"
	payload := getPayloadFromState(light.State, useXY)
	payloadReader := strings.NewReader(payload)
	req, err := http.NewRequest("PUT", url, payloadReader)
	if err != nil {
		return err
	}
	client := http.Client{}
	client.Do(req)

	return nil
}

func getAPIBaseURL(ip, username string) string {
	return "http://" + ip + "/api/" + username + "/"
}

func getPayloadFromState(s State, useXY bool) string {
	var onState string
	if s.On {
		onState = "true"
	} else {
		onState = "false"
	}
	var colorValue string
	if useXY {
		x := strconv.FormatFloat(s.Xy[0], 'f', 4, 64)
		y := strconv.FormatFloat(s.Xy[1], 'f', 4, 64)
		colorValue = "\"xy\":[" + x + "," + y + "]"
	} else {
		colorValue = "\"hue\":" + strconv.Itoa(s.Hue)
	}
	return "{\"on\":" + onState + ", \"sat\":" + strconv.Itoa(s.Saturation) + ", \"bri\":" + strconv.Itoa(s.Brightness) + "," + colorValue + "}"
}
