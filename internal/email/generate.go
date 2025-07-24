package email

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateOtpCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%04d", rand.Intn(10000))
}
