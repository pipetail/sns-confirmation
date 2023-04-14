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
	u, err := url.Parse(input)
	if err != nil {
		return fmt.Errorf("could not parse URL: %s", err)
	}

	hostnameParts := strings.Split(u.Hostname(), ".")

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

func ValidateURLAccounts(input string, accounts []string) error {
	u, err := url.Parse(input)
	if err != nil {
		return fmt.Errorf("could not parse URL: %s", err)
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return fmt.Errorf("could not parse query string: %s", err)
	}

	topicArn := q.Get("TopicArn")
	if topicArn == "" {
		return fmt.Errorf("topicArn can' be empty")
	}

	arnParts := strings.Split(topicArn, ":")

	if !util.Contains(accounts, arnParts[4]) {
		return fmt.Errorf("arn does not contain allowed accounts: %s vs %v", arnParts[4], accounts)
	}

	return nil
}

func ValidateURL(input string, regions []string, accounts []string) error {
	err := ValidateURLRegions(input, regions)
	if err != nil {
		return fmt.Errorf("could not validate regions: %s", err)
	}

	err = ValidateURLAccounts(input, accounts)
	if err != nil {
		return fmt.Errorf("could not validate accounts: %s", err)
	}

	return nil
}
