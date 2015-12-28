package gohue

import (
	"net/http"
	"strings"
)

// Light represents one light in the network of a Client.
type Light struct {
	id     string
	name   string
	State  State
	client *Client
}

// State represents the current or desired state of a light
type State struct {
	Hue, Brightness, Saturation int
	On                          bool
	Xy                          []float64
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
