package upload

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const (
	defaultR2Region        = "auto"
	defaultR2PresignExpiry = 5 * time.Minute
)

type R2UploadProviderConfig struct {
	AccountID       string
	EndpointURL     string
	BucketName      string
	AccessKeyID     string
	SecretAccessKey string
	PublicBaseURL   string
	Region          string
	PresignExpiry   time.Duration
}

type R2UploadProvider struct {
	bucketName    string
	publicBaseURL string
	presignClient *s3.PresignClient
	presignExpiry time.Duration
	s3Client      *s3.Client
}

func NewR2UploadProvider(ctx context.Context, cfg R2UploadProviderConfig) (*R2UploadProvider, error) {
	endpointURL, err := resolveR2EndpointURL(cfg)
	if err != nil {
		return nil, err
	}

	bucketName := strings.TrimSpace(cfg.BucketName)
	if bucketName == "" {
		return nil, fmt.Errorf("invalid R2 config: bucket name is required")
	}

	accessKeyID := strings.TrimSpace(cfg.AccessKeyID)
	if accessKeyID == "" {
		return nil, fmt.Errorf("invalid R2 config: access key ID is required")
	}

	secretAccessKey := strings.TrimSpace(cfg.SecretAccessKey)
	if secretAccessKey == "" {
		return nil, fmt.Errorf("invalid R2 config: secret access key is required")
	}

	publicBaseURL := strings.TrimSpace(cfg.PublicBaseURL)
	if publicBaseURL == "" {
		return nil, fmt.Errorf("invalid R2 config: public base URL is required")
	}

	region := strings.TrimSpace(cfg.Region)
	if region == "" {
		region = defaultR2Region
	}

	presignExpiry := cfg.PresignExpiry
	if presignExpiry <= 0 {
		presignExpiry = defaultR2PresignExpiry
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(
		ctx,
		awsconfig.WithRegion(region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, ""),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("load AWS config for R2: %w", err)
	}

	s3Client := s3.NewFromConfig(awsCfg, func(options *s3.Options) {
		options.BaseEndpoint = aws.String(endpointURL)
		options.UsePathStyle = true
	})

	return &R2UploadProvider{
		bucketName:    bucketName,
		publicBaseURL: publicBaseURL,
		presignClient: s3.NewPresignClient(s3Client),
		presignExpiry: presignExpiry,
		s3Client:      s3Client,
	}, nil
}

func (p *R2UploadProvider) CreateSignedUploadURL(ctx context.Context, req CreateSignedUploadURLRequest) (*CreateSignedUploadURLResponse, error) {
	if p == nil || p.presignClient == nil {
		return nil, ErrProviderNotConfigured
	}
	now := time.Now().UTC()

	objectKey := strings.TrimSpace(req.ObjectKey)
	if objectKey == "" {
		return nil, fmt.Errorf("invalid provider request: object key is required")
	}

	fileType := strings.TrimSpace(req.FileType)
	if fileType == "" {
		return nil, ErrInvalidFileType
	}

	presignedRequest, err := p.presignClient.PresignPutObject(
		ctx,
		&s3.PutObjectInput{
			Bucket:      aws.String(p.bucketName),
			Key:         aws.String(objectKey),
			ContentType: aws.String(fileType),
		},
		func(options *s3.PresignOptions) {
			options.Expires = p.presignExpiry
		},
	)
	if err != nil {
		return nil, fmt.Errorf("presign R2 upload URL: %w", err)
	}

	return &CreateSignedUploadURLResponse{
		UploadURL: presignedRequest.URL,
		FileURL:   p.FileURLForObjectKey(objectKey),
		ExpiresAt: now.Add(p.presignExpiry),
	}, nil
}

func (p *R2UploadProvider) FileURLForObjectKey(objectKey string) string {
	if p == nil {
		return ""
	}

	return joinURL(p.publicBaseURL, objectKey)
}

// Delete objects with given keys in batches of up to 1000
func (p *R2UploadProvider) DeleteObjects(ctx context.Context, objectKeys []string) (map[string]error, error) {
	failed := map[string]error{}

	if p == nil || p.s3Client == nil {
		return nil, ErrProviderNotConfigured
	}

	cleaned := make([]string, 0, len(objectKeys))
	for _, key := range objectKeys {
		if strings.TrimSpace(key) != "" {
			cleaned = append(cleaned, key)
		}
	}

	const batchSize = 1000
	for i := 0; i < len(cleaned); i += batchSize {
		end := i + batchSize
		if end > len(cleaned) {
			end = len(cleaned)
		}

		ids := make([]types.ObjectIdentifier, 0, end-i)
		for _, key := range cleaned[i:end] {
			ids = append(ids, types.ObjectIdentifier{Key: aws.String(key)})
		}

		out, err := p.s3Client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
			Bucket: aws.String(p.bucketName),
			Delete: &types.Delete{
				Objects: ids,
				Quiet:   aws.Bool(false),
			},
		})
		if err != nil {
			for _, key := range cleaned[i:end] {
				failed[key] = err
			}
			continue
		}

		for _, e := range out.Errors {
			if e.Key != nil {
				failed[aws.ToString(e.Key)] = fmt.Errorf("R2 delete error: %s", aws.ToString(e.Message))
			}
		}
	}

	if len(failed) > 0 {
		return failed, fmt.Errorf("failed to delete %d objects from R2", len(failed))
	}
	return failed, nil
}

func resolveR2EndpointURL(cfg R2UploadProviderConfig) (string, error) {
	endpointURL := strings.TrimSpace(cfg.EndpointURL)
	if endpointURL != "" {
		return endpointURL, nil
	}

	accountID := strings.TrimSpace(cfg.AccountID)
	if accountID == "" {
		return "", fmt.Errorf("invalid R2 config: account ID or endpoint URL is required")
	}

	return fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID), nil
}
