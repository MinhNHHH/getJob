package app

import (
	"bytes"
	"io"
	"strings"

	"github.com/ledongthuc/pdf"
)

func ExtractTextFromBytes(pdfBytes []byte) (string, error) {
	// Use pdf.NewReader to read from bytes instead of a file
	r, err := pdf.NewReader(bytes.NewReader(pdfBytes), int64(len(pdfBytes)))
	if err != nil {
		return "", err
	}

	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	_, err = io.Copy(&sb, b)
	if err != nil {
		return "", err
	}

	rawText := sb.String()
	return rawText, nil
}
