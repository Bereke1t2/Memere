/**
 * k6 3,000+ Virtual User Monster Stress Test
 *
 * Simulates 3,000 simultaneous clients hitting the Go backend.
 * Tests operating system file descriptor limits, socket queues, and thread pools.
 */

import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

const errorRate    = new Rate('errors');
const successRate  = new Rate('success_rate');
const courseListP  = new Trend('course_list_ms');
const systemStatsP = new Trend('system_stats_ms');

export const options = {
  stages: [
    { duration: '10s', target: 500 },   // Step 1: Rapid ramp to 500 VUs
    { duration: '20s', target: 1500 },  // Step 2: Ramp to 1,500 VUs
    { duration: '30s', target: 3000 },  // Step 3: PUSH TO 3,000 VUs!
    { duration: '15s', target: 0 },     // Step 4: Rapid cool down
  ],
  thresholds: {
    // Collect data without hard failure
  },
};

export default function () {
  const headers = {
    'Accept': 'application/json',
    'Accept-Encoding': 'gzip',
  };

  const roll = Math.random();

  if (roll < 0.60) {
    const res = http.get(`${BASE_URL}/api/v1/courses`, { headers });
    courseListP.add(res.timings.duration);
    handle(res);
  } else if (roll < 0.85) {
    const res = http.get(`${BASE_URL}/health/system`, { headers });
    systemStatsP.add(res.timings.duration);
    handle(res);
  } else {
    const res = http.get(`${BASE_URL}/health`, { headers });
    handle(res);
  }

  // Realistic human pacing
  sleep(0.1 + Math.random() * 0.2);
}

function handle(res) {
  const isOk = res.status === 200;
  successRate.add(isOk);
  errorRate.add(!isOk);
  check(res, { '200 OK': (r) => r.status === 200 });
}
