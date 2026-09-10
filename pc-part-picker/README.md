# AI PC Part Picker

This is a full-stack web application that generates price-optimized, compatible PC builds based on user-defined workloads, utilizing live AI-scraped pricing data.

## Project Structure

*   **`frontend/`**: Next.js (App Router) with React and Tailwind CSS. Provides the user interface for configuring builds.
*   **`backend/`**: Go using the Echo framework. Handles the core compatibility logic, budget allocation, and communicates with the database and scraper.
*   **`scraper/`**: Python FastAPI microservice utilizing `crawl4ai`. Fetches live pricing data for components.
*   **`schema.sql`** (in `backend/`): Configures the MySQL database schema.

---

## Prerequisites

*   Node.js (v18+) and npm
*   Go (v1.21+)
*   Python (3.10+)
*   MySQL (or Docker to run the database)

---

## Local Development Setup

### 1. Database Configuration

You can run the MySQL database locally or using Docker. Assuming you have a local MySQL instance:

Import the schema and seed data:

mysql -h 127.0.0.1 -u root -p < backend/schema.sql


### 2. Scraper Service (Python)

Open a new terminal and start the AI Web Scraper:

cd pc-part-picker/scraper
python3 -m venv venv
source venv/bin/activate
pip install fastapi uvicorn crawl4ai pydantic
python main.py &


*The scraper will run on `http://localhost:8000`.*

### 3. Backend Service (Go)

Open a new terminal and start the Go API:

cd pc-part-picker/backend
# Set database environment variables if different from defaults
export DB_USER=root
export DB_PASS=root
export DB_HOST=127.0.0.1:3306
export DB_NAME=pcpartpicker

go mod tidy
go run main.go &


*The backend will run on `http://localhost:8080`.*

### 4. Frontend Service (Next.js)

Open a new terminal and start the frontend UI:

cd pc-part-picker/frontend
npm install
npm run dev &


*The frontend will be accessible at `http://localhost:3000`.*

---

## Production Deployment Guide

For production, it is highly recommended to containerize all services and deploy them behind a reverse proxy (like Nginx or Traefik).

### Frontend (Next.js)
1. Build the application: `npm run build`
2. Start the production server: `npm start`
*Note: In production, ensure the API endpoint in the Next.js components points to your production Go API domain, not `localhost`.*

### Backend (Go)
1. Build the binary for your target OS: `GOOS=linux GOARCH=amd64 go build -o pcpartpicker-backend main.go`
2. Set the environment variables (`DB_USER`, `DB_PASS`, `DB_HOST`, `DB_NAME`) securely on your host.
3. Run the binary as a background service (e.g., using `systemd`).

### Scraper (Python FastAPI)
1. Run the FastAPI application using a production ASGI server like Gunicorn with Uvicorn workers:

gunicorn -k uvicorn.workers.UvicornWorker main:app --bind 0.0.0.0:8000 --workers 4 &


### Database (MySQL)
* Use a managed database service (e.g., AWS RDS, Google Cloud SQL) for better reliability and backups, or ensure your production Docker volumes are properly mounted and backed up.
