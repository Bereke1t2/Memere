

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

