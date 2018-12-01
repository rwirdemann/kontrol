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
	// Jahresüberschuss ist nun ermittelt

	// Buchung auf Verrechnungskonto Jahresüberschuss
	jue,okay := as.Get(accountSystem.ErgebnisNachSteuern.Id)
	if !okay {
		log.Panic("in GuV, there is no accountSystem.ErgebnisNachSteuern.Id")
	}
	soll := booking.NewBooking(0,"Jahresüberschuss "+strconv.Itoa(util.Global.FinancialYear), "", "", "", "",nil,  -jahresueberschuss, "Buchung Jahresüberschuss", int(now.Month()), now.Year(), now)
	jue.Book(*soll)

	// und Buchung auf Verbindlichkeitenkonto
	verb,okay := as.Get(accountSystem.SKR03_920_Gesellschafterdarlehen.Id)
	if !okay {
		log.Panic("in GuV, there is no accountSystem.SKR03_920_Gesellschafterdarlehen.Id")
	}
	haben := booking.NewBooking(0,"Jahresüberschuss "+strconv.Itoa(util.Global.FinancialYear), "", "", valueMagnets.StakeholderKM.Id, "", nil,  jahresueberschuss, "Buchung Jahresüberschuss", int(now.Month()), now.Year(), now)
	verb.Book(*haben)

}

