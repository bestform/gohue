package gohue_test

import (
	"log"

	"github.com/bestform/gohue"
)

func Example() {
	client := gohue.NewClient("<Your API Username>", "<You API IP Address")
	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	for _, light := range client.Lights {
		light.State.Hue = 10000
		err := light.UpdateState()
		if err != nil {
			log.Fatal(err)
		}
	}
}
