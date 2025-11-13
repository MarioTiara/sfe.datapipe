# 📦 SFE Datapipe

## 📝 Project Description
This project is a Go-based data pipeline designed to load, transform, and store data from multiple sources into a database.  
It includes features like:
- Reading Excel 
- Mapping data into domain entities
- Saving data into a PostgreSQL database
- Database migration management using `migrate` tool

---

## ⚙️ Prerequisites

Before running the project, make sure you have installed:

- [Go 1.25.3](https://go.dev/dl/)
- [PostgreSQL](https://www.postgresql.org/download/)
- [golang-migrate CLI](https://github.com/golang-migrate/migrate)
- Create a `data/.env` file in the project root:


## 🗄️ Run Migration
Make sure you already created the database before running migrations.

```bash
./migrate_up.ps1

```
## 🚀 Run Project

```bash
go run ./cmd/console
```