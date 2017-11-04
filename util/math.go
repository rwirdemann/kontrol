package util

import (
	"bitbucket.org/rwirdemann/kontrol/kontrol"
)

func Net(gross float64) float64 {
	return gross / (1.0 + kontrol.Umsatzsteuer)
}
