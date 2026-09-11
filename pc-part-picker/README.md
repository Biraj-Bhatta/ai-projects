# AI PC Part Picker

This is a full-stack web application that generates price-optimized, compatible PC builds based on user-defined workloads, utilizing live AI-scraped pricing data.

## Project Structure

*   **`frontend/`**: Next.js (App Router) with React and Tailwind CSS. Provides the user interface for configuring builds.
*   **`backend/`**: Go using the Echo framework. Handles the core compatibility logic, budget allocation, and communicates with the database and scraper.
*   **`scraper/`**: Python FastAPI microservice utilizing `crawl4ai`. Fetches live pricing data for components.
*   **`docker-compose.yml`**: Configures Docker for spinning up the complete environment (MySQL, Scraper, Backend, Frontend).

---

## Prerequisites

*   Docker and Docker-Compose installed on your local machine.

---

## Getting Started (Fully Dockerized)

This project is configured to run entirely via Docker.

### 1. Start the Environment

From the root `pc-part-picker` directory, run the following command to start all services:

```bash
docker-compose up -d --build
```

Docker will build and spin up the following containers:
*   **`pc-part-picker-db` (MySQL)**: Available internally on port `3306`. (Initialized automatically via `backend/schema.sql`).
*   **`pc-part-picker-scraper` (Python)**: Available on `http://localhost:8000`.
*   **`pc-part-picker-backend` (Go API)**: Available on `http://localhost:8080`.
*   **`pc-part-picker-frontend` (Next.js)**: Available on `http://localhost:3000`.

### 2. Access the Application

Once the containers are successfully running, open your web browser and navigate to:

**http://localhost:3000**

You can now configure your AI PC build.

### 3. Stopping the Environment

To stop the running containers:

```bash
docker-compose down
```

---

## Architecture Details

*   **Database**: The MySQL container uses a volume `db_data` to persist data. Upon first creation, it automatically executes the `schema.sql` to setup tables and seed data.
*   **Networking**: Services communicate via the default Docker bridge network. The Go backend references the database at `db:3306` and the Python scraper at `scraper:8000`.
*   **Frontend API Config**: The Next.js frontend calls the backend API at `http://localhost:8080`. Note that in a true production environment, `localhost` calls from the client browser need to be routed through an API Gateway, Nginx, or explicitly defined via the `NEXT_PUBLIC_API_URL` environment variable.
