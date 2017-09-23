package twilio

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
)

type SMS struct {
	// Required
	To string // required in posting messages for sending
	client         *client

	// Conditional Parameters
	From                string // present either with MessagingServiceSid  - for sending
	MessagingServiceSid string // present either with From  - for sending

	Body     string // at least this or MediaUrl
	MediaUrl string // url to an image for sending, only available in US and Canada

	// Optional Parameters

	// Receives standard request parameters
	//   - https://www.twilio.com/docs/api/twiml/sms/twilio_request#request-parameters
	//   together with :
	//		MessageStatus - possible values in https://www.twilio.com/docs/api/messaging/message#message-status-values
	//		ErrorCode  - possible values in https://www.twilio.com/docs/api/messaging/message#delivery-related-errors
	StatusCallback string // needed if MessageServiceId is the one sent, will receive delivery errors if any

	// For implementation
	// ApplicationSid string
	// MaxPrice	string
	// ProviderFeedback string
	// ValidityPeriod string
}

type SMSReponse struct {
	SID              string           `json:"sid,omitempty"`
	CreatedAt        string           `json:"date_created,omitempty"`
	UpdatedAt        string           `json:"date_updated,omitempty"`
	SentAt           string           `json:"date_sent,omitempty"`
	AccountSID       string           `json:"account_sid,omitempty"`
	To               string           `json:"to,omitempty"`
	From             string           `json:"from,omitempty"`
	Body             string           `json:"body,omitempty"`
	Status           string           `json:"status,omitempty"`
	NumberOfSegments string           `json:"num_segments,omitempty"`
	NumberOfMedia    string           `json:"num_media,omitempty"`
	Direction        string           `json:"direction,omitempty"`
	APIVersion       string           `json:"api_version,omitempty"`
	Price            string           `json:"price,omitempty"`
	PriceUnit        string           `json:"price_unit,omitempty"`
	ErrorCode        string           `json:"error_code,omitempty"`
	ErrorMessage     string           `json:"error_message,omitempty"`
	URI              string           `json:"uri,omitempty"`
	SubResourceURIS  []SubResourceURI `json:"subresource_uris,omitempty"`
	Code             int64            `json:"code,omitempty"`
	Message          string           `json:"message,omitempty"`
	MoreInfo         string           `json:"more_info,omitempty"`
}

type SubResourceURI struct {
	Media string `json:"media"`
}

func (s *SMS) SetClient(c *client) error {
	s.client = c
	return nil
}

func (s *SMS) getSMSPayload() string {
	v := url.Values{}
	v.Set("To", s.To)
	v.Set("From", s.From)
	v.Set("Body", s.Body)
	return v.Encode()
}

func (s *SMS) Send() (SMSReponse, error) {
	payload := s.getSMSPayload()
	fmt.Println(payload)
	log.Println("Payload", payload)
	response, _ := s.client.do(payload)
	var result SMSReponse
	err := json.Unmarshal(response, &result)
	return result, err
}
