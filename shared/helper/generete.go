package helper

import (
	"strings"

	"github.com/google/uuid"
	"github.com/thanhpk/randstr"
)

func GenerateId() string {
	return strings.ToLower(randstr.String(15))
}

func GenerateOrderId() string {
	return uuid.New().String()
}
