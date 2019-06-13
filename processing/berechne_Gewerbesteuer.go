package processing


const Gewerbesteuer_Freibetrag 	= 24.500
const Gewerbesteuer_Hebesatz	= 4.70
const Gewerbesteuer_Messbetrtag = 0.035


func berechne_Gewerbesteuer (gewinn float64) float64 {
	if gewinn < Gewerbesteuer_Freibetrag {
		return 0.0
	}

	return (gewinn - Gewerbesteuer_Freibetrag) * Gewerbesteuer_Hebesatz * Gewerbesteuer_Messbetrtag
}