package owner

// Stakeholder types
const (
	StakeholderTypeBank            = "bank"
	StakeholderTypeEmployee        = "employee"
	StakeholderTypePartner         = "partner"
	StakeholderTypeCompany         = "company"
	StakeholderTypeExtern          = "extern"
	StakeholderTypeOthers          = "Rest"
	StakeholderTypeInternalAccount = "internalAccount"
)

type Stakeholder struct {
	Id   string `json:",omitempty"`
	Name string
	Type string
}
