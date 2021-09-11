package storage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"io/ioutil"
	"path"
)

type S3Provider struct {
	client *s3.S3
	bucket string
}

func NewS3Provider(session *session.Session, bucket string) *S3Provider {
	return &S3Provider{bucket: bucket, client: s3.New(session)}
}

func (s S3Provider) Put(namespace string, filename string, body io.Reader) (string, error) {

	file, err := ioutil.TempFile("", "")

	if err != nil {
		return "", err
	}

	if _, err := io.Copy(file, body); err != nil {
		return "", err
	}

	_, err = s.client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path.Join(namespace, filename)),
		Body:   file,
	})

	if err != nil {
		return "", err
	}

	return filename, nil
}

func (s S3Provider) Get(namespace string, id string) (io.ReadCloser, error) {

	obj, err := s.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path.Join(namespace, id)),
	})

	if err != nil {
		return nil, err
	}

	return obj.Body, nil
}
