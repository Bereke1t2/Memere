

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

├── main.dart              \# App entry point, ProviderScope

├── app.dart               \# MaterialApp, GoRouter setup

├── core/

│   ├── constants/         \# API endpoints, colors, text styles

│   ├── errors/            \# Failure classes (ServerFailure, CacheFailure)

│   ├── network/           \# Dio client, interceptors, connectivity check

│   ├── storage/           \# Hive/SharedPrefs wrappers

│   ├── theme/             \# AppTheme, dark/light modes

│   ├── router/            \# GoRouter config, route guards

│   └── di/                \# Dependency injection providers

├── features/

│   ├── auth/

│   │   ├── data/

│   │   │   ├── datasources/   \# AuthRemoteDataSource (Dio calls)

│   │   │   │                  \# AuthLocalDataSource (Hive tokens)

│   │   │   ├── models/        \# UserModel (from JSON, extends UserEntity)

│   │   │   └── repositories/  \# AuthRepositoryImpl

│   │   ├── domain/

│   │   │   ├── entities/      \# UserEntity (pure Dart, no JSON)

│   │   │   ├── repositories/  \# AuthRepository (abstract)

│   │   │   └── usecases/      \# LoginUseCase, RegisterUseCase

│   │   └── presentation/

│   │       ├── providers/     \# authStateProvider, loginProvider

│   │       ├── screens/       \# LoginScreen, RegisterScreen

│   │       └── widgets/       \# LoginForm, SocialLoginButtons

│   ├── courses/

│   │   └── ... (same structure)

│   ├── video\_player/

│   │   └── ... (custom HLS player wrapper)

│   ├── quiz/

│   ├── exam/

│   ├── payment/

│   └── notifications/

└── shared/

    ├── widgets/           \# AppButton, AppTextField, LoadingOverlay

    ├── extensions/        \# String, DateTime, int extensions

    └── utils/             \# Formatters, validators, helpers

## **6.3 Riverpod State Management**

| Provider Type | Use Case | Example |
| :---- | :---- | :---- |
| Provider | Simple computed values, no state | currentUserProvider (reads auth state) |
| StateProvider | Simple mutable state | selectedSubjectProvider (String) |
| StateNotifierProvider | Complex state with methods | CourseListNotifier (filtering, pagination) |
| FutureProvider | One-time async fetch | courseDetailProvider(courseId) |
| StreamProvider | Real-time data | notificationsStreamProvider |
| AsyncNotifierProvider | Async state with methods | ExamAttemptNotifier (start, submit) |

## **6.4 Offline-First Strategy**

| Data Type | Storage | Strategy |
| :---- | :---- | :---- |
| Auth tokens | Flutter Secure Storage | Encrypted at rest; refresh on 401 |
| Course metadata | Hive (NoSQL local DB) | Cache API response; TTL \= 1 hour; stale-while-revalidate |
| Video files | Device filesystem (path\_provider) | Explicit user download; encrypted with device key |
| PDF/Notes | Device filesystem | Download on demand; indexed in Hive |
| Quiz state | Hive | Persist in-progress attempt; sync on reconnect |
| User preferences | SharedPreferences | Immediate write; no sync needed |

| Offline Sync Strategy 1\. On app launch: check connectivity via connectivity\_plus 2\. Offline: serve all data from Hive cache; queue writes to a SyncQueue 3\. On reconnect: flush SyncQueue to backend (idempotent endpoints) 4\. Conflict resolution: server-wins for all data except quiz answers (client-wins) 5\. Background sync: WorkManager (Android) / BGTaskScheduler (iOS) |
| :---- |

| PHASE 7 Authentication & Security |
| :---: |

## **7.1 JWT Authentication Flow**

**JWT Authentication Sequence (Mermaid)**

