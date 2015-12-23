package gohue

type jsonLightMap map[string]jsonLight

type jsonLight struct {
	jsonLightState   `json:"state"`
	Type             string `json:"type"`
	Name             string `json:"name"`
	Modelid          string `json:"modelid"`
	Manufacturername string `json:"manufacturername"`
	Uniqueid         string `json:"uniqueid"`
	Swversion        string `json:"swversion"`
}

type jsonLightState struct {
	On        bool      `json:"on"`
	Bri       int       `json:"bri"`
	Hue       int       `json:"hue"`
	Sat       int       `json:"sat"`
	Effect    string    `json:"effect"`
	Xy        []float64 `json:"xy"`
	Ct        int       `json:"ct"`
	Alert     string    `json:"alert"`
	Colormode string    `json:"colormode"`
	Reachable bool      `json:"reachable"`
}
