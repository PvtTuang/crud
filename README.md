# CRUD + Authentication Project

โปรเจคนี้เป็นระบบ **CRUD พร้อมระบบ Authentication** ที่พัฒนาโดยใช้ **Go (Gin Framework)** สำหรับ backend และ **Frontend** สำหรับ interface ของผู้ใช้ โดยรองรับการเชื่อมต่อกับ **PostgreSQL** และ **Redis** สำหรับจัดการข้อมูลและ session/cache

---

## 📂 โครงสร้างโปรเจค

backend/
├── auth/ # Auth service
│ ├── cmd/server # Main server ของ Auth service
│ └── pkg # Logic ของระบบ Auth
├── crud/ # CRUD service
│ ├── cmd/server # Main server ของ CRUD service
│ └── pkg # Logic ของ CRUD
└── database/ # การเชื่อมต่อฐานข้อมูลและ migration

---

## ⚙️ เทคโนโลยีที่ใช้

- **Backend:** Go, Gin Framework
- **Database:** PostgreSQL
- **Cache/Session:** Redis
- **Frontend:** Vite / React (หรือรายละเอียด frontend ตามจริง)
- **Containerization:** Docker, Docker Compose

---

## 🚀 การติดตั้งและรันโปรเจค

1. **Clone โปรเจค**
- git clone https://github.com/PvtTuang/crud.git
- cd crud

2. สร้างไฟล์ .env
# PostgreSQL
- POSTGRES_HOST=...
- POSTGRES_PORT=...
- POSTGRES_USER=...
- POSTGRES_PASSWORD=...
- POSTGRES_DB=...

# Redis
-REDIS_HOST=...     
-REDIS_PORT=...
-REDIS_PASSWORD=...
-REDIS_DB=...

# Go backend
-AUTH_HTTP_PORT=...
-CRUD_HTTP_PORT=...

#JWT
-JWT_SECRET=...

# Frontend
-REACT_APP_API_URL=...

# Backend
-VITE_AUTH_API_URL=...
-VITE_CRUD_API_URL=...

3. รัน Docker Compose
-docker-compose up -d --build

4. ตรวจสอบ container
-docker ps
*auth_service → รันที่พอร์ต 80..
*crud_service → รันที่พอร์ต 80..
*frontend_service → รันที่พอร์ต ...

🔧 ฟีเจอร์หลัก
Backend
*ระบบลงทะเบียนผู้ใช้ (Register)
*ระบบเข้าสู่ระบบ (Login)
*ระบบจัดการข้อมูล (CRUD) สำหรับแหล่งข้อมูลต่างๆ
*เชื่อมต่อ PostgreSQL และ Redis
Frontend
*หน้าจอสำหรับ Register/Login
*หน้าจอ CRUD สำหรับผู้ใช้งาน
*เชื่อมต่อ API กับ backend

📝 หมายเหตุ

โปรเจคนี้ยังอยู่ในระหว่างการพัฒนา
ข้อมูลสำคัญใน .env ห้ามใส่ใน GitHub
