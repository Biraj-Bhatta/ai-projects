# AI PC Part Picker

This is a full-stack web application that generates price-optimized, compatible PC builds based on user-defined workloads, utilizing live AI-scraped pricing data.

## Project Structure

*   **`frontend/`**: Next.js (App Router) with React and Tailwind CSS. Provides the user interface for configuring builds.
*   **`backend/`**: Go using the Echo framework. Handles the core compatibility logic, budget allocation, and communicates with the database and scraper.
*   **`scraper/`**: Python FastAPI microservice utilizing `crawl4ai`. Fetches live pricing data for components.
*   **`docker-compose.yml`**: Configures Docker for spinning up the MySQL Database and Python AI Scraper.

---

## Local Hybrid Setup (Docker + Localhost)

To run the MySQL Database and Python AI Scraper in Docker containers, while running the Frontend (Next.js) and Backend (Go) natively on your local machine, follow these steps:

### 1. Start Docker Containers (Scraper + Database)

From the `pc-part-picker` directory, run the following command to start both the MySQL database and the Python AI scraper in detached mode:

docker-compose up -d --build


*   **Database:** Runs on `localhost:3306`. (Credentials: User: `root`, Pass: `root`, DB: `pcpartpicker`). The schema will be automatically initialized.
*   **Scraper:** Runs on `http://localhost:8000`.

### 2. Start the Backend Service (Go on Localhost)

Open a new terminal, navigate to the backend directory, and run the Go API locally:

cd pc-part-picker/backend
go mod tidy
go run main.go &


*   The backend will run on `http://localhost:8080`.
*   It is pre-configured to connect to the MySQL database at `127.0.0.1:3306` and the Scraper at `http://localhost:8000`.

### 3. Start the Frontend Service (Next.js on Localhost)

Open a new terminal, navigate to the frontend directory, install dependencies, and run the Next.js server locally:

cd pc-part-picker/frontend
npm install
npm run dev &


*   The frontend will be accessible at `http://localhost:3000`.
*   It is pre-configured to make requests to your local Go backend at `http://localhost:8080/api/build`.

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
1. The `docker-compose.yml` already contains a robust configuration for the scraper, but in production, you might modify its Dockerfile to use an ASGI server like Gunicorn with Uvicorn workers.

### Database (MySQL)
* Use a managed database service (e.g., AWS RDS, Google Cloud SQL) for better reliability and backups, or ensure your production Docker volumes are properly mounted and backed up.
