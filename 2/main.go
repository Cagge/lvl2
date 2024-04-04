package main

import (
	"fmt"
	"strconv"

	"github.com/labstack/gommon/log"
)

func main() {
	BaseString := ""
	fmt.Println(BaseString)
	var SliceRune []rune
	multi := 0
	for _, v := range BaseString {
		if '0' <= v && v <= '9' {
			if multi != 0 {
				log.Error("brake string")
			}
			multi, _ = strconv.Atoi(string(v))
		} else if multi != 0 {
			for i := 0; i < multi; i++ {
				SliceRune = append(SliceRune, v)
			}
			multi = 0
		} else {
			SliceRune = append(SliceRune, v)
		}
	}
	fmt.Print(string(SliceRune))
}
