package kontrol

const (
	PartnerShare = 0.7
)

var SH_RW = Stakeholder{Id: "RW", Type: STAKEHOLDER_TYPE_PARTNER}
var SH_AN = Stakeholder{Id: "AN", Type: STAKEHOLDER_TYPE_PARTNER}
var SH_JM = Stakeholder{Id: "JM", Type: STAKEHOLDER_TYPE_PARTNER}
var SH_BW = Stakeholder{Id: "BW", Type: STAKEHOLDER_TYPE_EMPLOYEE}
var SH_KM = Stakeholder{Id: "KM", Type: STAKEHOLDER_TYPE_COMPANY}
var SH_EX = Stakeholder{Id: "EX", Type: STAKEHOLDER_TYPE_EXTERN}

var AllStakeholder = []Stakeholder{SH_RW, SH_AN, SH_JM, SH_BW, SH_EX}

// Beschreibt, dass die netto (Rechnungs-)Position in Spalte X der CSV-Datei dem Stakeholder Y gehört
type NetBookingColumn struct {
	Owner  Stakeholder
	Column int
}

// Liste aller Spalten-Stateholder Positions-Mappings
var NetBookings = []NetBookingColumn{
	NetBookingColumn{Owner: SH_RW, Column: 21},
	NetBookingColumn{Owner: SH_AN, Column: 20},
	NetBookingColumn{Owner: SH_JM, Column: 22},
	NetBookingColumn{Owner: SH_BW, Column: 19},
	NetBookingColumn{Owner: SH_EX, Column: 23},
}
