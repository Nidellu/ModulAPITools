package controller

import (
	"fmt"
	"time"

	"github.com/jasonlvhit/gocron"
)

func startScheduler() {
	startTime := time.Now()

	scheduler := gocron.NewScheduler()

	scheduler.Every(1).Minutes().Do(func() {
		elapsed := time.Since(startTime)
		minutes := int(elapsed.Minutes())
		fmt.Printf("%d minutes left for token to expired", 5-minutes)
		if minutes >= 5 {
			fmt.Println("Token is expired")
		}
	})

	scheduler.Start()

	time.Sleep(6 * time.Minute)

	scheduler.Clear()
}