func Bilanz (as accountSystem.AccountSystem) {

	var konto *account.Account
	var okay bool
	var bk *booking.Booking
	now := time.Now().AddDate(0, 0, 0)


	// Aktiva
	for rownr, acc := range as.All() {
		if acc.Description.Type == account.KontenartAktiv {
			// Buchung auf SummeAktiva
			konto, okay = as.Get(accountSystem.SummeAktiva.Id)
			if !okay {
				log.Panic("in Bilanz, could not get account SummeAktiva")
			}

			bk = booking.NewBooking(
				rownr,
				acc.Description.Name+strconv.Itoa(util.Global.FinancialYear),
				"",
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
			konto,okay = as.Get(accountSystem.SummePassiva.Id)
			if !okay {
				log.Panic("in Bilanz, could not get account SummePassiva")
			}
			bk = booking.NewBooking(
				rownr,
				acc.Description.Name+strconv.Itoa(util.Global.FinancialYear),
				"",
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

func GenerateProjectControlling  (as accountSystem.AccountSystem) {

	// concat KontenartErtrag and KontenartAufwand into accList
	accList := as.GetByType(account.KontenartErtrag)
	for k, v := range as.GetByType(account.KontenartAufwand) {
		accList[k] = v
	}


	for _, acc := range accList {
		for _, bk := range acc.Bookings {

			// handle empty projects
			if bk.Project == "" {
				bk.Project = "emptyProject"
			}
			// check if there is an project account bk.Projects
			acc, exists := as.Get(bk.Project)
			if !exists {
				// create a new projects account
				acc = account.NewAccount(account.AccountDescription{Id: bk.Project, Name: bk.Project, Type: account.KontenartProject})
				as.Add(acc)
			}

			// subtract "ER" from "AR"
			sign := +1.0
			if (bk.Type == "ER") {
				sign = -1.0
			}

			// now create a booking in the appropriate projectAccount
			clonedBooking := booking.Booking{
				RowNr:       bk.RowNr,
				Amount:      sign*bk.Amount,
				Project:     bk.Project,
				Type:        bk.Type,
				CostCenter:  bk.CostCenter,
				Month:		 bk.Month,
				Year:        bk.Year,
				Text:        bk.Text,
				FileCreated: bk.FileCreated,
				BankCreated: bk.BankCreated,
			}
			// and book it to tha account
			acc.Book(clonedBooking)
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
			case accountSystem.SKR03_1900.Id: // Privatentnahmen
				BookFromKtoKommitmensch{AccSystem: as, Booking: bk}.run()
			case accountSystem.SKR03_920_Gesellschafterdarlehen.Id,
			     accountSystem.ErgebnisNachSteuern.Id:
				bk.Type = booking.CC_KommitmenschDarlehen
				BookFromKtoKommitmensch{AccSystem: as, Booking: bk}.run()
			default:
			}
		}
		acc.UpdateSaldo()
	}
}

func DistributeKTopf (as accountSystem.AccountSystem) {

	// now calculate, what is left in the k-box
	kaccount,_ := as.Get("K")
	kaccount.UpdateSaldo()
	ErgebnisNachSteuernKonto ,_ := as.Get(accountSystem.ErgebnisNachSteuern.Id)
	ErgebnisNachSteuernKonto.UpdateSaldo()
	ergebnisNS := -1.0*ErgebnisNachSteuernKonto.Saldo
	log.Println("    GuV saldo:", ergebnisNS)
	log.Println("   ",100*ShareHoldersShare,"% fairshare: ",ergebnisNS*ShareHoldersShare)

	shrepo := valueMagnets.StakeholderRepository{}

	// distribute Shareholders-Share
	for _,sh := range shrepo.All(util.Global.FinancialYear) {
		if sh.Type == valueMagnets.StakeholderTypePartner {
			log.Print("    stakeholder: ", sh.Id)
			arbeit,_ := strconv.ParseFloat(sh.Arbeit, 64)
			log.Print("      Arbeit: ", arbeit)
			fairshares,_ := strconv.ParseFloat(sh.Fairshares, 64)
			log.Print("      Fairshares: ", fairshares)
			fairshareAnteil :=  math.Round(100*arbeit * fairshares*ergebnisNS*ShareHoldersShare)/100


			// Fairshare Anteil buchen
			log.Print(sh.Id,"      Anteil aus Fairshares: ", fairshareAnteil,"€")
			now := time.Now().AddDate(0, 0, 0)
			sollacc,_  := as.Get(valueMagnets.StakeholderKM.Id)
			habenacc,_  := as.Get(sh.Id)

			anteil_fairshares := booking.Booking{
				Amount:      -fairshareAnteil,
				Type:        booking.CC_AnteilAusFairshares,
				CostCenter:  sh.Id,
				Text:         "Anteil aus fairshares("+strconv.FormatFloat(fairshares, 'f', 2, 64)+")",
				FileCreated: now,
				BankCreated: now,
			}
			sollacc.Book(anteil_fairshares)
			anteil_fairshares.Amount *= -1.0
			habenacc.Book(anteil_fairshares)
		}
	}
	rest := ergebnisNS*(1-ShareHoldersShare)
	log.Println("    rest to distribute: ", rest)

	// now distribute internal hours
	for _,sh := range shrepo.All(util.Global.FinancialYear) {
		if sh.Type == valueMagnets.StakeholderTypePartner {
			//now := time.Now().AddDate(0, 0, 0)
			//sollacc,_  := as.Get(valueMagnets.StakeholderKM.Id)
			habenacc,_  := as.Get(sh.Id)


			// Interne Stunden buchen
			habenacc.UpdateSaldo()
			log.Print(sh.Id,"      Anteil aus internen Stunden: ", habenacc.Internals,"€")
			/* die sind schon gebucht
			anteil_interneStunden := booking.Booking{
				Amount:      -habenacc.Internals,
				Type:        booking.CC_InterneStunden,
				CostCenter:  sh.Id,
				Text:         "Anteil aus Internen Stunden",
				FileCreated: now,
				BankCreated: now,
			}
			sollacc.Book(anteil_interneStunden)
			anteil_interneStunden.Amount *= -1.0
			habenacc.Book(anteil_interneStunden)
			*/
			rest -= habenacc.Internals
		}
	}
	log.Println("    rest after internal hours: ", rest)

	// now distribute Vertriebsprovision
	for _,sh := range shrepo.All(util.Global.FinancialYear) {
		if sh.Type == valueMagnets.StakeholderTypePartner {

			now := time.Now().AddDate(0, 0, 0)
			sollacc, _ := as.Get(valueMagnets.StakeholderKM.Id)
			habenacc, _ := as.Get(sh.Id)

			// Vertriebsprovision buchen
			provisions := sumOfProvisonsForStakeholder(*kaccount, sh) // sum the partners revenue
			habenacc.UpdateSaldo()
			log.Print(sh.Id,"      Anteil Vertriebsprovision: ", math.Round(provisions*100)/100,"€")
			anteil_Vertriebsprovision := booking.Booking{
				Amount:      -provisions,
				Type:        booking.CC_Vertriebsprovision,
				CostCenter:  sh.Id,
				Text:         "Anteil aus Vertiebsprovision",
				FileCreated: now,
				BankCreated: now,
			}
			sollacc.Book(anteil_Vertriebsprovision)
			anteil_Vertriebsprovision.Amount *= -1.0
			habenacc.Book(anteil_Vertriebsprovision)
			rest -= provisions
		}
	}
	log.Println("    rest after Vertriebsprovision: ", rest)

	sumPartnerFaktura := kaccount.KommitmenschNettoFaktura
	log.Println("    Sum of Partnerfaktura", sumPartnerFaktura)

	restToDistribute := rest
	// now determine the Partners Contribution
	for _,sh := range shrepo.All(util.Global.FinancialYear) {
		if sh.Type == valueMagnets.StakeholderTypePartner {

			now := time.Now().AddDate(0, 0, 0)
			sollacc,_  := as.Get(valueMagnets.StakeholderKM.Id)
			habenacc,_  := as.Get(sh.Id)

			// Erlösanteil buchen
			rev := sumOfBookingsForStakeholder(*kaccount, sh) // sum the partners revenue
			log.Print(sh.Id,"      this partners revenue: ", math.Round(rev),"€")
			log.Print(sh.Id,"      this partners revenue: ", math.Round(10000*rev/sumPartnerFaktura)/100,"%")
			erloesAnteil := math.Round(restToDistribute*100*rev/sumPartnerFaktura )/100
			log.Print(sh.Id,"      Anteil aus Erlösen: ", erloesAnteil,"€")
			anteil_erloese := booking.Booking{
				Amount:      -erloesAnteil,
				Type:        booking.CC_AnteilAusFaktura,
				CostCenter:  sh.Id,
				Text:        "Anteil aus Erloesen: ",
				FileCreated: now,
				BankCreated: now,
			}
			sollacc.Book(anteil_erloese)
			anteil_erloese.Amount *= -1.0
			habenacc.Book(anteil_erloese)

			habenacc.UpdateSaldo()
			log.Print("      Summe:  ", habenacc.Salesprv+habenacc.Internals+erloesAnteil+habenacc.AnteilAusFairshares,"€")
			habenacc.YearS = habenacc.Salesprv+habenacc.Internals+erloesAnteil+habenacc.AnteilAusFairshares
			rest -= erloesAnteil
		}
	}
	log.Println("    rest (should be zero): ", rest)

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

func sumOfProvisonsForStakeholder (ac account.Account, sh valueMagnets.Stakeholder) float64 {
	saldo := 0.0
	for _,bk := range ac.Bookings {
		if sh.Id == bk.CostCenter && bk.Type == booking.CC_Vertriebsprovision {
			saldo += bk.Amount
		}
	}
	return saldo
}
