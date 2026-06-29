package domain

import "time"

type Role string

const (
	RoleRegular Role = "REGULAR"
	RoleWalkin  Role = "WALKIN"
	RoleStaff   Role = "STAFF"
)

type BookingStatus string

const (
	StatusExpected  BookingStatus = "EXPECTED"
	StatusAvailable BookingStatus = "AVAILABLE"
	StatusBooked    BookingStatus = "BOOKED"
	StatusAbsent    BookingStatus = "ABSENT"
)

type AttendanceStatus string

const (
	AttendancePending AttendanceStatus = "PENDING"
	AttendancePresent AttendanceStatus = "PRESENT"
	AttendanceNoShow  AttendanceStatus = "NOSHOW"
)

type User struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	NationalID string    `gorm:"uniqueIndex;size:13" json:"-"`
	LineUserID *string   `gorm:"uniqueIndex;size:64" json:"lineUserId,omitempty"`
	Name       string    `gorm:"size:100" json:"name"`
	Phone      string    `gorm:"size:20;index" json:"phone"`
	Role       Role      `gorm:"size:20;index" json:"role"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type Stall struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Code          string    `gorm:"uniqueIndex;size:20" json:"code"`
	Zone          string    `gorm:"size:50" json:"zone"`
	OwnerVendorID *uint     `gorm:"index" json:"ownerVendorId,omitempty"`
	Owner         *User     `gorm:"foreignKey:OwnerVendorID" json:"-"`
	PosX          float64   `json:"posX"`
	PosY          float64   `json:"posY"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type DailyBooking struct {
	ID               uint             `gorm:"primaryKey" json:"id"`
	StallID          uint             `gorm:"not null;uniqueIndex:idx_stall_date" json:"stallId"`
	MarketDate       time.Time        `gorm:"type:date;not null;uniqueIndex:idx_stall_date" json:"marketDate"`
	Status           BookingStatus    `gorm:"size:20;index" json:"status"`
	VendorID         *uint            `gorm:"index" json:"vendorId,omitempty"`
	ShopName         string           `gorm:"size:100" json:"shopName"`
	ProductType      string           `gorm:"size:100;index" json:"productType"`
	Attendance       AttendanceStatus `gorm:"size:20;default:'PENDING';index" json:"attendance"`
	CheckedByStaffID *uint            `json:"checkedByStaffId,omitempty"`
	CheckedAt        *time.Time       `json:"checkedAt,omitempty"`
	BookedAt         *time.Time       `json:"bookedAt,omitempty"`
	CreatedAt        time.Time        `json:"createdAt"`
	UpdatedAt        time.Time        `json:"updatedAt"`

	Stall *Stall `gorm:"foreignKey:StallID" json:"stall,omitempty"`
}

type AuditLog struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ActorUserID    uint      `gorm:"index" json:"actorUserId"`
	DailyBookingID *uint     `gorm:"index" json:"dailyBookingId,omitempty"`
	Action         string    `gorm:"size:30;index" json:"action"`
	Detail         string    `gorm:"type:text" json:"detail"`
	CreatedAt      time.Time `json:"createdAt"`
}
