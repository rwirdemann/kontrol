package processing

import (
	"math"
)

const Gewerbesteuer_Freibetrag 	= 24500.0
const Gewerbesteuer_Hebesatz	= 4.70
const Gewerbesteuer_Messbetrtag = 0.035


func berechne_Gewerbesteuer (gewinn float64) float64 {
	if gewinn <= Gewerbesteuer_Freibetrag {
		return 0.0
	}
	gewinn -= Gewerbesteuer_Freibetrag

	gewinn = math.Round(gewinn /100)*100

	return  gewinn * Gewerbesteuer_Messbetrtag * Gewerbesteuer_Hebesatz
}