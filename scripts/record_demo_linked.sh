#!/bin/bash
# Record demo_linked.gif - Linked navigation mode
# Terminal size: 120x24

SESSION="demo_linked"
CAST_FILE="/tmp/demo_linked.cast"
GIF_FILE=".github/demo_linked.gif"
GOFUL_BIN="$(pwd)/bin/goful"

# Kill any existing session
tmux kill-session -t $SESSION 2>/dev/null

# Create new tmux session with specific dimensions (120x24)
tmux new-session -d -s $SESSION -x 120 -y 24

# Start asciinema INSIDE tmux with proper environment
# CRITICAL: unset NO_COLOR and set TERM/COLORTERM before asciinema starts
tmux send-keys -t $SESSION "unset NO_COLOR && export TERM=xterm-256color COLORTERM=truecolor && asciinema rec --overwrite --cols 120 --rows 24 $CAST_FILE" Enter
sleep 3

# Start goful with three demo directories
tmux send-keys -t $SESSION "$GOFUL_BIN /tmp/demo/alpha /tmp/demo/beta /tmp/demo/gamma" Enter
sleep 4

# Linked mode is on by default, but let's ensure it's on
# Press L to toggle (if off) or just demonstrate it
# Navigate down to show cursor sync across panes
tmux send-keys -t $SESSION "j"
sleep 1
tmux send-keys -t $SESSION "j"
sleep 1
tmux send-keys -t $SESSION "j"
sleep 1

# Navigate into subdir to show directory sync
tmux send-keys -t $SESSION "l"
sleep 2

# Navigate back up
tmux send-keys -t $SESSION "u"
sleep 2

# Move cursor to show sync
tmux send-keys -t $SESSION "k"
sleep 1
tmux send-keys -t $SESSION "k"
sleep 1

# Quit goful
tmux send-keys -t $SESSION "q"
sleep 1
tmux send-keys -t $SESSION "y"
sleep 2

# Exit asciinema recording
tmux send-keys -t $SESSION "exit" Enter
sleep 2

# Cleanup
tmux kill-session -t $SESSION 2>/dev/null

# Convert to GIF
echo "Converting to GIF..."
agg --theme asciinema $CAST_FILE $GIF_FILE

echo "Demo recorded: $GIF_FILE"
