package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CourseDuration struct {
	Duration int
}

func (d *CourseDuration) Format() (formated string) {
	// Convert to x seconds
	timeStr := strconv.Itoa(d.Duration) + "s"
	// Convert x seconds to 1h:2m:3s
	duration, _ := time.ParseDuration(timeStr)
	durationStr := duration.String()
	// Convert to 1:2:3
	totalHours := strings.ReplaceAll(durationStr, "h", ":")
	totalHours = strings.ReplaceAll(totalHours, "m", ":")
	totalHours = strings.ReplaceAll(totalHours, "s", "")
	// Format using hh:mm:ss
	timeArr := strings.Split(totalHours, ":")

	index := 0
	for i := 0; i < 3; i++ {
		if i != 0 {
			formated = formated + ":"
		}
		if (len(timeArr) < 3 && i == 0) || (len(timeArr) == 1) {
			formated = formated + "00"
			continue
		}
		if len(timeArr[index]) == 1 {
			timeArr[index] = "0" + timeArr[index]
		}
		formated = formated + timeArr[index]
		index++
	}

	return formated + " Hours"
}

func (d *CourseDuration) Minute() (formated string) {
	// Convert to x seconds
	timeStr := strconv.Itoa(d.Duration) + "s"
	// Convert x seconds to 1h:2m:3s
	duration, _ := time.ParseDuration(timeStr)
	formated = fmt.Sprintf("%.f", duration.Minutes())

	return formated + " min"
}
