

**SOFTWARE ARCHITECTURE DOCUMENT**

**Grade 12 University Entrance Exam**

**Mobile Learning Platform**

| Platform | ExamPrep Mobile (Flutter \+ Go) |
| :---- | :---- |
| **Version** | 1.0.0 — Initial Architecture |
| **Date** | June 2025 |
| **Status** | Pre-Development — Architecture Phase |
| **Author** | Senior Software Architect |
| **Target Users** | Grade 12 Students, Teachers, Admins |

**CONFIDENTIAL — FOR INTERNAL USE ONLY**

# **Table of Contents**

**Phase 1  —**  Product Requirements Document (PRD)

**Phase 2  —**  Feature Planning — MVP & Roadmap

**Phase 3  —**  System Design & Architecture

**Phase 4  —**  Database Design

**Phase 5  —**  Backend Architecture (Go)

**Phase 6  —**  Flutter Architecture

**Phase 7  —**  Authentication & Security

**Phase 8  —**  Video Learning System

**Phase 9  —**  Quiz & Mock Exam Engine

**Phase 10  —**  Payment System

**Phase 11  —**  Notifications System

**Phase 12  —**  Scalability Planning

**Phase 13  —**  DevOps & CI/CD

**Phase 14  —**  Project Roadmap

**Phase 15  —**  Technical Learning Roadmap

| PHASE 1 Product Requirements Document |
| :---: |

## **1.1 Executive Summary**

ExamPrep is a mobile-first learning platform built specifically for Grade 12 students in Ethiopia preparing for university entrance examinations. The platform delivers structured video lessons, study notes, adaptive quizzes, and timed mock examinations — mirroring the Coursera/Udemy model but specialized for high-stakes national exam preparation.

## **1.2 Business Goals**

| Goal | Description | Success Metric |
| :---- | :---- | :---- |
| Market Entry | Launch a functional MVP targeting Grade 12 students in Ethiopia | 500 registered users in Month 1 |
| Revenue | Generate recurring revenue via course sales and subscriptions | ETB 50,000 MRR by Month 6 |
| Student Outcomes | Improve exam scores for enrolled students | Average 15% score improvement vs baseline |
| Content Scale | Build a library of exam-relevant video and text content | 50+ courses covering all subjects by Month 6 |
| Retention | Build habit-forming daily study patterns | DAU/MAU ratio \> 40% |

## **1.3 User Personas**

### **Persona 1 — The Focused Student**

| Attribute | Detail |
| :---- | :---- |
| Name | Hana Tesfaye |
| Age | 17 years old, Grade 12 |
| Goal | Score in top 5% on national university entrance exam |
| Pain Points | No access to quality tutors; textbooks are insufficient; no practice exams |
| Behavior | Studies 3–4 hours daily; uses phone for everything; motivated but anxious |
| Device | Android mid-range smartphone (Tecno/Samsung) |
| Connectivity | Intermittent mobile data; needs offline support |

### **Persona 2 — The Teacher/Content Creator**

| Attribute | Detail |
| :---- | :---- |
| Name | Mr. Abebe Girma |
| Age | 34, experienced Grade 12 teacher |
| Goal | Reach more students and supplement income |
| Pain Points | No platform to publish materials digitally; tracking student progress is manual |
| Behavior | Creates lessons using Word/PDF; wants simple upload tools |
| Tech Level | Moderate; comfortable with smartphones and basic web tools |

### **Persona 3 — The Platform Admin**

| Attribute | Detail |
| :---- | :---- |
| Name | Sara Bekele |
| Age | 28, operations manager |
| Goal | Manage platform operations, payments, and content quality |
| Responsibilities | User management, payment reconciliation, content approval, analytics |
| Tech Level | High; comfortable with dashboards and data |

## **1.4 Functional Requirements**

### **Student-Facing Features**

