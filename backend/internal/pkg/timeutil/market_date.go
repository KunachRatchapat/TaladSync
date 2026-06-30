package timeutil

import (
	"fmt"
	"time"
)

var bangkok *time.Location

func init() {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(fmt.Sprintf("load Asia/Bangkok: %v", err))
	}
	bangkok = loc
}

// Now คืนเวลาปัจจุบันใน Asia/Bangkok (ใช้กับ booked_at, checked_at)
func Now() time.Time {
	return time.Now().In(bangkok)
}

// Today คืน market_date ของวันนี้ (เที่ยงคืน Asia/Bangkok) สำหรับ GORM type:date
func Today() time.Time {
	now := time.Now().In(bangkok)
	y, m, d := now.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, bangkok)
}

// MarketDate normalize เวลาใด ๆ ให้เป็นวันที่ตลาด 
func MarketDate(t time.Time) time.Time {
	t = t.In(bangkok)
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, bangkok)
}
