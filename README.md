# Bank API

## Description
Bank API is a simple application that provides various banking functionalities. It uses PostgreSQL as the database backend.

## Tasks

### 1. **Build** - `task build`
This command builds the application and outputs an executable at `bin/bank.exe`.

To build the application, run the following command:

```bash
task build



### 1. **Run** - `task run`
This command builds the application and runs it in a single step `bin/bank.exe`.

To build the application, run the following command:

```bash
task run


### 1. **Test** - `task test`
This command runs the tests for the application to ensure everything is functioning as expected..

To run the tests, execute:

```bash
task test



Setup
Install Dependencies
Make sure you have the necessary dependencies installed (e.g., Go, PostgreSQL).

Database Configuration
Set up a PostgreSQL database with the necessary schema. You can configure the connection settings in the application configuration file.

Requirements
PostgreSQL for the database.

Go (Golang) for building and running the application.

Docker Setup
To run the PostgreSQL database in Docker, you can use the provided docker-compose.yml file. This will spin up a PostgreSQL container with the necessary configurations.

Running PostgreSQL using Docker Compose
Ensure you have Docker and Docker Compose installed.

Run the following command to start the PostgreSQL container:


```bash
docker-compose up -d