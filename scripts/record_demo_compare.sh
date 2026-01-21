#!/bin/bash
#
# record_demo_compare.sh — Record demo_compare.gif for goful directory comparison.
#
# Records a terminal session (tmux + asciinema) showing goful comparing three
# directories, then post-processes the cast (insert timing rows via Ruby),
# converts to GIF with agg, and overlays FFmpeg drawtext/drawbox annotations.
#
# The demo shows 9 navigation steps through differences, with UI highlights
# (drawbox overlays) marking the active file items and status lines.
#
# Prerequisites (macOS):
#   - tmux, asciinema, agg, ffmpeg, bc
#   - Ruby script: insert_timing_rows_standalone.rb (path set below)
#   - goful built at ./bin/goful
#   - Demo dirs: /tmp/demo/alpha, /tmp/demo/beta, /tmp/demo/gamma
#
# Output: .github/demo_compare.gif
#
# Terminal size: 120x18
#
# Naming (repeated patterns):
#   ff_*         — FFmpeg filter builder: emits one filter string with enable='between(t,start,end)'
#   vf_*         — Video filter chain: vf_chain (the -vf value), vf_append, vf_append_user_app, vf_append_ui_highlight_*
#   t_curr,t_prev— Elapsed seconds; tick(sec)= sleep + advance both
#   demo_step    — One interaction: send keys → tick → vf_append_user_app(user_label, app_label)
#   overlay_scheme — Set active COLOR_USER_CURRENT, COLOR_APP_CURRENT for subsequent ff_drawtext
#   VF_POS_*     — FFmpeg overlay position expressions (x/y for drawtext)
#   COLOR_{USER|APP}_{PRIMARY|ALT|CURRENT)}, FONT_SIZE_{USER|APP} — overlay style
#   PAUSE_*      — Timing constants for tick() calls
#   GOFUL_ITEM_TOP_* — Y positions of goful list items (for UI highlights)
#   UI_HIGHLIGHT_* — Constants for UI highlight boxes (HL_*, STATUS_*)

# --- Paths and session ---
TMUX_SESSION="demo_compare"
CAST_RAW_PATH="/tmp/demo_compare.1.cast"
CAST_FINAL_PATH="/tmp/demo_compare.cast"
GIF_TEMP_PATH="/tmp/demo_compare.gif"
GIF_OUTPUT_PATH=".github/demo_compare.gif"
GOFUL_BIN="$(pwd)/bin/goful"
INSERT_TIMING_SCRIPT="/Users/fareed/Documents/dev/ruby/markdown_exec/demo/lib/insert_timing_rows_standalone.rb"

# --- Overlay style: COLOR_*, FONT_SIZE_*, VF_POS_* (and BORDER_* for ff_drawtext) ---
COLOR_USER_PRIMARY=white
COLOR_USER_ALT=pink
COLOR_APP_PRIMARY=yellow
COLOR_APP_ALT=orange
COLOR_BOX=orange
BORDER_WIDTH=4
BORDER_COLOR=black
FONT_SIZE_USER=48
FONT_SIZE_APP=36

# VF overlay position expressions (FFmpeg -vf). USER fixed; APP updated per demo_step.
VF_CENTER_X='w/2-text_w/2'
VF_POS_USER_X="$VF_CENTER_X"
VF_POS_USER_Y='h*3/4-text_h'
VF_POS_APP_X=40
VF_POS_APP_Y=10

# --- Timing constants (seconds) ---
PAUSE_SETUP=3
PAUSE_OPEN=3
PAUSE_DIFF=2
PAUSE_EXIT=2.5
PAUSE_BEAT=1

# --- goful UI layout: list item Y positions (for UI highlight boxes) ---
GOFUL_ITEM_TOP_1=44
GOFUL_ITEM_TOP_2=64
GOFUL_ITEM_TOP_3=84
GOFUL_ITEM_TOP_4=104
GOFUL_ITEM_TOP_5=124
GOFUL_ITEM_TOP_6=143
GOFUL_ITEM_TOP_7=163

# --- UI highlight box constants ---
HL_HEIGHT=28
HL_X=4
HL_W=1016
STATUS_Y=300
STATUS2_Y=320
UI_HIGHLIGHT_BORDER_WIDTH=2

# --- State: t_curr/t_prev for enable='between(t,...)' in ff_* ---
t_curr=0
t_prev=0

# --- vf_chain: accumulated -vf filter string (comma‑separated) ---
vf_chain=""

