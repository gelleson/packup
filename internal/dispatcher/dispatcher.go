package dispatcher

import (
	"github.com/sirupsen/logrus"
	"os"
)

type Dispatcher struct {
	config Config
	logger *logrus.Logger
}

func New(config Config) (*Dispatcher, error) {

	out, err := config.LoggerConfig.GetOut()

	if err != nil {
		return nil, err
	}

	return &Dispatcher{
		config: config,
		logger: &logrus.Logger{
			Out:   out,
			Level: config.LoggerConfig.Level,
			Formatter: &logrus.JSONFormatter{
				PrettyPrint:      config.LoggerConfig.PrettyPrint,
				DisableTimestamp: config.LoggerConfig.DisableTimestamp,
			},
			ExitFunc:     os.Exit,
			Hooks:        make(logrus.LevelHooks),
			ReportCaller: config.LoggerConfig.PrintFuncName,
		},
	}, nil
}
