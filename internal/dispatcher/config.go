package dispatcher

import (
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
