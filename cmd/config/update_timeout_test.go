package config

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"go-pass/model"
	"go-pass/utils"
)

func TestUpdateConfigTimeout(t *testing.T) {
	cfg := &model.Config{
		Timeout: utils.THIRTY_MINUTES,
	}

	cmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {},
	}
	cmd.Flags()
	cmd.Flags().IntP("hours", "q", 0, "")
	cmd.Flags().IntP("minutes", "m", 30, "")

	err := cmd.Flags().Set("hours", "1")
	assert.NoError(t, err)

	err = cmd.Flags().Set("minutes", "45")
	assert.NoError(t, err)

	err = updateConfigTimeout(cfg, cmd)
	assert.NoError(t, err)
	assert.NotEqual(t, utils.THIRTY_MINUTES, cfg.Timeout)
	assert.Equal(t, int64(6300000), cfg.Timeout)
}

func TestGetHoursAndMinutes(t *testing.T) {
	cmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {},
	}
	cmd.Flags()
	cmd.Flags().IntP("hours", "q", 0, "")
	cmd.Flags().IntP("minutes", "m", 30, "")

	err := cmd.Flags().Set("hours", "1")
	assert.NoError(t, err)
	err = cmd.Flags().Set("minutes", "45")
	assert.NoError(t, err)

	hM, mM, err := getHoursAndMinutes(cmd)
	assert.NoError(t, err)
	assert.Equal(t, int64(3600000), hM)
	assert.Equal(t, int64(2700000), mM)
}
