package wppsenderlib

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_sendMessage(t *testing.T) {

	t.Run("make payload json", func(t *testing.T) {
		payloadObj := &Payload{
			From: "5511911111111",
			To:   "5522922222222",
			PayloadBody: &PayloadBody{
				TemplateId: "template_id_xpto",
			},
		}
		payload, err := buildPayload(payloadObj)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `{"from":"5511911111111","to":"5522922222222","body":{"parameters":null,"templateId":"template_id_xpto"}}`
		got := payload

		if got != want {
			t.Errorf(`got "%v", want "%v"`, got, want)
		}
	})

	t.Run("should send correct payload and headers", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected POST method, got %s", r.Method)
			}

			if ct := r.Header.Get("content-type"); !strings.Contains(ct, "application") {
				t.Errorf("missing or wrong content-type header: %s", ct)
			}

			body, _ := io.ReadAll(r.Body)
			expected := `{"from":"5511911111111","to":"5522922222222","body":{"parameters":null,"templateId":"template_id_xpto"}}`
			if string(body) != expected {
				t.Errorf("unexpected payload: got %s, want %s", body, expected)
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"status":"sent"}`)
		}))
		defer testServer.Close()

		wpp := &WppSender{
			Url: testServer.URL,
			Payload: &Payload{
				From: "5511911111111",
				To:   "5522922222222",
				PayloadBody: &PayloadBody{
					TemplateId: "template_id_xpto",
				},
			},
			Headers: map[string]string{
				"X-Custom-Header": "123",
			},
		}

		resp, err := SendSynchronousMessage(wpp)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if resp != `{"status":"sent"}` {
			t.Errorf("unexpected response: got %s", resp)
		}
	})
}
