package llm

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/MinhNHHH/get-job/pkg/cfgs"
)

type LLMPayload struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type LLMResponse struct {
	Response string `json:"response"`
}

type LLM struct {
	Cfgs cfgs.Configs
}

func NewLLM(config cfgs.Configs) LLM {
	return LLM{
		Cfgs: config,
	}
}

func (llm *LLM) LLMAPI(payload LLMPayload) (LLMResponse, error) {
	client := &http.Client{}
	method := "POST"
	// Convert payload to JSON string
	payloadStr, err := json.Marshal(payload)
	if err != nil {
		return LLMResponse{}, err
	}

	// Create a new request with the specified method
	req, err := http.NewRequest(method, llm.Cfgs.LLM_URI, strings.NewReader(string(payloadStr)))
	if err != nil {
		return LLMResponse{}, err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return LLMResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return LLMResponse{}, err
	}

	var response LLMResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return LLMResponse{}, err
	}
	return response, nil
}
