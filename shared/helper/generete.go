package helper

import (
	"strings"

	"github.com/thanhpk/randstr"
)

func GenerateId() string {
	return strings.ToLower(randstr.String(15))
}
