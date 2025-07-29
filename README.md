````markdown
# SETA Golang Training

🏗 **Training Project:** Microservices for User, Team & Asset Management

This repository is part of the SETA Golang training series. It demonstrates how to build a microservices-based system using Go, focusing on user management, team collaboration, and asset sharing.

---

## 📚 Overview

This system is split into three main services with clear domain boundaries:

1. **User Service (GraphQL)** – manages users, roles, authentication  
2. **Team Service (REST)** – manages teams, members, and managers  
3. **Asset Service (REST)** – manages folders, notes, and sharing between users  

---

## 🔧 Architecture

```text
             ┌────────────┐
             │  Client    │
             └────┬───────┘
                  │
    ┌─────────────┼──────────────┐
    │             │              │
┌───▼────────┐ ┌───▼────────┐ ┌───▼────────────┐
│ User Svc   │ │ Team Svc   │ │ Asset Svc      │
│ (GraphQL)  │ │ (REST API) │ │ (REST API)     │
└────────────┘ └────────────┘ └────────────────┘
````

---

## 🧩 Microservices Breakdown

### 1. 🧑‍💼 **User Service (GraphQL)**

> Handles user creation, login, logout, and role management.

* **Stack:** Go, gqlgen, JWT, PostgreSQL
* **Endpoints:**

  * `createUser(username, email, role)`
  * `login(email, password)`
  * `logout`
  * `fetchUsers`
  * `assignRole(userId, role)`
* **Roles:**

  * `manager`: can create/manage teams
  * `member`: can only be added to teams

---

### 2. 👥 **Team Service (REST)**

> RESTful service for team creation and member management.

* **Stack:** Go, Gin, PostgreSQL
* **Features:**

  * Create a team
  * Add/remove members
  * Add/remove other managers (only by main manager)
* **Entities:**

  * `teamId`, `teamName`
  * `managers[]`, `members[]`

---

### 3. 🗂 **Asset Service (REST)**

> Folder and note management with access control and sharing.

* **Entities:**

  * **Folders** – owned by users, contain notes
  * **Notes** – belong to folders, have content
* **Permissions:**

  * Share folders/notes (read or write)
  * Revoke access at any time
  * Shared folder = all notes inside are also shared
  * Managers can view assets from team members (read-only)

---

## 🛠️ Project Structure

```text
.
├── services/
│   ├── user-service/    # GraphQL-based auth & user logic
│   ├── team-service/    # RESTful team management
│   └── asset-service/   # RESTful asset sharing
├── migration/           # DB migration logic and models
├── shared/              # Shared packages: logger, db, error handling
├── docker-compose.yml   # Multi-service orchestration
├── promtail-config.yml  # Log shipping config
├── Taskfile.yaml        # Task runner for development
└── token.json           # JWT or token storage
```

---

## 🐳 Run with Docker

```bash
docker-compose up --build
```

> 📦 Ensure you’ve set up `.env` files for each service if needed.

---

## 🚀 Development

Use [`Taskfile`](https://taskfile.dev) to simplify common operations:

```bash
task migrate       # Run DB migrations
task user          # Start user-service
task team          # Start team-service
```

---

## 📈 Logging & Monitoring

Integrated with **Promtail + Loki + Grafana**:

* `promtail-config.yml` configures log scraping
* Logs go to Loki and can be queried via Grafana dashboard

---

## ✨ Contribution Guide

* Keep each service isolated and reusable
* Use gRPC or REST depending on service role
* Favor clear error handling with `shared/apperror`
* Use consistent logging (`shared/logger`)

---

## 📄 License

MIT – Free to learn, use, and modify.

---

## 🙌 Credits

Project built as part of **SETA Training Program** to learn production-grade Golang microservice architecture.

---

