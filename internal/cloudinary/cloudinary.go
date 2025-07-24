package cloudinary

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"image/png"
	"rent-application/configs"
	"rent-application/shared/web"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

func init() {
	_ = jpeg.DefaultQuality
	_ = png.BestCompression
}

type Cloudinary struct {
	Cloud   *cloudinary.Cloudinary
	IsError error
}

func NewCloudinary(cloud, apiKey, apiSecret string) Cloudinary {
	c, err := cloudinary.NewFromParams(cloud, apiKey, apiSecret)
	return Cloudinary{
		Cloud:   c,
		IsError: err,
	}
}

func NewCloudinaryCfg(cfg configs.CloudinaryConfig) Cloudinary {
	return NewCloudinary(cfg.CloudName, cfg.CloudApiKey, cfg.CloudApiSecret)
}

func (c *Cloudinary) UploadImage(ctx context.Context, file []byte) (string, error) {
	if c.IsError != nil {
		return "", c.IsError
	}

	const maxFileSize = 10 * 1024 * 1024 // 10 MB
	if len(file) > maxFileSize {
		return "", web.ErrBadRequest("file size exceeds 10 Mb limit")
	}

	img, format, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		return "", web.ErrBadRequest("failed to decode image: " + err.Error())
	}
	if format != "jpeg" && format != "png" && format != "jpg" {
		return "", web.ErrBadRequest("unsupported image format: " + format)
	}

	// Resize if needed (example: max width 2000px)
	if img.Bounds().Dx() > 2000 {
		img = imaging.Resize(img, 2000, 0, imaging.Lanczos)
	}

	var buf bytes.Buffer
	// Try WebP first, fallback to JPEG
	err = tryEncodeWebP(&buf, img)
	if err != nil {
		buf.Reset() // Clear the buffer if WebP failed
		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85}); err != nil {
			return "", web.ErrBadRequest("failed to encode image: " + err.Error())
		}
	}

	filename := uuid.NewString()
	formatToUpload := "webp"
	if buf.Len() == 0 {
		formatToUpload = "jpg"
	}

	res, err := c.Cloud.Upload.Upload(ctx, bytes.NewReader(buf.Bytes()), uploader.UploadParams{
		PublicID: "dev/" + filename,
		Format:   formatToUpload,
	})

	if err != nil {
		return "", err
	}

	return res.SecureURL, nil
}

// tryEncodeWebP attempts to encode as WebP without C dependencies
func tryEncodeWebP(w *bytes.Buffer, img image.Image) error {
	// This will only work if you have the pure Go version of a WebP encoder
	// For actual implementation, you'd need to use a pure Go WebP package
	// Currently returning error to trigger fallback
	return web.ErrBadRequest("webp encoding not available")
}
