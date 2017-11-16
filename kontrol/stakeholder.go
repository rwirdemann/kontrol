package kontrol

// Stakeholder types
const (
	StakeholderTypeEmployee = "employee"
	StakeholderTypePartner  = "partner"
	StakeholderTypeCompany  = "company"
	StakeholderTypeExtern   = "extern"
)

type Stakeholder struct {
	Id   string
	Name string
	Type string
}
