package main

import (
	"encoding/json"
	"fmt"
	"github.com/chrisvasquez/go-twilio"
	"flag"
	"os"
)

var (
	apiKey = flag.String("apiKey", "", "Required : Twilio API Key")
	apiSecret = flag.String("apiSecret", "", "Required : Twilio API Secret")
	to = flag.String("to", "", "Required : Number to send sms to")
	from = flag.String("from", "", "Required : Twilio number to send message from")
	message = flag.String("message", "Hello World!", "Message to send")
)

func main() {
	required := []string{"apiKey", "apiSecret", "to", "from"}
	flag.Parse()

	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	var success = true
	for _, req := range required {
		if !seen[req] {
			// or possibly use `log.Fatalf` instead of:
			fmt.Fprintf(os.Stderr, "missing required -%s argument/flag\n", req)
			success = false
		}
	}
	if success != true {
		os.Exit(2) // the same exit code flag.Parse uses
	}
	fmt.Println("Starting go-twilio example app")
	client := twilio.NewClient(*apiKey, *apiSecret)
	message := &twilio.SMS{To: *to, From: *from, Body: *message}
	message.SetClient(client)
	res, err := message.Send()
	fmt.Println("Error:", err)
	smsRespose, _ := json.Marshal(res)
	fmt.Println(string(smsRespose))
}