| ID | Feature | Description | Priority |
| :---- | :---- | :---- | :---- |
| FR-01 | Video Lessons | Stream and download HD video lessons per subject | P0 |
| FR-02 | Study Notes | Access PDF/document notes per lesson | P0 |
| FR-03 | Quizzes | Per-lesson quizzes with immediate feedback | P0 |
| FR-04 | Mock Exams | Full timed mock exams simulating national exam format | P0 |
| FR-05 | Progress Tracking | Dashboard showing completion %, scores, weak areas | P0 |
| FR-06 | Course Purchase | Buy individual courses or bundles | P0 |
| FR-07 | Subscriptions | Monthly/annual all-access plan | P1 |
| FR-08 | Offline Mode | Download videos and notes for offline use | P1 |
| FR-09 | Notifications | Push notifications for new content, reminders | P1 |
| FR-10 | Certificates | Downloadable completion certificates | P2 |
| FR-11 | Leaderboard | Class ranking based on quiz/exam scores | P2 |
| FR-12 | AI Tutor | Conversational AI to answer subject questions | P3 |

### **Teacher/Admin Features**

| ID | Feature | Description | Priority |
| :---- | :---- | :---- | :---- |
| FA-01 | Video Upload | Upload and manage video lessons | P0 |
| FA-02 | Notes Upload | Upload PDFs and documents | P0 |
| FA-03 | Quiz Builder | Create multiple-choice and short-answer quizzes | P0 |
| FA-04 | Exam Builder | Build timed mock exams with randomization | P0 |
| FA-05 | Course Manager | Create, edit, publish courses | P0 |
| FA-06 | Student Analytics | View per-student performance data | P1 |
| FA-07 | Payment Dashboard | Track revenue, payouts, subscriptions | P1 |
| FA-08 | Announcements | Send push/in-app notifications to students | P1 |
| FA-09 | Coupon Manager | Create discount coupons | P2 |

## **1.5 Non-Functional Requirements**

| Category | Requirement | Target |
| :---- | :---- | :---- |
| Performance | API response time (95th percentile) | \< 200ms |
| Performance | Video load time to first frame | \< 3 seconds on 3G |
| Availability | Platform uptime SLA | 99.9% (\< 8.7 hrs downtime/year) |
| Scalability | Concurrent active users at launch | 1,000 users |
| Scalability | Target concurrent users at 12 months | 10,000 users |
| Security | Authentication standard | JWT with refresh tokens, HTTPS only |
| Security | Payment data handling | PCI-DSS compliant via Stripe/Chapa |
| Offline | Offline content availability | Videos \+ notes downloadable, quizzes cached |
| Compliance | Data residency | Primary data in Africa (AWS af-south-1 or GCP africa-south1) |
| Accessibility | Language support | Amharic and English |

## **1.6 Monetization Strategy**

| Model | Description | Price Range (ETB) | Notes |
| :---- | :---- | :---- | :---- |
| Per-Course Purchase | One-time payment for lifetime course access | 150–500 per course | Best for subject-specific prep |
| Monthly Subscription | All-access subscription plan | 200–400/month | Best for full exam prep |
| Annual Subscription | Discounted all-year plan | 1,500–3,000/year | \~40% discount vs monthly |
| Bundle Packs | Group of related subject courses | 400–800 per bundle | e.g., "Science Bundle" |
| Premium Mock Exams | Standalone mock exam packages | 50–150 per pack | With detailed analytics |

| Revenue Projection Month 3:   500 users × ETB 250 avg spend \= ETB 125,000 Month 6:   2,000 users × ETB 300 avg spend \= ETB 600,000 Month 12:  10,000 users × ETB 350 avg spend \= ETB 3,500,000 Teacher revenue share: 70% to teacher, 30% platform fee |
| :---- |

| PHASE 2 Feature Planning — MVP & Roadmap |
| :---: |

## **2.1 MVP Features (Month 1–3)**

The MVP must validate core learning value: students can enroll, watch videos, read notes, take quizzes, and purchase courses.

