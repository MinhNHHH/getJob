package llm

import (
	"fmt"
	"strings"
)

const (
	Model = "deepseek-r1:1.5b"
)

func removeThinkTags(response string) string {
	// Find the start and end positions of <think> tags
	startIdx := strings.Index(response, "<think>")
	endIdx := strings.Index(response, "</think>")

	// If both tags are found, remove the content between them including the tags
	if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
		beforeThink := response[:startIdx]
		afterThink := response[endIdx+len("</think>"):]
		return beforeThink + afterThink
	}

	// If tags aren't found or are malformed, return the original response
	return response
}

func GenerateCoverLetter(jobTitle, companyName, location, description string) (string, error) {
	prompt := fmt.Sprintf(`
	You are a cover letter generator.
	You are given a job title, company name, location, and description.
	You need to generate a cover letter for the job.
	
	Job Title: %s
	Company Name: %s
	Location: %s
	Description: %s
	`, jobTitle, companyName, location, description)

	payload := LLMPayload{
		Model:  Model,
		Prompt: prompt,
		Stream: false,
	}

	response, err := LLMApi("http://localhost:11434/api/generate", "POST", payload)
	if err != nil {
		return "", err
	}

	// Remove <think> tags from the response
	cleanedResponse := removeThinkTags(response.Response)
	return cleanedResponse, nil
}
