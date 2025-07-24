package users

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
)

type UsersRequest struct {
	Name string `json:"username"`
	Img  string `json:"img"`
}

var (
	ErrUsernameTooShort    = errors.New("username must be at least 5 characters")
	ErrInvalidImageURL     = errors.New("invalid image URL format")
	ErrUnapprovedImageHost = errors.New("image must come from an approved source")
	ErrInvalidImageFormat  = errors.New("image must be in webp format")

	// Approved image sources (can be expanded)
	approvedImageHosts = []string{
		"res.cloudinary.com",
		"storage.googleapis.com",
		"aws.amazon.com",
	}

	// Regex to validate image URLs from approved sources
	imageURLRegex = regexp.MustCompile(`^https?://([a-zA-Z0-9-]+\.)*(` +
		strings.Join(approvedImageHosts, "|") + `)/.+\.webp(\?.*)?$`)
)

func (r *UsersRequest) Validate() error {
	// Validate username
	if len(strings.TrimSpace(r.Name)) < 5 {
		return ErrUsernameTooShort
	}

	// Skip image validation if empty
	if r.Img == "" {
		return nil
	}

	// Basic URL validation
	_, err := url.ParseRequestURI(r.Img)
	if err != nil {
		return ErrInvalidImageURL
	}

	// Check against approved sources and format
	if !imageURLRegex.MatchString(r.Img) {
		// Determine which part failed
		u, _ := url.Parse(r.Img)
		if !contains(approvedImageHosts, u.Host) {
			return ErrUnapprovedImageHost
		}
		if !strings.HasSuffix(strings.ToLower(u.Path), ".webp") {
			return ErrInvalidImageFormat
		}
		return ErrInvalidImageURL
	}

	return nil
}

// Helper function to check if host is in approved list
func contains(sources []string, host string) bool {
	for _, source := range sources {
		if strings.Contains(host, source) {
			return true
		}
	}
	return false
}
