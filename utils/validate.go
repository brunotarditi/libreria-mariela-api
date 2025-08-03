package utils

import (
	"fmt"
	"regexp"
)

func ValidatePassword(password string) error {
	var (
		minLen    = 8
		hasUpper  = regexp.MustCompile(`[A-Z]`)
		hasNumber = regexp.MustCompile(`[0-9]`)
		hasSymbol = regexp.MustCompile(`[^a-zA-Z0-9]`)
	)

	switch {
	case len(password) < minLen:
		return fmt.Errorf("la contraseña debe tener al menos 8 caracteres")
	case !hasUpper.MatchString(password):
		return fmt.Errorf("la contraseña debe incluir una letra mayúscula")
	case !hasNumber.MatchString(password):
		return fmt.Errorf("la contraseña debe incluir un número")
	case !hasSymbol.MatchString(password):
		return fmt.Errorf("la contraseña debe incluir un carácter especial")
	default:
		return nil
	}
}
