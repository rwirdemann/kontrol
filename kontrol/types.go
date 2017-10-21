package kontrol

const (
	SA_KM = "Kommitment"
	SA_BW = "BW"
	SA_RW = "RW"
	SA_JM = "JM"
	SA_AN = "AN"
	SA_EX = "EX"

	NET_COL_BW = 19
	NET_COL_AN = 20
	NET_COL_RW = 21
	NET_COL_JM = 22
	NET_COL_EX = 23
)

type NetPosition struct {
	Stakeholder string
	Column      int
}

var NetPositions = []NetPosition{
	NetPosition{Stakeholder: SA_RW, Column: NET_COL_RW},
	NetPosition{Stakeholder: SA_AN, Column: NET_COL_AN},
	NetPosition{Stakeholder: SA_JM, Column: NET_COL_JM},
	NetPosition{Stakeholder: SA_BW, Column: NET_COL_BW},
	NetPosition{Stakeholder: SA_EX, Column: NET_COL_EX},
}

type Position struct {
	Typ        string
	CostCenter string
	Subject    string
	Amount     float64
	Year       int
	Month      int

	Net map[string]float64
}
