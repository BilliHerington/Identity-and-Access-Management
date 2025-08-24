#  IAM (Identity and Access Management)
## Description
  IAM is a Go application for managing user registration, authentication, and access control.
The system is built using a REST API architecture and is designed for secure and scalable user interaction.
##  🔧 Features
  - Standard registration and login via email and password

  - Authentication using JWT (with user version control)

  - Sending email messages (e.g. confirmation codes) via SMTP

  - Login via Google OAuth 2.0

  - Flexible role assignment and access level management

  - Redis for fast session and temporary data storage

  - Clear 3-layer architecture: handlers, services, repositories

  - Custom error handling and validations

  - Quick deployment and launch using Docker

##  🛠️ Technologies Used

  - Go + Gin

  - Redis

  - Gmail SMTP

  - Google OAuth 2.0

  - JWT

##  🚀 Quick Start
**1. Clone the repository**

        git clone https://github.com/BilliHerington/Identity-and-Access-Management
        cd Identity-and-Access-Management

**2. Create .env based on example.env**

      cp config/example.env config/.env

  🔧 Fill in the required environment variables, including:
  
      JWT_SECRET_KEY
    
      EMAIL_SENDER, EMAIL_PASSWORD (if not using test mode)
    
      REDIS_PASSWORD, ROOT_EMAIL, ROOT_PASSWORD, etc.
  
  📌 If you are running the project for the first time, make sure to set:
  
      INIT_REDIS_IN_FIRST_TIME=true
  
  🧪 To run without Google OAuth, set:
  
      USE_TEST_MODE_WITHOUT_GOOGLE=true

**3. Run the project with Docker**

        docker compose up --build

The service will be available at:
📍 http://localhost:8080

📂 Endpoints

All REST API endpoints can be found in the directory:
/pkg/routes
## ✉️ Contact

**Author: PepeTheProger**

**GitHub: [GitHub](https://github.com/PepeTheProger)**
