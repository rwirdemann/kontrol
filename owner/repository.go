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

var StakeholderRW = Stakeholder{Id: "RW", Name: "Ralf Wirdemann", Type: StakeholderTypePartner}
var StakeholderAN = Stakeholder{Id: "AN", Name: "Anke Nehrenberg", Type: StakeholderTypePartner}
var StakeholderJM = Stakeholder{Id: "JM", Name: "Johannes Mainusch", Type: StakeholderTypePartner}
var StakeholderBW = Stakeholder{Id: "BW", Name: "Ben Wiedenmann", Type: StakeholderTypeEmployee}
var StakeholderKR = Stakeholder{Id: "KR", Name: "Katja Roth", Type: StakeholderTypeEmployee}
var StakeholderKM = Stakeholder{Id: "K", Name: "Kommitment", Type: StakeholderTypeCompany}
var StakeholderEX = Stakeholder{Id: "EX", Name: "Extern", Type: StakeholderTypeExtern}
var StakeholderRR = Stakeholder{Id: "RR", Name: "Rest", Type: StakeholderTypeOthers}
var StakeholderRueckstellung = Stakeholder{Id: "Rückstellung", Name: "Rückstellung", Type: StakeholderTypeInternalAccount}
var KontoJUSVJ = Stakeholder{Id: "JahresüberschussVJ", Name: "JahresüberschussVJ", Type: StakeholderTypeInternalAccount}
var SKR03_1400 = Stakeholder{Id: "1400", Name: "1400 OPOS-Kunde", Type: StakeholderTypeInternalAccount}
var SKR03_1600 = Stakeholder{Id: "1600", Name: "1600 OPOS-Lieferant", Type: StakeholderTypeInternalAccount}
var SKR03_4100_4199 = Stakeholder{Id: "4100_4199", Name: "4100-4199 Löhne und Gehälter", Type: StakeholderTypeInternalAccount}

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
		StakeholderRueckstellung,
		KontoJUSVJ,
		SKR03_1400,
		SKR03_1600,
		SKR03_4100_4199,
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
