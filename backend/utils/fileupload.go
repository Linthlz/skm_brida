package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const MaxSize = 10 * 1024 * 1024
const MaxImageSize = 5 * 1024 * 1024
const MaxVideoURLLength = 500
const MaxImagesPerDocumentation = 10
const MaxDocumentSize = 20 * 1024 * 1024

var FileFormatAllowed = []string{".pdf"}
var AllowedMimes = []string{"application/pdf"}

var allowedImageTypes = []string{".jpg", ".jpeg", ".png", ".webp"}
var allowedImageMimes = []string{"image/jpeg", "image/png", "image/webp", "image/jpg"}

var allowedDocumentExts = []string{".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".zip"}
var allowedImageExtsForBatch = []string{".jpg", ".jpeg", ".png", ".webp"}

var allowedDocumentMimes = []string{"application/pdf", "application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document", "application/vnd.ms-excel", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "application/vnd.ms-powerpoint", "application/vnd.openxmlformats-officedocument.presentationml.presentation", "application/zip", "application/x-zip-compressed"}

func isAllowedExt(ext string, list []string) bool {
	for _, a := range list {
		if ext == a {
			return true
		}
	}
	return false
}

func isAllowedMime(mime string, list []string) bool {
	for _, a := range list {
		if mime == a {
			return true
		}
	}
	return false
}

func FileValidation(fileHeader multipart.FileHeader) string {
	if fileHeader.Size > MaxSize {
		return fmt.Sprintf("File too large, Max Size Allowed is %d MB", MaxSize/(1024*1024))
	}
	extension := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !isAllowedExt(extension, FileFormatAllowed) {
		return fmt.Sprintf("Invalid file format, Allowed format: %s", strings.Join(FileFormatAllowed, ", "))
	}
	contentType := fileHeader.Header.Get("Content-Type")
	if !isAllowedMime(contentType, AllowedMimes) {
		return "File does not match with allowed format"
	}
	return ""
}

func ImageValidation(fileHeader multipart.FileHeader) string {
	if fileHeader.Size > MaxImageSize {
		return fmt.Sprintf("Image too large. Maximum is %d MB", MaxImageSize/(1024*1024))
	}
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !isAllowedExt(ext, allowedImageTypes) {
		return fmt.Sprintf("Invalid image extension. Allowed: %s", strings.Join(allowedImageTypes, ", "))
	}
	mime := fileHeader.Header.Get("Content-Type")
	if !isAllowedMime(mime, allowedImageMimes) {
		return fmt.Sprintf("Invalid image type. Got: %s", mime)
	}
	return ""
}

func BatchFileValidation(fileHeader multipart.FileHeader) string {
	if fileHeader.Size > MaxDocumentSize {
		return fmt.Sprintf("File too large. Maximum is %d MB", MaxDocumentSize/(1024*1024))
	}
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowedExt := isAllowedExt(ext, allowedDocumentExts) || isAllowedExt(ext, allowedImageExtsForBatch)
	if !allowedExt {
		return fmt.Sprintf("Invalid file extension. Allowed: %s, %s",
			strings.Join(allowedImageExtsForBatch, ", "),
			strings.Join(allowedDocumentExts, ", "))
	}
	mime := fileHeader.Header.Get("Content-Type")
	allowedMime := isAllowedMime(mime, allowedDocumentMimes) || isAllowedMime(mime, allowedImageMimes)
	if !allowedMime {
		return fmt.Sprintf("Invalid file type. Got: %s", mime)
	}
	return ""
}

func FileUpload(c *fiber.Ctx, Filename string, uploadPath string) (string, error) {
	file, err := c.FormFile(Filename)
	if err != nil {
		if strings.Contains(err.Error(), "there is no uploaded file associated with the given key") {
			return "", nil
		}
		return "", err
	}
	if errs := FileValidation(*file); errs != "" {
		return "", errors.New(errs)
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	newName := uuid.New().String() + ext
	Path := fmt.Sprintf("%s/%s", uploadPath, newName)
	if err := os.MkdirAll(uploadPath, 0o755); err != nil {
		return "", err
	}

	if err := c.SaveFile(file, Path); err != nil {
		return "", err
	}
	return Path, nil
}

func FileDeletion(path string) error {
	return os.Remove(path)
}

func saveFileToDisk(fileHeader *multipart.FileHeader, savePath string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	if err := os.MkdirAll(filepath.Dir(savePath), 0755); err != nil {
		return err
	}
	dst, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	return err
}

func SaveImageBatch(form *multipart.Form, fieldName string, diskPath string, urlPrefix string) ([]string, error) {
	files := form.File[fieldName]
	if len(files) == 0 {
		return nil, nil
	}
	paths := make([]string, 0, len(files))
	for _, fileHeader := range files {
		if errMsg := ImageValidation(*fileHeader); errMsg != "" {
			return nil, errors.New(errMsg)
		}
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		filename := uuid.New().String() + ext
		savePath := filepath.Join(diskPath, filename)
		if err := saveFileToDisk(fileHeader, savePath); err != nil {
			return nil, err
		}
		paths = append(paths, strings.TrimRight(urlPrefix, "/")+"/"+filename)
	}
	return paths, nil
}

func SaveFileBatch(form *multipart.Form, fieldName string, diskPath string, urlPrefix string) ([]string, error) {
	files := form.File[fieldName]
	if len(files) == 0 {
		return nil, nil
	}
	paths := make([]string, 0, len(files))
	for _, fileHeader := range files {
		if errMsg := BatchFileValidation(*fileHeader); errMsg != "" {
			return nil, errors.New(errMsg)
		}
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		filename := uuid.New().String() + ext
		savePath := filepath.Join(diskPath, filename)
		if err := saveFileToDisk(fileHeader, savePath); err != nil {
			return nil, err
		}
		paths = append(paths, strings.TrimRight(urlPrefix, "/")+"/"+filename)
	}
	return paths, nil
}

func NormalizeVideoURL(url *string) (*string, error) {
	if url == nil {
		return nil, nil
	}
	clean := strings.TrimSpace(*url)
	if clean == "" {
		return nil, nil
	}
	if len(clean) > MaxVideoURLLength {
		return nil, errors.New("video URL is too long")
	}
	if !strings.HasPrefix(clean, "http://") && !strings.HasPrefix(clean, "https://") {
		return nil, errors.New("video URL must start with http:// or https://")
	}
	return &clean, nil
}
