package dispatcher

import (
	"github.com/gelleson/packup/internal/core/service"
	"github.com/gelleson/packup/pkg/database"
	"github.com/sirupsen/logrus"
	"os"
)

type Dispatcher struct {
	config          Config
	logger          *logrus.Logger
	dispatcherInit  bool
	db              *database.Database
	backupService   *service.BackupService
	aclService      *service.AclService
	groupService    *service.GroupService
	snapshotService *service.SnapshotService
	userService     *service.UserService
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

	if err := dispatcher.initServices(); err != nil {
		return nil, err
	}

	return dispatcher, nil
}

func (d *Dispatcher) init() error {

	d.logger.Debugln("dispatcher init process started")

	if err := d.config.DatabaseConfig.Validate(); err != nil {
		d.logger.Error(err)
		return err
	}

	d.logger.Debugln("database config validated")

	d.db = database.NewDatabase(database.Config{
		DSN: d.config.DatabaseConfig.URL,
	})

	d.logger.Debugln("database pre connect")

	if err := d.db.Connect(); err != nil {
		d.logger.Error(err)
		return err
	}

	d.logger.Debugln("database post connect")

	d.dispatcherInit = true

	return nil
}

func (d *Dispatcher) initServices() error {

	d.logger.Debugln("init services process started")

	storage, err := d.config.StorageConfig.GetProvider()

	if err != nil {
		return err
	}

	d.backupService = service.NewBackupService(d.db)
	d.groupService = service.NewGroupService(d.db)
	d.aclService = service.NewAclService(d.db, d.groupService)
	d.snapshotService = service.NewSnapshotService(d.db, d.backupService, storage)
	d.userService = service.NewUserService(d.db, d.groupService)

	d.logger.Debugln("init services process finished")

	return nil
}
