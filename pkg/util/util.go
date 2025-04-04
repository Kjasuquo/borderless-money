package util

import (
	"math"
	"strings"
)

func RoundToTwoDecimalPlaces(value float64) float64 {
	multiplier := math.Pow(10, 2) // 100
	return math.Round(value*multiplier) / multiplier
}

func ParseUserAgent(userAgent string) (device, browser string) {

	if strings.Contains(userAgent, "Mobile") {
		device = "Mobile"
	} else {
		device = "Desktop"
	}

	if strings.Contains(userAgent, "Chrome") {
		browser = "Chrome"
	} else if strings.Contains(userAgent, "Firefox") {
		browser = "Firefox"
	} else if strings.Contains(userAgent, "Safari") {
		browser = "Safari"
	} else {
		browser = "Unknown"
	}

	return device, browser
}
