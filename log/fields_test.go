package log

import (
	"errors"
	"testing"
)

func TestFields(t *testing.T) {
	logf := Fields{"job_id": "100001"}
	logf.Fields(
		"mch_id", 1000067869,
		"error", errors.New("balabala"),
		"hahah", "hehehhe",
	)

	logf.Info("test")
}