| \# | Feature | Rationale |
| :---- | :---- | :---- |
| 1 | User registration & login (JWT) | Gate all content behind auth |
| 2 | Course listing & detail pages | Browse available courses |
| 3 | Video streaming (HLS adaptive) | Core learning delivery |
| 4 | PDF/note viewer | Study material access |
| 5 | Per-lesson multiple-choice quiz | Test comprehension |
| 6 | Mock exam engine (timed) | Simulate real exam conditions |
| 7 | Basic progress tracking | Show completion per course |
| 8 | Course purchase (Chapa/Stripe) | Monetization from day one |
| 9 | Teacher content upload (video, PDF) | Enable content creation |
| 10 | Admin dashboard (basic) | Manage users and content |
| 11 | Push notifications (FCM) | Engagement and reminders |
| 12 | Offline download (video \+ notes) | Critical for low-connectivity users |

## **2.2 Phase 2 Features (Month 4–6)**

| \# | Feature | Why After MVP |
| :---- | :---- | :---- |
| 1 | Leaderboard and ranking system | Needs baseline data to be meaningful |
| 2 | Monthly subscription plans | Requires payment infrastructure maturity |
| 3 | Advanced student analytics dashboard | Needs accumulated data |
| 4 | Coupon and discount system | Marketing activation after launch |
| 5 | Email notifications | Lower priority than push |
| 6 | Completion certificates | Nice-to-have after core loop works |
| 7 | Course review/rating system | Needs user base first |
| 8 | Teacher earnings dashboard | Needs payment history data |

## **2.3 Future Features (Month 7–12+)**

| \# | Feature | Notes |
| :---- | :---- | :---- |
| 1 | AI-powered tutoring chatbot | LLM integration (OpenAI/Claude API) |
| 2 | Adaptive learning paths | Personalized curriculum based on performance |
| 3 | Live video sessions (Zoom/Agora SDK) | Real-time teacher-student interaction |
| 4 | Peer study groups | Social/collaborative learning |
| 5 | Multi-language content (Amharic-first) | Localization for broader reach |
| 6 | Parent monitoring portal | Track child's study habits and scores |
| 7 | School/institution licensing | B2B revenue stream |
| 8 | Predictive score analytics (ML) | Predict exam outcomes, identify risk |

| PHASE 3 System Design & Architecture |
| :---: |

## **3.1 High-Level Architecture Overview**

The platform follows a microservices-oriented, API-first architecture. The Flutter mobile app communicates exclusively through a single API Gateway. Each bounded domain (courses, quizzes, payments, etc.) is an independently deployable service. All services share a PostgreSQL cluster (with per-service schemas) and Redis for caching and session management.

**Architecture Diagram (Mermaid — High-Level System)**

| graph TB   subgraph Mobile\["Flutter Mobile App"\]     APP\[Flutter App\<br/\>Riverpod \+ GoRouter\]   end   subgraph Edge\["Edge Layer"\]     CDN\[CloudFront / GCP CDN\]     GW\[API Gateway\<br/\>Nginx \+ Rate Limiter\]   end   subgraph Services\["Backend Microservices (Go)"\]     AUTH\[Auth Service\<br/\>JWT \+ Refresh Tokens\]     COURSE\[Course Service\<br/\>Videos, Notes, Sections\]     QUIZ\[Quiz Service\<br/\>Questions, Attempts\]     EXAM\[Exam Service\<br/\>Mock Exams, Timer\]     PAY\[Payment Service\<br/\>Stripe / Chapa / Telebirr\]     NOTIF\[Notification Service\<br/\>FCM \+ Email \+ In-App\]     PROG\[Progress Service\<br/\>Tracking \+ Analytics\]   end   subgraph Data\["Data Layer"\]     PG\[(PostgreSQL\<br/\>Primary DB)\]     REDIS\[(Redis\<br/\>Cache \+ Sessions)\]     S3\[(S3 / GCS\<br/\>Object Storage)\]     MQ\[RabbitMQ / SQS\<br/\>Message Queue\]   end   APP \--\> CDN   APP \--\> GW   GW \--\> AUTH   GW \--\> COURSE   GW \--\> QUIZ   GW \--\> EXAM   GW \--\> PAY   GW \--\> NOTIF   GW \--\> PROG   AUTH \--\> PG   AUTH \--\> REDIS   COURSE \--\> PG   COURSE \--\> S3   COURSE \--\> MQ   QUIZ \--\> PG   QUIZ \--\> REDIS   EXAM \--\> PG   EXAM \--\> REDIS   PAY \--\> PG   NOTIF \--\> MQ   PROG \--\> PG   PROG \--\> REDIS   CDN \--\> S3 |
| :---- |

