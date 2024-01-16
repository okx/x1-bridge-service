package utils

import (
	"crypto/sha256"
	"math/rand"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"golang.org/x/exp/constraints"
)

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))] //nolint:gosec
	}
	return string(b)
}

// GenerateRandomHash generates a random hash.
func GenerateRandomHash() [sha256.Size]byte {
	rand.Seed(time.Now().UnixNano())
	rs := generateRandomString(10) //nolint:gomnd
	return sha256.Sum256([]byte(rs))
}

// GenerateTraceID generates a random trace ID.
func GenerateTraceID() string {
	return generateRandomString(traceIDLen)
}

func Min[T constraints.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}

// LoggerWithTraceID returns a wrapping logger with a specific trace id
func LoggerWithTraceID(logger *log.Logger, traceID string) *log.Logger {
	if logger == nil {
		logger = log.GetDefaultLog()
	}
	return logger.WithFields(TraceID, traceID)
}

// LoggerWithRandomTraceID returns a wrapping logger with a random trace id
func LoggerWithRandomTraceID(logger *log.Logger) *log.Logger {
	return LoggerWithTraceID(logger, GenerateTraceID())
}