# -----------------------------------------------------------------------------
# ff_drawtext t_start t_end text fontcolor fontsize borderw bordercolor x y
#
# Emits one drawtext filter with enable='between(t,t_start,t_end)'.
# x,y may be FFmpeg expressions (e.g. VF_CENTER_X).
# -----------------------------------------------------------------------------
ff_drawtext() {
  local t_start="${1:-0}"
  local t_end="${2:-0}"
  printf "drawtext=text='%s':fontcolor=%s:fontsize=%s:borderw=%s:bordercolor=%s:x=%s:y=%s:enable='between(t,%s,%s)'" \
    "${3:-?}" "${4:-red}" "${5:-24}" "${6:-4}" "${7:-black}" "${8:-0}" "${9:-0}" "$t_start" "$t_end"
}

# -----------------------------------------------------------------------------
# ff_drawbox x y w h color border_thickness t_start t_end
#
# Emits one drawbox filter with enable='between(t,t_start,t_end)'.
# -----------------------------------------------------------------------------
ff_drawbox() {
  local x="${1:-24}"
  local y="${2:-12}"
  local w="${3:-200}"
  local h="${4:-60}"
  local color="${5:-red@0.6}"
  local border_thickness="${6:-4}"
  local t_start="${7:-1}"
  local t_end="${8:-40}"
  printf "drawbox=x=%s:y=%s:w=%s:h=%s:color=%s:t=%s:enable='between(t,%s,%s)'" \
    "$x" "$y" "$w" "$h" "$color" "$border_thickness" "$t_start" "$t_end"
}

# -----------------------------------------------------------------------------
# ff_drawbox_inner x y color border_thickness t_start t_end
#
# Full‑frame rect (margin x,y; width=iw-2*x, height=ih-2*y).
# -----------------------------------------------------------------------------
ff_drawbox_inner() {
  local x="${1:-24}"
  local y="${2:-12}"
  local color="${3:-red@0.6}"
  local border_thickness="${4:-4}"
  local t_start="${5:-1}"
  local t_end="${6:-40}"
  printf "drawbox=x=%s:w=iw-2*%s:y=%s:h=ih-2*%s:color=%s:t=%s:enable='between(t,%s,%s)'" \
    "$x" "$x" "$y" "$y" "$color" "$border_thickness" "$t_start" "$t_end"
}

# -----------------------------------------------------------------------------
# bc_eval expr
#
# Evaluates numeric expression with bc -l. Used for t_* math in overlay timings.
# -----------------------------------------------------------------------------
bc_eval() {
  echo "$*" | bc -l
}

# -----------------------------------------------------------------------------
# tick seconds
#
# Sleep(seconds), then: t_prev=t_curr, t_curr=t_curr+seconds. Overlay filters
# use (t_prev,t_curr) for the interval just waited.
# -----------------------------------------------------------------------------
tick() {
  t_prev="$t_curr"
  t_curr=$(bc_eval "$t_curr + $1")
  sleep "$1"
}

# -----------------------------------------------------------------------------
# overlay_scheme user_color app_color
#
# Sets COLOR_USER_CURRENT and COLOR_APP_CURRENT for subsequent ff_drawtext.
# -----------------------------------------------------------------------------
overlay_scheme() {
  COLOR_USER_CURRENT=$1
  COLOR_APP_CURRENT=$2
}

# -----------------------------------------------------------------------------
# vf_append filter_spec
#
# Appends one filter to vf_chain (comma‑separated).
# -----------------------------------------------------------------------------
vf_append() {
  vf_chain="$vf_chain,$1"
}

# -----------------------------------------------------------------------------
# vf_append_user_app user_label app_label
#
# Appends two ff_drawtext filters for the last (t_prev,t_curr) interval:
#   user:  (t_prev-1, t_curr-1) at VF_POS_USER_*, FONT_SIZE_USER, COLOR_USER_CURRENT
#   app:   (t_prev+1, t_curr)   at VF_POS_APP_*,  FONT_SIZE_APP,  COLOR_APP_CURRENT
# app_label may be empty (omits app overlay).
# -----------------------------------------------------------------------------
vf_append_user_app() {
  local t_user_start t_user_end t_app_start t_app_end
  t_user_start=$(bc_eval "$t_prev - 1")
  t_user_end=$(bc_eval "$t_curr - 1")
  t_app_start=$(bc_eval "$t_prev + 1")
  t_app_end=$t_curr
  vf_append "$(ff_drawtext $t_user_start $t_user_end "$1" "$COLOR_USER_CURRENT" $FONT_SIZE_USER $BORDER_WIDTH $BORDER_COLOR $VF_POS_USER_X $VF_POS_USER_Y)"
  [[ -n $2 ]] && vf_append "$(ff_drawtext $t_app_start $t_app_end "$2" "$COLOR_APP_CURRENT" $FONT_SIZE_APP $BORDER_WIDTH $BORDER_COLOR $VF_POS_APP_X $VF_POS_APP_Y)"
}