## **3.2 Service Responsibilities**

| Service | Port | Responsibilities |
| :---- | :---- | :---- |
| API Gateway | 8080 | Rate limiting, routing, auth validation, SSL termination |
| Auth Service | 8081 | Registration, login, JWT issuance, refresh tokens, password reset |
| Course Service | 8082 | Course CRUD, section management, lesson management, video/note metadata |
| Quiz Service | 8083 | Quiz creation, question bank, attempt recording, auto-grading |
| Exam Service | 8084 | Mock exam engine, timer management, exam sessions, score calculation |
| Payment Service | 8085 | Purchase flow, subscription management, webhook handling, receipts |
| Notification Service | 8086 | Push (FCM), email (SendGrid), in-app notifications, scheduling |
| Progress Service | 8087 | Completion tracking, streak calculation, analytics aggregation |

## **3.3 Infrastructure Components**

| Component | Technology | Purpose |
| :---- | :---- | :---- |
| Object Storage | AWS S3 / GCS | Store video files, PDFs, images, certificates |
| CDN | CloudFront / GCP CDN | Deliver video and static content with low latency globally |
| Video Processing | AWS MediaConvert / FFmpeg | Transcode uploaded videos to HLS adaptive bitrate streams |
| Database | PostgreSQL 15 (RDS/Cloud SQL) | Primary relational data store |
| Cache | Redis 7 (ElastiCache) | Session tokens, quiz states, rate limit counters |
| Message Queue | SQS / RabbitMQ | Async tasks: video processing, email sending, notifications |
| Container Registry | ECR / GCR | Docker image storage |
| Orchestration | Kubernetes (EKS/GKE) | Container deployment and scaling |
| Load Balancer | AWS ALB / GCP LB | Traffic distribution across service instances |
| Monitoring | Prometheus \+ Grafana | Metrics collection and visualization |
| Logging | ELK Stack / CloudWatch | Centralized log aggregation and search |

| PHASE 4 Database Design |
| :---: |

## **4.1 Database Design Principles**

* One PostgreSQL cluster, schemas per domain (auth, courses, payments, notifications)

* UUID primary keys for all entities (prevents enumeration attacks)

* Soft deletes (deleted\_at timestamp) — never hard delete user data

* created\_at / updated\_at on every table

* JSONB for flexible metadata; typed columns for queryable fields

## **4.2 Core Tables**

### **4.2.1 Users Table**

| Column | Type | Constraints | Description |
| :---- | :---- | :---- | :---- |
| id | UUID | PK, default gen\_random\_uuid() | Unique user identifier |
| email | VARCHAR(255) | UNIQUE, NOT NULL | Login email address |
| phone | VARCHAR(20) | UNIQUE, nullable | Phone number (Telebirr payments) |
| password\_hash | VARCHAR(255) | NOT NULL | bcrypt hashed password |
| role | ENUM | NOT NULL (student/teacher/admin) | User role for RBAC |
| first\_name | VARCHAR(100) | NOT NULL | First name |
| last\_name | VARCHAR(100) | NOT NULL | Last name |
| avatar\_url | TEXT | nullable | Profile picture URL (S3) |
| is\_active | BOOLEAN | default true | Account active state |
| is\_email\_verified | BOOLEAN | default false | Email verification status |
| email\_verification\_token | VARCHAR(255) | nullable | Token for email verification |
| password\_reset\_token | VARCHAR(255) | nullable | Token for password reset |
| password\_reset\_expires\_at | TIMESTAMPTZ | nullable | Token expiry time |
| last\_login\_at | TIMESTAMPTZ | nullable | Last successful login |
| created\_at | TIMESTAMPTZ | NOT NULL, default now() | Record creation time |
| updated\_at | TIMESTAMPTZ | NOT NULL, default now() | Last update time |
| deleted\_at | TIMESTAMPTZ | nullable | Soft delete timestamp |

