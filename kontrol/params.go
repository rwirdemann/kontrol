package kontrol

const (
	NET_COL_BW = 19
	NET_COL_AN = 20
	NET_COL_RW = 21
	NET_COL_JM = 22
	NET_COL_EX = 23

	PartnerShare  = 0.7
	ExternalShare = 0.8
)

var SH_RW = Stakeholder{Id: "RW", Type: STAKEHOLDER_TYPE_PARTNER}
var SH_AN = Stakeholder{Id: "AN", Type: STAKEHOLDER_TYPE_PARTNER}
var SH_JM = Stakeholder{Id: "JM", Type: STAKEHOLDER_TYPE_PARTNER}
var SH_BW = Stakeholder{Id: "BW", Type: STAKEHOLDER_TYPE_EMPLOYEE}
var SH_KM = Stakeholder{Id: "KM", Type: STAKEHOLDER_TYPE_COMPANY}
var SH_EX = Stakeholder{Id: "EX", Type: STAKEHOLDER_TYPE_EXTERN}
