# Itinerary Builder API

A production-ready RESTful API service built in Golang for managing travel itineraries with PDF generation capabilities.

## Features

**Core Features**
- Complete CRUD operations for itineraries
- Day-wise activity management with time slots (morning, afternoon, evening)
- Hotel, flight, and transfer management
- Payment plan tracking with installments
- Inclusions and exclusions management
- Professional PDF generation with formatted layouts
- User-specific itinerary filtering

**Architecture**
- Clean, modular architecture with separation of concerns
- In-memory storage (easily extendable to database)
- Comprehensive input validation
- RESTful API design
- Structured error handling

## Project Structure

```
vigovia-itenary-api/
├── main.go                 # Application entry point
├── config/
│   └── config.go          # Configuration file
├── models/
│   └── itinerary.go       # Data models 
├── repository/
│   └── itinerary_repo.go  # Data access layer
├── service/
│   ├── itinerary_service.go     # Business logic
│   └── pdf_service.go           # PDF generation Service
├── controllers/
│   └── route_controller.go     # HTTP handlers
├── routes/
│   └── routes.go                # Route configuration
├── output/                      # Generated PDFs
├── sample.json                  # Sample API request
└── README.md
```

## Prerequisites

- Go 1.21 or higher
- Git

## Installation

1. **Clone or create the project:**
```bash
mkdir vigovia-itenary-api
cd vigovia-itenary-api
```

2. **Initialize Go module:**
```bash
go mod init vigovia-itenary-api
```

3. **Copy all the provided files into their respective directories**

4. **Install dependencies:**
```bash
go mod tidy
```

5. **Create output directory:**
```bash
mkdir output
```

## Running the Application

**Start the server:**
```bash
go run main.go
```

The server will start on `http://localhost:8080`

**Configure port (optional):**
```bash
export SERVER_ADDRESS=":9000"
go run main.go
```

## API Endpoints

### Health Check
```http
GET /health
```

### Create Itinerary
```http
POST /api/v1/itineraries
Content-Type: application/json

{
  "user_id": "user-12345",
  "title": "Romantic Paris & Rome Getaway",
  "destination": "Paris & Rome, Europe",
  "start_date": "2025-06-15T00:00:00Z",
  "end_date": "2025-06-22T00:00:00Z",
  "days": [...],
  "hotels": [...],
  "flights": [...],
  "transfers": [...],
  "payment_plan": {...},
  "inclusions": [...],
  "exclusions": [...]
}
```

### Get All Itineraries
```http
GET /api/v1/itineraries
```

### Get Itineraries by User
```http
GET /api/v1/itineraries?user_id=user-12345
```

### Get Single Itinerary
```http
GET /api/v1/itineraries/{id}
```

### Update Itinerary
```http
PUT /api/v1/itineraries/{id}
Content-Type: application/json

{
  "title": "Updated Title",
  "destination": "New Destination"
}
```

### Delete Itinerary
```http
DELETE /api/v1/itineraries/{id}
```

### Generate PDF
```http
POST /api/v1/itineraries/{id}/pdf
```

### Download PDF
```http
GET /api/v1/itineraries/{id}/pdf/download
```

## Testing with cURL

### 1. Create an Itinerary
```bash
curl -X POST http://localhost:8080/api/v1/itineraries \
  -H "Content-Type: application/json" \
  -d @sample_request.json
```

### 2. Get All Itineraries
```bash
curl http://localhost:8080/api/v1/itineraries
```

### 3. Get Specific Itinerary
```bash
curl http://localhost:8080/api/v1/itineraries/{itinerary-id}
```

### 4. Generate PDF
```bash
curl -X POST http://localhost:8080/api/v1/itineraries/{itinerary-id}/pdf
```

### 5. Download PDF
```bash
curl -O -J http://localhost:8080/api/v1/itineraries/{itinerary-id}/pdf/download
```

### 6. Update Itinerary
```bash
curl -X PUT http://localhost:8080/api/v1/itineraries/{itinerary-id} \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Vacation Plan"
  }'
```

### 7. Delete Itinerary
```bash
curl -X DELETE http://localhost:8080/api/v1/itineraries/{itinerary-id}
```

## Sample Request

A complete sample request is provided in `sample_request.json` featuring:
- 8-day European vacation (Paris & Rome)
- Multiple daily activities organized by time slots
- Hotel accommodations in both cities
- Round-trip flights
- Airport and inter-city transfers
- 3-installment payment plan
- Comprehensive inclusions and exclusions

## PDF Output

Generated PDFs include:
- **Title Page**: Trip title, destination, dates, and duration
- **Trip Overview**: Summary of hotels, flights, and costs
- **Day-wise Itinerary**: Detailed daily activities by time slot
- **Hotel Details**: Complete accommodation information
- **Flight Details**: All flight information with timings
- **Transfer Details**: Ground transportation arrangements
- **Payment Plan**: Installment schedule and status
- **Inclusions & Exclusions**: Complete package details

PDFs are saved in the `output/` directory with timestamp-based filenames.

## Validation Rules

The API enforces strict validation:
- All required fields must be present
- End date must be after start date
- Number of days must match date range
- Payment installments must sum to total amount
- All nested objects (activities, hotels, flights) are validated
- Preferred payment statuses: `pending`, `paid`, `cancelled`

## Error Handling

The API returns structured error responses:
```json
{
  "error": "validation_error",
  "message": "Invalid request data",
  "details": {
    "error": "field validation error details"
  }
}
```

## Extending the Application

### Adding Database Support

Replace `InMemoryRepo` with a database implementation:

```go
// repository/postgres_repository.go
type PostgresRepository struct {
    db *sql.DB
}

func (r *PostgresRepository) Create(itinerary *models.Itinerary) error {
    // Database implementation
}
```

### Adding Authentication

Add middleware in `routes/routes.go`:

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Authentication logic
        c.Next()
    }
}

// Apply to routes
itineraries.Use(AuthMiddleware())
```

### Custom PDF Templates

Modify `service/pdf_service.go` to customize:
- Colors and fonts
- Layout and spacing
- Additional sections
- Company branding

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_ADDRESS` | `:8080` | Server address and port |
| `OUTPUT_DIR` | `./output` | PDF output directory |

## Code Quality Features

✅ **Modular Design**: Separate packages for different concerns  
✅ **Interface-based**: Repository pattern for easy testing  
✅ **Comprehensive Validation**: Input validation at multiple layers  
✅ **Error Handling**: Structured error responses  
✅ **Concurrency Safe**: Thread-safe in-memory storage  
✅ **Clean Code**: Well-documented and readable  
✅ **RESTful**: Follows REST API best practices  

## Dependencies

- **gin-gonic/gin**: Fast HTTP web framework
- **google/uuid**: UUID generation
- **jung-kurt/gofpdf**: PDF generation library

## Future Enhancements

Potential improvements:
- Database integration (PostgreSQL, MongoDB)
- Authentication and authorization
- Email notifications for itinerary updates
- Multi-language support for PDFs
- Image uploads for activities
- Real-time flight and hotel availability
- Payment gateway integration
- Version control for itineraries
- Export to other formats (Excel, Word)
