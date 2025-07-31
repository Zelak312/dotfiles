#!/bin/bash

# Direction: "up" or "down"
direction=$1

# Get current brightness (as percentage)
current=$(brightnessctl -m | awk -F, '{print $4}' | tr -d '%')

# Set step sizes
small_step=1
large_step=10

# Logic: if below 10%, use 1%, else 10%
if [[ "$direction" == "up" ]]; then
    if [ "$current" -lt 10 ]; then
        brightnessctl set ${small_step}%+
    else
        brightnessctl set ${large_step}%+
    fi
else
    if [ "$current" -le 10 ]; then
        brightnessctl set ${small_step}%-
    else
        brightnessctl set ${large_step}%-
    fi
fi

