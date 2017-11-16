package util

import (
	"os"
)

func GetHostname() string {
	if name, err := os.Hostname(); err == nil {
		if name == "Ubuntu-1704-zesty-64-minimal" {
			return "94.130.79.196"
		}
	}
	return "localhost"
}
