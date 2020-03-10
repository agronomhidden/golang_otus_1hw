package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	ntpTime, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		log.Fatalf("unable to get remote time: %s", err)
	}
	fmt.Println(fmt.Sprintf("current time: %s", time.Now().Format("2006-01-02 15:04:05 -0700 MST")))
	fmt.Println(fmt.Sprintf("exact time: %s", ntpTime.Format("2006-01-02 15:04:05 -0700 MST")))
}
