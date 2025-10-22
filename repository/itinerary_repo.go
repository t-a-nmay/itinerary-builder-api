package repository

import (
	"errors"
	"example/vigovia-itenary-api/models"
)

var(
	ErrNotFound=errors.New("itinerary not found")
)

type ItineraryRepository interface {
	Create(itinerary *models.Itinerary) error
	GetAll()([]*models.Itinerary,error)
	GetByID(id string)(*models.Itinerary,error)
	GetByUserID(userID string)([]*models.Itinerary,error)
	Update(id string, itinerary *models.Itinerary) error
	Delete(id string) error
}

//implements the InMemoryRepo using an in-memory map
type InMemoryRepo struct {
	itineraries map[string]*models.Itinerary
}

//creates and returns a new instance of InMemoryRepo
func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		itineraries: make(map[string]*models.Itinerary),
	}
}

//adds a new itinerary to the in-memory db (map)
func(r *InMemoryRepo) Create(itinerary *models.Itinerary) error {
	if _,exists:=r.itineraries[itinerary.ID];exists{
		return errors.New("itinerary already exists")
	}
	r.itineraries[itinerary.ID] = itinerary
	return nil
}

//gets all the itineraries from the in-memory db
func(r *InMemoryRepo) GetAll()([]*models.Itinerary,error){
	itineraries:=make([]*models.Itinerary,0,len(r.itineraries))
	for _,itinerary:=range r.itineraries{
		itineraries=append(itineraries,itinerary)
	}
	return itineraries,nil
}

//gets itinerary by ID 
func(r*InMemoryRepo) GetByID(id string)(*models.Itinerary,error){
	itinerary,exists:=r.itineraries[id]
	if(!exists){
		return  nil, errors.New("itinerary not found")
	}
	return itinerary,nil
}

//gets itineraries by UserID
func(r *InMemoryRepo) GetByUserID(userID string)([]*models.Itinerary,error){
	var userItineraries []*models.Itinerary
	for _,itinerary:=range r.itineraries{
		if itinerary.UserID==userID{
			userItineraries=append(userItineraries,itinerary)
		}
	}
	if(len(userItineraries)==0){}
	return userItineraries,nil
}

//update itinerary by ID
func(r *InMemoryRepo) Update(id string, itinerary *models.Itinerary) error {
	if _,exists:=r.itineraries[id];!exists{
		return errors.New("itinerary not found")
	}
	r.itineraries[id]=itinerary
	return nil
}

//delete itinerary by ID
func(r *InMemoryRepo) Delete(id string) error {
	if _,exists :=r.itineraries[id];!exists{
		return errors.New("itinerary not found")
	}
	delete(r.itineraries,id)
	return nil
}



