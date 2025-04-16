# CBCExams Backend 🧠📚

This is the backend API for **CBCExams**, a learning and tutoring platform that offers access to categorized resources, tutoring jobs, payment integration, feedback, and more — tailored for students and tutors under the CBC curriculum.

Built with **Go (Golang)** and **Gin**, and integrates with **PostgreSQL**, **Pesapal** for payments, and email services.

---

## 📦 Project Structure

```bash
cbcexams-backend/
├── config/          # App configuration
├── controllers/     # HTTP handlers for various features
├── database/        # PostgreSQL connection and initialization
├── middleware/      # Authentication middleware (JWT)
├── models/          # GORM models for database tables
├── pesapal/         # Payment integration with Pesapal API
├── routes/          # Route registration for API endpoints
├── utils/           # Helpers for JWT, email, and tokens
├── uploads/         # Directory for file uploads (e.g., resumes)
├── main.go          # Application entry point
├── common.sql       # SQL schema (optional)
├── go.mod/go.sum    # Go dependencies
```

---

## 🚀 Features

- ✅ User authentication (register, login, JWT)
- ✅ Tutor and job listing system
- ✅ Educational resource uploads/bookmarks
- ✅ Categorized content (web development, jobs, tutoring)
- ✅ Feedback collection
- ✅ Secure file upload handling
- ✅ Pesapal payment integration
- ✅ Email notifications (e.g., password reset)
- ✅ Middleware for protected routes
- ✅ CORS & secure host filtering

---

## ⚙️ Setup Instructions

### 1. Clone the repository

```bash
git clone https://github.com/yourusername/cbcexams-backend.git
cd cbcexams-backend
```

### 2. Create `.env`

Set your environment variables (DB, JWT secret, email config, Pesapal keys, etc.):

```env
PORT=8080
DB_URL=postgres://user:password@localhost:5432/cbcbackend
JWT_SECRET=your_secret
EMAIL_USER=your_email@example.com
EMAIL_PASS=your_password
PESAPAL_CONSUMER_KEY=your_key
PESAPAL_CONSUMER_SECRET=your_secret
```

> Make sure `.env` is in your `.gitignore`.

### 3. Install dependencies

```bash
go mod tidy
```

### 4. Run the app

```bash
go run main.go
```

Or compile and run:

```bash
go build -o cbcexams-backend
./cbcexams-backend
```

---

## 🛠️ Deployment

To deploy:

- Use **systemd** for process management
- Use **Nginx** as a reverse proxy
- Use **Certbot** for SSL (Let's Encrypt)

For more details, see the [Deployment Guide](#).

Server metrics can be monitored at: [http://102.209.68.145:19999/](http://102.209.68.145:19999/)

---

## 🔐 API Endpoints

API is organized into feature-based route groups:

| Feature        | Path Prefix     |
|----------------|------------------|
| Auth           | `/api/auth`      |
| Users          | `/api/users`     |
| Resources      | `/api/resources` |
| Bookmarks      | `/api/bookmark`  |
| Tutoring       | `/api/tutoring`  |
| Web Dev        | `/api/webdev`    |
| Feedback       | `/api/feedback`  |
| Payments       | `/api/payment`   |

You can find all the endpoints documented in the link provided: [Postman Documentation](https://documenter.getpostman.com/view/23285423/2sB2ca5eUP#7c8a41fe-cc56-4088-8f72-02da505472a9)

---

## 🧪 Testing

You can test endpoints using [Postman](https://postman.com) or cURL.

Example:

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com", "password":"secret"}'
```

---

## 📄 License

MIT License — free to use and modify.

---

## 👨‍💻 Author

Built with ❤️ by [@bot-on-tapwater](https://github.com/Bot-on-Tapwater)