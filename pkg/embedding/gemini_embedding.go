package embedding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// {"model": "models/gemini-embedding-001",
//      "content": {"parts":[{"text": "What is the meaning of life?"}]}
//     }

type EmbeddingRequestParts struct {
	Text string `json:"text"` // Tambahkan tag json
}

type EmbeddingRequestContent struct {
	Parts []EmbeddingRequestParts `json:"parts"` // Gunakan huruf kecil 'parts'
}

type EmbeddingRequest struct {
	Model    string                  `json:"model"`
	Content  EmbeddingRequestContent `json:"content"`
	TaskType string                  `json:"task_type,omitempty"` // Opsional
}

type EmbeddingResponse struct {
	Embedding EmbeddingResponseEmbedding `json:"embedding"`
}

type EmbeddingResponseEmbedding struct {
	Values []float32 `json:"values"`
}

func GetGemniniEmbedding(
	apiKey string,
	text string,
) (*EmbeddingResponse, error) {

	geminiReq := EmbeddingRequest{
		Model: "models/gemini-embedding-exp-03-07",
		Content: EmbeddingRequestContent{
			Parts: []EmbeddingRequestParts{
				{
					Text: text,
				},
			},
		},
		TaskType: "RETRIEVAL_DOCUMENT",
	}

	geminiReqJson, err := json.Marshal(geminiReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		"https://generativelanguage.googleapis.com/v1beta/models/gemini-embedding-exp-03-07:embedContent",
		bytes.NewBuffer(geminiReqJson),
	)

	if err != nil {
		return nil, err
	}

	req.Header.Set("x-goog-api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	resByte, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("error from response, code %d, body %s", res.StatusCode, string(resByte))
	}

	var resEmbedding EmbeddingResponse
	err = json.Unmarshal(resByte, &resEmbedding)
	if err != nil {
		panic(err)
	}

	return &resEmbedding, nil
}
