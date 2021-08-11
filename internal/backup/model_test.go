package backup

import (
	"github.com/gelleson/packup/pkg/compress"
	"github.com/gelleson/packup/pkg/validators"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBackup_Validate(t *testing.T) {

	tssc := Backup{
		Name:          "postgres",
		Compress:      compress.GzipType,
		Tag:           "alfa",
		Keystore:      "store",
		Bucket:        "domain",
		ExecutionType: RruleExecution,
		Rrule:         "hello",
	}

	if err := validators.Struct(tssc); err != nil {
		assert.Fail(t, err.Error())
	}
}
