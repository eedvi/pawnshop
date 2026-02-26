package logger

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/rs/zerolog"
)

// SensitiveFields lista de campos que deben ser sanitizados
var SensitiveFields = []string{
	"password",
	"token",
	"secret",
	"api_key",
	"apikey",
	"authorization",
	"credit_card",
	"ssn",
	"pin",
	"otp",
	"refresh_token",
	"access_token",
}

// SanitizedString es un tipo que siempre se loggea como redacted
type SanitizedString string

func (s SanitizedString) MarshalZerologObject(e *zerolog.Event) {
	e.Str("value", "***REDACTED***")
}

// Sanitize remueve datos sensibles de un string
func Sanitize(input string) string {
	if input == "" {
		return input
	}

	// Sanitizar passwords en URLs
	passwordInURLRegex := regexp.MustCompile(`://[^:]+:([^@]+)@`)
	input = passwordInURLRegex.ReplaceAllString(input, "://$user:***@")

	// Sanitizar tokens Bearer
	bearerRegex := regexp.MustCompile(`Bearer\s+[A-Za-z0-9\-_\.]+`)
	input = bearerRegex.ReplaceAllString(input, "Bearer ***")

	// Sanitizar números de tarjeta de crédito (formato básico)
	creditCardRegex := regexp.MustCompile(`\b\d{4}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}\b`)
	input = creditCardRegex.ReplaceAllString(input, "****-****-****-****")

	return input
}

// SanitizeMap remueve campos sensibles de un map
func SanitizeMap(data map[string]interface{}) map[string]interface{} {
	if data == nil {
		return nil
	}

	sanitized := make(map[string]interface{})
	for key, value := range data {
		if isSensitiveField(key) {
			sanitized[key] = "***REDACTED***"
		} else if strValue, ok := value.(string); ok {
			sanitized[key] = Sanitize(strValue)
		} else if mapValue, ok := value.(map[string]interface{}); ok {
			sanitized[key] = SanitizeMap(mapValue)
		} else {
			sanitized[key] = value
		}
	}

	return sanitized
}

// SanitizeJSON sanitiza un JSON string
func SanitizeJSON(jsonStr string) string {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		// Si no es JSON válido, sanitizar como string
		return Sanitize(jsonStr)
	}

	sanitized := SanitizeMap(data)
	result, err := json.Marshal(sanitized)
	if err != nil {
		return jsonStr
	}

	return string(result)
}

// isSensitiveField verifica si un campo es sensible
func isSensitiveField(fieldName string) bool {
	lowerField := strings.ToLower(fieldName)
	for _, sensitive := range SensitiveFields {
		if strings.Contains(lowerField, sensitive) {
			return true
		}
	}
	return false
}

// SanitizeSQL limpia queries SQL de valores sensibles
func SanitizeSQL(query string) string {
	// Remover valores de passwords en INSERTs/UPDATEs
	passwordRegex := regexp.MustCompile(`password\s*=\s*'[^']+'`)
	query = passwordRegex.ReplaceAllString(query, "password = '***'")

	// Limitar longitud de queries muy largas
	if len(query) > 500 {
		query = query[:500] + "... (truncated)"
	}

	return query
}

// SanitizationHook es un hook de zerolog para sanitizar automáticamente
type SanitizationHook struct{}

func (h SanitizationHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	// Este hook se puede usar para interceptar todos los logs
	// y sanitizar campos automáticamente
}