| sequenceDiagram   participant App as Flutter App   participant GW as API Gateway   participant Auth as Auth Service   participant DB as PostgreSQL   participant Cache as Redis   App-\>\>GW: POST /api/v1/auth/login {email, password}   GW-\>\>Auth: Forward request   Auth-\>\>DB: SELECT user WHERE email=?   DB--\>\>Auth: User record   Auth-\>\>Auth: bcrypt.Compare(password, hash)   Auth-\>\>Auth: Generate access\_token (15min TTL)   Auth-\>\>Auth: Generate refresh\_token (30 days TTL)   Auth-\>\>Cache: SET session:{user\_id} \= refresh\_token\_hash TTL=30d   Auth-\>\>DB: INSERT refresh\_tokens (hash, device\_info, expires\_at)   Auth--\>\>App: { access\_token, refresh\_token, user }   Note over App,Cache: Access token expired — silent refresh   App-\>\>GW: POST /api/v1/auth/refresh {refresh\_token}   GW-\>\>Auth: Forward   Auth-\>\>Cache: GET session:{user\_id} — verify not revoked   Auth-\>\>Auth: Issue new access\_token   Auth--\>\>App: { access\_token } |
| :---- |

## **7.2 Role-Based Access Control (RBAC)**

| Role | Permissions |
| :---- | :---- |
| student | Read courses (enrolled), submit quizzes/exams, view own progress, purchase |
| teacher | All student permissions \+ create/edit own courses, view student analytics for own courses |
| admin | All permissions \+ manage all users, all courses, payment reconciliation, platform config |

## **7.3 Attack Vectors & Mitigations**

| Attack Vector | Risk | Mitigation |
| :---- | :---- | :---- |
| Brute Force Login | Account takeover | Rate limit: 5 attempts/15min per IP. Account lockout after 10 failures. |
| JWT Token Theft | Unauthorized access | Short-lived access tokens (15 min). HttpOnly cookie option. Token rotation. |
| SQL Injection | Data breach | Parameterized queries via sqlx/pgx. No raw string interpolation. |
| XSS (Web Admin) | Script injection | CSP headers. React auto-escaping. Sanitize all user HTML input. |
| IDOR (Insecure Direct Object Ref) | Unauthorized data access | Always filter DB queries by authenticated user\_id. Never trust client-side IDs. |
| Payment Replay | Double charge | Idempotency keys on all payment requests. Webhook deduplication. |
| Video Hotlinking | Content theft | Pre-signed CDN URLs with 2-hour expiry per authenticated user. |
| Exam Answer Leaking | Academic dishonesty | Never send correct answers in API response. Grade server-side only. |
| DDoS | Platform downtime | CloudFlare / AWS Shield. Rate limiting at gateway. Auto-scaling. |
| Credential Stuffing | Mass account takeover | Email breach detection. Suspicious login alerting. 2FA option. |

| PHASE 8 Video Learning System |
| :---: |

## **8.1 Video Upload & Processing Pipeline**

**Video Processing Workflow (Mermaid)**

| graph TD   A\[Teacher uploads video\<br/\>via Flutter App\] \--\> B\[Request pre-signed S3 URL\<br/\>from Course Service\]   B \--\> C\[Flutter uploads directly\<br/\>to S3 via pre-signed URL\]   C \--\> D\[S3 Event Notification\<br/\>triggers SQS message\]   D \--\> E\[Video Processor Worker\<br/\>picks up SQS message\]   E \--\> F{Processing}   F \--\> G\[Transcode to HLS\<br/\>480p / 720p / 1080p\<br/\>using FFmpeg/MediaConvert\]   F \--\> H\[Generate thumbnail\<br/\>at 5-second mark\]   G \--\> I\[Upload HLS segments\<br/\>and .m3u8 manifest to S3\]   I \--\> J\[Update video record\<br/\>status \= ready\<br/\>in PostgreSQL\]   J \--\> K\[Invalidate CDN cache\<br/\>for new content\]   J \--\> L\[Push notification to teacher\<br/\>Video processing complete\]   H \--\> I |
| :---- |

## **8.2 Adaptive Bitrate Streaming (HLS)**

| Quality | Resolution | Bitrate | Use Case |
| :---- | :---- | :---- | :---- |
| Low | 360p | 400 kbps | 2G / very slow mobile data |
| Medium | 480p | 800 kbps | 3G mobile connection |
| High | 720p | 1.5 Mbps | 4G / good WiFi |
| HD | 1080p | 3 Mbps | Strong WiFi / broadband |

The HLS master manifest (.m3u8) lists all available quality levels. The video player (flutter\_hls\_parser \+ video\_player) automatically selects the best quality based on current bandwidth measurement, switching seamlessly without buffering.

## **8.3 Offline Video Download Strategy**

