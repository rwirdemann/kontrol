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


	// find OPOS
	switch {
	case booking.IsOpenPosition():
		booking.Text = "Achtung OPOS: "+booking.Text
	case booking.IsBeyondBudgetDate():
		booking.Text = "OPOS at BalanceDate: "+booking.Text
	default:
	}

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
	case "SKR03", "closingBalance", "openingBalance":
		command =  DontDoAnything{}
		sollAccount := accsystem.GetSKR03(booking.Soll)
		habenAccount := accsystem.GetSKR03(booking.Haben)
		bookFromTo(booking, sollAccount, habenAccount)
	default:
		log.Println("in Process: unknown command",booking.Type, " in row", booking.RowNr)
	}
	command.run()

}




// book all employees cost and revenues
func ErloesverteilungAnEmployees (as accountSystem.AccountSystem) {
	log.Println("in ErloesverteilungAnEmployees")

	// loop over employees subaccounts
	for _, acc := range as.All() {
		// loop through all accounts in accountSystem,
		// employees bookings only
		sh := valueMagnets.Stakeholder{}
		if  !sh.IsEmployee(acc.Description.Superaccount) {
			continue
		}

		// log.Println("	employee", acc.Description.Id)
		// beware: All() returns no bookings, so account here has no bookings[]
		a, _ := as.Get(acc.Description.Id)
		for _, bk := range a.Bookings  {

			// process bookings on GuV accounts
			switch acc.Description.Type {

			// alle Kosten
			case account.KontenartAufwand:
				bk.Text = "autom. Kostenvert.: " + bk.Text
				BookCostToCostCenter{AccSystem: as, Booking: bk}.run()

			// alle Ertröge
			case account.KontenartErtrag:
				bk.Text = "autom. Ertragsvert.: " + bk.Text
				BookRevenueToEmployeeCostCenter{AccSystem: as, Booking: bk}.run()

			default:
				log.Println("	ERROR: found booking which is neither Aufwand nor Ertrag", acc.Description.Id, acc.Description.Type)
			}

		}
	}

	vm := valueMagnets.Stakeholder{}
	for _,employee := range vm.AllEmployees(util.Global.FinancialYear) {

		acc,_ := as.Get(employee.Id)
		log.Println("	:", employee.Name, acc.Saldo )
	}
}









