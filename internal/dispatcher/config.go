package dispatcher

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gelleson/packup/pkg/storage"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

type Config struct {
	DatabaseConfig DatabaseConfig
	ServerConfig   ServerConfig
	EncryptConfig  EncryptConfig
	LoggerConfig   LoggerConfig
	StorageConfig  StorageConfig
}

type OutType string

const (
	Stdout OutType = "stdout"
	File   OutType = "file"
)

type LoggerConfig struct {
	Level            logrus.Level
	PrettyPrint      bool
	DisableTimestamp bool
	Out              OutType
	OutFileName      string
	PrintFuncName    bool
}

func (c LoggerConfig) GetOut() (io.Writer, error) {

	switch c.Out {
	case Stdout:
		return os.Stdout, nil
	case File:

		if c.OutFileName == "" {
			return nil, errors.New("filename of logger is zero value. Please, define filename")
		}

		file, err := os.OpenFile(c.OutFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend)

		if err != nil {
			return nil, err
		}

		return file, nil
	default:
		return os.Stdout, nil
	}
}

type DatabaseConfig struct {
	URL string
}

func (c DatabaseConfig) Validate() error {

	if c.URL == "" {
		return errors.New("url should be set up")
	}

	return nil
}

type ServerConfig struct {
}

type EncryptConfig struct {
}

type StorageConfig struct {
	Provider   string
	BaseFolder string

	AccessKey        string
	SecretKey        string
	S3ForcePathStyle *bool
	Bucket           string
	Endpoint         *string
	Region           *string
}

func (u StorageConfig) GetProvider() (storage.API, error) {

	switch u.Provider {
	case "file":
		return storage.NewFileProvider(u.BaseFolder), nil
	case "s3":
		config := aws.Config{
			S3ForcePathStyle: u.S3ForcePathStyle,
			Endpoint:         u.Endpoint,
			Region:           u.Region,
			Credentials:      credentials.NewStaticCredentials(u.AccessKey, u.SecretKey, ""),
		}

		if u.AccessKey != "" && u.SecretKey != "" {
			config.Credentials = credentials.NewStaticCredentials(u.AccessKey, u.SecretKey, "")
		}

		sess, err := session.NewSession(&config)

		if err != nil {
			return nil, err
		}

		return storage.NewS3Provider(sess, u.Bucket), nil
	default:
		return storage.NewFileProvider(u.BaseFolder), nil
	}
}
