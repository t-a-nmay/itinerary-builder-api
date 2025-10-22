package models

import (
	"time"
	"github.com/go-playground/validator/v10"
)

var validate=validator.New()

type Itinerary struct {
	ID		string   `json:"id" binding:"required" gorm:"primary_key;"`
	UserID  string   `json:"user_id" binding:"required"`
	Title   string      `json:"title" binding:"required"`
	Destination string  `json:"destination" binding:"required"`
	StartDate time.Time `json:"start_date" binding:"required" validate:"datetime=2006-01-02"`
	EndDate   time.Time `json:"end_date " binding:"required" validate:"datetime=2006-01-02"`
	Days      []Day	    `json:"days" binding:"required,dive" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Hotels    []Hotel    `json:"hotels" binding:"required,dive" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Flights   []Flight	`json:"flights" binding:"required,dive" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Transfers  []Transfer   `json:"transfers" binding:"dive" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	PaymentPlan PaymentPlan `json:"payment_plan" binding:"required" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Inclusions []string	`json:"inclusions" binding:"required" gorm:"type:text[]"`
	Exclusions []string	`json:"exclusions" binding:"required" gorm:"type:text[]"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type Day struct {
	// ItineraryID	  string `json:"itinerary_id" gorm:"foreignKey:ItineraryID"`
	DayNumber   int       `json:"day_number" binding:"required" min:"1" gorm:"uniqueIndex:idx_itinerary_daynumber"`
	Date		time.Time `json:"date" binding:"required" validate:"datetime=2006-01-02"`
	Title	   string    `json:"title" binding:"required"`
	Activities  Activities `json:"activities"`
}

type Activities struct {
	Morning   []Activity `json:"morning" binding:"dive"`
	Afternoon []Activity `json:"afternoon" binding:"dive"`
	Evening   []Activity `json:"evening" binding:"dive"`
}

type Activity struct{
	Name		string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Location	string `json:"location" binding:"required"`
	Duration	string `json:"duration"`
}

type Hotel struct {
	// ItineraryID	   string `json:"itinerary_id" gorm:"foreignKey:ItineraryID"`
	Name		string    `json:"name" binding:"required"`
	City 	 string    `json:"city" binding:"required"`
	CheckInDate  time.Time `json:"check_in_date" binding:"required" validate:"datetime=2006-01-02"`
	CheckOutDate time.Time `json:"check_out_date" binding:"required" validate:"datetime=2006-01-02"`
	Nights		int       `json:"nights" binding:"required" min:"1"`
	Address		string    `json:"address" binding:"required"`
}

type Flight struct {
	// ItineraryID	   string `json:"itinerary_id" gorm:"foreignKey:ItineraryID"`
	FlightNumber string    `json:"flight_number" binding:"required" gorm:"uniqueIndex:idx_itinerary_flightnumber"`
	Airline	 string    `json:"airline" binding:"required"`
	From string   `json:"from" binding:"required"`
	To   string   `json:"to" binding:"required"`
	Departure time.Time `json:"departure" binding:"required" validate:"datetime=2006-01-02"`
	Arrival   time.Time `json:"arrival" binding:"required" validate:"datetime=2006-01-02"`
}

type Transfer struct {
	// IItineraryID	   string `json:"itinerary_id" gorm:"foreignKey:ItineraryID"`
	From   string    `json:"from" binding:"required"`
	To     string    `json:"to" binding:"required"`
	Mode   string    `json:"mode" binding:"required"`
	Timing   time.Time `json:"time" binding:"required"`
}

type PaymentPlan struct {
	// ItineraryID	  string `json:"itinerary_id" gorm:"foreignKey:ItineraryID"`
	AmountDue   float64   `json:"amount_due" binding:"required"`
	DueDate     time.Time `json:"due_date" binding:"required"`
	Installments []Installment `json:"installments" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Installment struct {
	// ItineraryID	   string `json:"itinerary_id" gorm:"foreignKey:ItineraryID"`
	InstallmentNumber int       `json:"installment_number" binding:"required" min:"1"`
	Amount			float64   `json:"amount" binding:"required"`
	DueDate		time.Time `json:"due_date" binding:"required"`
	Status		string    `json:"status" binding:"required"`
}

type CreateItineraryReq struct {
	UserID     string   `json:"user_id" binding:"required"`
	Title      string      `json:"title" binding:"required"`
	Destination string     `json:"destination" binding:"required"`
	StartDate  time.Time   `json:"start_date" binding:"required"`
	EndDate    time.Time   `json:"end_date" binding:"required"`
	Days       []Day	    `json:"days" binding:"required,dive"`
	Hotels     []Hotel     `json:"hotels" binding:"dive"`
	Flights    []Flight	`json:"flights" binding:"dive"`
	Transfers  []Transfer   `json:"transfers" binding:"dive"`
	PaymentPlan PaymentPlan `json:"payment_plan" binding:"required"`
	Inclusions  []string	`json:"inclusions"`
	Exclusions  []string	`json:"exclusions"`
}

type UpdateItineraryReq struct {
	Title      *string      `json:"title"`
	Destination *string     `json:"destination"`
	StartDate  *time.Time   `json:"start_date"`
	EndDate    *time.Time   `json:"end_date"`
	Days	   []Day	    `json:"days"`
	Hotels     []Hotel     `json:"hotels"`
	Flights    []Flight	`json:"flights"`
	Transfers  []Transfer   `json:"transfers"`
	PaymentPlan *PaymentPlan `json:"payment_plan"`
	Inclusions  []string	`json:"inclusions"`
	Exclusions  []string	`json:"exclusions"`
}
