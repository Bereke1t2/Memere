#!/usr/bin/env bash
OUTPUT_FILE="${1:-/tmp/render_metrics_sample.csv}"
echo "timestamp,cpu_pct,alloc_mb,sys_mb,goroutines,db_acquired,db_idle,db_total" > "$OUTPUT_FILE"

# Sample every 2 seconds for 65 iterations (130 seconds total)
for i in {1..65}; do
  DATA=$(curl -s --connect-timeout 3 https://memere.onrender.com/health/system 2>/dev/null)
  if [ -n "$DATA" ]; then
    NOW=$(date +%T)
    CPU=$(echo "$DATA" | jq -r '.runtime.cpu_percent // 0')
    ALLOC=$(echo "$DATA" | jq -r '.runtime.alloc_mb // 0')
    SYS=$(echo "$DATA" | jq -r '.runtime.sys_mb // 0')
    GORO=$(echo "$DATA" | jq -r '.runtime.goroutines // 0')
    DB_ACQ=$(echo "$DATA" | jq -r '.db_pool.acquired_conns // 0')
    DB_IDLE=$(echo "$DATA" | jq -r '.db_pool.idle_conns // 0')
    DB_TOTAL=$(echo "$DATA" | jq -r '.db_pool.total_conns // 0')
    echo "$NOW,$CPU,$ALLOC,$SYS,$GORO,$DB_ACQ,$DB_IDLE,$DB_TOTAL" >> "$OUTPUT_FILE"
  fi
  sleep 2
done
