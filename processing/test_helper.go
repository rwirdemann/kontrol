package processing

import (
	"errors"

	"github.com/ahojsenn/kontrol/booking"
)

func findBookingByText(bookings []booking.Booking, text string) (*booking.Booking, error) {
	for _, b := range bookings {
		if b.Text == text {
			return &b, nil
		}
	}
	return nil, errors.New("booking with test '" + text + " not found")
}
