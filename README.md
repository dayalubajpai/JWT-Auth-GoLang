JWT Auth API
---

```markdown
# JWT Authentication API in Go

A simple RESTful API built using **Go (Golang)** and **MongoDB** that supports user authentication with **JWT (JSON Web Token)**. Users can sign up, log in, and fetch profile data securely using token-based authentication.

---

## 🚀 Technologies Used

- **Golang (Gin Framework)**
- **MongoDB (Official Go Driver)**
- **JWT (github.com/golang-jwt/jwt/v5)**
- **bcrypt (Password hashing)**
- **dotenv (Configuration)**

---

## 📦 Installation

1. Clone the repository:

```bash
git clone https://github.com/dayalubajpai/jwt-auth-go.git
cd jwt-auth-go
```

2. Create a `.env` file in the root directory:

```
MONGODB_URL=mongodb+srv://<username>:<password>@<cluster>.mongodb.net/?retryWrites=true&w=majority
```

3. Install dependencies:

```bash
go mod tidy
```

4. Run the project:

```bash
go run main.go
```

---

## 📂 API Endpoints

### 🔐 Sign Up

**URL:** `POST http://localhost:8080/users/signup/`

**Headers:**
- `Content-Type: application/json`

**Request Body:**

```json
{
    "First_name": "Malou",
    "Last_name": "David",
    "Password": "Maloudavid",
    "Phone": "9867326783",
    "Email": "maloudavid@gmail.com",
    "User_type": "ADMIN"
}
```

---

### 🔑 Log In

**URL:** `POST http://localhost:8080/users/login/`

**Headers:**
- `Content-Type: application/json`

**Request Body:**

```json
{
    "Email": "john@gmail.com",
    "password": "johncena"
}
```

**Response:**
```json
{
    "token": "JWT_TOKEN_HERE",
    "refresh_token": "REFRESH_TOKEN_HERE"
}
```

---

### 👤 Get Single User

**URL:** `GET http://localhost:8080/users/?user_id=67efecce2794311724402f55`

**Headers:**
- `token`: `<JWT token from login>`

---

### 👥 Get All Users (Paginated)

**URL:** `GET http://localhost:8080/users/`

**Headers:**
- `token`: `<JWT token from login>`

**Query Parameters:**
- `pageNumber`: (optional) Page number for pagination (default: 1)

---

## 🔒 Token-based Authentication

All protected routes require the `token` header. Tokens are generated on login and must be included in subsequent requests to access protected data.

---

## 📁 Folder Structure

```
.
├── controllers/
│   └── userController.go
├── database/
│   └── databaseConnection.go
├── routes/
│   └── authRoutes.go
│   └── userRoutes.go
├── models/
│   └── userModel.go
├── helpers/
│   └── tokenHelper.go
│   └── authHelper.go
├── middleware/
│   └── authMiddleware.go
├── main.go
├── go.mod
├── .env
```

---

## 🧪 Test Using Postman

- Import requests manually or use the above endpoints.
- Use **"raw JSON"** for body data.
- Don't forget to include the **`token`** in headers for protected routes.

---

```