Indexes: email (UNIQUE), phone (UNIQUE), role, deleted\_at, created\_at

### **4.2.2 Courses Table**

| Column | Type | Constraints | Description |
| :---- | :---- | :---- | :---- |
| id | UUID | PK | Course identifier |
| teacher\_id | UUID | FK → users.id, NOT NULL | Course creator/teacher |
| title | VARCHAR(200) | NOT NULL | Course display title |
| slug | VARCHAR(200) | UNIQUE, NOT NULL | URL-friendly identifier |
| description | TEXT | NOT NULL | Full course description |
| short\_description | VARCHAR(500) | nullable | Card preview text |
| subject | VARCHAR(100) | NOT NULL | e.g., Mathematics, Physics |
| grade | INTEGER | NOT NULL (e.g., 12\) | Target grade level |
| thumbnail\_url | TEXT | nullable | Cover image URL |
| price | DECIMAL(10,2) | NOT NULL default 0 | Price in ETB |
| currency | VARCHAR(3) | NOT NULL default ETB | Currency code |
| is\_free | BOOLEAN | default false | Free course flag |
| is\_published | BOOLEAN | default false | Visibility flag |
| language | VARCHAR(10) | default en | en or am (Amharic) |
| level | ENUM | beginner/intermediate/advanced | Difficulty level |
| total\_duration\_seconds | INTEGER | default 0 | Sum of all video durations |
| total\_lessons | INTEGER | default 0 | Denormalized lesson count |
| rating\_avg | DECIMAL(3,2) | default 0 | Average student rating |
| enrollment\_count | INTEGER | default 0 | Total enrolled students |
| metadata | JSONB | nullable | Flexible extra data |
| created\_at | TIMESTAMPTZ | NOT NULL, default now() |  |
| updated\_at | TIMESTAMPTZ | NOT NULL, default now() |  |
| deleted\_at | TIMESTAMPTZ | nullable | Soft delete |

Indexes: teacher\_id, subject, grade, is\_published, slug (UNIQUE), price, enrollment\_count

### **4.2.3 Course Sections & Lessons**

| Column | Type | Description |
| :---- | :---- | :---- |
| id | UUID | Primary key |
| course\_id | UUID FK | Parent course |
| title | VARCHAR(200) | Section title |
| description | TEXT | Section overview |
| order\_index | INTEGER | Display order within course |
| is\_published | BOOLEAN | Visibility toggle |
| created\_at / updated\_at | TIMESTAMPTZ | Audit timestamps |

| Column | Type | Description |
| :---- | :---- | :---- |
| id | UUID | Lesson primary key |
| section\_id | UUID FK | Parent section |
| course\_id | UUID FK | Denormalized for fast queries |
| title | VARCHAR(200) | Lesson title |
| type | ENUM | video / note / quiz / mixed |
| order\_index | INTEGER | Display order |
| is\_free\_preview | BOOLEAN | Allow non-enrolled preview |
| duration\_seconds | INTEGER | Video duration (if applicable) |
| is\_published | BOOLEAN | Visibility |
| created\_at / updated\_at | TIMESTAMPTZ | Audit timestamps |

### **4.2.4 Videos Table**

| Column | Type | Description |
| :---- | :---- | :---- |
| id | UUID | Primary key |
| lesson\_id | UUID FK UNIQUE | One video per lesson |
| original\_file\_key | TEXT | S3 key for raw upload |
| hls\_master\_key | TEXT | S3 key for HLS manifest (.m3u8) |
| processing\_status | ENUM | pending / processing / ready / failed |
| duration\_seconds | INTEGER | Video duration |
| resolution\_480p\_key | TEXT | HLS 480p stream key |
| resolution\_720p\_key | TEXT | HLS 720p stream key |
| resolution\_1080p\_key | TEXT | HLS 1080p stream key |
| thumbnail\_key | TEXT | Video thumbnail on S3 |
| file\_size\_bytes | BIGINT | Original file size |
| created\_at / updated\_at | TIMESTAMPTZ | Audit timestamps |

