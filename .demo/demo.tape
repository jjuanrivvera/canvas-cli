# canvas demo — rendered with VHS (https://github.com/charmbracelet/vhs).
# No account, signup, keyring or external network: only invented local API data.
# Regenerate from the repository root:
#   go build -o bin/canvas ./cmd/canvas && go build -o bin/demoserver ./tools/demoserver && vhs .demo/demo.tape
# The shell exits on a failed command rather than recording an error as a result.
# This CLI has no --jq flag; jq filters its actual JSON output instead.
# env -i removes inherited profiles; no home directory means the internal config/updater cannot read real settings.
# XDG_CACHE_HOME isolates version detection even with --no-cache.

Output assets/demo.gif

Set Shell bash
Set FontSize 15
Set Width 1200
Set Height 760
Set Padding 18
Set Theme "Dracula"
Set TypingSpeed 28ms

Hide
Type "set -e -o pipefail"
Enter
Type "export PATH=$PWD/bin:$PATH XDG_CONFIG_HOME=$(mktemp -d)"
Enter
Type "export CANVAS_URL=http://127.0.0.1:8645 CANVAS_TOKEN=demo-token"
Enter
Type `demoserver 127.0.0.1:8645 >"$XDG_CONFIG_HOME/demoserver.log" 2>&1 & demo_server_pid=$!`
Enter
Type `trap 'kill "$demo_server_pid" 2>/dev/null || true; rm -rf "$XDG_CONFIG_HOME"' EXIT`
Enter
Type `sleep 1; kill -0 "$demo_server_pid"`
Enter
Type `printf "settings: {}\n" >"$XDG_CONFIG_HOME/config.yaml"`
Enter
Type `alias canvas='env -i PATH="$PATH" XDG_CONFIG_HOME="$XDG_CONFIG_HOME" XDG_CACHE_HOME="$XDG_CONFIG_HOME/cache" CANVAS_URL="$CANVAS_URL" CANVAS_TOKEN="$CANVAS_TOKEN" canvas --config "$XDG_CONFIG_HOME/config.yaml" --no-cache --quiet'`
Enter
Type "clear"
Enter
Show

Sleep 500ms

Type "canvas courses list --columns id,name,course_code,workflow_state"
Sleep 250ms
Enter
Sleep 2400ms

Type "canvas assignments list --course-id 101 --columns id,name,points_possible,published"
Sleep 250ms
Enter
Sleep 2400ms

Type "clear"
Enter
Sleep 300ms

Type "canvas assignments list --course-id 101 -o json | jq '.[] | {name, points: .points_possible}'"
Sleep 250ms
Enter
Sleep 2400ms

Type "canvas assignments create --course-id 101 --name 'Demo constellation sketch' --points 15 --submission-types online_upload --published"
Sleep 250ms
Enter
Sleep 2400ms

Type "canvas assignments list --course-id 101 --columns id,name,points_possible,published"
Sleep 250ms
Enter
Sleep 2400ms
