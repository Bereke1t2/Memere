/**
 * k6 load test — Memere API hot paths (Phase 6 Skill 3 §1.5 SLO validation).
 *
 * Target SLO: p95 < 200 ms on read paths at 1 000 concurrent virtual users.
 *
 * Usage:
 *   k6 run --vus 100 --duration 30s scripts/load/hot_paths.js
 *   k6 run --vus 1000 --duration 60s scripts/load/hot_paths.js   # SLO run
 *
 * Environment variables:
 *   BASE_URL   — API base URL (default: http://localhost:8080)
 *   EMAIL      — seeded student email  (default: student@memere.et)
 *   PASSWORD   — seeded student password (default: Password1!)
 *
 * Scenarios (weight-split):
 *   40 % — GET /api/v1/courses          (published catalog — cache-aside)
 *   25 % — GET /api/v1/courses/:id      (course detail)
 *   15 % — GET /api/v1/me/dashboard     (progress dashboard — auth required)
 *   10 % — GET /api/v1/mock-exams       (exam catalog)
 *   10 % — GET /health                  (health probe baseline)
 */

import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

// ---- configuration ----------------------------------------------------------

const BASE_URL  = __ENV.BASE_URL  || 'http://localhost:8080';
const EMAIL     = __ENV.EMAIL     || 'student@memere.et';
const PASSWORD  = __ENV.PASSWORD  || 'Password1!';

// Custom metrics
const errorRate    = new Rate('errors');
const courseListP  = new Trend('course_list_duration', true);
const courseGetP   = new Trend('course_detail_duration', true);
const dashboardP   = new Trend('dashboard_duration', true);

// ---- SLO thresholds ---------------------------------------------------------

export const options = {
  vus: 100,
  duration: '30s',
  thresholds: {
    // Overall p95 < 200 ms (§1.5 NFR)
    'http_req_duration{scenario:default}': ['p(95)<200'],
    // Individual path SLOs
    course_list_duration:    ['p(95)<200'],
    course_detail_duration:  ['p(95)<200'],
    dashboard_duration:      ['p(95)<300'],  // auth + DB read — slightly wider
    // Error rate < 1 %
    errors:                  ['rate<0.01'],
  },
};

// ---- setup: obtain an access token ------------------------------------------

export function setup() {
  const res = http.post(`${BASE_URL}/api/v1/auth/login`, JSON.stringify({
    email: EMAIL,
    password: PASSWORD,
  }), { headers: { 'Content-Type': 'application/json' } });

  check(res, { 'login 200': r => r.status === 200 });
  const body = res.json();
  return { token: body.access_token, courseID: null };
}

// ---- default function -------------------------------------------------------

export default function (data) {
  const headers = {
    'Content-Type': 'application/json',
    'Accept-Encoding': 'gzip',
  };
  if (data.token) {
    headers['Authorization'] = `Bearer ${data.token}`;
  }

  const roll = Math.random();

  if (roll < 0.40) {
    // Course catalog
    const r = http.get(`${BASE_URL}/api/v1/courses`, { headers });
    courseListP.add(r.timings.duration);
    check(r, { 'course list 200': res => res.status === 200 });
    errorRate.add(r.status !== 200);

  } else if (roll < 0.65) {
    // Course detail — uses first course from a quick list fetch or a known ID
    const listRes = http.get(`${BASE_URL}/api/v1/courses?limit=1`, { headers });
    if (listRes.status === 200) {
      const items = listRes.json('courses') || [];
      if (items.length > 0) {
        const r = http.get(`${BASE_URL}/api/v1/courses/${items[0].id}`, { headers });
        courseGetP.add(r.timings.duration);
        check(r, { 'course detail 200': res => res.status === 200 });
        errorRate.add(r.status !== 200);
      }
    }

  } else if (roll < 0.80) {
    // Dashboard (requires auth)
    const r = http.get(`${BASE_URL}/api/v1/me/dashboard`, { headers });
    dashboardP.add(r.timings.duration);
    check(r, { 'dashboard 200 or 401': res => res.status === 200 || res.status === 401 });
    errorRate.add(r.status >= 500);

  } else if (roll < 0.90) {
    // Mock exam catalog
    const r = http.get(`${BASE_URL}/api/v1/mock-exams`, { headers });
    check(r, { 'mock exams 200': res => res.status === 200 });
    errorRate.add(r.status !== 200);

  } else {
    // Health probe — baseline / connection reuse check
    const r = http.get(`${BASE_URL}/health`);
    check(r, { 'health 200': res => res.status === 200 });
    errorRate.add(r.status !== 200);
  }

  sleep(0.1); // 100 ms think time between requests
}
