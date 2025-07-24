package helper

import "time"

func ConvertToJakarta(t time.Time) string {
	return t.Add(7 * time.Hour).Format(time.RFC3339)
}
