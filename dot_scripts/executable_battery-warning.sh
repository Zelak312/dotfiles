#!/bin/bash

LOCKFILE="/tmp/.low_battery_notified"
LOW_BATTERY_LEVEL=20
RESET_LEVEL=25

BATTERY_PATH="/sys/class/power_supply/BAT0"
BATTERY_LEVEL=$(cat "$BATTERY_PATH/capacity")
STATUS=$(cat "$BATTERY_PATH/status")

if [[ "$BATTERY_LEVEL" -le "$LOW_BATTERY_LEVEL" && "$STATUS" == "Discharging" ]]; then
    if [ ! -f "$LOCKFILE" ]; then
        notify-send -u critical "⚠️ Low Battery" "Battery level is at ${BATTERY_LEVEL}%"
        touch "$LOCKFILE"
    fi
elif [[ "$BATTERY_LEVEL" -ge "$RESET_LEVEL" ]]; then
    # Clear lock if battery is charging or has gone above threshold
    [ -f "$LOCKFILE" ] && rm "$LOCKFILE"
fi
