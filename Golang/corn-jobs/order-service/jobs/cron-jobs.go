package jobs

import (
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

// --- Job 1: Clean expired tokens ---
func CleanExpiredTokens() {
	defer recoverJob("CleanExpiredTokens")
	fmt.Println(" Cleaning expired tokens:", time.Now())
	// simulate work
	time.Sleep(2 * time.Second)
	// log.Println("Tokens cleaned successfully")
}

// --- Job 2: Send weekly report ---
func SendWeeklyReportToAdmin() {
	defer recoverJob("SendWeeklyReportToAdmin")
	fmt.Println(" Sending weekly report:", time.Now())
	time.Sleep(3 * time.Second)
	// log.Println("Report sent successfully")
}

// --- Recovery and logging ---
func recoverJob(jobName string) {
	if r := recover(); r != nil {
		log.Printf(" Panic recovered in %s: %v\n", jobName, r)
	}
}

// --- Register cron jobs ---
func StartCronJobs() {
	c := cron.New()

	// Clean expired tokens every hour
	c.AddFunc("@every 1h", func() {
		go withRetry(CleanExpiredTokens, 3)
	})

	// Send report every Monday 9 AM
	c.AddFunc("0 9 * * 1", func() {
		go withRetry(SendWeeklyReportToAdmin, 3)
	})

	c.Start()
	log.Println(" Cron jobs started...")
}


func withRetry(task func(), retries int) {
	for i := 0; i < retries; i++ {
		func() {
			defer recoverJob("withRetry")
			task()
		}()
		time.Sleep(2 * time.Second)
	}
}
