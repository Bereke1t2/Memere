/**
 * k6 Local Breakpoint & Stress Test
 *
 * Targets local backend (http://host.docker.internal:8080)
 * Pushes concurrency from 50 -> 250 -> 500 -> 1,000 -> 1,500 Virtual Users
 * to determine the exact breaking threshold.
 */

import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend, Counter } from 'k6/metrics';

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

const errorRate    = new Rate('errors');
const successRate  = new Rate('success_rate');
const courseListP  = new Trend('course_list_ms');
const courseDetailP= new Trend('course_detail_ms');
const mockExamsP   = new Trend('mock_exams_ms');
const systemStatsP = new Trend('system_stats_ms');

export const options = {
  stages: [
    { duration: '15s', target: 50 },    // Warm-up to 50 VUs
    { duration: '30s', target: 250 },   // Step 1: 250 VUs
    { duration: '30s', target: 500 },   // Step 2: 500 VUs
    { duration: '30s', target: 1000 },  // Step 3: 1,000 VUs (SLO target)
    { duration: '30s', target: 1500 },  // Step 4: 1,500 VUs (Breaking point push)
    { duration: '15s', target: 0 },     // Cooldown
  ],
  thresholds: {
    // Report errors if failure rate > 10%
    errors: ['rate<0.10'],
  },
};

export function setup() {
  const res = http.get(`${BASE_URL}/api/v1/courses?limit=1`);
  let courseId = null;
  if (res.status === 200) {
    try {
      const data = res.json();
      const items = data.courses || data.data || [];
      if (items.length > 0) courseId = items[0].id;
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
    const res = http.get(`${BASE_URL}/api/v1/courses`, { headers });
    courseListP.add(res.timings.duration);
    handle(res);
  } else if (roll < 0.75 && data.courseId) {
    const res = http.get(`${BASE_URL}/api/v1/courses/${data.courseId}`, { headers });
    courseDetailP.add(res.timings.duration);
    handle(res);
  } else if (roll < 0.90) {
    const res = http.get(`${BASE_URL}/api/v1/mock-exams`, { headers });
    mockExamsP.add(res.timings.duration);
    handle(res);
  } else {
    const res = http.get(`${BASE_URL}/health/system`, { headers });
    systemStatsP.add(res.timings.duration);
    handle(res);
  }

  // Pacing: small sleep between requests
  sleep(0.1 + Math.random() * 0.2);
}

function handle(res) {
  const isOk = res.status === 200;
  successRate.add(isOk);
  errorRate.add(!isOk);
  check(res, { '200 OK': (r) => r.status === 200 });
}
