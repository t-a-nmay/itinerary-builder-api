package service

import (
	"errors"
	"example/vigovia-itenary-api/models"
	"example/vigovia-itenary-api/repository"
	"time"
	"github.com/google/uuid"
	"fmt"
)	

var (
	ErrInvalidDateRange = errors.New("end date must be after start date")
	ErrInvalidDays      = errors.New("number of days doesn't match date range")
)

// ItineraryService handles business logic for itineraries
type ItineraryService struct {
	repo repository.ItineraryRepository
}

// NewItineraryService creates a new itinerary service
func NewItineraryService(repo repository.ItineraryRepository) *ItineraryService {
	return &ItineraryService{
		repo: repo,
	}
}

// CreateItinerary creates a new itinerary
func (s *ItineraryService) CreateItinerary(req *models.CreateItineraryReq) (*models.Itinerary, error) {
	// validates date range
	if req.EndDate.Before(req.StartDate) || req.EndDate.Equal(req.StartDate) {
		return nil, ErrInvalidDateRange
	}

	// Validate days
	if err := s.validateDays(req.Days, req.StartDate, req.EndDate); err != nil {
		return nil, err
	}

	// Validate payment plan
	if err := s.validatePaymentPlan(&req.PaymentPlan); err != nil {
		return nil, err
	}

	// Create itinerary
	now := time.Now()
	itinerary := &models.Itinerary{
		ID:          uuid.New().String(),
		UserID:      req.UserID,
		Title:       req.Title,
		Destination: req.Destination,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Days:        req.Days,
		Hotels:      req.Hotels,
		Flights:     req.Flights,
		Transfers:   req.Transfers,
		PaymentPlan: req.PaymentPlan,
		Inclusions:  req.Inclusions,
		Exclusions:  req.Exclusions,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(itinerary); err != nil {
		return nil, fmt.Errorf("failed to create itinerary: %w", err)
	}

	return itinerary, nil
}

// GetItinerary retrieves an itinerary by ID
func (s *ItineraryService) GetItinerary(id string) (*models.Itinerary, error) {
	itinerary, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("itinerary not found")
		}
		return nil, fmt.Errorf("failed to get itinerary: %w", err)
	}

	return itinerary, nil
}

// GetUserItineraries retrieves all itineraries for a user
func (s *ItineraryService) GetUserItineraries(userID string) ([]*models.Itinerary, error) {
	itineraries, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user itineraries: %w", err)
	}

	return itineraries, nil
}

// GetAllItineraries retrieves all itineraries
func (s *ItineraryService) GetAllItineraries() ([]*models.Itinerary, error) {
	itineraries, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all itineraries: %w", err)
	}

	return itineraries, nil
}

// UpdateItinerary updates an existing itinerary
func (s *ItineraryService) UpdateItinerary(id string, req *models.UpdateItineraryReq) (*models.Itinerary, error) {
	// Get existing itinerary
	existing, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("itinerary not found")
		}
		return nil, fmt.Errorf("failed to get itinerary: %w", err)
	}

	// Update fields
	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.Destination != nil {
		existing.Destination = *req.Destination
	}
	if req.StartDate != nil {
		existing.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		existing.EndDate = *req.EndDate
	}
	if req.Days != nil {
		existing.Days = req.Days
	}
	if req.Hotels != nil {
		existing.Hotels = req.Hotels
	}
	if req.Flights != nil {
		existing.Flights = req.Flights
	}
	if req.Transfers != nil {
		existing.Transfers = req.Transfers
	}
	if req.PaymentPlan != nil {
		existing.PaymentPlan = *req.PaymentPlan
	}
	if req.Inclusions != nil {
		existing.Inclusions = req.Inclusions
	}
	if req.Exclusions != nil {
		existing.Exclusions = req.Exclusions
	}

	existing.UpdatedAt = time.Now()

	// Validate updated data
	if existing.EndDate.Before(existing.StartDate) || existing.EndDate.Equal(existing.StartDate) {
		return nil, ErrInvalidDateRange
	}

	if err := s.validatePaymentPlan(&existing.PaymentPlan); err != nil {
		return nil, err
	}

	// Save changes
	if err := s.repo.Update(id, existing); err != nil {
		return nil, fmt.Errorf("failed to update itinerary: %w", err)
	}

	return existing, nil
}

// DeleteItinerary deletes an itinerary
func (s *ItineraryService) DeleteItinerary(id string) error {
	if err := s.repo.Delete(id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return fmt.Errorf("itinerary not found")
		}
		return fmt.Errorf("failed to delete itinerary: %w", err)
	}

	return nil
}

// validateDays validates that days match the date range
func (s *ItineraryService) validateDays(days []models.Day, startDate, endDate time.Time) error {
	duration := int(endDate.Sub(startDate).Hours()/24) + 1
	if len(days) != duration {
		return fmt.Errorf("%w: expected %d days, got %d", ErrInvalidDays, duration, len(days))
	}

	return nil
}

// validatePaymentPlan validates payment plan consistency
func (s *ItineraryService) validatePaymentPlan(plan *models.PaymentPlan) error {
	if plan.AmountDue <= 0 {
		return errors.New("total amount must be positive")
	}

	if len(plan.Installments) == 0 {
		return errors.New("at least one installment is required")
	}

	var sum float64
	for _, inst := range plan.Installments {
		sum += inst.Amount
	}

	// Allow small floating point difference
	if sum < plan.AmountDue-0.01 || sum > plan.AmountDue+0.01 {
		return fmt.Errorf("installment amounts (%.2f) don't match total amount (%.2f)", sum, plan.AmountDue)
	}

	return nil
}