| Step | Action | Implementation |
| :---- | :---- | :---- |
| 1\. Request | Student taps "Download" on lesson | Flutter app calls GET /api/v1/videos/{id}/download-url |
| 2\. Signed URL | Backend generates time-limited CDN URL | S3 pre-signed URL valid for 2 hours (single-use) |
| 3\. Download | App downloads HLS segments | Download all .ts segments \+ manifest |
| 4\. Storage | Save to app documents directory | path\_provider getApplicationDocumentsDirectory() |
| 5\. Encryption | Encrypt downloaded files | AES-256 with device-bound key from Secure Storage |
| 6\. Offline Play | Play from local files | Custom FileDataSource for video\_player |
| 7\. Expiry | Downloaded content valid 30 days | Check timestamp on each play; prompt to re-download |

| PHASE 9 Quiz & Mock Exam Engine |
| :---: |

## **9.1 Quiz Engine Design**

| Aspect | Design Decision | Rationale |
| :---- | :---- | :---- |
| Answer security | Correct answers NEVER sent to client | Prevents cheating via API inspection |
| Question randomization | Shuffle on attempt creation, snapshot stored in Redis | Consistent order for duration of attempt; different per attempt |
| Grading | Server-side only, triggered on submission | Client cannot manipulate score |
| Partial attempts | State saved to Redis every 30 seconds | No data loss on app crash or disconnect |
| Time enforcement | Timer managed server-side (started\_at in DB) | Client timer is display-only; server auto-submits at expiry |
| Attempt limits | Configurable per quiz (1–unlimited) | Teacher sets max attempts |

## **9.2 Mock Exam Flow**

**Exam Attempt State Machine (Mermaid)**

| stateDiagram-v2   \[\*\] \--\> NOT\_STARTED   NOT\_STARTED \--\> IN\_PROGRESS: Student clicks Start   IN\_PROGRESS \--\> IN\_PROGRESS: Save answer (auto-save every 30s)   IN\_PROGRESS \--\> SUBMITTED: Student submits manually   IN\_PROGRESS \--\> EXPIRED: Server timer fires at duration\_minutes   SUBMITTED \--\> GRADED: Auto-grading completes   EXPIRED \--\> GRADED: Auto-graded on expiry   GRADED \--\> \[\*\] |
| :---- |

## **9.3 Scoring & Analytics**

| Metric | Calculation | Storage |
| :---- | :---- | :---- |
| Raw Score | Sum of points for correct answers | exam\_attempts.score |
| Percentage | (raw\_score / total\_marks) \* 100 | exam\_attempts.percentage |
| Subject Breakdown | Points scored per subject tag on questions | JSONB in exam\_attempts |
| Weak Areas | Questions answered wrong grouped by topic | Aggregated in progress service |
| Percentile Rank | Student score vs all attempts for same exam | Calculated on-demand via Redis sorted set |
| Trend Analysis | Score over consecutive attempts for same subject | Time-series query on exam\_attempts |

| PHASE 10 Payment System |
| :---: |

## **10.1 Payment Providers**

| Provider | Use Case | Currency | Notes |
| :---- | :---- | :---- | :---- |
| Chapa | Primary ETB payments for local users | ETB | Mobile money \+ bank transfer; dominant in Ethiopia |
| Telebirr | Ethio Telecom mobile wallet | ETB | Largest mobile money in Ethiopia; \~30M users |
| Stripe | International cards; diaspora users | USD/EUR | For future international expansion |

## **10.2 Payment Flow**

**Payment Sequence Diagram (Mermaid)**

| sequenceDiagram   participant App   participant PayService as Payment Service   participant Provider as Chapa/Stripe   participant DB   App-\>\>PayService: POST /payments/initiate {course\_id, provider}   PayService-\>\>DB: INSERT payment {status=pending, idempotency\_key}   PayService-\>\>Provider: Create payment intent / checkout   Provider--\>\>PayService: payment\_url or client\_secret   PayService--\>\>App: { payment\_url, payment\_id }   App-\>\>Provider: User completes payment in WebView   Provider-\>\>PayService: POST /payments/webhook {event, transaction\_id}   PayService-\>\>PayService: Verify webhook signature   PayService-\>\>DB: UPDATE payment {status=completed, paid\_at}   PayService-\>\>DB: INSERT enrollment {student\_id, course\_id}   PayService-\>\>NotifService: Trigger enrollment confirmation notification   App-\>\>PayService: GET /payments/{id}/status (polling)   PayService--\>\>App: { status: "completed" }   App-\>\>App: Navigate to course player |
| :---- |

