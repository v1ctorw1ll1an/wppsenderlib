package wppsenderlib

import (
	"encoding/json"
	"fmt"
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
		return "", fmt.Errorf("falha ao enviar requisição para API wts.chat: %w", err)
	}

	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)

	if err != nil {
		return "", fmt.Errorf("falha ao ler resposta da API wts.chat: %w", err)
	}

	var jsonData map[string]any
	err = json.Unmarshal(bodyBytes, &jsonData)

	if err != nil {
		return "", fmt.Errorf("falha ao decodificar resposta da API wts.chat: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode > 299 || jsonData["status"] != "SENT" {
		return "", fmt.Errorf("API wts.chat retornou status %d: %s", res.StatusCode, jsonData["status"])
	}

	return fmt.Sprintf("Mensagem enviada com sucesso para %s (From: %s)", wppSender.Payload.To, wppSender.Payload.From), nil
}

func buildPayload(payload *Payload) (string, error) {
	jsonBytes, err := json.Marshal(payload)

	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}
