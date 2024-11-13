#!/bin/bash

# Create metrics directory if it doesn't exist
mkdir -p /var/log/metrics
chmod 755 /var/log/metrics

# Initialize log files if they don't exist
touch /var/log/metrics/daily_searches.log
touch /var/log/metrics/daily_users.log
touch /var/log/metrics/cpu_load.log
touch /var/log/metrics/disk_usage.log

# Set permissions
chmod 644 /var/log/metrics/*.log

# Timestamp for the log entries
TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')

# Collect CPU load using /proc/loadavg
CPU_LOAD=$(cat /proc/loadavg | awk '{print $1,$2,$3}')
echo "[$TIMESTAMP] $CPU_LOAD" >> /var/log/metrics/cpu_load.log

# Ensure access.log exists and is readable
if [ -f /var/log/nginx/access.log ]; then
    # Count unique IPs in nginx access log
    UNIQUE_USERS=$(awk '{print $1}' /var/log/nginx/access.log | sort -u | wc -l)
    echo "[$TIMESTAMP] $UNIQUE_USERS" >> /var/log/metrics/daily_users.log

    # Count searches
    DAILY_SEARCHES=$(grep "GET /?q=" /var/log/nginx/access.log | wc -l)
    echo "[$TIMESTAMP] $DAILY_SEARCHES" >> /var/log/metrics/daily_searches.log
else
    echo "[$TIMESTAMP] 0" >> /var/log/metrics/daily_users.log
    echo "[$TIMESTAMP] 0" >> /var/log/metrics/daily_searches.log
fi

# Disk usage
DISK_USAGE=$(df -h / | tail -n1 | awk '{print $5}')
echo "[$TIMESTAMP] $DISK_USAGE" >> /var/log/metrics/disk_usage.log

# Optional: Create a daily summary (runs at midnight)
if [ "$(date +%H:%M)" = "00:00" ]; then
    YESTERDAY=$(date -d "yesterday" '+%Y-%m-%d')
    
    # Calculate yesterday's total searches
    YESTERDAY_SEARCHES=$(grep "$YESTERDAY" /var/log/metrics/daily_searches.log | tail -n 1 | awk '{print $3}')
    
    # Calculate yesterday's unique users
    YESTERDAY_USERS=$(grep "$YESTERDAY" /var/log/metrics/daily_users.log | tail -n 1 | awk '{print $3}')
    
    # Create daily summary
    echo "=== Summary for $YESTERDAY ===" >> /var/log/metrics/daily_summary.log
    echo "Total Searches: $YESTERDAY_SEARCHES" >> /var/log/metrics/daily_summary.log
    echo "Unique Users: $YESTERDAY_USERS" >> /var/log/metrics/daily_summary.log
    echo "Final CPU Load: $CPU_LOAD" >> /var/log/metrics/daily_summary.log
    echo "Final Disk Usage: $DISK_USAGE" >> /var/log/metrics/daily_summary.log
    echo "===========================" >> /var/log/metrics/daily_summary.log
fi

# Optional: Rotate logs if they get too big (keep last 30 days)
find /var/log/metrics/ -name "*.log" -mtime +30 -delete