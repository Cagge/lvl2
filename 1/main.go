package main

import (
	"fmt"

	"github.com/beevik/ntp"
	"github.com/labstack/gommon/log"
)

func main() {
	time, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		log.Error(err)
	}
	fmt.Println(time)
}
