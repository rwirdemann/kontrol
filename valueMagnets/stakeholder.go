package valueMagnets

import (
	"encoding/json"
	"fmt"
	"github.com/ahojsenn/kontrol/util"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
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


var StakeholderKM = Stakeholder{Id: "K", Name: "k:  kommitment", Type: StakeholderTypeCompany, Arbeit: "1", Fairshares: "0"}
var StakeholderEX = Stakeholder{Id: "Extern", Name: "k:  Extern", Type: StakeholderTypeExtern, Arbeit: "1", Fairshares: "0"}
var StakeholderRR = Stakeholder{Id: "Rest", Name: "k:  Buchungsreste AR like Reisekosten etc.", Type: StakeholderTypeOthers, Arbeit: "1", Fairshares: "0"}



type Kommitmenschen struct {
	Id string `json:"Id"`
	Name string `json:"Name"`
	Type string `json:"Type"`
	Arbeit string `json:"Arbeit"`
	FairShares string `json:"Fairshares"`
}

type KommimtmentYear struct {
	Abrechenzeitpunkt string `json:"Abrechenzeitpunkt"`
	Liquiditaetsbedarf string `json:"Liquiditaetsbedarf"`
	JahresAbschluss_done bool `json:"JahresAbschluss_done"`
	Menschen []Kommitmenschen `json:"Kommitmenschen"`
}

var kommitmentHistory []KommimtmentYear

func (this KommimtmentYear) Init(year int)  {
	env := util.GetEnv()

	rawFile, err := ioutil.ReadFile(env.KommitmentFile)
	if err != nil {
		fmt.Println("in KommimtmentYear.Init(), file: ", env)
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = json.Unmarshal(rawFile, &kommitmentHistory)
	if err != nil {
		fmt.Println("in KommimtmentYear.Init(), cannot Unmarshal rawfile... ", env.KommitmentFile)
		fmt.Println(err.Error())
		os.Exit(1)
	}

}

func ReadJahresAbschluss_done (year int) bool {
	reval := false
	// find the right year
	for i,yrep := range kommitmentHistory {
		layout := "2006-01-02"
		t, err := time.Parse(layout, yrep.Abrechenzeitpunkt)

		if err != nil {
			fmt.Println(err)
		}
		if year == t.Year() {
			reval = kommitmentHistory[i].JahresAbschluss_done
		}
	}

	return reval
}

func (this KommimtmentYear) Liqui(year int) float64 {

	if len(kommitmentHistory) == 0 {
		this.Init(year)
	}

	var liqui = 0.0
	// find the right year
	for i,yrep := range kommitmentHistory {
		layout := "2006-01-02"
		t, err := time.Parse(layout, yrep.Abrechenzeitpunkt)

		if err != nil {
			fmt.Println(err)
		}
		if year == t.Year() {
			liqui,_ = strconv.ParseFloat(kommitmentHistory[i].Liquiditaetsbedarf, 64)
		}
	}

	return liqui
}


func (this KommimtmentYear) All(year int) []Kommitmenschen {

	if len(kommitmentHistory) == 0 {
		this.Init(year)
	}

	// find the right year
	for i,yrep := range kommitmentHistory {
		layout := "2006-01-02"
		t, err := time.Parse(layout, yrep.Abrechenzeitpunkt)

		if err != nil {
			fmt.Println(err)
		}
		if year == t.Year() {
			return kommitmentHistory[i].Menschen
		}

	}
	return kommitmentHistory[0].Menschen
}




var StakeholderRepository []Stakeholder


// generates the initial stakeholder
func (this *Stakeholder) Init(year int, shptr *[]Stakeholder) *[]Stakeholder {

	sh := *shptr

	// add kommitment company
	sh = append(sh, StakeholderKM)

	kmrepo := KommimtmentYear{}
	for _, mensch := range kmrepo.All(year) {
		s := Stakeholder{}
		s.Type = mensch.Type
		sh = append(sh, Stakeholder{Id: mensch.Id, Name: mensch.Name, Type: mensch.Type, Arbeit: mensch.Arbeit, Fairshares: mensch.FairShares})
	}

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
//	log.Println("in Stakeholder.Get: stakeholder '%s' not found", id, util.Global.FinancialYear)
//	log.Println("in Stakeholder.Get: returning ", StakeholderKM, "  instead...")
	return StakeholderKM
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


// loop over all employees
func (this *Stakeholder) AllEmployees(year int) []Stakeholder {
	if len(StakeholderRepository) == 0 {
		StakeholderRepository = *this.Init(year, &StakeholderRepository)
	}
	var employees []Stakeholder
	for _,s := range this.All(util.Global.FinancialYear) {
		if s.Type == StakeholderTypeEmployee {
			employees = append(employees, s)
		}
	}
	return employees
}



// loop over all kommanditisten
func (this *Stakeholder) AllPartners(year int) []Stakeholder {
	if len(StakeholderRepository) == 0 {
		StakeholderRepository = *this.Init(year, &StakeholderRepository)
	}
	var employees []Stakeholder
	for _,s := range this.All(util.Global.FinancialYear) {
		if s.Type == StakeholderTypePartner {
			employees = append(employees, s)
		}
	}
	return employees
}