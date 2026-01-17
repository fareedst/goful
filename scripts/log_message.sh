#!/bin/bash

# Log file path
: "${LOG_FILE:=script_log.log}"

# Function to log a message with timestamp and arguments
log_message() {
  local timestamp=$(date -u +"%Y-%m-%d %H:%M:%S")
  echo "[$timestamp] $* " >> "$LOG_FILE"
}
echo "$(date -u) : $*" >> /tmp/panel-demo/logs/t1.log

log_message "$*"

# Output the last 4 lines to stdout
tail -n 4 "$LOG_FILE"