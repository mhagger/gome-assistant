package gomeassistant

import (
	"time"

	"github.com/golang-module/carbon"
	i "github.com/saml-dev/gome-assistant/internal"
)

type conditionCheck struct {
	fail bool
}

func CheckWithinTimeRange(startTime, endTime string) conditionCheck {
	cc := conditionCheck{fail: false}
	// if betweenStart and betweenEnd both set, first account for midnight
	// overlap, then check if between those times.
	if startTime != "" && endTime != "" {
		parsedStart := i.ParseTime(startTime)
		parsedEnd := i.ParseTime(endTime)

		// check for midnight overlap
		if parsedEnd.Lt(parsedStart) { // example turn on night lights when motion from 23:00 to 07:00
			if parsedEnd.IsPast() { // such as at 15:00, 22:00
				parsedEnd = parsedEnd.AddDay()
			} else {
				parsedStart = parsedStart.SubDay() // such as at 03:00, 05:00
			}
		}

		// skip callback if not inside the range
		if !carbon.Now().BetweenIncludedStart(parsedStart, parsedEnd) {
			cc.fail = true
		}

		// otherwise just check individual before/after
	} else if startTime != "" && i.ParseTime(startTime).IsFuture() {
		cc.fail = true
	} else if endTime != "" && i.ParseTime(endTime).IsPast() {
		cc.fail = true
	}
	return cc
}

func CheckStatesMatch(listenerState, s string) conditionCheck {
	cc := conditionCheck{fail: false}
	// check if fromState or toState are set and don't match
	if listenerState != "" && listenerState != s {
		cc.fail = true
	}
	return cc
}

func CheckThrottle(throttle time.Duration, lastRan carbon.Carbon) conditionCheck {
	cc := conditionCheck{fail: false}
	// check if Throttle is set and that duration hasn't passed since lastRan
	if throttle.Seconds() > 0 &&
		lastRan.DiffAbsInSeconds(carbon.Now()) < int64(throttle.Seconds()) {
		cc.fail = true
	}
	return cc
}