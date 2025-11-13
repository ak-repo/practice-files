package main

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New()

	// Shedule every 10s
	c.AddFunc("@every 3s", func() {
		fmt.Println("Running shedule task : ", time.Now())
	})

	c.Start()

	time.Sleep(time.Second * 10)

}
