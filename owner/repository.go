package owner

import (
	"fmt"
		"io/ioutil"
	"os"
	"encoding/json"
		"time"
	"github.com/ahojsenn/kontrol/util"
		)

var StakeholderKM = Stakeholder{Id: "K", Name: "k: Kommitment", Type: StakeholderTypeCompany, Arbeit : "100%"}
var StakeholderEX = Stakeholder{Id: "EX", Name: "Extern", Type: StakeholderTypeExtern, Arbeit : "100%"}
var StakeholderRR = Stakeholder{Id: "RR", Name: "Buchungsreste AR like Reisekosten etc.", Type: StakeholderTypeOthers, Arbeit : "0%"}

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

var kmry []KommitmenschenRepository

func (this KommitmenschenRepository) Init(year int)  {

	env := util.GetEnv()
	rawFile, err := ioutil.ReadFile(env.KommitmenschenFile)
	if err != nil {
		fmt.Println("in KommitmenschenRepository.All(), file: ", env)
		fmt.Println(err.Error())
		os.Exit(1)
	}

	json.Unmarshal(rawFile, &kmry)

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



type StakeholderRepository struct {
}

func (this StakeholderRepository) All(year int) []Stakeholder {

	kmrepo := KommitmenschenRepository{}
	stakehr := []Stakeholder{}
	for _, mensch := range kmrepo.All(year) {
		s := Stakeholder{}
		s.Type = mensch.Type
		stakehr = append(stakehr, Stakeholder{Id: mensch.Id, Name: mensch.Name, Type: mensch.Type, Arbeit: mensch.Arbeit})
	}

	// add kommitment company
	stakehr = append(stakehr, StakeholderKM)

	// add externals
	stakehr = append(stakehr, StakeholderEX)

	// add Stakeholder for booking ressts RR
	stakehr = append(stakehr, StakeholderRR)

	return stakehr
}

func (this StakeholderRepository) TypeOf(id string) string {


	for _, s := range this.All(util.Global{}.FinancialYear) {
		if s.Id == id {
			return s.Type
		}
	}
	panic(fmt.Sprintf("stakeholder '%s' not found", id))
}

func (this StakeholderRepository) Get(id string) Stakeholder {

	for _,s := range this.All(util.Global{}.FinancialYear) {
		if s.Id == id {
			return s
		}
	}
	panic(fmt.Sprintf("stakeholder '%s' not found", id))
}