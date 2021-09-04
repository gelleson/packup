package upload

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"io/ioutil"
	"path"
)

type S3Provider struct {
	client *s3.S3
	bucket string
}

func (s S3Provider) Put(ctx context.Context, namespace string, filename string, body io.Reader) (string, error) {

	file, err := ioutil.TempFile("", "")

	if err != nil {
		return "", err
	}

	if _, err := io.Copy(file, body); err != nil {
		return "", err
	}

	_, err = s.client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path.Join(namespace, filename)),
		Body:   file,
	})

	if err != nil {
		return "", err
	}

	return filename, nil
}

func (s S3Provider) Get(ctx context.Context, namespace string, id string) (io.ReadCloser, error) {

	obj, err := s.client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path.Join(namespace, id)),
	})

	if err != nil {
		return nil, err
	}

	return obj.Body, nil
}
