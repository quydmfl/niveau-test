package helper

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateProductReference generate product reference
// Pattern: PROD-YYYY-MM-DD-XXXX
func GenerateProductReference() string {
	return fmt.Sprintf("PROD-%s-%04d", time.Now().Format("2006-01"), rand.Intn(10000))
}
