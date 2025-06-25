# Project: js-script-service

This microservice executes administratively configured JavaScript scripts from RESTful endpoints.

## 🚀 Features
- Each script has its own GET/POST endpoint
- Scripts stored in PostgreSQL (Aurora in production)
- JS executed using `goja` engine
- Resilient design with logging
- Supports local and AWS Lambda deployment

## 📦 Project Structure
```
├── Makefile
├── go.mod
├── go.sum
├── .env
├── docker-compose.yml
├── cmd/
│   ├── server/
│   │   └── main.go
│   └── lambda/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── db/
│   │   ├── postgres.go
│   │   └── models.go
│   ├── executor/
│   │   └── executor.go
│   ├── router/
│   │   └── router.go
│   └── logger/
│       └── logger.go
└── scripts/
    └── init.sql
```

## 🛠️ Setup Commands
```bash
git init
git remote add origin git@github.com:peter/js-script-service.git
git add .
git commit -m "Initial scaffold"
git push -u origin main
```

## 📄 scripts/init.sql
```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS scripts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    path TEXT UNIQUE NOT NULL,
    code TEXT NOT NULL,
    enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS execution_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    script_id UUID REFERENCES scripts(id) ON DELETE CASCADE,
    timestamp TIMESTAMP DEFAULT NOW(),
    input JSONB,
    output JSONB,
    error TEXT
);
```

## 🧪 Running Locally
```bash
docker-compose up -d
go run ./cmd/server
```

## 🧼 .env Example
```env
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=admin
POSTGRES_PASSWORD=secret
POSTGRES_DB=scriptdb
ENV=development
```

---

Now you're ready to implement the core service!
