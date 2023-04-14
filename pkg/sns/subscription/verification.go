package subscription

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/pipetail/sns_verification/pkg/http"
	"github.com/pipetail/sns_verification/pkg/util"
)

type Verification struct {
	Type             string `json:"Type"`
	MessageID        string `json:"MessageId"`
	Token            string `json:"Token"`
	TopicARN         string `json:"TopicArn"`
	Message          string `json:"Message"`
	SubscribeURL     string `json:"SubscribeURL"`
	Timestamp        string `json:"Timestamp"`
	SignatureVersion string `json:"SignatureVersion"`
	Signature        string `json:"Signature"`
}

func VerificationFromByte(input []byte) (Verification, error) {
	v := Verification{}
	err := json.Unmarshal(input, &v)
	if err != nil {
		return v, fmt.Errorf("could not unmarshal verification request: %s", err)
	}
	return v, err
}

func VerificationFromReadCloser(input io.ReadCloser) (Verification, error) {
	v := Verification{}
	err := json.NewDecoder(input).Decode(&v)
	if err != nil {
		return v, fmt.Errorf("could not unmarshal verification request: %s", err)
	}
	return v, err
}

func Verify(client http.Client, url string, code int) error {
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("could not process the request: %s", err)
	}

	if resp.StatusCode != code {
		return fmt.Errorf("status code %d does not match the expected %d", resp.StatusCode, code)
	}

	return nil
}

func ValidateURLRegions(input string, regions []string) error {
	url, err := url.Parse(input)
	if err != nil {
		return fmt.Errorf("could not parse URL: %s", err)
	}

	hostnameParts := strings.Split(url.Hostname(), ".")

	if hostnameParts[0] != "sns" {
		return fmt.Errorf("url does not contain required sns part")
	}

	if !util.Contains(regions, hostnameParts[1]) {
		return fmt.Errorf("url does not contain allowed region: %s vs %v", hostnameParts[1], regions)
	}

	if hostnameParts[2] != "amazonaws" && hostnameParts[3] != "com" {
		return fmt.Errorf("url does not contain allowed base aws domain")
	}

	return nil
}

// TODO: add validation for AWS accounts
// TODO: add validation for SNS ARNs
// TODO: add common validation function
