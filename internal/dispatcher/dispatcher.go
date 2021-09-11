package dispatcher

import (
	"github.com/gelleson/packup/pkg/database"
	"github.com/sirupsen/logrus"
	"os"
)

type Dispatcher struct {
	config         Config
	logger         *logrus.Logger
	dispatcherInit bool
	db             *database.Database
}

func New(config Config) (*Dispatcher, error) {

	out, err := config.LoggerConfig.GetOut()

	if err != nil {
		return nil, err
	}

	dispatcher := &Dispatcher{
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
	}

	if err := dispatcher.init(); err != nil {
		return nil, err
	}

	return dispatcher, nil
}

func (d *Dispatcher) init() error {

	if err := d.config.DatabaseConfig.Validate(); err != nil {
		return err
	}

	d.db = database.NewDatabase(database.Config{
		DSN: d.config.DatabaseConfig.URL,
	})

	if err := d.db.Connect(); err != nil {
		return err
	}

	d.dispatcherInit = true

	return nil
}
