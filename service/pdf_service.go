package service

import (
	"example/vigovia-itenary-api/models"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// PDFService handles PDF generation
type PDFService struct {
	outputDir string
}

// NewPDFService creates a new PDF service
func NewPDFService(outputDir string) *PDFService {
	return &PDFService{
		outputDir: outputDir,
	}
}

// GeneratePDF creates a PDF from an itinerary
func (s *PDFService) GeneratePDF(itinerary *models.Itinerary) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(true, 15)

	// Add first page
	pdf.AddPage()

	// Title page
	s.addTitlePage(pdf, itinerary)

	// Trip overview
	pdf.AddPage()
	s.addTripOverview(pdf, itinerary)

	// Day-wise itinerary
	for _, day := range itinerary.Days {
		pdf.AddPage()
		s.addDayDetails(pdf, &day)
	}

	// Hotels
	pdf.AddPage()
	s.addHotels(pdf, itinerary.Hotels)

	// Flights
	pdf.AddPage()
	s.addFlights(pdf, itinerary.Flights)

	// Transfers
	if len(itinerary.Transfers) > 0 {
		pdf.AddPage()
		s.addTransfers(pdf, itinerary.Transfers)
	}

	// Payment plan
	pdf.AddPage()
	s.addPaymentPlan(pdf, &itinerary.PaymentPlan)

	// Inclusions and Exclusions
	pdf.AddPage()
	s.addInclusionsExclusions(pdf, itinerary.Inclusions, itinerary.Exclusions)

	// Generate filename
	filename := fmt.Sprintf("itinerary_%s_%s.pdf", itinerary.ID, time.Now().Format("20060102_150405"))
	filepath := filepath.Join(s.outputDir, filename)

	if err:=os.MkdirAll(s.outputDir, 0755);err!=nil{
		return "", fmt.Errorf("failed to create output directory : %w", err)
	}

	// Save PDF
	if err := pdf.OutputFileAndClose(filepath); err != nil {
		return "", fmt.Errorf("failed to save PDF: %w", err)
	}

	return filepath, nil
}

