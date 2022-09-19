package fate

import (
	"time"

	"github.com/6tail/lunar-go/calendar"
)

type Calendar struct {
	solar *calendar.Solar
	lunar *calendar.Lunar
}

func NewCalendar(time time.Time) *Calendar {
	s := calendar.NewSolarFromDate(time)
	return &Calendar{
		solar: s,
		lunar: s.GetLunar(),
	}
}
