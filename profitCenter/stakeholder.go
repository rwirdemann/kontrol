package profitCenter

// Stakeholder types
const (
	StakeholderTypeBank            = "bank"
	StakeholderTypeEmployee        = "employee"
	StakeholderTypePartner         = "partner"
	StakeholderTypeCompany         = "company"
	StakeholderTypeExtern          = "extern"
	StakeholderTypeOthers          = "Rest"
	SKR03                          = "SKR03"
)

type Stakeholder struct {
	Id   string `json:",omitempty"`
	Name string
	Type string
	Arbeit string
}