### **4.2.5 Quizzes, Questions & Answers**

| Column | Type | Description |
| :---- | :---- | :---- |
| quizzes.id | UUID | Quiz primary key |
| quizzes.lesson\_id | UUID FK | Linked lesson (nullable for standalone) |
| quizzes.course\_id | UUID FK | Linked course |
| quizzes.title | VARCHAR(200) | Quiz title |
| quizzes.time\_limit\_seconds | INTEGER | Nullable — timed or untimed |
| quizzes.pass\_percentage | DECIMAL(5,2) | Minimum pass score |
| quizzes.randomize\_questions | BOOLEAN | Shuffle question order |
| questions.id | UUID | Question PK |
| questions.quiz\_id | UUID FK | Parent quiz |
| questions.text | TEXT | Question content |
| questions.type | ENUM | multiple\_choice / true\_false / short\_answer |
| questions.points | INTEGER | Point value |
| questions.explanation | TEXT | Explanation shown after submission |
| questions.order\_index | INTEGER | Display order |
| answers.id | UUID | Answer option PK |
| answers.question\_id | UUID FK | Parent question |
| answers.text | TEXT | Answer option text |
| answers.is\_correct | BOOLEAN | Correct answer flag (server-only) |
| answers.order\_index | INTEGER | Display order |

### **4.2.6 Exams & Attempts**

| Column | Type | Description |
| :---- | :---- | :---- |
| exams.id | UUID | Exam primary key |
| exams.course\_id | UUID FK | Parent course (nullable for standalone) |
| exams.title | VARCHAR(200) | Exam title |
| exams.subject | VARCHAR(100) | Subject area |
| exams.grade | INTEGER | Target grade |
| exams.duration\_minutes | INTEGER | Total allowed time |
| exams.total\_marks | INTEGER | Maximum possible score |
| exams.pass\_marks | INTEGER | Minimum passing marks |
| exams.instructions | TEXT | Exam instructions text |
| exams.is\_published | BOOLEAN | Visibility |
| exam\_attempts.id | UUID | Attempt PK |
| exam\_attempts.exam\_id | UUID FK | Which exam |
| exam\_attempts.student\_id | UUID FK | Who took it |
| exam\_attempts.started\_at | TIMESTAMPTZ | Timer start |
| exam\_attempts.submitted\_at | TIMESTAMPTZ | Submission time (nullable \= in progress) |
| exam\_attempts.score | DECIMAL(10,2) | Final score |
| exam\_attempts.percentage | DECIMAL(5,2) | Score as percentage |
| exam\_attempts.answers\_snapshot | JSONB | Snapshot of student answers at submission |
| exam\_attempts.status | ENUM | in\_progress / submitted / graded / expired |

### **4.2.7 Enrollments & Payments**

| Column | Type | Description |
| :---- | :---- | :---- |
| enrollments.id | UUID | Enrollment PK |
| enrollments.student\_id | UUID FK | Enrolled student |
| enrollments.course\_id | UUID FK | Enrolled course |
| enrollments.enrolled\_at | TIMESTAMPTZ | Enrollment timestamp |
| enrollments.expires\_at | TIMESTAMPTZ | Nullable (subscription-based expiry) |
| enrollments.source | ENUM | purchase / subscription / free / coupon |
| payments.id | UUID | Payment PK |
| payments.student\_id | UUID FK | Payer |
| payments.course\_id | UUID FK | Nullable (null \= subscription) |
| payments.amount | DECIMAL(10,2) | Amount charged |
| payments.currency | VARCHAR(3) | ETB or USD |
| payments.provider | ENUM | stripe / chapa / telebirr |
| payments.provider\_transaction\_id | VARCHAR(255) | External payment reference |
| payments.status | ENUM | pending / completed / failed / refunded |
| payments.coupon\_id | UUID FK | Applied coupon (nullable) |
| payments.metadata | JSONB | Provider response payload |
| payments.paid\_at | TIMESTAMPTZ | Payment completion time |

