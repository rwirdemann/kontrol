package owner

import (
	"fmt"
)

const (
	PartnerShare             = 0.7
	KommmitmentShare         = 0.25
	KommmitmentExternShare   = 0.95
	KommmitmentOthersShare   = 1.00
	KommmitmentEmployeeShare = 0.95
	PartnerProvision         = 0.05
)

var StakeholderRW = Stakeholder{Id: "RW", Name: "k: Ralf Wirdemann", Type: StakeholderTypePartner}
var StakeholderAN = Stakeholder{Id: "AN", Name: "k: Anke Nehrenberg", Type: StakeholderTypePartner}
var StakeholderJM = Stakeholder{Id: "JM", Name: "k: Johannes Mainusch", Type: StakeholderTypePartner}
var StakeholderBW = Stakeholder{Id: "BW", Name: "k: Ben Wiedenmann", Type: StakeholderTypeEmployee}
var StakeholderKR = Stakeholder{Id: "KR", Name: "k: Katja Roth", Type: StakeholderTypeEmployee}
var StakeholderKM = Stakeholder{Id: "K", Name: "k: Kommitment", Type: StakeholderTypeCompany}
var StakeholderEX = Stakeholder{Id: "EX", Name: "Extern", Type: StakeholderTypeExtern}
var StakeholderRR = Stakeholder{Id: "RR", Name: "Rest", Type: StakeholderTypeOthers}



type StakeholderRepository struct {
}

func (this StakeholderRepository) All() []Stakeholder {
	return []Stakeholder{
		StakeholderRW,
		StakeholderAN,
		StakeholderJM,
		StakeholderBW,
		StakeholderEX,
		StakeholderKM,
		StakeholderKR,
		StakeholderRR,
	}
}

func (this StakeholderRepository) TypeOf(id string) string {
	for _, s := range this.All() {
		if s.Id == id {
			return s.Type
		}
	}
	panic(fmt.Sprintf("stakeholder '%s' not found", id))
}