// Distribute Revenues and Costs according to the costcenters provided in the
// booking sheet
func ErloesverteilungAnKommanditisten(as accountSystem.AccountSystem) {

	for _, acc := range as.All() {
		// loop through all accounts in accountSystem,
		// skip employees bookings
		sh := valueMagnets.Stakeholder{}
		if  sh.IsEmployee(acc.Description.Superaccount) {
			continue
		}
		// beware: All() returns no bookings, so account here has no bookings[]
		a, _ := as.Get(acc.Description.Id)
		for _, bk := range a.Bookings {

			// process bookings on GuV accounts
			switch acc.Description.Type {

			// alle Kosten
			case account.KontenartAufwand:
				bk.Text = "autom. Kostenvert.: " + bk.Text
				BookCostToCostCenter{AccSystem: as, Booking: bk}.run()

			// alle Ertröge
			case account.KontenartErtrag:
				bk.Text = "autom. Ertragsvert.: " + bk.Text
				BookRevenueToEmployeeCostCenter{AccSystem: as, Booking: bk}.run()

			// alle Anlagen und Abschreibungen
			case account.KontenartAktiv:
				// now process other accounts like accountSystem.SKR03_1900.Id
				switch acc.Description.Id {
				case accountSystem.SKR03_Anlagen.Id, accountSystem.SKR03_Anlagen25_35.Id:
					switch bk.Type {
					case "openingBalance":
						BookToValuemagnetsByShares{AccSystem: as, Booking: bk, SubAcc: accountSystem.UK_AnteilAnAnlagen.Id}.run()
					case "SKR03":
						BookToValuemagnetsByShares{AccSystem: as, Booking: bk, SubAcc: accountSystem.UK_VeraenderungAnlagen.Id}.run()
					// Schlussbilanz ausnehmen
					case "closingBalance":
					default:
					}
				case accountSystem.SKR03_Abschreibungen.Id:
					//
					BookToValuemagnetsByShares{AccSystem: as, Booking: bk, SubAcc: accountSystem.UK_VeraenderungAnlagen.Id}.run()
				default:
				}

			// alle Darlehen und Entnahmen
			case account.KontenartPassiv:
				switch acc.Description.Id {
				case accountSystem.SKR03_920_Gesellschafterdarlehen.Id:
					bk.Type = booking.CC_KommitmenschDarlehen
					// book on kommitment.UK_Darlehen auf Hauptkonto kommanditist
					debit,_ := as.GetSubacc(valueMagnets.StakeholderKM.Id, accountSystem.UK_Darlehen.Id)
					credit,_ := as.Get(bk.CostCenter)
					bookFromTo(bk, debit, credit)

					// von Hauptkonto kommanditist an kommanditist.UK_Darlehen
					debit,_ = as.Get(bk.CostCenter)
					credit,_ = as.GetSubacc(bk.CostCenter, accountSystem.UK_Darlehen.Id)
					bookFromTo(bk, debit, credit)

				case accountSystem.SKR03_1900.Id: // Privatentnahmen
					bk.Type = booking.CC_Entnahme

					// von kommitment.Entnahmen auf Hauptkonto kommanditist
					debit,_ := as.GetSubacc(valueMagnets.StakeholderKM.Id, accountSystem.UK_Entnahmen.Id)
					credit,_ := as.Get(bk.CostCenter)
					bookFromTo(bk, debit, credit)

					// von Hauptkonto kommanditist an kommanditist.UK_Entnahmen
					debit,_ = as.Get(bk.CostCenter)
					credit,_ = as.GetSubacc(bk.CostCenter, accountSystem.UK_Entnahmen.Id)
					bookFromTo(bk, debit, credit)
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
		anteil_fairshares := booking.Booking{
			Amount:      -fairshareAnteil,
			Type:        booking.CC_AnteilAusFairshares,
			CostCenter:  sh.Id,
			Text:         "Anteil aus fairshares("+strconv.FormatFloat(fairshares, 'f', 2, 64)+")",
			FileCreated: now,
			BankCreated: now,
		}
		// 1. Buchung von StakeholderKM auf sh.Id
		sollacc,_  := as.GetSubacc(valueMagnets.StakeholderKM.Id, accountSystem.UK_AnteileausFairshare.Id)
		habenacc,_  := as.Get(sh.Id)
		bookFromTo(anteil_fairshares, habenacc, sollacc)

		// 2. Buchung von sh.Id auf subaccf UK_AnteileausFairshare
		sollacc,_  = as.Get(sh.Id)
		habenacc,_  = as.GetSubacc(sh.Id, accountSystem.UK_AnteileausFairshare.Id)
		bookFromTo(anteil_fairshares, habenacc, sollacc)
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


	arbeitPercentage := (1.00 - 0.2 - provisionPercentage) / 2
	// now distribute 50% of the rest according to factor "Arbeit"
	// use account interne stunden fpr now...
	sumOfArbeit := 0.0
	for _,sh := range shrepo.GetAllOfType (valueMagnets.StakeholderTypePartner) {
		shArbeit,_ := strconv.ParseFloat(sh.Arbeit, 64)
		sumOfArbeit += shArbeit
	}
	log.Printf("      Sum Arbeit =  %2.2f years", sumOfArbeit)
	//sumOfArbeitShare := 0.0
	restToDistributeByArbeit := rest*0.5
	log.Printf("      ArbeitsShare %2.2f€ = %2.2f%%", restToDistributeByArbeit, arbeitPercentage)
	for _,sh := range shrepo.GetAllOfType (valueMagnets.StakeholderTypePartner) {
		now := time.Now().AddDate(0, 0, 0)

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
		// 1. Buchung von StakeholderKM auf sh.Id
		sollacc,_  := as.GetSubacc(valueMagnets.StakeholderKM.Id, accountSystem.UK_AnteilMitmachen.Id)
		habenacc,_  := as.Get(sh.Id)
		bookFromTo(anteil_erloese, habenacc, sollacc)

		// 2. Buchung von sh.Id auf subaccf UK_AnteilMitmachen
		sollacc,_  = as.Get(sh.Id)
		habenacc,_  = as.GetSubacc(sh.Id, accountSystem.UK_AnteilMitmachen.Id)
		bookFromTo(anteil_erloese, habenacc, sollacc)

		log.Printf("      %s Anteil ArbeitShare: %2.2f€", sh.Id, habenacc.Saldo)
		rest -= habenacc.Saldo
		//		habenacc.YearS = accountIfYearlyIncome(*habenacc)
	}

	//	log.Printf("      ArbeitShare: %2.2f€ = %2.2f%%", sumOfArbeitShare, sumOfArbeitShare/totalSumToDistribute)
	log.Println("    rest after ArbeitShare: ", math.Round(100*rest)/100)

	// Erlösanteile
	k_erloesAcc,_ :=  as.GetSubacc(kaccount.Description.Id, accountSystem.UK_Erloese.Id)
	sumPartnerFaktura := sumOfAllBookings(*k_erloesAcc)
	if (sumPartnerFaktura < 1.0) {
		// in the very unli9kely event that there is a mnonth withour partner faktura...
		// the distribution does not work sesibly anyway ...
		sumPartnerFaktura = 1.0
	}

	log.Println("    Sum of Partnerfaktura", sumPartnerFaktura)
	restToDistribute := rest
	//sumOfErloesAnteil := rest
	fakturaPercentage := 1.00 - 0.2 - provisionPercentage - arbeitPercentage
	log.Printf("      fakturaShare: %2.2f€ = %2.2f%%", restToDistribute, fakturaPercentage)

	// now determine the Partners Contribution
	for _,sh := range shrepo.GetAllOfType (valueMagnets.StakeholderTypePartner) {

		now := time.Now().AddDate(0, 0, 0)

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
		// 1. Buchung von StakeholderKM auf sh.Id
		sollacc,_  := as.GetSubacc(valueMagnets.StakeholderKM.Id, accountSystem.UK_AnteileAuserloesen.Id)
		habenacc,_  := as.Get(sh.Id)
		bookFromTo(anteil_erloese, habenacc, sollacc)

		// 2. Buchung von sh.Id auf subaccf UK_AnteileAuserloesen
		sollacc,_  = as.Get(sh.Id)
		habenacc,_  = as.GetSubacc(sh.Id, accountSystem.UK_AnteileAuserloesen.Id)
		bookFromTo(anteil_erloese, habenacc, sollacc)

		// now book habenacc.KommitmenschNettoFaktura  to anteil_erloese / fakturaPercentage
		habenacc.KommitmenschNettoFaktura = erloesAnteil / fakturaPercentage

		rest -= erloesAnteil
	}
	//log.Printf("      ErlösAnt.: %2.2f€ = %2.2f%%", sumOfErloesAnteil, sumOfErloesAnteil/totalSumToDistribute)
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
// book from partner main account to subaccount  "BookLiquidityNeedToPartners"
func BookLiquidityNeedToPartners (as accountSystem.AccountSystem, liquidityNeed float64) {

	// each partner takes his share of the needed Liquidity according to her fairshares
	shrepo := valueMagnets.Stakeholder{}
	for _,sh := range shrepo.GetAllOfType (valueMagnets.StakeholderTypePartner) {
		fairshares,_ := strconv.ParseFloat(sh.Fairshares, 64)
		bk := booking.Booking{
			RowNr:       0,
			Amount:      -1.0*fairshares*liquidityNeed,
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

		// book from kommitment.UK_UK_LiquidityReserve to stakeholder's main account
		debit,_ := as.GetSubacc(valueMagnets.StakeholderKM.Id, accountSystem.UK_LiquidityReserve.Id)
		credit,_ := as.Get(sh.Id)
		bookFromTo(bk, debit, credit)

		// book from stakeholders main account to stakeholders subaccount
		debit,_ = as.Get(sh.Id)
		credit,_ = as.GetSubacc(sh.Id, accountSystem.UK_LiquidityReserve.Id)
		bookFromTo(bk, debit, credit)

		log.Println("in bookLiquidityNeedToPartners: ", sh.Id, bk.Amount)
	}
}


func BookAmountAtDisposition (as accountSystem.AccountSystem) {
	// book the kommitment company Sado to Suacc Verfügungsrahmen/Bonus once
/*
    k,_ := as.Get(valueMagnets.StakeholderKM.Id)
	k_UK, _ := as.GetSubacc(valueMagnets.StakeholderKM.Id, accountSystem.UK_Verfuegungsrahmen.Id)
	bk := booking.Booking{
		RowNr:       0,
		Amount:      -1.0*k.Saldo,
		Soll:		 "",
		Haben: 		 "",
		Type:        booking.CC_LiquidityReserve,
		CostCenter:  valueMagnets.StakeholderKM.Id,
		Text:        fmt.Sprintf("Erlöse  %d",  util.Global.FinancialYear),
		Month:       12,
		Year:        util.Global.FinancialYear,
		FileCreated: time.Now().AddDate(0, 0, 0),
		BankCreated: time.Now().AddDate(0, 0, 0),
	}
	bookFromTo(bk,k, k_UK)
*/
// I think this is not korrekt... 2019-12-02



	// book from stakeholder UK_Verfuegungsrahmen back to K_UK_Verfuegungsrahmen
	// this should give equal Salden at both sides of balance sheet and zero this account
	shrepo := valueMagnets.Stakeholder{}
	for _,sh := range shrepo.GetAllOfType (valueMagnets.StakeholderTypePartner) {
		stakeholder_Saldo := 0.0

		a,_ := as.GetSubacc(sh.Id, accountSystem.UK_Kosten.Id)
		stakeholder_Saldo += a.Saldo
		a,_ = as.GetSubacc(sh.Id, accountSystem.UK_AnteileausFairshare.Id)
		stakeholder_Saldo += a.Saldo
		a,_ = as.GetSubacc(sh.Id, accountSystem.UK_AnteileAuserloesen.Id)
		stakeholder_Saldo += a.Saldo
		a,_ = as.GetSubacc(sh.Id, accountSystem.UK_AnteilMitmachen.Id)
		stakeholder_Saldo += a.Saldo
		a,_ = as.GetSubacc(sh.Id, accountSystem.UK_Vertriebsprovision.Id)
		stakeholder_Saldo += a.Saldo

		debit,_  := as.GetSubacc(valueMagnets.StakeholderKM.Id, accountSystem.UK_Verfuegungsrahmen.Id)
		credit,_ := as.GetSubacc(sh.Id, accountSystem.UK_Verfuegungsrahmen.Id)
		bk := booking.Booking{
			RowNr:       0,
			Amount:      stakeholder_Saldo,
			Soll:		 "",
			Haben: 		 "",
			Type:        booking.CC_J_Bonus,
			CostCenter:  sh.Id,
			Text:        fmt.Sprintf("Jahresüberschuss/Bonus %s %d", sh.Id, util.Global.FinancialYear),
			Month:       12,
			Year:        util.Global.FinancialYear,
			FileCreated: time.Now().AddDate(0, 0, 0),
			BankCreated: time.Now().AddDate(0, 0, 0),
		}
		bookFromTo(bk,debit, credit)

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
		log.Printf("in CalculateEmployeeBonus: %s:  %7.2f€",sh.Id, bonus)

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
				Text:        fmt.Sprintf("in kontrol kalkulierte Bonusrückstellung für %s in %d", sh.Id, util.Global.FinancialYear),
				Month:       12,
				Year:        util.Global.FinancialYear,
				FileCreated: time.Now().AddDate(0, 0, 0),
				BankCreated: time.Now().AddDate(0, 0, 0),
			}
			// create new booking
			sollAccount := as.GetSKR03(bk.Soll)
			habenAccount := as.GetSKR03(bk.Haben)
			bookFromTo(bk, sollAccount, habenAccount)


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

	a1,_  := as.GetSubacc(stkhldr, accountSystem.UK_Kosten.Id)
	a2,_  := as.GetSubacc(stkhldr, accountSystem.UK_AnteileausFairshare.Id)
	a3,_  := as.GetSubacc(stkhldr, accountSystem.UK_Vertriebsprovision.Id)
	a4,_  := as.GetSubacc(stkhldr, accountSystem.UK_AnteilMitmachen.Id)
	a5,_  := as.GetSubacc(stkhldr, accountSystem.UK_AnteileAuserloesen.Id)
	a6,_  := as.GetSubacc(stkhldr, accountSystem.UK_Erloese.Id)

	log.Printf("in StakeholderYearlyIncome: %s %7.2f€ %7.2f€ %7.2f€ %7.2f€ %7.2f€ %7.2f€", stkhldr, a1.Saldo, a2.Saldo, a3.Saldo, a4.Saldo, a5.Saldo, a6.Saldo)

	return	a1.Saldo +
		a2.Saldo +
		a3.Saldo +
		a4.Saldo +
		a5.Saldo +
		a6.Saldo
}


