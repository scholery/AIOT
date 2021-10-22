package test

import (
	"testing"

	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

func TestCron(t *testing.T) {
	i := 0
	c := cron.New()
	spec := "*/5 * * * * ?"
	c.AddFunc(spec, func() {
		i++
		logrus.Info("cron running:", i)
	})
	c.AddFunc("@every 1s", func() {
		i++
		logrus.Info("cron running every:", i)
	})
	c.Start()
	for {

	}
}
