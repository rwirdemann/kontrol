package processing

import (
	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/util"
	"github.com/ahojsenn/kontrol/valueMagnets"
	"log"
	"math"
	"strconv"
	"time"
)

const (
	ShareHoldersShare        = 0.20
)

type Command interface {
	run()
}

// Implementiert den Kommitment-Verteilungsalgorithmus
func Process(accsystem accountSystem.AccountSystem, booking booking.Booking) {

	// Assign booking GuV and Bilanz accounts
	var command Command

	switch booking.Typ {
	case "GV":
		command = BookPartnerEntnahmeCommand{AccSystem: accsystem, Booking: booking}
	case "GV-Vorjahr":
		command = BookPartnerEntnahmeVorjahrCommand{AccSystem: accsystem, Booking: booking}
	case "AR":
		command = BookAusgangsrechnungCommand{AccSystem: accsystem, Booking: booking}
	case "ER":
		command = BookEingangsrechnungCommand{AccSystem: accsystem, Booking: booking}
	case "IS":
		command = BookInterneStundenCommand{AccSystem: accsystem, Booking: booking}
	case "SV-Beitrag":
		command = BookSVBeitragCommand{AccSystem: accsystem, Booking: booking}
	case "GWSteuer":
		command = BookGWSteuerCommand{AccSystem: accsystem, Booking: booking}
	case "Gehalt":
		command = BookGehaltCommand{AccSystem: accsystem, Booking: booking}
	case "LNSteuer":
		command = BookLNSteuerCommand{AccSystem: accsystem, Booking: booking}
	case "UstVZ":
		command = BookUstCommand{AccSystem: accsystem, Booking: booking}
	case "SKR03":
		command = BookSKR03Command{AccSystem: accsystem, Booking: booking}
	default:
		log.Println("in Process: unknown command",booking.Type, " in row", booking.RowNr)
	}

	command.run()

}

func GuV (as accountSystem.AccountSystem) {
	log.Println("in GuV")

	var jahresueberschuss float64

	for _, acc := range as.All() {
		if acc.Description.Type == account.KontenartAufwand ||  acc.Description.Type == account.KontenartErtrag {
			jahresueberschuss += acc.Saldo
		}
	}


	now := time.Now().AddDate(0, 0, 0)
	// Jarhesüberschuss ist nun ermittelt

	// Buchung auf Verrechnungskonto Jahresüberschuss
	jue,_ := as.Get(accountSystem.ErgebnisNachSteuern.Id)
	soll := booking.NewBooking(0,"Jahresüberschuss "+strconv.Itoa(util.Global.FinancialYear), "", "", "", nil,  jahresueberschuss, "Buchung Jahresüberschuss", int(now.Month()), now.Year(), now)
	jue.Book(*soll)

	// und Buchung auf Verbindlichkeitenkonto
	verb,_ := as.Get(accountSystem.SKR03_920_Gesellschafterdarlehen.Id)
	haben := booking.NewBooking(0,"Jahresüberschuss "+strconv.Itoa(util.Global.FinancialYear), "", "", valueMagnets.StakeholderKM.Id, nil,  jahresueberschuss, "Buchung Jahresüberschuss", int(now.Month()), now.Year(), now)
	verb.Book(*haben)

}

func Bilanz (as accountSystem.AccountSystem) {

	var konto *account.Account
	var bk *booking.Booking
	now := time.Now().AddDate(0, 0, 0)


	// Aktiva
	for rownr, acc := range as.All() {
		if acc.Description.Type == account.KontenartAktiv {
			// Buchung auf SummeAktiva
			konto,_ = as.Get(accountSystem.SummeAktiva.Id)
			bk = booking.NewBooking(
				rownr,
				acc.Description.Name+strconv.Itoa(util.Global.FinancialYear),
				"",
				"",
				"",
				nil,
				acc.Saldo,
				"SummeAktiva "+strconv.Itoa(util.Global.FinancialYear),
				int(now.Month()),
				now.Year(),
				now)
			konto.Book(*bk)
		}
	}


	// Passiva
	for rownr, acc := range as.All() {
		if acc.Description.Type == account.KontenartPassiv {
			// Buchung auf SummePassiva
			konto,_ = as.Get(accountSystem.SummePassiva.Id)
			bk = booking.NewBooking(
				rownr,
				acc.Description.Name+strconv.Itoa(util.Global.FinancialYear),
				"",
				"",
				"",
				nil,
				acc.Saldo,
				"SummePassiva "+strconv.Itoa(util.Global.FinancialYear),
				int(now.Month()),
				now.Year(),
				now)
			konto.Book(*bk)
		}
	}
}



