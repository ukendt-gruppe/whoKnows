#!/bin/bash

# Create metrics directory if it doesn't exist
mkdir -p /var/log/metrics
chmod 755 /var/log/metrics

# Initialize log files if they don't exist
touch /var/log/metrics/daily_searches.log
touch /var/log/metrics/daily_users.log
touch /var/log/metrics/cpu_load.log
touch /var/log/metrics/disk_usage.log
touch /var/log/metrics/weekly_cost.log

# Set permissions
chmod 644 /var/log/metrics/*.log

# Timestamp for the log entries
TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')

# Collect CPU load using /proc/loadavg
CPU_LOAD=$(cat /proc/loadavg | awk '{print $1,$2,$3}')
echo "[$TIMESTAMP] $CPU_LOAD" >> /var/log/metrics/cpu_load.log

# Get today's date for filtering
TODAY=$(date '+%Y-%m-%d')

# Ensure access.log exists and is readable
if [ -f /var/log/nginx/access.log ]; then
    # Count unique IPs in nginx access log for today
    UNIQUE_USERS=$(grep "$TODAY" /var/log/nginx/access.log | awk '{print $1}' | sort -u | wc -l)
    echo "[$TIMESTAMP] $UNIQUE_USERS" >> /var/log/metrics/daily_users.log

    # Count searches for today
    DAILY_SEARCHES=$(grep "$TODAY" /var/log/nginx/access.log | grep "GET /?q=" | wc -l)
    echo "[$TIMESTAMP] $DAILY_SEARCHES" >> /var/log/metrics/daily_searches.log
else
    echo "[$TIMESTAMP] 0" >> /var/log/metrics/daily_users.log
    echo "[$TIMESTAMP] 0" >> /var/log/metrics/daily_searches.log
fi

# Disk usage
DISK_USAGE=$(df -h / | tail -n1 | awk '{print $5}')
echo "[$TIMESTAMP] $DISK_USAGE" >> /var/log/metrics/disk_usage.log

# Weekly cost
WEEKLY_COST="2.80"
echo "[$TIMESTAMP] $WEEKLY_COST" >> /var/log/metrics/weekly_cost.log