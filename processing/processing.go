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
		//log.Println("in Process: skipping internal hours",booking.Type, " in row", booking.RowNr)
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
	// only do that for the current year!
	log.Println("in GuV, Gewerbesteuer gebucht:", gwsteuer)
	log.Println("in GuV, Gewinn vor Steuer:", jahresueberschuss-gwsteuer)
	if  util.Global.FinancialYear == time.Now().Year() {
		log.Println("in GuV, GWsteuer:", berechne_Gewerbesteuer(jahresueberschuss-gwsteuer))
		gwsRück := math.Round( 100* (berechne_Gewerbesteuer(jahresueberschuss-gwsteuer) + gwsteuer ) /100 )



		// ermittelte GWSteuer Rückstellung verbuchen
		gwsKonto,_ := as.Get(accountSystem.SKR03_Steuern.Id)
		gwsSoll := booking.NewBooking(0,"in kontrol ermittelte Gewerbesteuer-Rückstellung "+strconv.Itoa(util.Global.FinancialYear), "4320", "956", "", "",nil,  -gwsRück, ("in kontrol ermittelte Gewerbesteuer-Rückstellung "+strconv.Itoa(util.Global.FinancialYear)), int(now.Month()), now.Year(), now)
		gwsKonto.Book(*gwsSoll)
		//
		gwsGegenKonto,_ := as.Get(accountSystem.SKR03_Rueckstellungen.Id)
		gwsHaben := booking.NewBooking(0,"in kontrol ermittelte Gewerbesteuer-Rückstellung "+strconv.Itoa(util.Global.FinancialYear), "", "", "", "",nil,  gwsRück, ("in kontrol ermittelte Gewerbesteuer-Rückstellung "+strconv.Itoa(util.Global.FinancialYear)), int(now.Month()), now.Year(), now)
		gwsGegenKonto.Book(*gwsHaben)
		// ermittelte GWSteuer Rückstellung von jahresueberschuss abziehen
		log.Println("in GuV, Gewerbesteuer-Rückstellung", gwsRück)
		jahresueberschuss -= gwsRück

	}

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



