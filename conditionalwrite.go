package conditionalwrite

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

// ETag is an HTTP ETag.
type ETag string

// None is the empty ETag.
const None ETag = ""

// Client is an S3 client for a single object.
type Client struct {
	client *s3.Client
	bucket string
	key    string
}

// NewClient configures the underlying S3 client and returns an object-specific
// Client.
func NewClient(endpoint, user, pw, region, bucket, key string) *Client {
	c := s3.New(s3.Options{
		Region:       region,
		BaseEndpoint: aws.String(endpoint),
		DefaultsMode: aws.DefaultsModeStandard,
		Credentials: credentials.NewStaticCredentialsProvider(
			user,
			pw,
			"", /* session */
		),
		UsePathStyle:               true,
		RequestChecksumCalculation: aws.RequestChecksumCalculationWhenSupported,
		ResponseChecksumValidation: aws.ResponseChecksumValidationWhenSupported,
		HTTPClient: &http.Client{
			Transport: &http.Transport{},
		},
	})
	return &Client{client: c, bucket: bucket, key: key}
}

// Set writes to object storage, creating a new object or overwriting an
// existing object. If the supplied ETag is the empty string, Set fails if
// there's already an existing object. Otherwise, Set fails if the existing
// object doesn't match the supplied ETag.
//
// Use IsPreconditionFailed to distinguish ETag mismatches from other errors.
func (c *Client) Set(ctx context.Context, r io.Reader, previous ETag) (ETag, error) {
	input := &s3.PutObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(c.key),
		Body:   r,
	}
	if previous == "" {
		input.IfNoneMatch = aws.String("*")
	} else {
		input.IfMatch = aws.String(string(previous))
	}
	res, err := c.client.PutObject(ctx, input)
	if err != nil {
		return "", err
	}
	if res == nil || res.ETag == nil {
		return "", errors.New("no ETag")
	}
	return ETag(*res.ETag), nil
}

// CreateBucket attempts to create the configured bucket. If the bucket already
// exists and is writable by this client, this method does not return an error.
func (c *Client) CreateBucket(ctx context.Context) error {
	_, err := c.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(c.bucket),
	})
	if getSmithyCode(err) == "BucketAlreadyOwnedByYou" {
		return nil
	}
	return err
}

// IsPreconditionFailed checks whether an error indicates a conditional write
// failure.
func IsPreconditionFailed(err error) bool {
	return getSmithyCode(err) == "PreconditionFailed"
}

func getSmithyCode(err error) string {
	if err == nil {
		return ""
	}
	var e smithy.APIError
	if errors.As(err, &e) {
		return e.ErrorCode()
	}
	return ""
}
