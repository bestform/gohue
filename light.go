package gohue

import (
	"net/http"
	"strconv"
	"strings"
)

type updateType int

const (
	updateColorHue updateType = iota
	updateColorXy
	updateBrightness
	updateOnOff
	updateSaturation
	updateEffect
	updateAlert
)

// Effect represents an effect type used by SetEffect()
type Effect string

const (
	EffectNone      Effect = "none"
	EffectColorloop Effect = "colorloop"
)

// Alert represents an effect type used by SetAlert()
type Alert string

const (
	AlertNone   Alert = "none"
	AlertSelect Alert = "select"
)

// Light represents one light in the network of a Client.
type Light struct {
	id     string
	name   string
	state  state
	client *client
}

// State represents the current or desired state of a light
type state struct {
	Hue, Brightness, Saturation int
	On                          bool
	Xy                          []float64
	Effect                      string
	Alert                       string
}

func newLight(id string) *Light {
	l := Light{}
	l.id = id
	l.state = state{10000, 254, 254, true, []float64{0.0, 0.0}, "none", "none"}

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

// SetEffect sets the effect value to one of the possible Effect constants. It has an immediate effect.
func (l *Light) SetEffect(e Effect) error {
	l.state.Effect = string(e)
	return l.updateState(updateEffect)
}

// SetAlert sets the alert value to one of the possible Alert constants. It has an immediate effect.
func (l *Light) SetAlert(a Alert) error {
	l.state.Alert = string(a)
	return l.updateState(updateAlert)
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

func getPayloadFromState(s state, uType updateType) string {
	var payload string
	switch uType {
	case updateColorHue:
		payload = "{\"hue\":" + strconv.Itoa(s.Hue) + "}"
	case updateColorXy:
		x := strconv.FormatFloat(s.Xy[0], 'f', 4, 64)
		y := strconv.FormatFloat(s.Xy[1], 'f', 4, 64)
		payload = "{\"xy\":[" + x + "," + y + "]}"
	case updateBrightness:
		payload = "{\"bri\":" + strconv.Itoa(s.Brightness) + "}"
	case updateSaturation:
		payload = "{\"sat\":" + strconv.Itoa(s.Saturation) + "}"
	case updateOnOff:
		var onState string
		if s.On {
			onState = "true"
		} else {
			onState = "false"
		}
		payload = "{\"on\":" + onState + "}"
	case updateEffect:
		payload = "{\"effect\": \"" + s.Effect + "\"}"
	case updateAlert:
		payload = "{\"alert\": \"" + s.Alert + "\"}"
	}

	return payload
}
