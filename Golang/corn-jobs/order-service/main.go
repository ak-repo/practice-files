package main

import (
	"log"
	"order-service/jobs"
	"time"
)

func main() {
	log.Println(" Order service started...")

	jobs.StartCronJobs()

	// Keep service alive
	for {
		time.Sleep(10 * time.Minute)
	}
}
