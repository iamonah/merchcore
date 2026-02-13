# MerchCore SaaS E-Commerce Platform

A modern SaaS e-commerce platform that enables small businesses to effortlessly create and manage online stores. It provides secure payment options using both stablecoins and traditional **FIAT**, giving merchants flexibility in how they operate and customers freedom in how they pay.

Built with **Go** and following Clean Architecture principles, Storefront separates its core domain logic from infrastructure, ensuring scalability, maintainability, and testability.

***

## ✨ Key Features

- **🏪 Multi-Tenant Store Management**: Each business gets its own isolated store with customizable domains.
- **💸 Dual Payment Support**: Accepts both crypto (stablecoins) and fiat payments.
- **👤 User Authentication & Authorization**: Secure JWT-based access with role-based permissions.
- **🔄 Transactional Safety**: Consistent database operations in PostgreSQL.
- **⚙️ Modular Architecture**: Clearly defined layers (transport, service, repository) for easy maintenance.
- **🧰 Observability & Logging**: Structured logging using Zerolog.
- **🧪 Testable Design**: Repositories and services are fully mockable for isolated unit testing.

***

## 🏗️ Tech Stack

- **Backend**: Go
- **Database**: PostgreSQL
- **Cache / Jobs**: Redis
- **Auth**: JWT (JSON Web Tokens)
- **Architecture**: Clean Architecture + Domain-Driven Design principles
- **Containerization**: Docker & Docker Compose
- **Logging**: Zerolog

***
