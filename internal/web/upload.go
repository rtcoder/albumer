package web

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var allowedCoverTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

func (s *Server) saveCover(file multipart.File, header *multipart.FileHeader) (string, error) {
	if header == nil || header.Filename == "" {
		return "", nil
	}

	var sniff [512]byte
	n, err := file.Read(sniff[:])
	if err != nil && err != io.EOF {
		return "", err
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	contentType := http.DetectContentType(sniff[:n])
	ext, ok := allowedCoverTypes[contentType]
	if !ok {
		return "", fmt.Errorf("cover must be JPEG, PNG, or WebP")
	}

	if err := os.MkdirAll(s.cfg.UploadDir, 0o755); err != nil {
		return "", err
	}

	base := newToken()
	originalName := base + ext
	originalPath := uploadPath(s.cfg.UploadDir, originalName)
	out, err := os.OpenFile(originalPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(out, file); err != nil {
		_ = out.Close()
		return "", err
	}
	if err := out.Close(); err != nil {
		return "", err
	}

	thumbName := base + "-thumb" + ext
	thumbPath := uploadPath(s.cfg.UploadDir, thumbName)
	if contentType == "image/webp" {
		if err := copyFile(originalPath, thumbPath); err != nil {
			return "", err
		}
	} else if err := createThumbnail(originalPath, thumbPath, ext); err != nil {
		return "", err
	}

	return "/uploads/" + filepath.Base(originalName), nil
}

func createThumbnail(source, target, ext string) error {
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()

	img, _, err := image.Decode(in)
	if err != nil {
		return err
	}

	thumb := resizeNearest(img, 360)
	out, err := os.Create(target)
	if err != nil {
		return err
	}
	defer out.Close()

	if strings.EqualFold(ext, ".png") {
		return png.Encode(out, thumb)
	}
	return jpeg.Encode(out, thumb, &jpeg.Options{Quality: 82})
}

func resizeNearest(src image.Image, maxSide int) image.Image {
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	if width <= maxSide && height <= maxSide {
		return src
	}

	newWidth := maxSide
	newHeight := height * maxSide / width
	if height > width {
		newHeight = maxSide
		newWidth = width * maxSide / height
	}
	if newWidth < 1 {
		newWidth = 1
	}
	if newHeight < 1 {
		newHeight = 1
	}

	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			srcX := bounds.Min.X + x*width/newWidth
			srcY := bounds.Min.Y + y*height/newHeight
			dst.Set(x, y, src.At(srcX, srcY))
		}
	}
	return dst
}

func copyFile(source, target string) error {
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(target)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
