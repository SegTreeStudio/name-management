# Name Management System

A simple full-stack application for adding names, checking whether a
name exists, and viewing all saved names.

## Tech stack

-   **Backend:** Go (`net/http`)
-   **Database:** PostgreSQL 17
-   **Frontend:** React, TypeScript, Tailwind CSS, and Vite
-   **API testing:** Postman
-   **Database administration:** pgAdmin 4

## Project structure

``` text
name-management/
├── backend/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── database/
│   │   ├── handler/
│   │   ├── model/
│   │   └── repository/
│   ├── .env
│   ├── .env.example
│   ├── go.mod
│   └── go.sum
└── frontend/
    ├── src/
    ├── .env
    ├── package.json
    └── ...
```

> The exact files may differ as the project evolves. Keep secrets such
> as database passwords in local environment files, not in source
> control.

## Prerequisites

Install the following before running the project:

1.  **Go** --- install a current stable version from
    [go.dev/dl](https://go.dev/dl/).
2.  **PostgreSQL 17** --- install PostgreSQL locally from
    [postgresql.org](https://www.postgresql.org/download/).
3.  **pgAdmin 4** --- useful for creating and inspecting the database.
    It is available from
    [pgadmin.org](https://www.pgadmin.org/download/).
4.  **Node.js LTS and npm** --- install from
    [nodejs.org](https://nodejs.org/).
5.  **Git** --- optional, but recommended.
6.  **Postman** --- optional, for testing the API. Download from
    [postman.com](https://www.postman.com/downloads/).

Docker is **not required** for this setup because PostgreSQL runs as a
local Windows service.

## 1. Get the project

If you have the repository URL:

``` powershell
git clone <REPOSITORY_URL>
cd name-management
```

Alternatively, download and extract the project, then open the
`name-management` folder in your terminal or editor.

## 2. Set up PostgreSQL

### 2.1 Start PostgreSQL

On Windows:

1.  Press `Win + R`.
2.  Enter `services.msc`.
3.  Find the PostgreSQL 17 service (often named `postgresql-x64-17`).
4.  Start it if it is stopped.

The default PostgreSQL connection settings are commonly:

  Setting    Value
  ---------- -------------------------------------------------------
  Host       `localhost`
  Port       `5432`
  Username   `postgres`
  Database   `name_management`
  Password   The password you chose during PostgreSQL installation

If you selected a different port, use that port in the backend
configuration.

### 2.2 Create the database

In pgAdmin 4:

1.  Connect to your local PostgreSQL server.
2.  Right-click **Databases** and choose **Create → Database**.
3.  Set the database name to `name_management`.
4.  Set the owner to `postgres` (or your PostgreSQL user).
5.  Save.

### 2.3 Create the names table

Select the `name_management` database, open **Tools → Query Tool**, and
run:

``` sql
CREATE TABLE IF NOT EXISTS names (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_names_name_lower
ON names (LOWER(name));
```

The unique index prevents duplicate names even when their letter casing
differs (for example, `Rahim` and `RAHIM`).

To verify the table:

``` sql
SELECT * FROM names ORDER BY id;
```

## 3. Configure and run the Go backend

Open a terminal in the backend folder:

``` powershell
cd backend
```

### 3.1 Configure environment variables

Create a file named `.env` inside `backend/`:

``` env
PORT=8080
DATABASE_URL=postgres://postgres:YOUR_PASSWORD@localhost:5432/name_management?sslmode=disable
FRONTEND_URL=http://localhost:5173
```

Replace `YOUR_PASSWORD` with the password for your local PostgreSQL
user.

If your password contains URL-reserved characters (such as `@`, `#`,
`/`, or `:`), URL-encode those characters in the connection URL. Do not
share or commit your `.env` file.

The repository should include an `.env.example` with placeholder values,
not real credentials.

### 3.2 Install Go dependencies

From the `backend` directory:

``` powershell
go mod tidy
```

### 3.3 Format, build, and run

``` powershell
go fmt ./...
go build ./...
go run ./cmd/api
```

Keep this terminal open while using the application.

The API should be listening at:

``` text
http://localhost:8080
```

Check the health endpoint in a browser or API client:

``` text
GET http://localhost:8080/health
```

A successful response should return HTTP `200 OK` and `OK`.

## 4. Test the API

You can test the backend with Postman before starting the frontend.

### 4.1 Health check

-   Method: `GET`
-   URL: `http://localhost:8080/health`

Expected status: `200 OK`.

### 4.2 Add a name

-   Method: `POST`
-   URL: `http://localhost:8080/api/names`
-   In Postman, select **Body → raw → JSON**.

Request body:

``` json
{
  "name": "Rahim"
}
```

A successful creation should return HTTP `201 Created`. The response
body follows the `NameResponse` structure, for example:

``` json
{
  "message": "Name added successfully",
  "name": "Rahim"
}
```

The exact message depends on the backend handler implementation.

### 4.3 Check whether a name exists

-   Method: `GET`
-   URL: `http://localhost:8080/api/names/exists?name=Rahim`

Example response:

``` json
{
  "name": "Rahim",
  "exists": true
}
```

Try a name that has not been added to check the `false` case.

### 4.4 Get all names

-   Method: `GET`
-   URL: `http://localhost:8080/api/names`

The response should be a JSON array of saved records, including their
IDs and creation timestamps.

### API endpoint summary

  ---------------------------------------------------------------------------------
  Method                  Endpoint                          Purpose
  ----------------------- --------------------------------- -----------------------
  `GET`                   `/health`                         Check whether the API
                                                            is running

  `POST`                  `/api/names`                      Add a name

  `GET`                   `/api/names/exists?name=<name>`   Check whether a name
                                                            exists

  `GET`                   `/api/names`                      Get all saved names
  ---------------------------------------------------------------------------------

**Note:** Opening a URL directly in a browser sends a `GET` request. To
add a name, send a `POST` request with a JSON body using Postman or the
frontend.

## 5. Configure and run the React frontend

Open a **new terminal** (leave the backend running).

``` powershell
cd frontend
```

Install dependencies:

``` powershell
npm install
```

Create `frontend/.env`:

``` env
VITE_API_URL=http://localhost:8080/api
```

Start the Vite development server:

``` powershell
npm run dev
```

Vite will print the local development URL, usually:

``` text
http://localhost:5173
```

Open that URL in your browser.

If you change a Vite environment variable, restart the frontend
development server for the change to take effect.

## 6. Run the project locally

For normal development, use two terminals:

**Terminal 1 --- backend**

``` powershell
cd name-management\backend
go run ./cmd/api
```

**Terminal 2 --- frontend**

``` powershell
cd name-management\frontend
npm run dev
```

Also make sure the PostgreSQL 17 Windows service is running.

Open the frontend URL printed by Vite, usually `http://localhost:5173`.

## 7. Troubleshooting

### `404 page not found`

-   Confirm that the backend is running on port `8080`.
-   Check the exact path and HTTP method. For example, adding a name
    requires `POST /api/names`, not `GET /api/names`.
-   Confirm that the routes are registered in `backend/cmd/api/main.go`.

### Database connection refused

-   Make sure the PostgreSQL Windows service is running.
-   Check the host and port in `DATABASE_URL`.
-   Confirm that the database `name_management` exists.

### Password authentication failed

-   Check the PostgreSQL username and password in `DATABASE_URL`.
-   Remember that the pgAdmin master password is not necessarily the
    PostgreSQL server password.
-   URL-encode special characters in the password when placing it in the
    connection URL.

### `go build` reports an import or module error

-   Run commands from the `backend` directory.
-   Ensure internal imports use the same module path declared at the top
    of `backend/go.mod`.
-   Run:

``` powershell
go mod tidy
go fmt ./...
go build ./...
```

### Frontend cannot connect to the backend

-   Confirm that the backend is running.
-   Check `VITE_API_URL` in `frontend/.env`.
-   Check that the backend CORS configuration allows
    `http://localhost:5173`.
-   Restart Vite after editing `.env`.

### Port already in use

Another process may already be using port `8080` or `5173`. Stop that
process or configure the application to use another available port, and
update the corresponding environment variable.

## 8. Security and contribution notes

-   Never commit `.env` files, database passwords, API keys, or other
    secrets.
-   Keep `.env.example` populated with placeholders so other developers
    know which variables are required.
-   Do not expose the PostgreSQL server publicly for local development.
-   Validate user input in the backend; frontend validation alone is not
    sufficient.
-   Before making database schema changes, consider how existing local
    databases will be migrated.
-   If you contribute code, format Go files with `go fmt ./...` and
    verify the backend with `go build ./...`.

## License

Add the project's license and copyright information here.
