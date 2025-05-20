package schedule

import (
	"github.com/robfig/cron/v3"
	"sonarqube-ouath-async/async"
)

func Run() {
	c := cron.New()
	// Synchronize once every day at 3:30 AM
	c.AddFunc("CRON_TZ=Asia/Shanghai 30 03 * * *", async.RunSyncToSonar)

	c.Start()
}
