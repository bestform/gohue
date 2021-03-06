package gohue_test

import (
	"log"

	"github.com/bestform/gohue"
)

func Example() {
	client := gohue.NewClient("<Your API Username>", "<Your API IP Address>")
	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	for _, light := range client.GetLights() {
		err := light.SwitchOn()
		if err != nil {
			log.Fatal(err)
		}
		err = light.SetColorRGB(110, 110, 210)
		if err != nil {
			log.Fatal(err)
		}
	}
}
