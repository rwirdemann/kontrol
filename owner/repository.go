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
var StakeholderRueckstellung = Stakeholder{Id: "Rückstellung", Name: "Rückstellung", Type: SKR03}
var KontoJUSVJ = Stakeholder{Id: "JahresüberschussVJ", Name: "JahresüberschussVJ", Type: SKR03}
var SKR03_1400 = Stakeholder{Id: "1400", Name: "OPOS-Kunde 1400", Type: SKR03}
var SKR03_1600 = Stakeholder{Id: "1600", Name: "OPOS-Lieferant 1600", Type: SKR03}
var SKR03_Anlagen = Stakeholder{Id: "SKR03_Anlagen", Name: "Zugang Anlagen", Type: SKR03}
var SKR03_Anlagen25 = Stakeholder{Id: "SKR03_Anlagen25", Name: "Zugang Anlagen Ähnl.R&W", Type: SKR03}
var SKR03_Kautionen = Stakeholder{Id: "SKR03_Kautionen", Name: "SKR03_Kautionen 1525", Type: SKR03}

var SKR03_Umsatzerloese = Stakeholder{Id: "SKR03_Umsatzerloese", Name: "1 SKR03_Umsatzerloese 8100-8402", Type: SKR03}
var SKR03_4100_4199 = Stakeholder{Id: "4100_4199", Name: "3 Löhne und Gehälter 4100-4199", Type: SKR03}
var SKR03_Abschreibungen = Stakeholder{Id: "SKR03_Abschreibungen", Name: "4 Abschreibungen auf Anlagen 4822-4855", Type: SKR03}
var SKR03_sonstigeAufwendungen = Stakeholder{Id: "SKR03_sonstigeAufwendungen", Name: "5 sonstige Aufwendungen", Type: SKR03}
var SKR03_Steuern = Stakeholder{Id: "SKR03_Steuern", Name: "6 SKR03_Steuern 4320", Type: SKR03}

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
		SKR03_sonstigeAufwendungen,
		SKR03_Anlagen,
		SKR03_Anlagen25,
		SKR03_Abschreibungen,
		SKR03_Kautionen,
		SKR03_Umsatzerloese,
		SKR03_Steuern,
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
