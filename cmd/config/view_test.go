package config

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var templateString = "%d hours %d minutes"

func TestConvertTimeMsToDuration(t *testing.T) {
	hours := int64(2)
	minutes := int64(32)
	total := (time.Hour.Milliseconds() * hours) + (time.Minute.Milliseconds() * minutes)

	out := convertTimeMsToDuration(total)
	assert.Equal(t, fmt.Sprintf(templateString, hours, minutes), out)
}
