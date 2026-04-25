package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/isa0-gh/reader/internal/config"
)

type S3Client struct {
	client       *s3.Client
	presigner    *s3.PresignClient
	bucket       string
	cdnPrefix    string
	endpoint     string
	usePathStyle bool
}

func NewS3Client(cfg *config.Config) *S3Client {
	awsCfg := aws.Config{
		Region:      cfg.S3.Region,
		Credentials: credentials.NewStaticCredentialsProvider(cfg.S3.AccessKeyID, cfg.S3.SecretAccessKey, ""),
	}

	opts := []func(*s3.Options){
		func(o *s3.Options) { o.UsePathStyle = cfg.S3.UsePathStyle },
	}
	if cfg.S3.Endpoint != "" {
		opts = append(opts, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.S3.Endpoint)
		})
	}

	client := s3.NewFromConfig(awsCfg, opts...)
	return &S3Client{
		client:       client,
		presigner:    s3.NewPresignClient(client),
		bucket:       cfg.S3.Bucket,
		cdnPrefix:    cfg.CDN,
		endpoint:     cfg.S3.Endpoint,
		usePathStyle: cfg.S3.UsePathStyle,
	}
}

// PresignPut returns a presigned PUT URL for direct browser upload.
func (s *S3Client) PresignPut(ctx context.Context, key string) (string, error) {
	req, err := s.presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		return "", fmt.Errorf("presign: %w", err)
	}
	return req.URL, nil
}

// PublicURL returns the CDN or S3 URL for a key.
func (s *S3Client) PublicURL(key string) string {
	if s.cdnPrefix != "" {
		return strings.TrimSuffix(s.cdnPrefix, "/") + "/" + key
	}

	if s.endpoint != "" {
		endpoint := strings.TrimSuffix(s.endpoint, "/")
		if s.usePathStyle {
			return fmt.Sprintf("%s/%s/%s", endpoint, s.bucket, key)
		}
		// Virtual-host style with custom endpoint
		protocol := "https://"
		if strings.HasPrefix(endpoint, "http://") {
			protocol = "http://"
			endpoint = strings.TrimPrefix(endpoint, "http://")
		} else if strings.HasPrefix(endpoint, "https://") {
			endpoint = strings.TrimPrefix(endpoint, "https://")
		}
		return fmt.Sprintf("%s%s.%s/%s", protocol, s.bucket, endpoint, key)
	}

	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucket, key)
}

// Bucket returns the configured bucket name.
func (s *S3Client) Bucket() string { return s.bucket }

// KeyFromURL extracts the S3 key from a public URL, or returns "" if unrecognized.
func (s *S3Client) KeyFromURL(url string) string {
	if s.cdnPrefix != "" && strings.HasPrefix(url, strings.TrimSuffix(s.cdnPrefix, "/")) {
		return strings.TrimPrefix(url, strings.TrimSuffix(s.cdnPrefix, "/")+"/")
	}

	if s.endpoint != "" {
		endpoint := strings.TrimSuffix(s.endpoint, "/")
		pathStylePrefix := fmt.Sprintf("%s/%s/", endpoint, s.bucket)
		if strings.HasPrefix(url, pathStylePrefix) {
			return strings.TrimPrefix(url, pathStylePrefix)
		}

		// Virtual host style check
		protocol := "https://"
		pureEndpoint := endpoint
		if strings.HasPrefix(endpoint, "http://") {
			protocol = "http://"
			pureEndpoint = strings.TrimPrefix(endpoint, "http://")
		} else if strings.HasPrefix(endpoint, "https://") {
			pureEndpoint = strings.TrimPrefix(endpoint, "https://")
		}
		vhPrefix := fmt.Sprintf("%s%s.%s/", protocol, s.bucket, pureEndpoint)
		if strings.HasPrefix(url, vhPrefix) {
			return strings.TrimPrefix(url, vhPrefix)
		}
	}

	prefix := fmt.Sprintf("https://%s.s3.amazonaws.com/", s.bucket)
	if strings.HasPrefix(url, prefix) {
		return strings.TrimPrefix(url, prefix)
	}
	return ""
}
func (s *S3Client) DeleteObject(ctx context.Context, bucket, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return err
}
