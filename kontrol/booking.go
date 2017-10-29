package kontrol

import (
	"fmt"
)

type Booking struct {
	Amount float64
	Text   string
	Year   int
	Month  int
}

func (b Booking) Print(account string) {
	fmt.Printf("[%s: %2d-%d %-40s %9.2f]\n", account, b.Month, b.Year, b.Text, b.Amount)
}

type ByMonth []Booking

func (a ByMonth) Len() int           { return len(a) }
func (a ByMonth) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByMonth) Less(i, j int) bool { return a[i].Month < a[j].Month }