## **10.3 Coupon System**

| Column | Type | Description |
| :---- | :---- | :---- |
| code | VARCHAR(50) UNIQUE | Human-readable coupon code (e.g., EXAM50) |
| discount\_type | ENUM | percentage or fixed\_amount |
| discount\_value | DECIMAL(10,2) | e.g., 50 \= 50% off or ETB 50 off |
| max\_uses | INTEGER | Total redemption limit (null \= unlimited) |
| used\_count | INTEGER | Current redemption count |
| expires\_at | TIMESTAMPTZ | Coupon expiry date |
| applicable\_to | ENUM | all / specific\_courses / subscription\_only |
| course\_ids | UUID\[\] | Applicable course IDs (if specific) |

| PHASE 11 Notifications System |
| :---: |

## **11.1 Notification Channels**

| Channel | Provider | Use Cases | Priority |
| :---- | :---- | :---- | :---- |
| Push (Mobile) | Firebase Cloud Messaging (FCM) | New lesson, exam reminder, score ready, announcement | P0 |
| In-App | PostgreSQL \+ WebSocket/polling | Unread count badge, notification center in app | P0 |
| Email | SendGrid | Registration confirmation, purchase receipt, weekly progress report | P1 |
| SMS | Africa's Talking | Critical: password reset, payment confirmation (if no email) | P2 |

## **11.2 Notification Event Catalog**

| Event | Trigger | Channels |
| :---- | :---- | :---- |
| welcome\_email | User registration completed | Email |
| email\_verification | Registration or resend request | Email |
| exam\_reminder | 24h before a scheduled exam | Push \+ Email |
| lesson\_published | Teacher publishes new lesson in enrolled course | Push \+ In-App |
| exam\_graded | Exam attempt grading complete | Push \+ In-App |
| purchase\_confirmed | Payment completed successfully | Email \+ In-App |
| streak\_warning | User hasn't studied in 2 days | Push |
| certificate\_ready | Course 100% completed | Push \+ Email |
| announcement | Admin broadcasts platform news | Push \+ In-App |

| PHASE 12 Scalability Planning |
| :---: |

## **12.1 Scaling Tiers**

| User Scale | Architecture | Infrastructure | Monthly AWS Cost (Est.) |
| :---- | :---- | :---- | :---- |
| 1K users | Monolith (single Go binary) | EC2 t3.medium (2vCPU/4GB), RDS db.t3.micro, S3, CloudFront | \~$80–120/month |
| 10K users | Modular monolith \+ Redis \+ CDN | EC2 t3.large x2 (LB), RDS db.t3.medium, ElastiCache t3.micro, S3 | \~$250–350/month |
| 100K users | Microservices \+ Kubernetes | EKS 3-node cluster, RDS db.r5.large (Multi-AZ), ElastiCache r6g.large, ALB | \~$1,200–1,800/month |
| 1M users | Full microservices \+ read replicas | EKS multi-region, RDS r5.2xlarge \+ 2 read replicas, ElastiCache cluster, CDN edge caching | \~$8,000–15,000/month |

## **12.2 Scaling Decision Triggers**

| Component | Introduce When | Signal |
| :---- | :---- | :---- |
| Redis Cache | 500+ concurrent users | DB CPU \> 60% or p95 query time \> 100ms |
| CDN | First video upload | Always — even at 100 users, CDN pays for itself |
| Load Balancer | 2,000+ concurrent users or deployment downtime is unacceptable | Single instance bottleneck |
| Message Queue | Any async task (video processing, emails) | Always — decouple processing from API response |
| Read Replicas | 100K+ users OR analytics queries impacting write performance | DB replication lag or read query p99 \> 500ms |
| Microservices split | Individual services need to scale independently | One service consuming disproportionate resources |
| Multi-region | 1M+ users or regulatory data requirements | Latency \> 300ms for significant user segment |

| PHASE 13 DevOps & CI/CD |
| :---: |

## **13.1 CI/CD Pipeline**

**GitHub Actions CI/CD Pipeline (Mermaid)**

