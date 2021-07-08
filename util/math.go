package util

func Net(gross float64) float64 {
	return gross / (1.0 + 0.19)
}

// in 2020 teh Umsatzsteuer became a function of month and year...
func Net2020(gross float64, year int, month int) float64 {
	var ust float64

	ust = 0.19
	if (year == 2020) && (month > 6) {
		ust = 0.16
	}
	return gross / (1.0 + ust)
}
