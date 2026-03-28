# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Billadm-GO is a desktop personal finance application built with Electron + Vue.js + Go. It manages transaction records across multiple ledgers with categories and tags.

## Architecture

```
kernel/          # Go backend (Gin HTTP server)
  api/           # HTTP handlers/controllers
  dao/           # Data access objects (SQL)
  models/        # Domain models and DTOs
  service/       # Business logic layer (Logger interface for logging)
  workspace/     # Workspace management (one SQLite DB per workspace)
  pkg/operator/  # Query builder for filtering/sorting/paging
  util/          # Utilities (DB, config, logging, UUID, etc.)

app/             # Vue.js frontend
  src/
    backend/     # API client wrappers
    components/  # Vue components
    stores/      # Pinia state management

electron/        # Electron main process
  src/           # Electron entry (main.js, preload.js)
```

### Go Backend Layers

- **api**: HTTP handlers, parse requests, call services, return responses
- **service**: Business logic, transactions, logging via Logger interface
- **dao**: Database CRUD, no business logic
- **workspace**: Database lifecycle, transaction support (`Workspace.Transaction(fn)`)

## Database Schema

Each workspace contains its own SQLite database (`billadm.db`) with tables:
- `tbl_billadm_ledger` - ledgers
- `tbl_billadm_transaction_record` - transaction records (expense/income/transfer)
- `tbl_billadm_transaction_record_tag` - tags for transactions
- `tbl_billadm_category` - expense categories
- `tbl_billadm_tag` - tags organized by category

## Key Commands

**Backend (Go kernel):**
```bash
cd kernel && go build -ldflags '-s -w -extldflags "-static"' -o Billadm-Kernel.exe
# Runs on 127.0.0.1:31943
```

**Frontend (Vue dev):**
```bash
cd app && npm run dev
```

**Electron:**
```bash
cd electron && npm start
```

**Build full application (Windows):**
```powershell
./build/build.ps1
```

## Development (Hot Reload)

Three processes must run simultaneously:
1. Go backend: run `kernel/main.go` in GoLand/IDE
2. Vue dev server: `npm run dev` in `app/`
3. Electron: `npm start` in `electron/`

## API Conventions

- Base URL: `http://127.0.0.1:31943`
- All endpoints use POST with JSON body
- Response format: `{"code": 0, "msg": "", "data": ...}`
- Non-zero code indicates error

Key endpoints:
- `/api/v1/ledger/*` - ledger CRUD
- `/api/v1/tr/*` - transaction record operations
- `/api/v1/category/query/:type` - query categories by transaction type
- `/api/v1/tag/query/:category` - query tags by category
- `/api/v1/workspace/open` - open workspace directory

## Configuration

Kernel accepts flags: `--mode`, `--port`, `--log-level`, `--workspace`. See `kernel/util/config.go`.
