# Data Ingestion Service

A cloud-native data ingestion pipeline that collects logs from public APIs, processes them, and stores them in BigQuery using Go services and DBT transformations. Also provides and API to read data. At present, I have used AI to write this as per my instructions and design. THe ideas and completely mine with AI only writing it down for me. I wasn't able to test the system since I lack a proper GCP account which have increased verification to enable services like BigQuery.

## Architecture

```
External API → Go Ingester → BigQuery (Raw) → DBT → BigQuery (Processed) → Go API Fetcher
```

## Components

- **Go Ingester**: Fetches data from external APIs and stores raw JSON in BigQuery
- **Go API**: REST API to retrieve processed data
- **DBT**: Transforms raw data into structured, analytics-ready format
- **BigQuery**: Cloud-native storage for both raw and processed data
- **Terraform**: Infrastructure as Code for GCP resources
- **Docker**: Containerization for consistent deployments

## Setup Instructions

### Prerequisites

- Go 1.24+
- Docker & Docker Compose
- GCP Account with BigQuery API enabled
- Terraform 1.0+
- Python

### Local Development

1. **Clone the repository**
```bash
git clone git@github.com:ihsan-96/data-ingesion-service.git
cd data-ingestion-service
```

2. **Set up environment variables**
```bash
cp .env.example .env
cp service-account.json.example service-account.json
cp go-service/.env.example .env
# Need to update .env with proper GCP_PROJECT_ID and service-account.json with gcp creds.
```

3. **Deploy infrastructure**
```bash
cd terraform
terraform init
terraform plan -var="project_id=gcp-project-id"
terraform apply
```

5. **Start services**
```bash
docker-compose up -d
```

## Running the Application

### Local with Docker Compose

```bash
# Start all services
docker-compose up -d

# Check logs
docker-compose logs -f ingester
docker-compose logs -f api

# Run DBT transformations
docker-compose run dbt dbt run
```

### Individual Services

```bash
# Go Ingester
cd go-service
go run cmd/ingester/main.go

# Go API
cd go-service
go run cmd/api/main.go

# DBT
cd dbt
pip install dbt-bigquery
source ../.env # Use proper env
dbt run --profiles-dir ./profiles-local
```

## Running Tests

### Go Tests
```bash
cd go-service
go test ./tests/...
```

## API Documentation

### Endpoints

#### GET /health
Health check endpoint
```json
{
  "status": "healthy"
}
```

#### GET /logs
Retrieve processed logs with pagination
```
Query Parameters:
- limit: Number of records (default: 100, max: 1000)
- offset: Pagination offset (default: 0)
```

Response:
```json
{
  "data": [
    {
      "log_id": 1,
      "user_id": 1,
      "title": "sunt aut facere",
      "body": "quia et suscipit",
      "word_count": 4,
      "ingested_at": "2024-01-01T12:00:00Z",
      "source": "placeholder_api",
      "process_date": "2024-01-01",
      "created_at": "2024-01-01T12:05:00Z"
    }
  ],
  "limit": 100,
  "offset": 0,
  "count": 1
}
```

## Transformation Logic

1. Data fetched from API is dumped to GCP with ingesion time and source
2. DBT extracts the data into a stage stable
3. DBT then puts it to a prod table with correct types and any other data extraction we need.

## Deploying to Cloud Environment

### GCP Cloud Run Deployment

1. **Set up secrets in GitHub**
```
GCP_PROJECT_ID: GCP project ID
GCP_SA_KEY: Service account key JSON (base64 encoded)
API_ENDPOINT: External API endpoint URL to get the data (if we are setting it in env vars of github actions/ we can use gcp secrets to store it also but that is not implemented)
```

2. **Deploy via GitHub Actions**
```bash
# Push to main branch triggers deployment
git push origin main
```


## Trade-offs and Design Decisions

### Storage Choice: BigQuery
**Pros:**
- Serverless and fully managed
- Excellent for analytics
- Built-in DBT integration
- Cost-effective
- Schema evolution support

**Cons:**
- Not ideal for real-time updates

### Trade-offs
Thought to use Postgres/Mongo/S3/BigQuery for the use case. BigQuery seems apt for a data ingesion like ours. It is more cost effective for the use case and give more analytical capabilities and handles data at scale.
Made it micro services that handles different parts even though it has some deployment overhead.
Scheduling part is a bit more open for business decision, wether we need a periodic data or more freequent updates. ingester is currently a service which gets data every 5 minutes. we can make it a job that triggers every 3-4 hours using cloud scheduler also. DBT is currently implemented like that for demonstration.


## Implementation Challenges

### Hardest Parts to Implement
Big Query connection issues and Error Handling.
Also the Go Service since I am new to Go



## Future Improvements

- Add metrics to some observable stack like Prometheus
- Use Airflow to orchestrate the system
- Rate limiting for API endpoints (Both External and Internal)
- Alerting for ingestion failures
- Data quality monitoring and checks
