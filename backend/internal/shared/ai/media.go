package ai

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

// Part is one piece of a multi-part user message: text, an image or a PDF.
type Part struct {
	Type     string     `json:"type"`
	Text     string     `json:"text,omitempty"`
	ImageURL *imagePart `json:"image_url,omitempty"`
	File     *filePart  `json:"file,omitempty"`
}

type imagePart struct {
	URL string `json:"url"`
}

type filePart struct {
	Filename string `json:"filename"`
	FileData string `json:"file_data"`
}

func TextPart(text string) Part { return Part{Type: "text", Text: text} }

// ImagePart embeds an image (png, jpeg, webp, gif) as a data URL.
func ImagePart(mimeType string, data []byte) Part {
	return Part{Type: "image_url", ImageURL: &imagePart{
		URL: "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(data),
	}}
}

// PDFPart embeds a PDF file.
func PDFPart(filename string, data []byte) Part {
	return Part{Type: "file", File: &filePart{
		Filename: filename,
		FileData: "data:application/pdf;base64," + base64.StdEncoding.EncodeToString(data),
	}}
}

type partsMessage struct {
	Role    string `json:"role"`
	Content any    `json:"content"`
}

type partsRequest struct {
	Model          string          `json:"model"`
	Messages       []partsMessage  `json:"messages"`
	ResponseFormat *responseFormat `json:"response_format,omitempty"`
	Temperature    *float64        `json:"temperature,omitempty"`
}

// CompleteJSONWithParts sends a system prompt and a multi-part user message
// (images, PDFs, text) and decodes a strict JSON-schema reply into out.
func (c *Client) CompleteJSONWithParts(ctx context.Context, system string, parts []Part, schema Schema, out any) error {
	if !c.Configured() {
		return ErrNotConfigured
	}
	temperature := 0.0
	payload, err := json.Marshal(partsRequest{
		Model: c.model,
		Messages: []partsMessage{
			{Role: "system", Content: system},
			{Role: "user", Content: parts},
		},
		ResponseFormat: &responseFormat{
			Type:       "json_schema",
			JSONSchema: &jsonSchemaFormat{Name: schema.Name, Strict: true, Schema: schema.Schema},
		},
		Temperature: &temperature,
	})
	if err != nil {
		return fmt.Errorf("marshal openai request: %w", err)
	}
	content, err := c.completeRaw(ctx, payload)
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(content), out); err != nil {
		return fmt.Errorf("decode openai json reply: %w", err)
	}
	return nil
}

// Transcription is Whisper's verbose_json result.
type Transcription struct {
	Text     string                 `json:"text"`
	Language string                 `json:"language"`
	Duration float64                `json:"duration"`
	Segments []TranscriptionSegment `json:"segments"`
}

type TranscriptionSegment struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Text  string  `json:"text"`
}

// Transcribe sends audio to the transcription endpoint (Whisper) and returns
// the text with segment timestamps. One retry on 5xx.
func (c *Client) Transcribe(ctx context.Context, filename string, data []byte) (Transcription, error) {
	if !c.Configured() {
		return Transcription{}, ErrNotConfigured
	}

	var lastErr error
	for attempt := 0; attempt < 2; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return Transcription{}, ctx.Err()
			case <-time.After(c.retryDelay):
			}
		}
		result, retry, err := c.transcribeOnce(ctx, filename, data)
		if err == nil {
			return result, nil
		}
		lastErr = err
		if !retry {
			break
		}
	}
	return Transcription{}, lastErr
}

func (c *Client) transcribeOnce(ctx context.Context, filename string, data []byte) (Transcription, bool, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	fields := [][2]string{
		{"model", c.transcribeModel},
		{"response_format", "verbose_json"},
		{"timestamp_granularities[]", "segment"},
		{"language", "en"},
	}
	for _, field := range fields {
		if err := writer.WriteField(field[0], field[1]); err != nil {
			return Transcription{}, false, err
		}
	}
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return Transcription{}, false, err
	}
	if _, err := part.Write(data); err != nil {
		return Transcription{}, false, err
	}
	if err := writer.Close(); err != nil {
		return Transcription{}, false, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/audio/transcriptions", &body)
	if err != nil {
		return Transcription{}, false, fmt.Errorf("build transcription request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	client := &http.Client{Timeout: 3 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return Transcription{}, ctx.Err() == nil, fmt.Errorf("call transcription: %w", err)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return Transcription{}, true, fmt.Errorf("read transcription: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		var parsed chatResponse
		_ = json.Unmarshal(raw, &parsed)
		message := fmt.Sprintf("status %d", resp.StatusCode)
		if parsed.Error != nil && parsed.Error.Message != "" {
			message += ": " + parsed.Error.Message
		}
		return Transcription{}, resp.StatusCode >= 500, errors.New("transcription failed (" + message + ")")
	}
	var result Transcription
	if err := json.Unmarshal(raw, &result); err != nil {
		return Transcription{}, false, fmt.Errorf("parse transcription: %w", err)
	}
	return result, false, nil
}
