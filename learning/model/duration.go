package model

import (
	"strconv"
	"strings"
	"time"
)

type CourseDuration struct {
	duration int
}

func NewCourseDuration(seconds int) *CourseDuration {
	return &CourseDuration{duration: seconds}
}

func (d *CourseDuration) Format() (formated string) {
	// Convert to x seconds
	timeStr := strconv.Itoa(d.duration) + "s"
	// Convert x seconds to 1h:2m:3s
	duration, _ := time.ParseDuration(timeStr)
	durationStr := duration.String()
	// Convert to 1:2:3
	totalHours := strings.ReplaceAll(durationStr, "h", ":")
	totalHours = strings.ReplaceAll(totalHours, "m", ":")
	totalHours = strings.ReplaceAll(totalHours, "s", "")
	// Format using hh:mm:ss
	timeArr := strings.Split(totalHours, ":")
	/*for i, t := range timeArr {
		if len(t) == 1{
			t = "0"+t
		}
		if i != 0 {
			formated = formated + ":"
		}
		formated = formated + t
	}*/

	for i := 0; i <= 3; i++ {
		if i != 0 {
			formated = formated + ":"
		}
		if len(timeArr) != 3 {
			formated = formated + "00"
			continue
		}
		if len(timeArr) == 1 {
			timeArr[i] = "0" + timeArr[i]
		}
		formated = formated + timeArr[i]
	}

	return formated + " Hours"
}
