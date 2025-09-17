# OpenEdu Technical Architecture & Stack

## 📌 Overview
This document describes the technical architecture and technology stack used in the **OpenEdu Platform**.  
It explains the rationale behind choosing each technology to ensure scalability, performance, and maintainability.

---

## 🖥️ Frontend
- **React + Next.js**
  - Provides SSR (Server-Side Rendering) and SSG (Static Site Generation) for SEO and fast performance.
  - Strong community support and ecosystem.
  - React: popular library, strong community, easy to extend with many UI libraries.
  - Next.js: supports SSR (Server-Side Rendering) and SSG (Static Site Generation) → improves SEO, performance, and       time-to-first-byte.
  - A very suitable combination for online learning systems that need to load quickly and have good SEO (course pages, lectures, public content).
- **TailwindCSS**
  - Utility-first CSS framework for faster, consistent, and responsive UI development.
- **SST (Serverless Stack)**
  - Simplifies serverless frontend deployment with AWS Lambda integration.

---

## ⚙️ Backend
- **Go + Gin**
  - High-performance, concurrent-friendly programming language.
  - Gin is a lightweight framework with middleware support (auth, logging, etc.).
- **Nginx (Load Balancer)**
  - Handles reverse proxy and SSL termination, distributes requests across backend services.
- **Redis (Cache Layer)**
  - Caches API responses and sessions to reduce load on databases.
- **RabbitMQ (Message Queue)**
  - Provides asynchronous communication for background jobs (emails, notifications, AI tasks).
- **Grafana + Prometheus**
  - Prometheus collects metrics; Grafana visualizes system health and performance.

---

## 🗄️ Database
- **PostgreSQL (RDS)**
  - Relational database for transactional data (users, courses, payments).
  - Supports JSONB for semi-structured data.
- **MongoDB**
  - NoSQL database for unstructured or flexible schema data (logs, learning content metadata).

---

## ☁️ Infrastructure & Services
- **Cloudflare CDN**
  - Accelerates content delivery globally and provides DDoS protection.
- **Blockchain Service**
  - Ensures transparency and verifiable learning certificates using smart contracts.

---

## ✅ Summary
- **Frontend**: SEO-friendly and fast rendering (React + Next.js + TailwindCSS).  
- **Backend**: High-performance Go microservices with caching, messaging, and monitoring.  
- **Database**: PostgreSQL for structured data + MongoDB for unstructured data.  
- **Infra**: CDN, load balancing, monitoring, and blockchain integration for scalability and trust.
