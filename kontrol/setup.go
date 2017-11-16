package kontrol

const (
	PartnerShare             = 0.7
	KommmitmentShare         = 0.25
	KommmitmentExternShare   = 0.95
	KommmitmentEmployeeShare = 0.95
	PartnerProvision         = 0.05
)

var SH_RW = Stakeholder{Id: "RW", Name: "Ralf Wirdemann", Type: StakeholderTypePartner}
var SH_AN = Stakeholder{Id: "AN", Name: "Anke Nehrenberg", Type: StakeholderTypePartner}
var SH_JM = Stakeholder{Id: "JM", Name: "Johannes Mainusch", Type: StakeholderTypePartner}
var SH_BW = Stakeholder{Id: "BW", Name: "Ben Wiedenmann", Type: StakeholderTypeEmployee}
var SH_KM = Stakeholder{Id: "K", Name: "Kommitment", Type: StakeholderTypeCompany}
var SH_EX = Stakeholder{Id: "EX", Name: "Extern", Type: StakeholderTypeExtern}

var AllStakeholder = []Stakeholder{SH_RW, SH_AN, SH_JM, SH_BW, SH_EX, SH_KM}

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
