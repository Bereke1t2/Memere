#!/usr/bin/env bash
# smoke_phase6.sh — Phase 6 Skill 5 deploy validation runbook.
#
# Usage (against a running cluster / kind):
#   NAMESPACE=staging BASE_URL=https://api.memere.app ./scripts/smoke_phase6.sh
#
# Requires: kubectl, curl, jq
set -euo pipefail

NAMESPACE="${NAMESPACE:-staging}"
BASE_URL="${BASE_URL:-http://localhost:8080}"

pass() { echo "  [PASS] $1"; }
fail() { echo "  [FAIL] $1"; exit 1; }
section() { echo; echo "=== $1 ==="; }

# ---------------------------------------------------------------------------
section "Manifest validation (dry-run)"
# ---------------------------------------------------------------------------
kubectl apply -k k8s/ --dry-run=server --namespace="$NAMESPACE" \
  && pass "kustomize manifests valid" \
  || fail "manifest validation failed"

# ---------------------------------------------------------------------------
section "Migration Job"
# ---------------------------------------------------------------------------
kubectl apply -f k8s/jobs/migrate-job.yaml -n "$NAMESPACE"
kubectl wait --for=condition=complete job/migrate -n "$NAMESPACE" --timeout=120s \
  && pass "migrate Job completed" \
  || fail "migrate Job did not complete in 120s"

# ---------------------------------------------------------------------------
section "Deployment rollout"
# ---------------------------------------------------------------------------
kubectl rollout status deployment/api    -n "$NAMESPACE" --timeout=120s \
  && pass "api Deployment rolled out" \
  || fail "api rollout timed out"
kubectl rollout status deployment/worker -n "$NAMESPACE" --timeout=120s \
  && pass "worker Deployment rolled out" \
  || fail "worker rollout timed out"

# ---------------------------------------------------------------------------
section "Pod readiness"
# ---------------------------------------------------------------------------
API_READY=$(kubectl get deployment api -n "$NAMESPACE" \
  -o jsonpath='{.status.readyReplicas}' 2>/dev/null || echo 0)
[ "$API_READY" -ge 2 ] \
  && pass "api: $API_READY ready replicas (≥ 2)" \
  || fail "api: only $API_READY ready replicas"

# ---------------------------------------------------------------------------
section "HPA exists"
# ---------------------------------------------------------------------------
kubectl get hpa api -n "$NAMESPACE" &>/dev/null \
  && pass "api HPA present" \
  || fail "api HPA not found"

kubectl get hpa worker -n "$NAMESPACE" &>/dev/null \
  && pass "worker HPA present" \
  || fail "worker HPA not found"

# ---------------------------------------------------------------------------
section "PodDisruptionBudget"
# ---------------------------------------------------------------------------
kubectl get pdb api-pdb -n "$NAMESPACE" &>/dev/null \
  && pass "api PDB present" \
  || fail "api PDB not found"

# ---------------------------------------------------------------------------
section "Ingress + TLS"
# ---------------------------------------------------------------------------
INGRESS_IP=$(kubectl get ingress api -n "$NAMESPACE" \
  -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null || echo "")
[ -n "$INGRESS_IP" ] \
  && pass "Ingress has IP: $INGRESS_IP" \
  || echo "  [WARN] Ingress IP not yet assigned (may be pending LB)"

# ---------------------------------------------------------------------------
section "Health probes"
# ---------------------------------------------------------------------------
HEALTHZ=$(curl -sf --max-time 5 "${BASE_URL}/healthz" | jq -r '.status' 2>/dev/null || echo "")
[ "$HEALTHZ" = "ok" ] \
  && pass "/healthz → ok" \
  || fail "/healthz returned: '${HEALTHZ}'"

READYZ=$(curl -sf --max-time 5 "${BASE_URL}/readyz" | jq -r '.status' 2>/dev/null || echo "")
[ "$READYZ" = "ready" ] \
  && pass "/readyz → ready" \
  || fail "/readyz returned: '${READYZ}'"

# ---------------------------------------------------------------------------
section "Version endpoint"
# ---------------------------------------------------------------------------
VERSION=$(curl -sf --max-time 5 "${BASE_URL}/version" | jq -r '.version' 2>/dev/null || echo "")
[ -n "$VERSION" ] \
  && pass "/version → $VERSION" \
  || fail "/version unreachable"

# ---------------------------------------------------------------------------
section "Prometheus scrape (internal)"
# ---------------------------------------------------------------------------
# Port-forward metrics port and scrape a known metric.
kubectl port-forward svc/api 9091:9090 -n "$NAMESPACE" &>/dev/null &
PF_PID=$!
sleep 2
METRIC=$(curl -sf --max-time 5 "http://localhost:9091/metrics" \
  | grep -c "^http_requests_total" || echo 0)
kill "$PF_PID" 2>/dev/null || true
[ "$METRIC" -ge 1 ] \
  && pass "Prometheus: http_requests_total metric present" \
  || echo "  [WARN] http_requests_total not found — app may not have received any traffic yet"

# ---------------------------------------------------------------------------
echo
echo "Phase 6 Skill 5 smoke checks complete."
