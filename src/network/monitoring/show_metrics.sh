#!/bin/bash

echo "=== WhoKnows Metrics ==="
echo "Last 24 Hours:"
echo "Searches: $(tail -n 288 /var/log/metrics/daily_searches.log | awk -F']' '{sum += $2} END {print sum}')"
echo "Unique Users: $(tail -n 1 /var/log/metrics/daily_users.log | awk -F']' '{print $2}')"
echo "Current CPU Load: $(tail -n 1 /var/log/metrics/cpu_load.log | awk -F']' '{print $2}')"
echo "Disk Usage: $(tail -n 1 /var/log/metrics/disk_usage.log | awk -F']' '{print $2}')"

# Show latest daily summary if it exists
if [ -f /var/log/metrics/daily_summary.log ]; then
    echo -e "\nLatest Daily Summary:"
    tail -n 6 /var/log/metrics/daily_summary.log
fi
