package upload

import (
	"context"
	"rent-application/internal/cloudinary"
)

type UploadService interface {
	UploadImage(ctx context.Context, file []byte) (string, error)
}

type uploadService struct {
	cloudinary cloudinary.Cloudinary
}

func NewUploadService(cld cloudinary.Cloudinary) UploadService {
	return &uploadService{
		cloudinary: cld,
	}
}

func (s *uploadService) UploadImage(ctx context.Context, file []byte) (string, error) {
	return s.cloudinary.UploadImage(ctx, file)
}