func ErloesverteilungAnValueMagnets (as accountSystem.AccountSystem) {

	for _, acc := range as.All() {
		// loop through all accounts in accountSystem,
		// beware: All() returns no bookings, so account here has no bookings[]
		a, _ := as.Get(acc.Description.Id)
		for _, bk := range a.Bookings {

			// process bookings on GuV accounts
			switch acc.Description.Type {
			case account.KontenartAufwand:
				BookCostToCostCenter{AccSystem: as, Booking: bk}.run()
			case account.KontenartErtrag:
				BookRevenueToEmployeeCostCenter{AccSystem: as, Booking: bk}.run()
			}

			// now process other accounts like accountSystem.SKR03_1900.Id
			switch acc.Description.Id {
			case accountSystem.SKR03_1900.Id:
				BookCostToCostCenter{AccSystem: as, Booking: bk}.run()
			}
		}
	}
}

func DistributeKTopf (as accountSystem.AccountSystem) {

	// now calculate, what is left in the k-box
	kaccount,_ := as.Get("K")
	kaccount.UpdateSaldo()
	log.Println("    kommitment saldo:", kaccount.Saldo)
	log.Println("   ",100*ShareHoldersShare,"% fairshare: ",kaccount.Saldo*ShareHoldersShare)
	rest := kaccount.Saldo*(1-ShareHoldersShare)
	log.Println("    rest to distribute: ", rest)

	shrepo := valueMagnets.StakeholderRepository{}
	// now determine the sum of all partners revenue
	sum := 0.0
	for _,sh := range shrepo.All(util.Global.FinancialYear) {
		if sh.Type == valueMagnets.StakeholderTypePartner {
			// sum the partners revenue
			sum += sumOfBookingsForStakeholder(*kaccount, sh)
		}
	}
	log.Println("    allPartnerRevenue =", sum)
	log.Println("    PartnerShare =", rest / sum)

	// now determine the Partners Contribution
	for _,sh := range shrepo.All(util.Global.FinancialYear) {
		if sh.Type == valueMagnets.StakeholderTypePartner {
			// sum the partners revenue
			rev := sumOfBookingsForStakeholder(*kaccount, sh)

			log.Print("    stakeholder: ", sh.Id)
			arbeit,_ := strconv.ParseFloat(sh.Arbeit, 64)
			log.Print("      Arbeit: ", arbeit)
			fairshares,_ := strconv.ParseFloat(sh.Fairshares, 64)
			log.Print("      Fairshares: ", fairshares)
			fairshareAnteil :=  math.Round(100*arbeit * fairshares*kaccount.Saldo*ShareHoldersShare)/100
			log.Print("      revenue: ", math.Round(rev),"€")
			log.Print("      ", math.Round(100*rev/sum),"%")
			erloesAnteil := math.Round(100*rev/sum*kaccount.Saldo*(1-ShareHoldersShare) )/100
			log.Print("      Anteil aus Fairshares: ", fairshareAnteil,"€")
			log.Print("      Anteil aus Erlösen: ", erloesAnteil,"€")

			// Fairshare Anteil buchen
			now := time.Now().AddDate(0, 0, 0)

			soll := booking.Booking{
				Amount:      -fairshareAnteil,
				Type:        booking.CC_AnteilAusFairshares,
				CostCenter:  sh.Id,
				Text:        "Anteil aus fairshares",
				FileCreated: now,
				BankCreated: now,
			}
			sollacc,_  := as.Get("K")
			sollacc.Book(soll)

			haben := booking.Booking{
				Amount:      fairshareAnteil,
				Type:        booking.CC_AnteilAusFairshares,
				CostCenter:  sh.Id,
				Text:        "Anteil aus fairshares",
				FileCreated: now,
				BankCreated: now,
			}
			habenacc,_  := as.Get(sh.Id)
			habenacc.Book(haben)

			// Erlösanteil buchen
			soll = booking.Booking{
				Amount:      -erloesAnteil,
				Type:        booking.CC_AnteilAusFaktura,
				CostCenter:  sh.Id,
				Text:        "Anteil aus fairshares",
				FileCreated: now,
				BankCreated: now,
			}
			sollacc,_  = as.Get("K")
			sollacc.Book(soll)

			haben = booking.Booking{
				Amount:      erloesAnteil,
				Type:        booking.CC_AnteilAusFaktura,
				CostCenter:  sh.Id,
				Text:        "Anteil aus fairshares",
				FileCreated: now,
				BankCreated: now,
			}
			habenacc,_  = as.Get(sh.Id)
			habenacc.Book(haben)

		}
	}




}

func sumOfBookingsForStakeholder (ac account.Account, sh valueMagnets.Stakeholder) float64 {
 	saldo := 0.0
	for _,bk := range ac.Bookings {
		if sh.Id == bk.CostCenter && bk.Type == booking.CC_PartnerNettoFaktura {
			saldo += bk.Amount
		}
	}

	return saldo
}
