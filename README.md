<div align="center">

# TRY-IT
*Interactive Learning Through Real-Time Gamification*

![last-commit](https://img.shields.io/github/last-commit/IU-Capstone-Project-2025/Vanguard?style=flat&logo=git&logoColor=white&color=0080ff)
![repo-top-language](https://img.shields.io/github/languages/top/IU-Capstone-Project-2025/Vanguard?style=flat&color=0080ff)
![repo-language-count](https://img.shields.io/github/languages/count/IU-Capstone-Project-2025/Vanguard?style=flat&color=0080ff)

*Built with:*

![Redis](https://img.shields.io/badge/Redis-FF4438.svg?style=flat&logo=Redis&logoColor=white)
![RabbitMQ](https://img.shields.io/badge/RabbitMQ-FF6600.svg?style=flat&logo=RabbitMQ&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1.svg?style=flat&logo=PostgreSQL&logoColor=white)
![Grafana](https://img.shields.io/badge/Grafana-F46800.svg?style=flat&logo=Grafana&logoColor=white)
![Prometheus](https://img.shields.io/badge/Prometheus-E6522C.svg?style=flat&logo=Prometheus&logoColor=white)
![Loki](https://img.shields.io/badge/Loki-2C2D72?style=flat&logo=Loki&logoColor=white)
![Alertmanager](https://img.shields.io/badge/Alertmanager-FF6D00?style=flat)
![Go](https://img.shields.io/badge/Go-00ADD8.svg?style=flat&logo=Go&logoColor=white)
![React](https://img.shields.io/badge/React-61DAFB.svg?style=flat&logo=React&logoColor=black)
![Docker](https://img.shields.io/badge/Docker-2496ED.svg?style=flat&logo=Docker&logoColor=white)
![Python](https://img.shields.io/badge/Python-3776AB.svg?style=flat&logo=Python&logoColor=white)
![FastAPI](https://img.shields.io/badge/FastAPI-009688.svg?style=flat&logo=FastAPI&logoColor=white)
![SQLAlchemy](https://img.shields.io/badge/SQLAlchemy-D71F00.svg?style=flat&logo=SQLAlchemy&logoColor=white)
![Cypress](https://img.shields.io/badge/Cypress-69D3A7.svg?style=flat&logo=Cypress&logoColor=white)
![Pytest](https://img.shields.io/badge/Pytest-0A9EDC.svg?style=flat&logo=Pytest&logoColor=white)

</div>

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Live Deployment](#live-deployment)
- [Getting Started](#getting-started)
  - [Development](#development)
  - [Production](#production)
- [Testing](#testing)
- [Monitoring Stack](#monitoring-stack)

---

## Overview

**TRY-IT** is an interactive learning platform that allows users to create and participate in real-time quizzes. Designed for teachers, trainers, and event hosts, it brings engagement and gamification to learning and assessments.

---

## Features

- 🎨 **Dynamic Quiz Engine**: Create and host quizzes with real-time synchronization
- 📱 **Cross-Platform Access**: Join sessions from any device via browser
- 🧠 **Gamification**: Points system with live feedback (timers in development)
- 📊 **Advanced Analytics**: Track participant performance and quiz statistics
- 🏆 **Competitive Leaderboards**: Real-time ranking during sessions
- 🔔 **Telegram Alerts**: Instant notifications for critical system events
- 📈 **Business Metrics**: Track user registrations and quiz lifecycle events

---

## Live Deployment

A production instance is already available at [https://tryit.selnastol.ru](https://tryit.selnastol.ru). You may use this hosted version or deploy your own instance for customization/development.

---

## Getting Started

### Development

```sh
git clone https://github.com/IU-Capstone-Project-2025/Vanguard.git
cd Vanguard
cp .env.dev.example .env

# Without monitoring:
docker compose --env-file .env -f docker-compose.yaml -f docker-compose.dev.yaml up -d frontend

# With monitoring (Grafana/Prometheus):
docker compose --env-file .env -f docker-compose.yaml -f docker-compose.dev.yaml up -d --build
```

**Access:**
- Frontend: [http://localhost:3000](http://localhost:3000)
- Grafana: [http://localhost:3001](http://localhost:3001)

### Production

```sh
cp .env.prod.example .env
chmod +x setup.sh
./setup.sh  # SSL setup
docker compose --env-file .env -f docker-compose.yaml -f docker-compose.prod.yaml up -d --build
```

---

## Testing

**Test Coverage:**
- Python: 80% integration tests (pytest)
- Go: Unit tests for critical components + integration/e2e tests
- Frontend: Cypress post-build and session flow tests

All run tests through CI pipeline.

---

## Monitoring Stack

**Observability Suite:**
- 📊 **Grafana**: Unified dashboards for metrics/logs
- 🔍 **Prometheus**: Metrics collection with custom business KPIs
- 📜 **Loki + Promtail**: Centralized logging (backend services only)
- 🚨 **Alertmanager**: Telegram alerts for:
  - System resources (CPU/RAM thresholds)
  - Service throughput degradation
  - Container failures

**Planned Enhancements:**
- Frontend error logging
- Quiz process metrics (completion rates, latency)
- User behavior analytics

*Monitoring access is restricted to developers for security.*