// Distribute Revenues and Costs according to the costcenters provided in the
// booking sheet
func ErloesverteilungAnStakeholder (as accountSystem.AccountSystem) {

	for _, acc := range as.All() {
		// loop through all accounts in accountSystem,
		// beware: All() returns no bookings, so account here has no bookings[]
		a, _ := as.Get(acc.Description.Id)
		for _, bk := range a.Bookings {
			// process bookings on GuV accounts
			switch acc.Description.Type {
			case account.KontenartAufwand:
				bk.Text = "autom. Kostenvert.: " + bk.Text
				BookCostToCostCenter{AccSystem: as, Booking: bk}.run()
			case account.KontenartErtrag:
				bk.Text = "autom. Ertragsvert.: " + bk.Text
				BookRevenueToEmployeeCostCenter{AccSystem: as, Booking: bk}.run()
			case account.KontenartAktiv:
				// now process other accounts like accountSystem.SKR03_1900.Id
				// this applies only to kommanditisten
				switch acc.Description.Id {
				case accountSystem.SKR03_Anlagen.Id,
					accountSystem.SKR03_Anlagen25_35.Id:
					BookToValuemagnetsByShares{AccSystem: as, Booking: bk, SubAcc: accountSystem.UK_VeraenderungAnlagen.Id}.run()
				case accountSystem.SKR03_Abschreibungen.Id:
					//
					BookToValuemagnetsByShares{AccSystem: as, Booking: bk, SubAcc: accountSystem.UK_VeraenderungAnlagen.Id}.run()
				default:
				}
			case account.KontenartPassiv:
				switch acc.Description.Id {
				case accountSystem.SKR03_920_Gesellschafterdarlehen.Id:
					bk.Type = booking.CC_KommitmenschDarlehen
					debit,_ := as.GetSubacc(bk.CostCenter, accountSystem.UK_Darlehen.Id)
					credit,_ := as.Get(bk.CostCenter)
					bookFromTo(bk,debit, credit)
				case accountSystem.SKR03_1900.Id: // Privatentnahmen
					bk.Type = booking.CC_Entnahme
					debit,_ := as.GetSubacc(bk.CostCenter, accountSystem.UK_Entnahmen.Id)
					credit,_ := as.Get(bk.CostCenter)
					bookFromTo(bk,debit, credit)
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
		sollacc,_  := as.GetSubacc(sh.Id, accountSystem.UK_Kosten.Id)
		log.Print("      sh: ", sh.Id, ", hat Kosten: ", math.Round(100*sollacc.Saldo)/100)
//		sollacc.YearS = accountIfYearlyIncome(*sollacc)
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
//		sollacc,_  := as.GetSubacc(valueMagnets.StakeholderKM.Id, accountSystem.UK_AnteileausFairshare)
		sollacc,_  := as.Get(sh.Id)
		habenacc,_  := as.GetSubacc(sh.Id, accountSystem.UK_AnteileausFairshare.Id)

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
//		habenacc.YearS = accountIfYearlyIncome(*habenacc)

	}
	rest -= shareHoldersShare
	log.Println("    rest to distribute: ", rest)

	// now care for Vertriebsprovisionen, die sind schon verteilt...
	sumOfProvisions := 0.0
	for _,sh := range shrepo.GetAllOfType (valueMagnets.StakeholderTypePartner) {

		// Vertriebsprovision buchen
		sollacc,_ := as.GetSubacc(sh.Id, accountSystem.UK_Vertriebsprovision.Id)
		provisions := sumOfProvisonsForStakeholder(*sollacc, sh) // sum all the partners revenue
		sumOfProvisions += provisions

		log.Printf("      %s Anteil Vertriebsprov: %2.2f€", sh.Id, provisions)
//		sollacc.YearS = accountIfYearlyIncome(*sollacc)

		rest -= provisions
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
		sollacc,_  := as.Get(sh.Id)
		habenacc,_  := as.GetSubacc(sh.Id, accountSystem.UK_AnteilMitmachen.Id)

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
//		habenacc.YearS = accountIfYearlyIncome(*habenacc)
	}


	log.Printf("      ArbeitShare: %2.2f€ = %2.2f%%", sumOfArbeitShare, sumOfArbeitShare/totalSumToDistribute)
	log.Println("    rest after ArbeitShare: ", math.Round(100*rest)/100)


	// Erlösanteile
	k_erloesAcc,_ :=  as.GetSubacc(kaccount.Description.Id, accountSystem.UK_Erloese.Id)
	sumPartnerFaktura := sumOfAllBookings(*k_erloesAcc)
	log.Println("    Sum of Partnerfaktura", sumPartnerFaktura)
	restToDistribute := rest
	sumOfErloesAnteil := rest

	// now determine the Partners Contribution
	for _,sh := range shrepo.GetAllOfType (valueMagnets.StakeholderTypePartner) {

		now := time.Now().AddDate(0, 0, 0)
		sollacc,_  := as.Get(sh.Id)
		habenacc,_  := as.GetSubacc(sh.Id, accountSystem.UK_AnteileAuserloesen.Id)

		// Erlösanteil buchen
		partnersRev := sumOfBookingsForStakeholder(*k_erloesAcc, sh) // sum the partners revenue
		erloesAnteil := math.Round(restToDistribute*100*partnersRev/sumPartnerFaktura )/100
		log.Printf("      %s revenue %2.2f%% = %2.0f€ / %2.0f€", sh.Id, partnersRev/sumPartnerFaktura, partnersRev, sumPartnerFaktura)
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
//		habenacc.YearS = habenacc.Saldo
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


// distribute need for Liquidity between partners
// book from partner main account to subaccount  "Darlehen"
func BookLiquidityNeedToPartners (as accountSystem.AccountSystem, liquidityNeed float64) {

	// each partner takes his share of the needed Liquidity according to her fairshares
	shrepo := valueMagnets.Stakeholder{}
	for _,sh := range shrepo.GetAllOfType (valueMagnets.StakeholderTypePartner) {
		fairshares,_ := strconv.ParseFloat(sh.Fairshares, 64)
		bk := booking.Booking{
			RowNr:       0,
			Amount:      fairshares*liquidityNeed,
			Soll:		 "",
			Haben: 		 "",
			Type:        booking.CC_LiquidityReserve,
			CostCenter:  sh.Id,
			Text:        fmt.Sprintf("Liquiditätsbeitrag %s %d", sh.Id, util.Global.FinancialYear),
			Month:       12,
			Year:        util.Global.FinancialYear,
			FileCreated: time.Now().AddDate(0, 0, 0),
			BankCreated: time.Now().AddDate(0, 0, 0),
		}
		debit,_ := as.Get(sh.Id)
		credit,_ := as.GetSubacc(sh.Id, accountSystem.UK_LiquidityReserve.Id)
		bookFromTo(bk,debit, credit)
/*		BookFromCreditToDebit {
			AccSystem: as,
			Booking: bk,
			Debit: debit,
			Credit: credit,
			Reason: "costcenter booking: ",
		}.run()*/

		log.Println("in bookLiquidityNeedToPartners: ", sh.Id, bk.Amount)
	}
}


func BookAmountAtDisposition (as accountSystem.AccountSystem) {
	shrepo := valueMagnets.Stakeholder{}
	for _,sh := range shrepo.GetAllOfType (valueMagnets.StakeholderTypePartner) {
		debit,_ := as.Get(sh.Id)
		credit,_ := as.GetSubacc(sh.Id, accountSystem.UK_Verfuegungsrahmen.Id)
		bk := booking.Booking{
			RowNr:       0,
			Amount:      -debit.Saldo,
			Soll:		 "",
			Haben: 		 "",
			Type:        booking.CC_LiquidityReserve,
			CostCenter:  sh.Id,
			Text:        fmt.Sprintf("Liquiditätsbeitrag %s %d", sh.Id, util.Global.FinancialYear),
			Month:       12,
			Year:        util.Global.FinancialYear,
			FileCreated: time.Now().AddDate(0, 0, 0),
			BankCreated: time.Now().AddDate(0, 0, 0),
		}
		bookFromTo(bk,debit, credit)
/*		BookFromCreditToDebit {
			AccSystem: as,
			Booking: bk,
			Debit: debit,
			Credit: credit,
			Reason: "costcenter booking: ",
		}.run()*/

		log.Println("in bookAmountAtDisposition: ", sh.Id, math.Round(100*credit.Saldo)/100)
	}
}



func sumOfBookingsForStakeholder (ac account.Account, sh valueMagnets.Stakeholder) float64 {
	saldo := 0.0
	for _,bk := range ac.Bookings {
		if sh.Id == bk.CostCenter && bk.Type == booking.CC_RevDistribution_1 {
			saldo += bk.Amount
		}
	}
	return saldo
}


func sumOfAllBookings (ac account.Account) float64 {
	saldo := 0.0
	shrepo := valueMagnets.Stakeholder{}
	for _,sh := range shrepo.GetAllOfType (valueMagnets.StakeholderTypePartner) {
		saldo += sumOfBookingsForStakeholder (ac, sh)
	}
	return saldo
}




// loop through the employees, calculates their Bonusses and book them from employee to kommitment cost.
func CalculateEmployeeBonus (as accountSystem.AccountSystem) accountSystem.AccountSystem {
	log.Println("in CalculateEmployeeBonus: ")
	shrepo := valueMagnets.Stakeholder{}
	for _,sh := range shrepo.GetAllOfType (valueMagnets.StakeholderTypeEmployee) {
		bonus := StakeholderYearlyIncome(as, sh.Id)
		log.Println("in CalculateEmployeeBonus: ",sh.Id, bonus)

		// take care, this is idempotent, i.e. that the next bonus calculation overwrites the last one...
		if bonus > 0.0 {
			// only book positive bonusses of valuemagnets
			// log.Println("in CalculateEmployeeBonus: ", sh.Id, math.Round(bonus*100)/100)
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
				FileCreated: time.Now().AddDate(0, 0, 0),
				BankCreated: time.Now().AddDate(0, 0, 0),
			}
			BookSKR03Command{AccSystem: as, Booking: bk}.run()

			// book from company cost to valuemagnets Hauptaccount
			k_subacc_costs,_ := as.GetSubacc(valueMagnets.StakeholderKM.Id, accountSystem.UK_Kosten.Id)
			sh_mainacc,_ := as.Get(bk.CostCenter)
			bookFromTo(bk, k_subacc_costs, sh_mainacc)

			// book from valuemagnets Hauptaccount into valuemagnets bonus account
			credit,_ := as.GetSubacc(bk.CostCenter, accountSystem.UK_Verfuegungsrahmen.Id)
			bookFromTo(bk, sh_mainacc, credit)
		}
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



// sum up whatever the stakeholder earned in the actual year
func StakeholderYearlyIncome (as accountSystem.AccountSystem, stkhldr string) float64 {

	yearsum := 0.0
	vm := valueMagnets.Stakeholder{}
	for _, acc := range as.GetAllAccountsOfStakeholder(vm.Get(stkhldr)) {
		if  (acc.Description.Type == "Aktiv") {
			yearsum += acc.Saldo
		}
	}
	return yearsum
}


