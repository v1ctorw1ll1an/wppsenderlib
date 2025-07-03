package wppsenderlib

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type PayloadBody struct {
	Parameters map[string]any
	TemplateId string
}

type Payload struct {
	From        string
	To          string
	PayloadBody *PayloadBody
}

type WppSender struct {
	Url     string
	Payload *Payload
	Headers map[string]string
}

func SendSynchronousMessage(wppSender *WppSender) (string, error) {

	payload, err := buildPayload(wppSender.Payload)

	if err != nil {
		return "", err
	}

	payloadReader := strings.NewReader(payload)

	req, err := http.NewRequest("POST", wppSender.Url, payloadReader)

	if err != nil {
		return "", err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/*+json")

	for key, value := range wppSender.Headers {
		req.Header.Add(key, value)
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	return string(body), nil
}

func buildPayload(payload *Payload) (string, error) {
	jsonBytes, err := json.Marshal(payload)

	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}
