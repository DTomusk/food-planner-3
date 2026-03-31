package upload

import (
	"context"
	"strings"
)

type CreateSignedUploadURLRequest struct {
	ObjectKey     string
	FileType      string
	FileSizeBytes int64
}

type CreateSignedUploadURLResponse struct {
	// UploadURL is the presigned URL that clients can use to upload the file.
	UploadURL string
	// FileURL is the public URL where the uploaded file will be accessible after upload.
	FileURL string
}

type UploadProvider interface {
	CreateSignedUploadURL(ctx context.Context, req CreateSignedUploadURLRequest) (*CreateSignedUploadURLResponse, error)
}

type StaticUploadProvider struct {
	uploadBaseURL string
	publicBaseURL string
}

func NewStaticUploadProvider(uploadBaseURL, publicBaseURL string) *StaticUploadProvider {
	return &StaticUploadProvider{
		uploadBaseURL: uploadBaseURL,
		publicBaseURL: publicBaseURL,
	}
}

func (p *StaticUploadProvider) CreateSignedUploadURL(_ context.Context, req CreateSignedUploadURLRequest) (*CreateSignedUploadURLResponse, error) {
	if p == nil {
		return nil, ErrProviderNotConfigured
	}

	return &CreateSignedUploadURLResponse{
		UploadURL: joinURL(p.uploadBaseURL, req.ObjectKey),
		FileURL:   joinURL(p.publicBaseURL, req.ObjectKey),
	}, nil
}

func joinURL(baseURL, path string) string {
	return strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(path, "/")
}
