
# Bank API

## Description
Bank API is a simple application that provides various banking functionalities. It uses PostgreSQL as the database backend.

---

## Tasks

### 1. **Build** - `task build`
This command builds the application and outputs an executable at `bin/bank.exe`.

To build the application, run the following command:

```bash
task build
```

### 2. **Run** - `task run`
This command builds the application and runs it in a single step.

To build and run the application, use:

```bash
task run
```

### 3. **Test** - `task test`
This command runs the tests for the application to ensure everything is functioning as expected.

To run the tests, execute:

```bash
task test
```

---

## Setup

### Install Dependencies
Make sure you have the necessary dependencies installed:

- **Go (Golang)**: Required for building and running the application.
- **PostgreSQL**: Required as the database backend.

### Database Configuration
Set up a PostgreSQL database with the necessary schema. You can configure the connection settings in the application configuration file.

---

## Requirements

- **PostgreSQL**: For the database.
- **Go (Golang)**: For building and running the application.

---

## Docker Setup

To run the PostgreSQL database in Docker, you can use the provided `docker-compose.yml` file. This will spin up a PostgreSQL container with the necessary configurations.

### Running PostgreSQL using Docker Compose
1. Ensure you have Docker and Docker Compose installed.
2. Run the following command to start the PostgreSQL container:

```bash
docker-compose up -d
```

3. This will start PostgreSQL on `localhost:5432` with the following credentials:
   - **Username**: `admin`
   - **Database**: `bank`
   - **Password**: `password`

---


## Notes:
- If you encounter any issues while setting up or running the application, please open an issue or contact the maintainers.
