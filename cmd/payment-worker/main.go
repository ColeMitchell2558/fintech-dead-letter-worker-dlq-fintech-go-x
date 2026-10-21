package main

import (
	lib "dlq-fintech-go/src"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	if os.Getenv("INFRAI_API_KEY") == "" {
		log.Fatal("set INFRAI_API_KEY")
	}
	c := lib.NewClient()
	msgs, err := c.QueueConsume(10, 30)
	if err != nil {
		log.Fatal(err)
	}
	for _, m := range msgs {
		d, err := lib.Handle(c, m)
		if err != nil {
			log.Print(err)
			continue
		}
		b, _ := json.Marshal(d)
		fmt.Printf("payment decision %s\n", b)
	}
}
