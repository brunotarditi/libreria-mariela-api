package security

import (
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/o1egl/paseto"
)

type PasetoManager struct {
	symmetricKey []byte
	v2           *paseto.V2
}

func NewPasetoManager() (*PasetoManager, error) {
	pasetoKey := os.Getenv("PASETO_SECRET_KEY")
	if pasetoKey == "" {
		return nil, fmt.Errorf("la variable de entorno PASETO_SECRET_KEY no está definida")
	}
	secretPaseto, err := base64.StdEncoding.DecodeString(pasetoKey)

	if err != nil {
		return nil, fmt.Errorf("error al descodificar PASETO_SECRET_KEY: %v", err)
	}

	if len(secretPaseto) != 32 {
		return nil, fmt.Errorf("la clave PASETO debe tener 32 bytes")
	}

	return &PasetoManager{
		symmetricKey: secretPaseto,
		v2:           paseto.NewV2(),
	}, nil
}

func (p *PasetoManager) GenerateToken(userID uint, username string, duration time.Duration) (string, error) {
	if duration <= 0 {
		return "", fmt.Errorf("la duración del token debe ser mayor a cero")
	}

	if userID == 0 {
		return "", fmt.Errorf("el userID debe ser un valor válido")
	}

	if username == "" {
		return "", fmt.Errorf("el username no puede estar vacío")
	}

	now := time.Now()
	exp := now.Add(duration)

	jsonToken := paseto.JSONToken{
		Expiration: exp,
		IssuedAt:   now,
		Subject:    fmt.Sprintf("%d", userID),
	}
	jsonToken.Set("username", username)

	token, err := p.v2.Encrypt(p.symmetricKey, jsonToken, nil)
	if err != nil {
		return "", fmt.Errorf("error al generar el token: %v", err)
	}

	return token, nil
}

func (p *PasetoManager) VerifyToken(token string) (*paseto.JSONToken, error) {
	var jsonToken paseto.JSONToken
	var footer string
	err := p.v2.Decrypt(token, p.symmetricKey, &jsonToken, &footer)
	if err != nil {
		return nil, fmt.Errorf("error al desencriptar el token: %v", err)
	}
	if err := jsonToken.Validate(); err != nil {
		return nil, fmt.Errorf("token inválido: %v", err)
	}
	return &jsonToken, nil
}
