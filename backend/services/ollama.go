package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const modelAPI = "http://localhost:11434/api/generate"

type generateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

func Generate(ctx context.Context, model, prompt string) (string, error) {
	reqBody := generateRequest{Model: model, Prompt: prompt, Stream: false}
	b, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, "POST", modelAPI, bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// Attempt to parse JSON; fall back to raw text
	var parsed map[string]any
	if err := json.Unmarshal(body, &parsed); err == nil {
		// common keys: "text", "output", "result"
		for _, k := range []string{"text", "output", "result"} {
			if v, ok := parsed[k]; ok {
				return fmt.Sprint(v), nil
			}
		}
		// otherwise return whole JSON string
		return string(body), nil
	}
	return string(body), nil
}

// StreamToWriter forwards streaming model API response (raw) to w as-is.
// Caller is responsible for setting SSE headers and framing.
func StreamToWriter(ctx context.Context, model, prompt string, w io.Writer) error {
	reqBody := generateRequest{Model: model, Prompt: prompt, Stream: true}
	b, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, "POST", modelAPI, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			chunk := buf[:n]
			// write as SSE data lines
			if _, werr := w.Write(append([]byte("data: "), chunk...)); werr != nil {
				return werr
			}
			if _, werr := w.Write([]byte("\n\n")); werr != nil {
				return werr
			}
		}
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}
