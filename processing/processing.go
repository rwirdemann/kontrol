package processing

import (
	"fmt"
	"github.com/ahojsenn/kontrol/account"
	"github.com/ahojsenn/kontrol/accountSystem"
	"github.com/ahojsenn/kontrol/booking"
	"github.com/ahojsenn/kontrol/util"
	"github.com/ahojsenn/kontrol/valueMagnets"
	"log"
	"math"
	"strconv"
	"strings"
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
		// ignore internal hours for this...
		// command = BookInterneStundenCommand{AccSystem: accsystem, Booking: booking}
		log.Println("in Process: skipping internal hours",booking.Type, " in row", booking.RowNr)
		command = DontDoAnything {AccSystem: accsystem, Booking: booking}
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

	var jahresueberschuss, gwsteuer float64
	now := time.Now().AddDate(0, 0, 0)

	for _, acc := range as.All() {

		switch  {
		case acc.Description.Type == account.KontenartAufwand, acc.Description.Type == account.KontenartErtrag:
			jahresueberschuss += acc.Saldo
			if (acc.Description.Id == accountSystem.SKR03_Steuern.Id) {
				gwsteuer += acc.Saldo
			}
		default:
		}
	}

	// calculate Gewerbesteuer
	log.Println("in GuV, Gewerbesteuer gebucht:", gwsteuer)
	log.Println("in GuV, Gewinn vor Steuer:", jahresueberschuss-gwsteuer)
	log.Println("in GuV, GWsteuer:", berechne_Gewerbesteuer(jahresueberschuss-gwsteuer))
	gwsRück := math.Round( 100* (berechne_Gewerbesteuer(jahresueberschuss-gwsteuer) + gwsteuer ) /100 )


	log.Println("in GuV, Gewerbesteuer-Rückstellung", gwsRück)

	// ermittelte GWSteuer Rückstellung verbuchen
	gwsKonto,_ := as.Get(accountSystem.SKR03_Steuern.Id)
	gwsSoll := booking.NewBooking(0,"in kontrol ermittelte Gewerbesteuer-Rückstellung "+strconv.Itoa(util.Global.FinancialYear), "4320", "956", "", "",nil,  -gwsRück, ("in kontrol ermittelte Gewerbesteuer-Rückstellung "+strconv.Itoa(util.Global.FinancialYear)), int(now.Month()), now.Year(), now)
	gwsKonto.Book(*gwsSoll)
	//
	gwsGegenKonto,_ := as.Get(accountSystem.SKR03_Rueckstellungen.Id)
	gwsHaben := booking.NewBooking(0,"in kontrol ermittelte Gewerbesteuer-Rückstellung "+strconv.Itoa(util.Global.FinancialYear), "", "", "", "",nil,  gwsRück, ("in kontrol ermittelte Gewerbesteuer-Rückstellung "+strconv.Itoa(util.Global.FinancialYear)), int(now.Month()), now.Year(), now)
	gwsGegenKonto.Book(*gwsHaben)


	// ermittelte GWSteuer Rückstellung von jahresueberschuss abziehen
	jahresueberschuss -= gwsRück
	log.Println("in GuV, Gewinn nach Steuer:", jahresueberschuss)


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

	log.Printf("in GuV, Jahresüberschuss: %6.2f€\n", math.Round(100*jahresueberschuss)/100)
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





func ErloesverteilungAnStakeholder (as accountSystem.AccountSystem) {
	log.Println("in ErloesverteilungAnStakeholder:")

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
			case account.KontenartAktiv:
				// now process other accounts like accountSystem.SKR03_1900.Id
				// this applies only to kommanditisten
				switch acc.Description.Id {
				case accountSystem.SKR03_920_Gesellschafterdarlehen.Id:
					bk.Type = booking.CC_KommitmenschDarlehen
					BookFromKtoKommitmensch{AccSystem: as, Booking: bk, SubAcc: accountSystem.UK_Darlehen}.run()
				case accountSystem.SKR03_Anlagen.Id,
					accountSystem.SKR03_Anlagen25_35.Id:
					//
					BookFromKtoKommitmenschenByShares{AccSystem: as, Booking: bk, SubAcc: accountSystem.UK_VeraenderungAnlagen}.run()
				case accountSystem.SKR03_Abschreibungen.Id:
					//
					BookFromKtoKommitmenschenByShares{AccSystem: as, Booking: bk, SubAcc: accountSystem.UK_VeraenderungAnlagen}.run()
				default:
				}
			case account.KontenartPassiv:
				switch acc.Description.Id {
				case accountSystem.SKR03_1900.Id: // Privatentnahmen
					BookFromKtoKommitmensch{AccSystem: as, Booking: bk, SubAcc: accountSystem.UK_Entnahmen}.run()

				}
			}
		}
	}
}