# -----------------------------------------------------------------------------
# demo_step keys seconds user_label app_label
#
# One recorded interaction: tmux send-keys → tick(seconds) → vf_append_user_app.
# -----------------------------------------------------------------------------
demo_step() {
  tmux send-keys -t $TMUX_SESSION "$1"
  tick "$2"
  vf_append_user_app "$3" "$4"
}

# -----------------------------------------------------------------------------
# vf_append_ui_highlight item_y
#
# Highlights one goful list item (at item_y) and the status line (STATUS_Y)
# for the last (t_prev,t_curr) interval. Used after demo_step to mark the
# active item.
# -----------------------------------------------------------------------------
vf_append_ui_highlight() {
  vf_append "$(ff_drawbox $HL_X "$1" $HL_W $HL_HEIGHT $COLOR_BOX $UI_HIGHLIGHT_BORDER_WIDTH $t_prev $t_curr)"
  vf_append "$(ff_drawbox $HL_X $STATUS_Y $HL_W $HL_HEIGHT $COLOR_BOX $UI_HIGHLIGHT_BORDER_WIDTH $t_prev $t_curr)"
}

# -----------------------------------------------------------------------------
# vf_append_ui_highlight_status2
#
# Highlights only the second status line (STATUS2_Y) for the last interval.
# Used when showing "All differences found" message.
# -----------------------------------------------------------------------------
vf_append_ui_highlight_status2() {
  vf_append "$(ff_drawbox $HL_X $STATUS2_Y $HL_W $HL_HEIGHT $COLOR_BOX $UI_HIGHLIGHT_BORDER_WIDTH $t_prev $t_curr)"
}

# -----------------------------------------------------------------------------
# vf_append_ui_highlight_3 item_y1 item_y2 item_y3
#
# Highlights three goful list items (split into 3 columns) and the status line.
# Used when multiple items are active (e.g., navigating through subdirectories).
# -----------------------------------------------------------------------------
vf_append_ui_highlight_3() {
  local cw x2 x3
  cw=$(bc_eval "$HL_W / 3")
  x2=$(bc_eval "$HL_X + $HL_W / 3")
  x3=$(bc_eval "$HL_X + (($HL_W / 3) * 2)")
  vf_append "$(ff_drawbox $HL_X "$1" $cw $HL_HEIGHT $COLOR_BOX $UI_HIGHLIGHT_BORDER_WIDTH $t_prev $t_curr)"
  vf_append "$(ff_drawbox $x2 "$2" $cw $HL_HEIGHT $COLOR_BOX $UI_HIGHLIGHT_BORDER_WIDTH $t_prev $t_curr)"
  vf_append "$(ff_drawbox $x3 "$3" $cw $HL_HEIGHT $COLOR_BOX $UI_HIGHLIGHT_BORDER_WIDTH $t_prev $t_curr)"
  vf_append "$(ff_drawbox $HL_X $STATUS_Y $HL_W $HL_HEIGHT $COLOR_BOX $UI_HIGHLIGHT_BORDER_WIDTH $t_prev $t_curr)"
}

# -----------------------------------------------------------------------------
# gif_apply_vf input_gif output_gif vf_chain
#
# ffmpeg -i in -vf "$vf_chain" -loop 0 -y out.
# -----------------------------------------------------------------------------
gif_apply_vf() {
  ffmpeg -i "$1" -vf "$3" -loop 0 -y "$2"
}

# --- Bootstrap: scheme and first vf (placeholder) ---
overlay_scheme $COLOR_USER_PRIMARY $COLOR_APP_PRIMARY
vf_chain="$(ff_drawtext 0 0 ' ' white 72 4 black 0 0)"

# --- Setup tmux session and start asciinema recording ---
tmux kill-session -t $TMUX_SESSION 2>/dev/null || true
tmux new-session -d -s $TMUX_SESSION -x 120 -y 18

tmux send-keys -t $TMUX_SESSION "unset NO_COLOR && export TERM=xterm-256color COLORTERM=truecolor && asciinema rec --overwrite --cols 120 --rows 18 $CAST_RAW_PATH" Enter
tick $PAUSE_BEAT

# --- Start goful with three demo directories ---
tmux send-keys -t $TMUX_SESSION "$GOFUL_BIN /tmp/demo/alpha /tmp/demo/beta /tmp/demo/gamma" Enter
tick $PAUSE_OPEN

# Opening title overlay
vf_append "$(ff_drawtext 0 $t_curr 'Compare three directories' $COLOR_USER_CURRENT 60 $BORDER_WIDTH $BORDER_COLOR '(w-text_w)/2' 'h*3/4-text_h')"
tick $PAUSE_BEAT

