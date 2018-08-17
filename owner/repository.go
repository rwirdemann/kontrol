package owner

import (
	"fmt"
		"io/ioutil"
	"os"
	"encoding/json"
		"time"
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




// environments and HTTPS certificate locations.
type KommitmenschenRepository struct {
	Abrechenzeitpunkt string `json:"Abrechenzeitpunkt"`
	Menschen []Kommitmenschen `json:"Kommitmenschen"`
}

type Kommitmenschen struct {
	Id string `json:"Id"`
	Name string `json:"Name"`
	Type string `json:"Type"`
	Arbeit string `json:"Arbeit"`
}

func (this KommitmenschenRepository) All(year int) []Kommitmenschen {
	rawFile, err := ioutil.ReadFile("./kommitmenschen.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var kme []KommitmenschenRepository
	json.Unmarshal(rawFile, &kme)

	// find the right year
	for i,yrep := range kme {
		layout := "2006-01-02"
		t, err := time.Parse(layout, yrep.Abrechenzeitpunkt)

		if err != nil {
			fmt.Println(err)
		}
		if year == t.Year() {
			return kme[i].Menschen
		}

	}

	return kme[0].Menschen
}



type StakeholderRepository struct {
}

func (this StakeholderRepository) All() []Stakeholder {

	kmrepo := KommitmenschenRepository{}
	stakehr := []Stakeholder{}
	for _, mensch := range kmrepo.All(2018) {
		s := Stakeholder{}
		s.Type = mensch.Type

		stakehr = append(stakehr, Stakeholder {Id: mensch.Id, Name: mensch.Name, Type: mensch.Type, Arbeit: mensch.Arbeit} )
	}

	return stakehr

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