| graph LR   A\[Git Push / PR\] \--\> B\[GitHub Actions Trigger\]   B \--\> C{Branch?}   C \--\>|feature branch| D\[CI Only\]   C \--\>|main branch| E\[CI \+ CD to Staging\]   C \--\>|release tag| F\[CI \+ CD to Production\]   D \--\> D1\[go test ./...\]   D \--\> D2\[flutter test\]   D \--\> D3\[golangci-lint\]   D \--\> D4\[docker build \--no-push\]   E \--\> E1\[All CI steps\]   E1 \--\> E2\[docker build \+ push to ECR\]   E2 \--\> E3\[kubectl apply to staging\]   E3 \--\> E4\[Run integration tests\]   E4 \--\> E5\[Notify Slack\]   F \--\> F1\[All CI steps\]   F1 \--\> F2\[docker build \+ push prod tag\]   F2 \--\> F3\[kubectl rollout production\]   F3 \--\> F4\[Smoke tests\]   F4 \--\> F5\[Notify team \+ Datadog deploy marker\] |
| :---- |

## **13.2 Kubernetes Deployment Structure**

k8s/

├── namespaces/

│   ├── production.yaml

│   └── staging.yaml

├── deployments/

│   ├── auth-service.yaml     \# replicas: 2 (prod), 1 (staging)

│   ├── course-service.yaml

│   ├── quiz-service.yaml

│   └── ... (one per service)

├── services/

│   └── ... (ClusterIP per service, LoadBalancer for gateway)

├── ingress/

│   └── api-ingress.yaml     \# nginx-ingress with TLS cert-manager

├── configmaps/

│   └── app-config.yaml      \# Non-secret config

├── secrets/

│   └── app-secrets.yaml     \# DB password, JWT secret (use Secrets Manager)

└── hpa/

    └── autoscaling.yaml     \# HorizontalPodAutoscaler per service

## **13.3 Monitoring & Alerting**

| Tool | Purpose | Key Metrics |
| :---- | :---- | :---- |
| Prometheus | Metrics scraping and storage | HTTP request rate, error rate, response time (p50/p95/p99), DB connections |
| Grafana | Metrics visualization | Service dashboards, SLO tracking, alert panels |
| ELK Stack | Centralized logging | Error aggregation, request tracing, security events |
| PagerDuty/OpsGenie | On-call alerting | 5xx error rate \> 1%, p99 latency \> 500ms, DB replication lag \> 30s |
| Sentry | Error tracking (Flutter \+ Go) | Crash-free session rate, unhandled exceptions |

| PHASE 14 Project Roadmap |
| :---: |

## **14.1 3-Month Roadmap — MVP Launch**

| Month | Focus | Milestones | Team Size |
| :---- | :---- | :---- | :---- |
| Month 1 | Foundation & Architecture | Go backend scaffolding, DB schema, Auth service, Flutter project setup, Design system, CI/CD pipeline | 2 backend, 1 Flutter, 1 DevOps |
| Month 2 | Core Features | Course \+ Video service, HLS streaming, Flutter video player, Basic quiz engine, Course purchase (Chapa) | 2 backend, 2 Flutter, 1 DevOps |
| Month 3 | MVP Completion | Mock exam engine, Offline download, Push notifications, Admin dashboard, UAT with 20 beta students | 2 backend, 2 Flutter, 1 QA, 1 PM |

## **14.2 6-Month Roadmap — Growth Phase**

| Month | Focus | Milestones |
| :---- | :---- | :---- |
| Month 4 | Analytics & Engagement | Student analytics dashboard, Leaderboard, Study streaks, Course ratings |
| Month 5 | Monetization Expansion | Subscription plans, Bundle packs, Coupon system, Teacher earnings dashboard |
| Month 6 | Scale & Quality | Performance optimization, Redis caching, iOS app release, 1,000+ active users target |

## **14.3 12-Month Roadmap — Scale Phase**

| Quarter | Focus | Milestones |
| :---- | :---- | :---- |
| Q3 (M7–9) | AI & Advanced Features | AI tutoring chatbot, Adaptive quizzes, Amharic language support, Live sessions (MVP) |
| Q4 (M10–12) | Scale & Expansion | Microservices migration, 10,000+ users, B2B school licensing, Multi-subject expansion |

## **14.4 Key Risks & Mitigations**

