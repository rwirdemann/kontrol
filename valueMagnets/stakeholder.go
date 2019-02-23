package valueMagnets

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"encoding/json"
	"time"
	"github.com/ahojsenn/kontrol/util"
		)

// Stakeholder types
const (
	StakeholderTypeEmployee = "Employee"
	StakeholderTypePartner  = "Partner"
	StakeholderTypeCompany  = "Company"
	StakeholderTypeExtern   = "Extern"
	StakeholderTypeOthers   = "Rest"
	StakeholderTypeKUA      = "Unterkonto"
)

type Stakeholder struct {
	Id   string `json:",omitempty"`
	Name string
	Type string
	Arbeit string
	Fairshares string
	YearlySaldo float64
}


var StakeholderKM = Stakeholder{Id: "K", Name: "k:  Kommitment", Type: StakeholderTypeCompany, Arbeit: "1", Fairshares: "0"}
var StakeholderEX = Stakeholder{Id: "EX", Name: "k:  Extern", Type: StakeholderTypeExtern, Arbeit: "1", Fairshares: "0"}
var StakeholderRR = Stakeholder{Id: "RR", Name: "k:  Buchungsreste AR like Reisekosten etc.", Type: StakeholderTypeOthers, Arbeit: "1", Fairshares: "0"}



type Kommitmenschen struct {
	Id string `json:"Id"`
	Name string `json:"Name"`
	Type string `json:"Type"`
	Arbeit string `json:"Arbeit"`
	FairShares string `json:"Fairshares"`
}

type KommitmenschenRepository struct {
	Abrechenzeitpunkt string `json:"Abrechenzeitpunkt"`
	Menschen []Kommitmenschen `json:"Kommitmenschen"`
}

var kmry []KommitmenschenRepository

func (this KommitmenschenRepository) Init(year int)  {
	env := util.GetEnv()

	rawFile, err := ioutil.ReadFile(env.KommitmenschenFile)
	if err != nil {
		fmt.Println("in KommitmenschenRepository.All(), file: ", env)
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = json.Unmarshal(rawFile, &kmry)
	if err != nil {
		fmt.Println("in KommitmenschenRepository.All(), cannot Unmarshal rawfile... ", env.KommitmenschenFile)
		fmt.Println(err.Error())
		os.Exit(1)
	}

}


func (this KommitmenschenRepository) All(year int) []Kommitmenschen {

	if len(kmry) == 0 {
		this.Init(year)
	}

	// find the right year
	for i,yrep := range kmry {
		layout := "2006-01-02"
		t, err := time.Parse(layout, yrep.Abrechenzeitpunkt)

		if err != nil {
			fmt.Println(err)
		}
		if year == t.Year() {
			return kmry[i].Menschen
		}

	}
	return kmry[0].Menschen
}




var StakeholderRepository []Stakeholder


// generates the initial stakeholder
func (this *Stakeholder) Init(year int, shptr *[]Stakeholder) *[]Stakeholder {

	sh := *shptr
	kmrepo := KommitmenschenRepository{}
	for _, mensch := range kmrepo.All(year) {
		s := Stakeholder{}
		s.Type = mensch.Type
		sh = append(sh, Stakeholder{Id: mensch.Id, Name: mensch.Name, Type: mensch.Type, Arbeit: mensch.Arbeit, Fairshares: mensch.FairShares})
	}

	// add kommitment company
	sh = append(sh, StakeholderKM)

	// add externals
	sh = append(sh, StakeholderEX)

	// add Stakeholder for booking rests like Fakturierte Reisekosten etc. RR
	sh = append(sh, StakeholderRR)

	return &sh
}



// returns an array with a copy of all stakeholders
func (this *Stakeholder) All(year int) []Stakeholder {

	if len(StakeholderRepository) == 0 {
		StakeholderRepository = *this.Init(year, &StakeholderRepository)
	}
	return StakeholderRepository

}


func (this *Stakeholder) IsValidStakeholder (stakeholderId string) bool {

	for _, sh := range this.All(util.Global.FinancialYear) {
		if sh.Id == stakeholderId  {
			return true
		}
	}
	log.Println("in IsValidStakeholder: Warning! Unknown Stakeholder", stakeholderId)
	return false
}

func (this *Stakeholder) TypeOf(id string) string {

	for _, s := range this.All(util.Global.FinancialYear) {
		if s.Id == id ||
			id == StakeholderEX.Id  ||
			id == StakeholderRR.Id ||
			id == StakeholderKM.Id {
			return s.Type
		}
	}
	panic(fmt.Sprintf("stakeholder '%s' not found", id))
}

func (this *Stakeholder) Get(id string) Stakeholder {

	for _,s := range this.All(util.Global.FinancialYear) {
		if s.Id == id {
			return s
		}
	}
	panic(fmt.Sprintf("in Stakeholder.Get: stakeholder '%s' not found", id))
}

// return a array of pointers to selected stakeholders
func (this *Stakeholder) GetAllOfType(typ string) []Stakeholder {
	var stakeholders []Stakeholder
	for _,s := range this.All(util.Global.FinancialYear) {
		if s.Type == typ {
			// fill it with a pointer to the original stakeholder
			stakeholders = append(stakeholders, s)
		}
	}
	return stakeholders
}

// check if this is an employee
func  (this *Stakeholder) IsEmployee (id string) bool {
	return (id != "" && this.Get(id).Type == StakeholderTypeEmployee)
}

// check if this is an kommanditist
func  (this *Stakeholder) IsPartner (id string) bool {
	return ( id != "" && this.Get(id).Type == StakeholderTypePartner)
}