func DistributeKTopf (as accountSystem.AccountSystem) accountSystem.AccountSystem {

	// now calculate, what is left in the k-box
	kaccount,_ := as.Get("K")
	ErgebnisNachSteuernKonto ,_ := as.Get(accountSystem.ErgebnisNachSteuern.Id)
	ergebnisNS := -1.0*ErgebnisNachSteuernKonto.Saldo
	log.Println("    GuV saldo:", ergebnisNS)

	shrepo := valueMagnets.Stakeholder{}

	rest := ergebnisNS
	log.Printf("    Amount to distrib. bef. cost: %2.2f€", rest)

	// Partner returns her/his personal cost
	log.Println ("    Partner tragen ihre Kosten:")
	for _,sh := range shrepo.GetAllOfType (valueMagnets.StakeholderTypePartner) {
		sollacc,_  := as.GetSubacc(sh.Id, accountSystem.UK_Kosten)
		log.Print("      sh: ", sh.Id, ", hat Kosten: ", math.Round(100*sollacc.Saldo)/100)
		sollacc.YearS = accountIfYearlyIncome(*sollacc)
		rest -= sollacc.Saldo
	}
	totalSumToDistribute := rest
	log.Printf("    rest: %2.2f€-%2.2f€=%2.2f€", ergebnisNS, (ergebnisNS - rest),totalSumToDistribute)


	// distribute Shareholders-Share
	shareHoldersShare := ShareHoldersShare*rest
	log.Printf("    dist. %2.2f%% ShareHoldersShare=%2.2f€",ShareHoldersShare, shareHoldersShare)
	for _,sh := range shrepo.GetAllOfType (valueMagnets.StakeholderTypePartner) {
		fairshares,_ := strconv.ParseFloat(sh.Fairshares, 64)
		fairshareAnteil :=  fairshares * shareHoldersShare

		// Fairshare Anteil buchen
		log.Printf("      %s %2.2f%% Fairshares: %2.2f€", sh.Id, 100*fairshares, fairshareAnteil)
		now := time.Now().AddDate(0, 0, 0)
		sollacc,_  := as.GetSubacc(valueMagnets.StakeholderKM.Id, accountSystem.UK_AnteileausFairshare)
		habenacc,_  := as.GetSubacc(sh.Id, accountSystem.UK_AnteileausFairshare)

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
		habenacc.YearS = accountIfYearlyIncome(*habenacc)

	}
	rest -= shareHoldersShare
	log.Println("    rest to distribute: ", rest)

	// now care for Vertriebsprovisionen, die sind schon verteilt...
	sumOfProvisions := 0.0
	for _,sh := range shrepo.GetAllOfType (valueMagnets.StakeholderTypePartner) {
		// now := time.Now().AddDate(0, 0, 0)
		// sollacc, _ := as.GetSubacc(valueMagnets.StakeholderKM.Id, accountSystem.UK_Vertriebsprovision)
		habenacc, _ := as.GetSubacc(sh.Id, accountSystem.UK_Vertriebsprovision)

		// Vertriebsprovision buchen
		sollacc,_ := as.GetSubacc(sh.Id, accountSystem.UK_Vertriebsprovision)
		provisions := sumOfProvisonsForStakeholder(*sollacc, sh) // sum the partners revenue
		sumOfProvisions += provisions

		log.Printf("      %s Anteil Vertriebsprov: %2.2f€", sh.Id, provisions)
		sollacc.YearS = accountIfYearlyIncome(*sollacc)

		rest -= provisions
		habenacc.YearS = provisions
	}
	provisionPercentage := sumOfProvisions/totalSumToDistribute
	log.Printf("      Vertriebsprov: %2.2f€ = %2.2f%%", sumOfProvisions, provisionPercentage)
	log.Println("    rest after Vertriebsprov.: ", math.Round(100*rest)/100)



	// now distribute 50% of the rest according to factor "Arbeit"
	// use account interne stunden fpr now...


	sumOfArbeit := 0.0
	for _,sh := range shrepo.GetAllOfType (valueMagnets.StakeholderTypePartner) {
		shArbeit,_ := strconv.ParseFloat(sh.Arbeit, 64)
		sumOfArbeit += shArbeit
	}
	log.Printf("      Sum Arbeit =  %2.2f years", sumOfArbeit)


	sumOfArbeitShare := 0.0
	restToDistributeByArbeit := rest*0.5
	for _,sh := range shrepo.GetAllOfType (valueMagnets.StakeholderTypePartner) {
		now := time.Now().AddDate(0, 0, 0)
		sollacc,_  := as.Get(valueMagnets.StakeholderKM.Id)
		habenacc,_  := as.GetSubacc(sh.Id, accountSystem.UK_AnteilMitmachen)

		// sumofArbeitShare buchen
		shArbeit,_ := strconv.ParseFloat(sh.Arbeit, 64)
		arbeitShare := restToDistributeByArbeit*shArbeit/sumOfArbeit
		log.Printf("      %s Anteil aus arbeitShare: %2.2f€", sh.Id, arbeitShare)
		anteil_erloese := booking.Booking{
			Amount:      -arbeitShare,
			Type:        booking.CC_AnteilAusFaktura,
			CostCenter:  sh.Id,
			Text:        "Anteil aus Erloesen: ",
			FileCreated: now,
			BankCreated: now,
		}
		sollacc.Book(anteil_erloese)
		anteil_erloese.Amount *= -1.0
		habenacc.Book(anteil_erloese)


		log.Printf("      %s Anteil ArbeitShare: %2.2f€", sh.Id, habenacc.Saldo)
		rest -= habenacc.Saldo
		habenacc.YearS = accountIfYearlyIncome(*habenacc)
	}


	log.Printf("      ArbeitShare: %2.2f€ = %2.2f%%", sumOfArbeitShare, sumOfArbeitShare/totalSumToDistribute)
	log.Println("    rest after ArbeitShare: ", math.Round(100*rest)/100)


	// Erlösanteile
	subacc,_ :=  as.GetSubacc(kaccount.Description.Id, accountSystem.UK_AnteileAuserloesen)
	sumPartnerFaktura := subacc.SumOfBookingType(booking.CC_PartnerNettoFaktura)
	log.Println("    Sum of Partnerfaktura", sumPartnerFaktura)
	restToDistribute := rest
	sumOfErloesAnteil := rest

	// now determine the Partners Contribution
	for _,sh := range shrepo.GetAllOfType (valueMagnets.StakeholderTypePartner) {

		now := time.Now().AddDate(0, 0, 0)
		sollacc,_  := as.Get(valueMagnets.StakeholderKM.Id)
		habenacc,_  := as.GetSubacc(sh.Id, accountSystem.UK_AnteileAuserloesen)

		// Erlösanteil buchen
		subacc,_ := as.GetSubacc(kaccount.Description.Id, accountSystem.UK_AnteileAuserloesen)
		rev := sumOfBookingsForStakeholder(*subacc, sh) // sum the partners revenue
		erloesAnteil := math.Round(restToDistribute*100*rev/sumPartnerFaktura )/100
		log.Printf("      %s revenue %2.2f%% = %2.0f€ / %2.0f€", sh.Id, rev/sumPartnerFaktura, rev, sumPartnerFaktura)
		log.Printf("      %s Anteil aus Erlösen: %2.2f€", sh.Id, erloesAnteil)
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

		// book the yearly sum from saldo to yearS
		habenacc.YearS = habenacc.Saldo
		rest -= erloesAnteil
	}
	log.Printf("      ErlösAnt.: %2.2f€ = %2.2f%%", sumOfErloesAnteil, sumOfErloesAnteil/totalSumToDistribute)
	log.Println("    rest (should be zero): ", math.Round(100*rest)/100)

	// now book the yearlySaldo on Stakeholder.YearlySaldo
	sum := 0.0
	for _,sh := range shrepo.GetAllOfType (valueMagnets.StakeholderTypePartner) {
		s := shrepo.Get(sh.Id)
		s.YearlySaldo =  StakeholderYearlyIncome (as, s.Id)
		sum += s.YearlySaldo
		log.Printf("DistributeKTopf: %s. YearlySaldo: %2.2f€", s.Id, s.YearlySaldo)
	}
	log.Printf("DistributeKTopf, sum: %2.2f€", sum)
	log.Printf("GuV saldo:            %2.2f€", ergebnisNS)

	return as

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



// loop through the employees, calculates their Bonusses and book them from employee to kommitment cost.
func CalculateEmployeeBonus (as accountSystem.AccountSystem) accountSystem.AccountSystem {

	shrepo := valueMagnets.Stakeholder{}
	for _,sh := range shrepo.GetAllOfType (valueMagnets.StakeholderTypeEmployee) {
		bonus := StakeholderYearlyIncome(as, sh.Id)

		// take care, this is idempotent, i.e. that the next bonus calculation overwrites the last one...
		if bonus > 0.0 {
			// only book positive bonusses of valuemagnets
			log.Println("in CalculateEmployeeBonus: ", sh.Id, math.Round(bonus*100)/100)

			now := time.Now().AddDate(0, 0, 0)

			// book in into GuV
			bk := booking.Booking{
				RowNr:       0,
				Amount:      bonus,
				Soll:		 "4120",
				Haben: 		 "965",
				Type:        booking.CC_Gehalt,
				CostCenter:  sh.Id,
				Text:        fmt.Sprintf("Bonusrückstellung für %s in %d", sh.Id, util.Global.FinancialYear),
				Month:       12,
				Year:        util.Global.FinancialYear,
				FileCreated: now,
				BankCreated: now,
			}
			BookSKR03Command{AccSystem: as, Booking: bk}.run()

			// book into kommitmentschen accountsystem
/*			habenAcc,_ := as.Get(sh.Id)
			sollAcc,_ := as.Get(valueMagnets.StakeholderKM.Id)
			bk.Amount *= -1.0
			habenAcc.Book(bk)
			bk.Amount *= -1.0
			sollAcc.Book(bk)
*/		}
	}
	return as
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


// a function that helps to return the yearly income of a kommitmensch...
func accountIfYearlyIncome (ac account.Account) float64 {
	switch ac.Description.Type {
	case accountSystem.UK_Darlehen, accountSystem.UK_Entnahmen:
		return 0.0
	default:
		return ac.Saldo
	}
}




// sum up whatever the stakeholder earned in the actual year
func StakeholderYearlyIncome (as accountSystem.AccountSystem, stkhldr string) float64 {

	yearsum := 0.0
	for _, acc := range as.All() {
		if strings.HasPrefix(acc.Description.Id, stkhldr) {
			// found the main or a subaccount of kommi
			yearsum += accountIfYearlyIncome (acc)
		}
	}
	return yearsum
}