| Risk | Probability | Impact | Mitigation |
| :---- | :---- | :---- | :---- |
| Poor mobile internet in Ethiopia | High | High | Offline-first architecture; HLS adaptive bitrate; minimal API payloads |
| Payment provider API changes | Medium | High | Payment abstraction layer; multiple provider support |
| Low content quality | Medium | High | Teacher onboarding program; content review process; student ratings |
| Flutter release delays | Medium | Medium | Focus MVP on Android first; iOS in Month 6 |
| Scaling costs exceed projections | Low | High | Cost monitoring from day 1; AWS budget alerts; optimize queries early |

| PHASE 15 Technical Learning Roadmap |
| :---: |

## **15.1 Flutter**

| Skill | Why Needed | Priority | Difficulty | Learn Order |
| :---- | :---- | :---- | :---- | :---- |
| Dart language fundamentals | Core language; null safety, async/await, streams | P0 | Easy | 1st |
| Flutter widgets & layouts | Build all UI screens | P0 | Easy-Medium | 2nd |
| Riverpod state management | All state management in the app | P0 | Medium | 3rd |
| GoRouter navigation | Screen routing, auth guards, deep links | P0 | Medium | 4th |
| Dio HTTP client | API communication, interceptors, JWT refresh | P0 | Easy | 4th |
| Freezed code generation | Immutable models, union types for state | P1 | Medium | 5th |
| Hive local database | Offline caching, local data persistence | P0 | Easy | 5th |
| flutter\_secure\_storage | Secure token storage | P0 | Easy | 5th |
| video\_player \+ HLS | Core feature — video playback | P0 | Medium-Hard | 6th |
| Clean Architecture in Flutter | Structure for maintainability | P0 | Hard | Ongoing |
| Flutter testing (unit \+ widget \+ integration) | Quality assurance | P1 | Medium | Ongoing |

## **15.2 Go Backend**

| Skill | Why Needed | Priority | Difficulty | Learn Order |
| :---- | :---- | :---- | :---- | :---- |
| Go fundamentals (goroutines, channels, interfaces) | Core language; concurrency model | P0 | Medium | 1st |
| Gin or Echo HTTP framework | API routing and middleware | P0 | Easy | 2nd |
| sqlx / pgx PostgreSQL driver | Database queries | P0 | Easy-Medium | 2nd |
| Clean Architecture patterns in Go | Maintainable structure | P0 | Hard | 3rd |
| JWT implementation (golang-jwt) | Authentication | P0 | Easy | 3rd |
| Docker for Go | Containerization | P0 | Easy | 4th |
| golang-migrate | Database migrations | P0 | Easy | 4th |
| Redis client (go-redis) | Caching and sessions | P1 | Easy | 5th |
| Testing (testing package \+ testify) | Reliability | P1 | Medium | Ongoing |
| OpenAPI/Swagger documentation | API docs for Flutter team | P1 | Easy | Ongoing |

## **15.3 PostgreSQL**

| Skill | Why Needed | Priority | Difficulty |
| :---- | :---- | :---- | :---- |
| SQL fundamentals (SELECT, JOIN, GROUP BY, subqueries) | All data access | P0 | Easy |
| Schema design & normalization | Database structure decisions | P0 | Medium |
| Indexes (B-tree, partial, composite) | Query performance | P0 | Medium |
| JSONB column usage | Flexible metadata storage | P1 | Easy |
| Transactions and ACID guarantees | Payment and enrollment integrity | P0 | Medium |
| Query optimization (EXPLAIN ANALYZE) | Performance tuning | P1 | Hard |
| golang-migrate migrations | Schema version control | P0 | Easy |

## **15.4 System Design**

| Skill | Why Needed | Priority | Difficulty |
| :---- | :---- | :---- | :---- |
| REST API design principles | Backend API contract design | P0 | Easy |
| Microservices patterns (strangler fig, saga) | Multi-service architecture | P1 | Hard |
| Message queues (SQS/RabbitMQ) | Async video processing, notifications | P1 | Medium |
| CDN and object storage patterns | Video and file delivery | P0 | Medium |
| HLS adaptive bitrate streaming | Video lesson delivery | P0 | Medium |
| Redis patterns (cache-aside, session store) | Performance and sessions | P1 | Medium |
| Database indexing strategies | Query performance | P0 | Medium |

## **15.5 Cloud (AWS or GCP)**