func (s *PDFService) addTitlePage(pdf *gofpdf.Fpdf, itinerary *models.Itinerary) {
	// Title
	pdf.SetFont("Arial", "B", 28)
	pdf.SetTextColor(25, 25, 112)
	pdf.CellFormat(0, 20, itinerary.Title, "", 1, "C", false, 0, "")

	pdf.Ln(10)

	// Destination
	pdf.SetFont("Arial", "I", 18)
	pdf.SetTextColor(70, 130, 180)
	pdf.CellFormat(0, 10, itinerary.Destination, "", 1, "C", false, 0, "")

	pdf.Ln(20)

	// Dates
	pdf.SetFont("Arial", "", 14)
	pdf.SetTextColor(0, 0, 0)
	dateStr := fmt.Sprintf("%s to %s",
		itinerary.StartDate.Format("January 2, 2006"),
		itinerary.EndDate.Format("January 2, 2006"))
	pdf.CellFormat(0, 10, dateStr, "", 1, "C", false, 0, "")

	pdf.Ln(10)

	// User info
	pdf.SetFont("Arial", "", 12)
	pdf.SetTextColor(100, 100, 100)
	pdf.CellFormat(0, 8, fmt.Sprintf("Itinerary ID: %s", itinerary.ID), "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 8, fmt.Sprintf("User ID: %s", itinerary.UserID), "", 1, "C", false, 0, "")

	pdf.Ln(20)

	// Duration
	duration := int(itinerary.EndDate.Sub(itinerary.StartDate).Hours()/24) + 1
	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(220, 20, 60)
	pdf.CellFormat(0, 10, fmt.Sprintf("%d Days / %d Nights", duration, duration-1), "", 1, "C", false, 0, "")
}

func (s *PDFService) addTripOverview(pdf *gofpdf.Fpdf, itinerary *models.Itinerary) {
	s.addSectionTitle(pdf, "Trip Overview")

	pdf.SetFont("Arial", "", 11)
	pdf.SetTextColor(0, 0, 0)

	// Summary
	pdf.MultiCell(0, 6, fmt.Sprintf("Duration: %d days", len(itinerary.Days)), "", "L", false)
	pdf.MultiCell(0, 6, fmt.Sprintf("Hotels: %d accommodations", len(itinerary.Hotels)), "", "L", false)
	pdf.MultiCell(0, 6, fmt.Sprintf("Flights: %d flights booked", len(itinerary.Flights)), "", "L", false)
	pdf.MultiCell(0, 6, fmt.Sprintf("Transfers: %d transfers arranged", len(itinerary.Transfers)), "", "L", false)

	pdf.Ln(5)

	// Total amount
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(220, 20, 60)
	pdf.MultiCell(0, 8, fmt.Sprintf("Total Package Cost: %.2f", itinerary.PaymentPlan.AmountDue), "", "L", false)
}

func (s *PDFService) addDayDetails(pdf *gofpdf.Fpdf, day *models.Day) {
	// Day header
	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(25, 25, 112)
	pdf.CellFormat(0, 10, fmt.Sprintf("Day %d - %s", day.DayNumber, day.Title), "", 1, "L", false, 0, "")
	
	pdf.SetFont("Arial", "I", 10)
	pdf.SetTextColor(100, 100, 100)
	pdf.CellFormat(0, 6, day.Date.Format("Monday, January 2, 2006"), "", 1, "L", false, 0, "")
	
	pdf.Ln(5)

	// Morning activities
	if len(day.Activities.Morning) > 0 {
		s.addTimeSlot(pdf, "Morning", day.Activities.Morning)
	}

	// Afternoon activities
	if len(day.Activities.Afternoon) > 0 {
		s.addTimeSlot(pdf, "Afternoon", day.Activities.Afternoon)
	}

	// Evening activities
	if len(day.Activities.Evening) > 0 {
		s.addTimeSlot(pdf, "Evening", day.Activities.Evening)
	}
}

func (s *PDFService) addTimeSlot(pdf *gofpdf.Fpdf, timeSlot string, activities []models.Activity) {
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(70, 130, 180)
	pdf.CellFormat(0, 8, timeSlot, "", 1, "L", false, 0, "")

	for _, activity := range activities {
		pdf.SetFont("Arial", "B", 11)
		pdf.SetTextColor(0, 0, 0)
		pdf.MultiCell(0, 6, fmt.Sprintf("• %s", activity.Name), "", "L", false)

		pdf.SetFont("Arial", "", 10)
		pdf.SetTextColor(60, 60, 60)
		pdf.SetLeftMargin(20)
		pdf.MultiCell(0, 5, activity.Description, "", "L", false)
		
		pdf.SetFont("Arial", "I", 9)
		pdf.SetTextColor(100, 100, 100)
		pdf.MultiCell(0, 5, fmt.Sprintf("Location: %s", activity.Location), "", "L", false)
		
		if activity.Duration != "" {
			pdf.MultiCell(0, 5, fmt.Sprintf("Duration: %s", activity.Duration), "", "L", false)
		}
		
		pdf.SetLeftMargin(10)
		pdf.Ln(3)
	}

	pdf.Ln(2)
}

func (s *PDFService) addHotels(pdf *gofpdf.Fpdf, hotels []models.Hotel) {
	s.addSectionTitle(pdf, "Accommodation Details")

	for i, hotel := range hotels {
		pdf.SetFont("Arial", "B", 12)
		pdf.SetTextColor(0, 0, 0)
		pdf.CellFormat(0, 8, fmt.Sprintf("%d. %s", i+1, hotel.Name), "", 1, "L", false, 0, "")

		pdf.SetFont("Arial", "", 10)
		pdf.SetTextColor(60, 60, 60)
		pdf.MultiCell(0, 5, fmt.Sprintf("City: %s", hotel.City), "", "L", false)
		pdf.MultiCell(0, 5, fmt.Sprintf("Check-in: %s", hotel.CheckInDate.Format("January 2, 2006")), "", "L", false)
		pdf.MultiCell(0, 5, fmt.Sprintf("Check-out: %s", hotel.CheckOutDate.Format("January 2, 2006")), "", "L", false)
		pdf.MultiCell(0, 5, fmt.Sprintf("Nights: %d", hotel.Nights), "", "L", false)
		
		if hotel.Address != "" {
			pdf.SetFont("Arial", "I", 9)
			pdf.SetTextColor(100, 100, 100)
			pdf.MultiCell(0, 5, fmt.Sprintf("Address: %s", hotel.Address), "", "L", false)
		}
		
		pdf.Ln(5)
	}
}

func (s *PDFService) addFlights(pdf *gofpdf.Fpdf, flights []models.Flight) {
	s.addSectionTitle(pdf, "Flight Details")

	for i, flight := range flights {
		pdf.SetFont("Arial", "B", 12)
		pdf.SetTextColor(0, 0, 0)
		pdf.CellFormat(0, 8, fmt.Sprintf("%d. %s - %s", i+1, flight.FlightNumber, flight.Airline), "", 1, "L", false, 0, "")

		pdf.SetFont("Arial", "", 10)
		pdf.SetTextColor(60, 60, 60)
		pdf.MultiCell(0, 5, fmt.Sprintf("From: %s", flight.From), "", "L", false)
		pdf.MultiCell(0, 5, fmt.Sprintf("To: %s", flight.To), "", "L", false)
		pdf.MultiCell(0, 5, fmt.Sprintf("Departure: %s", flight.Departure.Format("January 2, 2006 at 3:04 PM")), "", "L", false)
		pdf.MultiCell(0, 5, fmt.Sprintf("Arrival: %s", flight.Arrival.Format("January 2, 2006 at 3:04 PM")), "", "L", false)
		
		pdf.Ln(5)
	}
}

func (s *PDFService) addTransfers(pdf *gofpdf.Fpdf, transfers []models.Transfer) {
	s.addSectionTitle(pdf, "Transfer Details")

	for i, transfer := range transfers {
		pdf.SetFont("Arial", "B", 11)
		pdf.SetTextColor(0, 0, 0)
		pdf.CellFormat(0, 8, fmt.Sprintf("%d. %s to %s", i+1, transfer.From, transfer.To), "", 1, "L", false, 0, "")

		pdf.SetFont("Arial", "", 10)
		pdf.SetTextColor(60, 60, 60)
		pdf.MultiCell(0, 5, fmt.Sprintf("Mode: %s", transfer.Mode), "", "L", false)
		pdf.MultiCell(0, 5, fmt.Sprintf("Timing: %s", transfer.Timing.Format("January 2, 2006 at 3:04 PM")), "", "L", false)
		
		pdf.Ln(4)
	}
}

func (s *PDFService) addPaymentPlan(pdf *gofpdf.Fpdf, plan *models.PaymentPlan) {
	s.addSectionTitle(pdf, "Payment Plan")

	// Total amount
	pdf.SetFont("Arial", "B", 14)
	pdf.SetTextColor(220, 20, 60)
	pdf.CellFormat(0, 10, fmt.Sprintf("Total Amount: %.2f", plan.AmountDue), "", 1, "L", false, 0, "")
	
	pdf.Ln(5)

	// Installments
	pdf.SetFont("Arial", "B", 11)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 8, "Installments:", "", 1, "L", false, 0, "")

	for _, inst := range plan.Installments {
		pdf.SetFont("Arial", "B", 10)
		pdf.SetTextColor(0, 0, 0)
		pdf.CellFormat(0, 7, fmt.Sprintf("Installment %d - %.2f", 
			inst.InstallmentNumber ,inst.Amount), "", 1, "L", false, 0, "")

		pdf.SetFont("Arial", "", 9)
		pdf.SetTextColor(60, 60, 60)
		pdf.MultiCell(0, 5, fmt.Sprintf("Due Date: %s", inst.DueDate.Format("January 2, 2006")), "", "L", false)
		pdf.MultiCell(0, 5, fmt.Sprintf("Status: %s", inst.Status), "", "L", false)
		
		pdf.Ln(3)
	}
}

func (s *PDFService) addInclusionsExclusions(pdf *gofpdf.Fpdf, inclusions, exclusions []string) {
	// Inclusions
	s.addSectionTitle(pdf, "Inclusions")

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(0, 100, 0)
	for _, item := range inclusions {
		pdf.MultiCell(0, 6, fmt.Sprintf("✓ %s", item), "", "L", false)
	}

	pdf.Ln(10)

	// Exclusions
	s.addSectionTitle(pdf, "Exclusions")

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(200, 0, 0)
	for _, item := range exclusions {
		pdf.MultiCell(0, 6, fmt.Sprintf("✗ %s", item), "", "L", false)
	}
}

func (s *PDFService) addSectionTitle(pdf *gofpdf.Fpdf, title string) {
	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(25, 25, 112)
	pdf.CellFormat(0, 12, title, "", 1, "L", false, 0, "")
	pdf.Ln(3)
}