### **4.2.8 Progress Tracking**

| Column | Type | Description |
| :---- | :---- | :---- |
| id | UUID | Progress record PK |
| student\_id | UUID FK | Student reference |
| lesson\_id | UUID FK | Lesson reference |
| course\_id | UUID FK | Denormalized for fast course-level queries |
| is\_completed | BOOLEAN | Lesson completion flag |
| completed\_at | TIMESTAMPTZ | When lesson was completed |
| video\_progress\_seconds | INTEGER | How far into the video |
| last\_accessed\_at | TIMESTAMPTZ | Last session timestamp |

**ER Diagram — Core Entities (Mermaid)**

| erDiagram   USERS ||--o{ ENROLLMENTS : "enrolls"   USERS ||--o{ QUIZ\_ATTEMPTS : "attempts"   USERS ||--o{ EXAM\_ATTEMPTS : "takes"   USERS ||--o{ PAYMENTS : "pays"   USERS ||--o{ PROGRESS : "tracks"   USERS ||--o{ COURSES : "teaches"   COURSES ||--o{ COURSE\_SECTIONS : "contains"   COURSES ||--o{ ENROLLMENTS : "has"   COURSES ||--o{ EXAMS : "includes"   COURSE\_SECTIONS ||--o{ LESSONS : "has"   LESSONS ||--o| VIDEOS : "may have"   LESSONS ||--o| NOTES : "may have"   LESSONS ||--o| QUIZZES : "may have"   LESSONS ||--o{ PROGRESS : "tracked in"   QUIZZES ||--o{ QUESTIONS : "contains"   QUESTIONS ||--o{ ANSWERS : "has options"   QUIZZES ||--o{ QUIZ\_ATTEMPTS : "has"   EXAMS ||--o{ EXAM\_QUESTIONS : "contains"   EXAMS ||--o{ EXAM\_ATTEMPTS : "has"   PAYMENTS ||--o| SUBSCRIPTIONS : "may create"   PAYMENTS }o--o| COUPONS : "may use" |
| :---- |

| PHASE 5 Backend Architecture — Go Clean Architecture |
| :---: |

## **5.1 Architectural Principles**

The backend follows Uncle Bob's Clean Architecture with strict dependency inversion. Outer layers depend on inner layers — never the reverse. The Domain layer has zero external dependencies.

**Clean Architecture Dependency Flow**

| graph LR   subgraph Core\["Core (No External Deps)"\]     E\[Domain\<br/\>Entities\]     UC\[Use Cases\<br/\>Business Logic\]   end   subgraph Adapters\["Adapters"\]     R\[Repositories\<br/\>DB Interfaces\]     H\[HTTP Handlers\<br/\>Delivery\]   end   subgraph Infra\["Infrastructure"\]     PG\[PostgreSQL\<br/\>Implementation\]     RD\[Redis\<br/\>Implementation\]     S3I\[S3\<br/\>Implementation\]   end   UC \--\> E   H \--\> UC   R \--\> UC   PG \--\> R   RD \--\> R   S3I \--\> R |
| :---- |

## **5.2 Folder Structure**

examprep-backend/

├── cmd/

│   ├── api/            \# main.go — entry point, dependency wiring

│   └── migrate/        \# Database migration runner

├── internal/

│   ├── domain/         \# ← INNERMOST LAYER (no external imports)

│   │   ├── entity/     \# User, Course, Quiz, Exam structs

│   │   ├── repository/ \# Repository INTERFACES (contracts)

│   │   └── service/    \# Domain service interfaces

│   ├── usecase/        \# ← BUSINESS LOGIC

│   │   ├── auth/       \# RegisterUser, LoginUser, RefreshToken

│   │   ├── course/     \# CreateCourse, EnrollStudent, GetCourseById

│   │   ├── quiz/       \# SubmitQuizAttempt, GradeQuiz, GetLeaderboard

│   │   ├── exam/       \# StartExam, SubmitExam, GetExamResults

│   │   ├── payment/    \# InitiatePayment, HandleWebhook, IssueRefund

│   │   └── progress/   \# MarkLessonComplete, GetCourseProgress

│   ├── repository/     \# ← REPOSITORY IMPLEMENTATIONS

│   │   ├── postgres/   \# SQL queries per entity

│   │   └── redis/      \# Cache operations

│   ├── delivery/       \# ← HTTP HANDLERS (outermost layer)

│   │   ├── http/       \# Gin/Echo route handlers

│   │   │   ├── auth\_handler.go

│   │   │   ├── course\_handler.go

│   │   │   ├── quiz\_handler.go

│   │   │   ├── exam\_handler.go

│   │   │   └── payment\_handler.go

│   │   └── middleware/ \# Auth, CORS, logging, rate-limit

│   └── infrastructure/ \# ← EXTERNAL CONCERNS

│       ├── database/   \# PostgreSQL connection, migrations

│       ├── cache/      \# Redis client wrapper

│       ├── storage/    \# S3/GCS client wrapper

│       ├── messaging/  \# SQS/RabbitMQ client

│       ├── payment/    \# Stripe, Chapa, Telebirr clients

│       └── notification/ \# FCM, SendGrid wrappers

├── pkg/

│   ├── jwt/            \# JWT creation and validation helpers

│   ├── password/       \# bcrypt helpers

│   ├── validator/      \# Input validation utilities

│   ├── pagination/     \# Cursor-based pagination helpers

│   └── errors/         \# Custom error types and codes

├── config/             \# Config structs, env loading

├── migrations/         \# SQL migration files (golang-migrate)

├── api/                \# OpenAPI/Swagger spec (generated)

├── scripts/

│   ├── seed.go         \# Test data seeder

│   └── generate.go     \# Code generation scripts

├── Dockerfile

├── docker-compose.yml

└── Makefile

## **5.3 Layer Responsibilities**

| Layer | Location | Responsibility | May Import |
| :---- | :---- | :---- | :---- |
| Domain Entities | internal/domain/entity | Pure Go structs — User, Course, Quiz. No business logic. No DB tags. | Standard library only |
| Repository Interfaces | internal/domain/repository | Go interfaces defining DB operations. e.g., UserRepository.FindByEmail() | Domain entities only |
| Use Cases | internal/usecase | Business logic. Orchestrates repos and services. Stateless functions. | Domain layer only |
| Repository Impl. | internal/repository | Concrete PostgreSQL/Redis implementations of domain interfaces | Domain \+ infrastructure |
| HTTP Handlers | internal/delivery/http | Parse HTTP requests, validate input, call use cases, return JSON | Use cases \+ pkg utilities |
| Middleware | internal/delivery/middleware | JWT auth, CORS, logging, rate limiting, request ID injection | pkg/jwt, pkg/errors |
| Infrastructure | internal/infrastructure | Thin wrappers around external clients: DB, S3, Redis, FCM | External SDKs only |

## **5.4 API Design Conventions**

| Convention | Rule | Example |
| :---- | :---- | :---- |
| Base URL versioning | All routes prefixed with /api/v1 | GET /api/v1/courses |
| Resource naming | Plural nouns, kebab-case | GET /api/v1/mock-exams |
| HTTP methods | REST-standard GET/POST/PUT/PATCH/DELETE | POST /api/v1/quiz-attempts |
| Pagination | Cursor-based with limit/after params | GET /api/v1/courses?limit=20\&after=\<cursor\> |
| Error format | { "code": "RESOURCE\_NOT\_FOUND", "message": "...", "details": {} } | Standard across all services |
| Auth header | Bearer token in Authorization header | Authorization: Bearer \<jwt\> |
| File uploads | Multipart form-data for uploads; pre-signed URLs for large files | POST /api/v1/videos/upload-url |

| PHASE 6 Flutter Architecture |
| :---: |

## **6.1 Clean Architecture in Flutter**

The Flutter app mirrors the backend's Clean Architecture. Features are self-contained modules. Riverpod manages all state. GoRouter handles declarative navigation with redirect guards.

## **6.2 Folder Structure**

lib/

