#!/bin/bash

echo "=== WhoKnows Metrics ==="
echo "Last 24 Hours:"
echo "Searches: $(tail -n 288 /var/log/metrics/daily_searches.log | awk '{sum += $3} END {print sum}')"
echo "Unique Users: $(tail -n 1 /var/log/metrics/daily_users.log | awk '{print $3}')"
echo "Current CPU Load: $(tail -n 1 /var/log/metrics/cpu_load.log | awk '{print $2,$3,$4}')"
echo "Disk Usage: $(tail -n 1 /var/log/metrics/disk_usage.log | awk '{print $2}')"

# Show latest daily summary if it exists
if [ -f /var/log/metrics/daily_summary.log ]; then
    echo -e "\nLatest Daily Summary:"
    tail -n 6 /var/log/metrics/daily_summary.log
fi
