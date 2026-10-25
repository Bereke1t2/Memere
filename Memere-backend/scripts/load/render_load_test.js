/**
 * k6 Load Test for Live Render Backend (https://memere.onrender.com)
 *
 * Tests real production endpoints:
 *   - 50% GET /api/v1/courses       (Course catalog - cache hit/miss)
 *   - 25% GET /api/v1/courses/:id   (Course details)
 *   - 15% GET /api/v1/mock-exams    (Exams catalog)
 *   - 10% GET /health/system        (Real-time server stats)
 *
 * Usage:
 *   # Quick run with 20 concurrent users for 30 seconds:
 *   docker run --rm -i grafana/k6 run --vus 20 --duration 30s - < Memere-backend/scripts/load/render_load_test.js
 *
 *   # Stepped ramp-up test (up to 50 concurrent users):
 *   docker run --rm -i grafana/k6 run - < Memere-backend/scripts/load/render_load_test.js
 */

import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend, Counter } from 'k6/metrics';

const BASE_URL = __ENV.BASE_URL || 'https://memere.onrender.com';

// Custom Metrics
const rateLimitCount = new Counter('rate_limited_429');
const successRate    = new Rate('success_rate');
const courseListP    = new Trend('course_list_ms');
const courseDetailP  = new Trend('course_detail_ms');
const mockExamsP     = new Trend('mock_exams_ms');
const systemStatsP   = new Trend('system_stats_ms');

export const options = {
  stages: [
    { duration: '15s', target: 10 },  // Phase 1: Warm-up (10 users)
    { duration: '30s', target: 25 },  // Phase 2: Moderate load (25 users)
    { duration: '30s', target: 50 },  // Phase 3: High load (50 users)
    { duration: '30s', target: 100 }, // Phase 4: Stress peak (100 users)
    { duration: '15s', target: 0 },   // Phase 5: Cooldown & Recovery (0 users)
  ],
  thresholds: {
    http_req_duration: ['p(95)<800'],
  },
};

export function setup() {
  // Pre-fetch a course ID from the live catalog for the detail test
  const res = http.get(`${BASE_URL}/api/v1/courses?limit=1`);
  let courseId = null;
  if (res.status === 200) {
    try {
      const data = res.json();
      const items = data.courses || data.data || [];
      if (items.length > 0) {
        courseId = items[0].id;
      }
    } catch (_) {}
  }
  return { courseId };
}

export default function (data) {
  const headers = {
    'Accept': 'application/json',
    'Accept-Encoding': 'gzip',
  };

  const roll = Math.random();

  if (roll < 0.50) {
    // 1. Course Catalog
    const res = http.get(`${BASE_URL}/api/v1/courses`, { headers });
    courseListP.add(res.timings.duration);
    handleResponse(res);

  } else if (roll < 0.75 && data.courseId) {
    // 2. Course Details
    const res = http.get(`${BASE_URL}/api/v1/courses/${data.courseId}`, { headers });
    courseDetailP.add(res.timings.duration);
    handleResponse(res);

  } else if (roll < 0.90) {
    // 3. Mock Exams Catalog
    const res = http.get(`${BASE_URL}/api/v1/mock-exams`, { headers });
    mockExamsP.add(res.timings.duration);
    handleResponse(res);

  } else {
    // 4. Live System Diagnostics Probe
    const res = http.get(`${BASE_URL}/health/system`, { headers });
    systemStatsP.add(res.timings.duration);
    handleResponse(res);
  }

  // Realistic human think time: 0.5s to 1s between clicks
  sleep(0.5 + Math.random() * 0.5);
}

function handleResponse(res) {
  if (res.status === 429) {
    rateLimitCount.add(1);
    successRate.add(false);
  } else {
    const isOk = res.status === 200;
    successRate.add(isOk);
    check(res, { 'status is 200': (r) => r.status === 200 });
  }
}
