package wppsenderlib

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type PayloadBody struct {
	Parameters map[string]any `json:"parameters"`
	TemplateId string         `json:"templateId"`
}

type Payload struct {
	From        string         `json:"from"`
	To          string         `json:"to"`
	Options     map[string]any `json:"options,omitempty"`
	PayloadBody *PayloadBody   `json:"body"`
}

type WppSender struct {
	Url     string
	Payload *Payload
	Headers map[string]string
}

func SendSynchronousMessage(wppSender *WppSender) (string, error) {
	// Validação de entrada
	if wppSender == nil || wppSender.Payload == nil {
		return "", fmt.Errorf("wppSender ou payload não podem ser nil")
	}

	if wppSender.Url == "" {
		return "", fmt.Errorf("URL não pode estar vazia")
	}

	if wppSender.Payload.From == "" || wppSender.Payload.To == "" {
		return "", fmt.Errorf("campos From e To são obrigatórios")
	}

	payload, err := buildPayload(wppSender.Payload)
	if err != nil {
		return "", err
	}

	payloadReader := strings.NewReader(payload)

	req, err := http.NewRequest("POST", wppSender.Url, payloadReader)
	if err != nil {
		return "", fmt.Errorf("erro ao criar requisição: %w", err)
	}

	// Headers corretos
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	// Adicionar headers customizados
	for key, value := range wppSender.Headers {
		req.Header.Add(key, value)
	}

	// Cliente HTTP com timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
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

	// Verificação mais segura do status
	status, ok := jsonData["status"].(string)
	if res.StatusCode < 200 || res.StatusCode > 299 || !ok || status != "SENT" {
		return "", fmt.Errorf("API wts.chat retornou status %d: bodyBytes: %v", res.StatusCode, string(bodyBytes))
	}

	return fmt.Sprintf("Mensagem enviada com sucesso para %s (From: %s)", wppSender.Payload.To, wppSender.Payload.From), nil
}

func buildPayload(payload *Payload) (string, error) {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("erro ao serializar payload: %w", err)
	}
	return string(jsonBytes), nil
}
