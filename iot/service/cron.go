package service

import (
	"koudai-box/iot/db"

	status "koudai-box/iot/gateway/status"

	"github.com/robfig/cron"

	"github.com/sirupsen/logrus"
)

var dayCalcJob cron.Cron

const dayCalc = "0 0/30 0 * * *"

func StartCron() {
	dayCalcJob = *cron.New()
	err := dayCalcJob.AddFunc(dayCalc, func() {
		devices, err := db.QueryDevices()
		if err != nil {
			return
		}
		for _, d := range devices {
			preday := CalcPredayAvg(d.Id)
			dStatus := status.GetDeviceStatus(d.Id)
			dStatus.PreDayStatus = preday
		}
	})
	if err != nil {
		logrus.Errorf("dayCalcJob cron[%s] is error", dayCalc)
	}
	dayCalcJob.Start()
}

func StopCron() {
	dayCalcJob.Stop()
}
