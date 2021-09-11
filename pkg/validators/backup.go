package validators

import (
	"github.com/gelleson/packup/internal/core/model"
	"github.com/pkg/errors"
	"github.com/robfig/cron"
	"github.com/teambition/rrule-go"
	"time"
)

func IsValidExecutionValue(b model.Backup) error {

	switch b.ExecutionType {
	case model.OnceExecution:

		now := time.Now()

		if !now.Before(b.ExecutionTime) {
			return errors.New("execution_time should be future time")
		}

		return nil

	case model.RruleExecution:

		if _, err := rrule.StrToRRule(b.Rrule); err != nil {
			return err
		}

		return nil

	case model.CronExecution:

		parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

		if _, err := parser.Parse(b.Cron); err != nil {
			return err
		}

		return nil

	default:

		return nil
	}
}

func IsValidTimezone(timezone string) error {

	_, err := time.LoadLocation(timezone)

	return err
}
