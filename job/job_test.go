package job

import (
	"context"
	"testing"
	"time"

	"github.com/gocraft/work"
	"github.com/stretchr/testify/assert"
)

type testJob struct {
	jobName    struct{} `Job:"job_name"`
	Int        int
	Int32      int32
	Int64      int64
	String     string
	Bool       bool
	SliceInt   []int
	SliceInt8  []int8
	SliceInt16 []int16
	SliceInt32 []int32
	SliceInt64 []int64
	Time       time.Time
}

func TestJob(t *testing.T) {
	job := testJob{
		Int:        1,
		Int32:      2,
		Int64:      3,
		String:     "4",
		Bool:       true,
		SliceInt:   []int{10, 20, 30},
		SliceInt8:  []int8{11, 21, 31},
		SliceInt16: []int16{12, 22, 32},
		SliceInt32: []int32{13, 23, 33},
		SliceInt64: []int64{14, 24, 34},
		Time:       time.Date(2020, 1, 1, 2, 2, 2, 0, time.UTC),
	}
	packArgs := map[string]interface{}{
		"Int":        float64(1),
		"Int32":      float64(2),
		"Int64":      float64(3),
		"String":     "4",
		"Bool":       true,
		"SliceInt":   []interface{}{float64(10), float64(20), float64(30)},
		"SliceInt8":  []interface{}{float64(11), float64(21), float64(31)},
		"SliceInt16": []interface{}{float64(12), float64(22), float64(32)},
		"SliceInt32": []interface{}{float64(13), float64(23), float64(33)},
		"SliceInt64": []interface{}{float64(14), float64(24), float64(34)},
		"Time":       "2020-01-01T02:02:02Z",
	}

	t.Run("Pack arguments", func(t *testing.T) {
		gotArgs, err := packArguments(context.Background(), job)

		assert.Nil(t, err)
		assert.Equal(t, packArgs, gotArgs)
	})

	t.Run("Pack arguments (pointer)", func(t *testing.T) {
		gotArgs, err := packArguments(context.Background(), &job)
		assert.Nil(t, err)
		assert.Equal(t, packArgs, gotArgs)
	})

	t.Run("Unpack arguments", func(t *testing.T) {
		gotJob := testJob{}
		err := UnpackArguments(context.Background(), &gotJob, &work.Job{Args: packArgs})
		assert.Nil(t, err)
		assert.Equal(t, job, gotJob)
	})
}