# --- Demo steps: navigate through differences ---
VF_POS_APP_X=$VF_CENTER_X
VF_POS_APP_Y=55

# Step 1: [ to first difference (file2)
overlay_scheme $COLOR_USER_PRIMARY $COLOR_APP_PRIMARY
demo_step '[' $PAUSE_SETUP '1. Press [ to start' 'First difference'
vf_append_ui_highlight $GOFUL_ITEM_TOP_3

# Step 2: ] to next difference (file3)
VF_POS_APP_Y=75
overlay_scheme $COLOR_USER_ALT $COLOR_APP_ALT
demo_step ']' $PAUSE_SETUP '2. Press ] to continue' 'Next difference'
vf_append_ui_highlight $GOFUL_ITEM_TOP_4

# Step 3: ] to next difference (file4)
VF_POS_APP_Y=95
overlay_scheme $COLOR_USER_PRIMARY $COLOR_APP_PRIMARY
demo_step ']' $PAUSE_DIFF '3. Press ] to continue' ''
vf_append_ui_highlight $GOFUL_ITEM_TOP_5

# Step 4: ] to next difference (file5)
overlay_scheme $COLOR_USER_ALT $COLOR_APP_ALT
demo_step ']' $PAUSE_DIFF '4. Press ] to continue' ''
vf_append_ui_highlight_3 $GOFUL_ITEM_TOP_6 $GOFUL_ITEM_TOP_6 $GOFUL_ITEM_TOP_7

# Step 5: ] to next difference (file6)
overlay_scheme $COLOR_USER_PRIMARY $COLOR_APP_PRIMARY
demo_step ']' $PAUSE_DIFF '5. Press ] to continue' ''
vf_append_ui_highlight_3 $GOFUL_ITEM_TOP_7 $GOFUL_ITEM_TOP_7 $GOFUL_ITEM_TOP_6

# Step 6: ] to next difference (subdir/subfile2)
overlay_scheme $COLOR_USER_ALT $COLOR_APP_ALT
demo_step ']' $PAUSE_DIFF '6. Press ] to continue' ''
vf_append_ui_highlight_3 $GOFUL_ITEM_TOP_3 $GOFUL_ITEM_TOP_3 $GOFUL_ITEM_TOP_5

# Step 7: ] to next difference (subdir/another/subfile3)
overlay_scheme $COLOR_USER_PRIMARY $COLOR_APP_PRIMARY
demo_step ']' $PAUSE_DIFF '7. Press ] to continue' ''
vf_append_ui_highlight_3 $GOFUL_ITEM_TOP_1 $GOFUL_ITEM_TOP_1 $GOFUL_ITEM_TOP_2

# Step 8: ] to next difference (subdir/another/subfile3b)
overlay_scheme $COLOR_USER_ALT $COLOR_APP_ALT
demo_step ']' $PAUSE_DIFF '8. Press ] to continue' ''
vf_append_ui_highlight_3 $GOFUL_ITEM_TOP_2 $GOFUL_ITEM_TOP_2 $GOFUL_ITEM_TOP_1

# Step 9: ] all differences found
overlay_scheme $COLOR_USER_PRIMARY $COLOR_APP_PRIMARY
demo_step ']' $PAUSE_DIFF '9. Press ] to continue' 'All differences found'
vf_append_ui_highlight_status2

# --- Quit goful ---
tmux send-keys -t $TMUX_SESSION "q"
tick $PAUSE_DIFF
vf_append "$(ff_drawtext $t_prev $t_curr 'Press q to quit' $COLOR_USER_CURRENT $FONT_SIZE_USER $BORDER_WIDTH $BORDER_COLOR $VF_POS_USER_X $VF_POS_USER_Y)"
vf_append "$(ff_drawtext $t_prev $t_curr 'Confirm' $COLOR_APP_CURRENT $FONT_SIZE_APP $BORDER_WIDTH $BORDER_COLOR $VF_POS_APP_X $STATUS_Y)"

tmux send-keys -t $TMUX_SESSION "y"
sleep $PAUSE_EXIT

# Exit asciinema recording
tmux send-keys -t $TMUX_SESSION "exit" Enter
sleep $PAUSE_EXIT

# Cleanup
tmux kill-session -t $TMUX_SESSION 2>/dev/null || true

# Post-process: insert timing rows, convert to GIF, apply overlays
"$INSERT_TIMING_SCRIPT" "$CAST_RAW_PATH" "$CAST_FINAL_PATH"
echo "Converting to GIF..."
agg --theme asciinema "$CAST_FINAL_PATH" "$GIF_TEMP_PATH"

gif_apply_vf "$GIF_TEMP_PATH" "$GIF_OUTPUT_PATH" "$vf_chain"
echo "Demo recorded: $GIF_OUTPUT_PATH"
