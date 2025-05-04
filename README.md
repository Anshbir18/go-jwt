# 🛡️ Role-Based Authentication API in Go

A secure, scalable RESTful API built with **Go**, featuring **JWT authentication**, **role-based access control**, **request validation**, and a **MongoDB (Dockerized)** backend. Designed with clean architecture principles and real-world middleware patterns.

---

## 🚀 Features

- ✅ JWT Authentication (Access + Refresh Tokens)  
- ✅ Role-Based Access Control (Admin/User)  
- ✅ Secure Password Hashing (bcrypt)  
- ✅ Input Validation for Clean Requests  
- ✅ MongoDB Integration via Docker  
- ✅ Pagination with MongoDB Aggregation Pipeline  
- ✅ Middleware for Authorization Handling  
- ✅ Built using `gin-gonic` for fast and clean routing

---

## 🧰 Tech Stack

- **Go (Golang)**  
- **Gin** – Routing & Middleware  
- **MongoDB** – Dockerized  
- **JWT** – Token-based Authentication  
- **bcrypt** – Password Hashing  
- **validator/v10** – Request Validation

---

## 📦 Setup

### 1. Clone the Repo
```bash
git clone https://github.com/yourusername/yourproject.git
cd yourproject
```
### 2. Start MongoDB via Docker
```bash
docker run -d -p 27017:27017 --name goauth-mongo mongo
```

### 3. Run the Server
```bash 
go mod tidy
go run main.go

```

## 🛠️ API Endpoints
| Method | Endpoint         | Description                              |
|--------|------------------|------------------------------------------|
| POST   | `/signup`        | Register a new user                      |
| POST   | `/login`         | Login and receive tokens                 |
| GET    | `/users`         | Get paginated user list *(Admin only)*   |
| GET    | `/users/:user_id`| Get details of a specific user *(Admin)* |
