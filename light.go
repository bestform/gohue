package gohue

import (
	"net/http"
	"strings"
)

type updateType int

const (
	updateColorHue updateType = iota
	updateColorXy
	updateBrightness
	updateOnOff
	updateSaturation
)

// Light represents one light in the network of a Client.
type Light struct {
	id     string
	name   string
	state  state
	client *Client
}

// State represents the current or desired state of a light
type state struct {
	Hue, Brightness, Saturation int
	On                          bool
	Xy                          []float64
}

func newLight(id string) *Light {
	l := Light{}
	l.id = id
	l.state = state{10000, 254, 254, true, []float64{0.0, 0.0}}

	return &l
}

// SetColorHue sets the hue value. It has an immediate effect.
func (l *Light) SetColorHue(hue int) error {
	l.state.Hue = hue
	return l.updateState(updateColorHue)
}

// SetColorRGB converts the given RGB values to the corresponding x/y values and sets the color based on those. It has an immediate effect.
func (l *Light) SetColorRGB(r, g, b int) error {
	x, y := convertRGBToXY(r, g, b)
	l.state.Xy[0] = x
	l.state.Xy[1] = y
	return l.updateState(updateColorXy)
}

// SetBrightness sets the brightness value. It has an immediate effect.
func (l *Light) SetBrightness(bri int) error {
	l.state.Brightness = bri
	return l.updateState(updateBrightness)
}

// SetSaturation sets the saturation value. It has an immediate effect.
func (l *Light) SetSaturation(sat int) error {
	l.state.Saturation = sat
	return l.updateState(updateSaturation)
}

// SwitchOn sets the on state to true. It has an immediate effect.
func (l *Light) SwitchOn() error {
	l.state.On = true
	return l.updateState(updateOnOff)
}

// SwitchOff sets the on state to false. It has an immediate effect.
func (l *Light) SwitchOff() error {
	l.state.On = false
	return l.updateState(updateOnOff)
}

// UpdateState sends the current state of the light to its hardware counterpart.
func (light *Light) updateState(uType updateType) error {
	url := getAPIBaseURL(light.client.ip, light.client.username) + "lights/" + light.id + "/state"
	payload := getPayloadFromState(light.state, uType)
	payloadReader := strings.NewReader(payload)
	req, err := http.NewRequest("PUT", url, payloadReader)
	if err != nil {
		return err
	}
	client := http.Client{}
	client.Do(req)

	return nil
}
