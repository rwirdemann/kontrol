package kontrol

const (
	STAKEHOLDER_TYPE_EMPLOYEE = "employee"
	STAKEHOLDER_TYPE_PARTNER  = "partner"
	STAKEHOLDER_TYPE_COMPANY  = "company"
	STAKEHOLDER_TYPE_EXTERN   = "extern"
)

type Stakeholder struct {
	Id   string
	Name string
	Type string
}
