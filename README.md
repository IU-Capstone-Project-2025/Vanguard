# 📚 Vanguard

Vanguard is an interactive learning platform that allows users to create and participate in real-time quizzes. Designed for teachers, trainers, and event hosts, it brings engagement and gamification to learning and assessments.

## ✨ Features

- 🎨 Create and host dynamic quizzes in real time
- 📱 Participants can join from any device via a simple interface
- 🧠 Gamified experience with timers, points, and live feedback
- 📊 Detailed performance analytics for insights and evaluation
- 🏆 Leaderboard for competitive fun
- 👥 Ideal for classrooms, corporate training, and social events

## 🚀 Getting Started with Docker Compose

### 1. Clone the repository

```sh
git clone https://github.com/IU-Capstone-Project-2025/Vanguard.git
cd Vanguard
```

---

## 🛠 Development Deployment

### 1. Copy or create the environment file

```sh
cp .env.dev.example .env
```

### 2. Start without monitoring

```sh
docker compose --env-file .env -f docker-compose.yaml -f docker-compose.dev.yaml build
docker compose --env-file .env -f docker-compose.yaml -f docker-compose.dev.yaml up -d frontend
```

### 3. Start with monitoring (Grafana)

```sh
docker compose --env-file .env -f docker-compose.yaml -f docker-compose.dev.yaml up -d --build
```

### 4. Access the platform

* Frontend available at: [http://localhost:3000](http://localhost:3000)
* Grafana available at: [http://localhost:3001](http://localhost:3001)

---

## 🔐 Production Deployment

### 1. Copy the production environment file

```sh
cp .env.prod.example .env
```

### 2. Run SSL setup (required before first deployment)

```sh
chmod +x setup.sh
./setup.sh
```

### 3. Start all services with monitoring

```sh
docker compose --env-file .env -f docker-compose.yaml -f docker-compose.prod.yaml up -d --build
```

### 4. Access the platform

* Frontend available at: [https://tryit.selnastol.ru](https://tryit.selnastol.ru)
* Grafana available at: [https://grafana.tryit.selnastol.ru](https://grafana.tryit.selnastol.ru)

---

## 🧹 Stopping Services

To stop and remove all running containers:

- For dev:
    ```sh
    docker compose -f docker-compose.yaml -f docker-compose.dev.yaml down
    ```

- For prod:
    ```sh
    docker compose -f docker-compose.yaml -f docker-compose.prod.yaml down
    ```
