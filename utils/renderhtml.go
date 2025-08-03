package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

func RenderVerificationEmail(path string, data any) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("no se pudo obtener el directorio actual: %w", err)
	}

	fullPath := filepath.Join(wd, path)

	tmpl, err := template.ParseFiles(fullPath)
	if err != nil {
		return "", fmt.Errorf("error al parsear template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("error al ejecutar template: %w", err)
	}

	return buf.String(), nil
}