| Skill | Service | Priority | Difficulty |
| :---- | :---- | :---- | :---- |
| Object Storage | S3 / GCS: store videos, PDFs, images | P0 | Easy |
| CDN | CloudFront / GCP CDN: fast content delivery | P0 | Easy |
| Compute | EC2 / Cloud Run: host Go services | P0 | Easy |
| Managed Database | RDS PostgreSQL / Cloud SQL | P0 | Easy |
| Container Orchestration | EKS / GKE: Kubernetes deployment | P1 | Hard |
| Message Queue | SQS / Cloud Pub/Sub | P1 | Medium |
| Video Transcoding | MediaConvert / Transcoder API | P1 | Medium |
| IAM & Security | Roles, policies, secrets management | P0 | Medium |

## **15.6 DevOps**

| Skill | Why Needed | Priority | Difficulty |
| :---- | :---- | :---- | :---- |
| Docker | Containerize all services | P0 | Easy |
| Docker Compose | Local development environment | P0 | Easy |
| GitHub Actions | CI/CD automation | P0 | Medium |
| Kubernetes basics | Production deployment (K8s) | P1 | Hard |
| Helm charts | K8s application packaging | P1 | Hard |
| Prometheus \+ Grafana | Monitoring and alerting | P1 | Medium |
| Nginx configuration | API gateway and reverse proxy | P0 | Medium |

## **15.7 Security**

| Skill | Why Needed | Priority | Difficulty |
| :---- | :---- | :---- | :---- |
| JWT & OAuth 2.0 concepts | Auth system design | P0 | Medium |
| bcrypt password hashing | Secure credential storage | P0 | Easy |
| HTTPS & TLS configuration | Encrypt all traffic | P0 | Easy |
| OWASP Top 10 vulnerabilities | Know what to defend against | P0 | Medium |
| SQL injection prevention | Data breach prevention | P0 | Easy |
| Rate limiting strategies | DDoS and brute-force mitigation | P0 | Easy |
| Payment security (PCI-DSS basics) | Handle payment data correctly | P1 | Medium |

## **15.8 Recommended Learning Order (Full Stack Path)**

| Week | Focus Area | Goal |
| :---- | :---- | :---- |
| W1–2 | Go fundamentals \+ Dart/Flutter basics | Build simple REST API \+ Flutter screen |
| W3–4 | PostgreSQL \+ Gin/Echo \+ sqlx | CRUD API with database |
| W5–6 | Flutter Clean Architecture \+ Riverpod | Feature module with state management |
| W7–8 | Authentication (JWT \+ refresh tokens) | Full auth flow end-to-end |
| W9–10 | Docker \+ GitHub Actions CI/CD | Automated build and deploy pipeline |
| W11–12 | AWS S3 \+ CloudFront \+ video upload | File upload and CDN delivery |
| W13–16 | Full feature build (quiz \+ payment) | Mini product with core features |
| W17–20 | Kubernetes \+ monitoring | Production-ready deployment |
| W21+ | Advanced (AI, scaling, performance) | Iterative improvement |

# **Appendix: Quick Reference**

| Key Technology Decisions Summary Mobile: Flutter (cross-platform Android \+ iOS, single codebase) Backend: Go (high performance, low memory, fast compile) Database: PostgreSQL (ACID, JSONB, excellent scaling path) Cache: Redis (sessions, exam state, leaderboard sorted sets) Video: HLS adaptive bitrate via S3 \+ CloudFront Payments: Chapa (primary) \+ Telebirr \+ Stripe (international) Push: Firebase Cloud Messaging (FCM) Containers: Docker \+ Kubernetes (EKS or GKE) CI/CD: GitHub Actions → ECR → kubectl rolling deploy Monitoring: Prometheus \+ Grafana \+ Sentry |
| :---- |

| Architecture Non-Negotiables 1\. Correct answers NEVER sent to client — always grade server-side 2\. Exam timer MUST be enforced server-side (client timer is display only) 3\. Pre-signed CDN URLs ONLY for video — no public S3 access 4\. All payments must use idempotency keys — no double-charge risk 5\. Soft deletes on all user-facing data — no hard DELETEs 6\. HTTPS-only — HTTP redirect to HTTPS at gateway 7\. Never log raw passwords, tokens, or payment card data 8\. All DB queries must filter by authenticated user\_id (prevent IDOR) |
| :---- |

