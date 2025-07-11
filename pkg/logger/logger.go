package logger

import (
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"os"
	"strings"
)

var Logger *slog.Logger

// Init initializes the logger with appropriate settings
func Init() {
	logLevel := slog.LevelInfo
	
	switch strings.ToUpper(os.Getenv("LOG_LEVEL")) {
	case "DEBUG":
		logLevel = slog.LevelDebug
	case "WARN":
		logLevel = slog.LevelWarn
	case "ERROR":
		logLevel = slog.LevelError
	}
	
	opts := &slog.HandlerOptions{
		Level: logLevel,
		ReplaceAttr: maskSensitiveData,
	}

	var handler slog.Handler
	if os.Getenv("ENVIRONMENT") == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}
	
	Logger = slog.New(handler)
}

func maskSensitiveData(groups []string, a slog.Attr) slog.Attr {
	if os.Getenv("DISABLE_LOG_MASKING") == "true" {
		return a
	}

	switch a.Key {
	case "zone_id", "record_id":
		if str, ok := a.Value.Any().(string); ok && len(str) > 8 {
			a.Value = slog.StringValue(MaskID(str))
		}
	case "record_name", "zone_name":
		if str, ok := a.Value.Any().(string); ok {
			a.Value = slog.StringValue(MaskDomain(str))
		}
	case "api_token", "api_email", "api_key":
		a.Value = slog.StringValue("***REDACTED***")
	case "email":
		if str, ok := a.Value.Any().(string); ok {
			a.Value = slog.StringValue(MaskEmail(str))
		}
	case "target_ip", "ip_address":
		if str, ok := a.Value.Any().(string); ok {
			a.Value = slog.StringValue(MaskIP(str))
		}
	}
	return a
}

// MaskID masks an ID showing only first and last 4 characters
func MaskID(id string) string {
	if len(id) <= 8 {
		return "****"
	}
	return id[:4] + "****" + id[len(id)-4:]
}

// MaskDomain masks a domain showing only TLD
func MaskDomain(domain string) string {
	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return "***"
	}

	// Keep only the TLD visible
	for i := 0; i < len(parts)-1; i++ {
		parts[i] = "***"
	}

	return strings.Join(parts, ".")
}

// MaskEmail masks an email address
func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "***@***"
	}
	
	localPart := parts[0]
	if len(localPart) > 2 {
		localPart = localPart[:1] + "***" + localPart[len(localPart)-1:]
	} else {
		localPart = "***"
	}
	
	return localPart + "@" + parts[1]
}

// MaskIP masks an IP address showing only the first octet
func MaskIP(ip string) string {
	parts := strings.Split(ip, ".")
	if len(parts) == 4 {
		return parts[0] + ".***.***.***"
	}
	return "***.***.***"
}

// HashString returns a SHA256 hash of the input string (first 8 chars)
func HashString(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	hash := hex.EncodeToString(h.Sum(nil))
	if len(hash) > 8 {
		return hash[:8]
	}
	return hash
}

// Debug logs a debug message
func Debug(msg string, args ...any) {
	if Logger != nil {
		Logger.Debug(msg, args...)
	}
}

// Info logs an info message
func Info(msg string, args ...any) {
	if Logger != nil {
		Logger.Info(msg, args...)
	}
}

// Warn logs a warning message
func Warn(msg string, args ...any) {
	if Logger != nil {
		Logger.Warn(msg, args...)
	}
}

// Error logs an error message
func Error(msg string, args ...any) {
	if Logger != nil {
		Logger.Error(msg, args...)
	}
}
