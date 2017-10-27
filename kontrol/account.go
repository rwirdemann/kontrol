package kontrol

type Account struct {
	Bookings []Booking
}

var Accounts map[string]*Account

func (a *Account) Book(booking Booking) {
	a.Bookings = append(a.Bookings, booking)
}

func (a Account) Saldo() float64 {
	saldo := 0.0
	for _, b := range a.Bookings {
		saldo += b.Amount
	}
	return saldo
}
