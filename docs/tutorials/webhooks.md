# Webhooks Tutorial

Learn how to receive real-time Canvas events using webhooks.

## Overview

Canvas CLI includes a webhook listener that receives real-time notifications from Canvas LMS. This enables you to:

- React to grade changes instantly
- Track assignment submissions in real-time
- Monitor enrollment changes
- Build integrations with external systems
- Automate workflows based on Canvas events

## Supported Events

Canvas CLI supports 19 event types:

| Category | Events |
|----------|--------|
| **Assignment** | `assignment_created`, `assignment_updated`, `assignment_deleted` |
| **Submission** | `submission_created`, `submission_updated`, `grade_change` |
| **Enrollment** | `enrollment_created`, `enrollment_updated`, `enrollment_deleted` |
| **User** | `user_created`, `user_updated` |
| **Course** | `course_created`, `course_updated`, `course_completed` |
| **Discussion** | `discussion_topic_created`, `discussion_entry_created` |
| **Quiz** | `quiz_submitted` |
| **Conversation** | `conversation_created`, `conversation_message_created` |

List all events with:

```bash
canvas webhook events
```

## Quick Start

### Step 1: Start the Webhook Listener

```bash
canvas webhook listen --log
```

This starts a server on `http://localhost:8080` that:

- Listens for Canvas events at `/webhook`
- Provides a health check at `/health`
- Logs all incoming events to stdout

### Step 2: Test the Listener

In another terminal, simulate an event:

```bash
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -d '{"event_type":"grade_change","body":{"user_id":123,"score":95}}'
```

You should see the event logged in the first terminal.

## Configuration Options

### Address

Change the listening address and port:

```bash
# Listen on a specific port
canvas webhook listen --addr :9000

# Listen on specific interface
canvas webhook listen --addr 127.0.0.1:8080

# Listen on all interfaces
canvas webhook listen --addr 0.0.0.0:8080
```

### Security Options

Canvas CLI supports two verification methods for webhook security:

#### Option 1: JWT Verification (Canvas Data Services)

For **Instructure-hosted Canvas** using Canvas Data Services, use JWT verification:

```bash
# Use official Canvas Data Services JWK endpoint
canvas webhook listen --canvas-data-services --log

# Or specify a custom JWK endpoint
canvas webhook listen --jwks-url "https://your-jwks-endpoint.com/jwks"
```

With JWT verification:

- Payloads are signed JWTs with event data in claims
- JWKs are fetched from Canvas and cached (1 hour TTL)
- Keys rotate monthly; the listener handles this automatically
- Invalid JWTs are rejected (401 Unauthorized)

!!! tip "Recommended for Canvas Data Services"
    Use `--canvas-data-services` for production deployments with Instructure-hosted Canvas.

#### Option 2: HMAC Signatures (Custom Integrations)

For custom integrations or self-hosted Canvas, use HMAC-SHA256:

```bash
canvas webhook listen --secret "your-canvas-webhook-secret"
```

With HMAC verification:

- Requests without valid signatures are rejected (401 Unauthorized)
- Prevents unauthorized systems from sending fake events
- The secret must match what's configured in Canvas

#### Combined Mode (Fallback)

You can enable both methods for maximum compatibility:

```bash
canvas webhook listen --canvas-data-services --secret "backup-secret" --log
```

The listener tries JWT first, then falls back to HMAC if JWT fails.

!!! warning "Always Use Verification in Production"
    Without `--canvas-data-services` or `--secret`, any system can send events to your listener.

### Filter Specific Events

Listen only for events you care about:

```bash
# Only grade changes
canvas webhook listen --events grade_change

# Multiple event types
canvas webhook listen --events grade_change,submission_created,assignment_updated

# All submission-related events
canvas webhook listen --events submission_created,submission_updated,grade_change
```

### Enable Logging

See detailed request information:

```bash
canvas webhook listen --log
```

Output includes event type, ID, timestamp, and full payload.

## Setting Up Canvas

### Canvas Cloud (Instructure Hosted)

1. Contact your Canvas administrator - webhook configuration requires admin access
2. Navigate to **Admin > Developer Keys**
3. Create a new Developer Key with webhook capabilities
4. Configure the webhook endpoint URL
5. Select the events to subscribe to
6. Copy the signing secret for `--secret` flag

### Self-Hosted Canvas

1. Go to **Admin > Settings > Webhooks** (or use Canvas API)
2. Add your webhook endpoint URL
3. Select events to receive
4. Save and note the signing secret

### Exposing Local Development Server

For testing with a real Canvas instance, expose your local server using ngrok:

```bash
# Terminal 1: Start webhook listener
canvas webhook listen --secret "test-secret" --log

# Terminal 2: Expose with ngrok
ngrok http 8080
```

Copy the ngrok HTTPS URL (e.g., `https://abc123.ngrok.io`) and configure it in Canvas as: `https://abc123.ngrok.io/webhook`

## Use Case: React to Grade Changes

### Option 1: Shell Script Integration

Create a handler script:

```bash
#!/bin/bash
# ~/scripts/on-grade-change.sh

# Read event JSON from stdin
EVENT=$(cat)

# Extract fields with jq
USER_ID=$(echo "$EVENT" | jq -r '.body.user_id')
SCORE=$(echo "$EVENT" | jq -r '.body.score')
ASSIGNMENT_ID=$(echo "$EVENT" | jq -r '.body.assignment_id')

echo "Grade received: User $USER_ID scored $SCORE on assignment $ASSIGNMENT_ID"

# Send Slack notification
curl -X POST "https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK" \
  -H "Content-Type: application/json" \
  -d "{\"text\": \"New grade: Student $USER_ID scored $SCORE\"}"

# Or send email, update database, etc.
```

Pipe webhook events to your script:

```bash
canvas webhook listen --events grade_change --log 2>&1 | \
  grep --line-buffered "Body:" | \
  sed 's/.*Body: //' | \
  while read line; do echo "$line" | ~/scripts/on-grade-change.sh; done
```

### Option 2: Custom Go Application (Recommended)

For production use, build a custom application using the webhook package.

#### Step 1: Create Project

```bash
mkdir grade-notifier && cd grade-notifier
go mod init grade-notifier
go get github.com/jjuanrivvera/canvas-cli/internal/webhook
```

#### Step 2: Write Handler

```go
// main.go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "bytes"
    "encoding/json"

    "github.com/jjuanrivvera/canvas-cli/internal/webhook"
)

func main() {
    secret := os.Getenv("WEBHOOK_SECRET")
    slackURL := os.Getenv("SLACK_WEBHOOK_URL")

    // Create listener
    listener := webhook.NewListener(webhook.Config{
        Addr:   ":8080",
        Secret: secret,
        Logger: log.Default(),
    })

    // Register grade change handler
    listener.On(webhook.EventGradeChange, func(ctx context.Context, event *webhook.Event) error {
        return handleGradeChange(event, slackURL)
    })

    // Graceful shutdown
    ctx, cancel := context.WithCancel(context.Background())
    go func() {
        sigCh := make(chan os.Signal, 1)
        signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
        <-sigCh
        fmt.Println("\nShutting down...")
        cancel()
    }()

    fmt.Printf("Listening for grade changes on :8080...\n")
    if err := listener.Start(ctx); err != nil && err != context.Canceled {
        log.Fatal(err)
    }
}

func handleGradeChange(event *webhook.Event, slackURL string) error {
    body := event.Body

    userID := body["user_id"]
    score := body["score"]
    assignmentID := body["assignment_id"]
    courseID := body["course_id"]

    log.Printf("Grade change: User %v scored %v on assignment %v (course %v)",
        userID, score, assignmentID, courseID)

    // Send Slack notification
    if slackURL != "" {
        msg := map[string]string{
            "text": fmt.Sprintf("New grade: Student %v scored %v on assignment %v",
                userID, score, assignmentID),
        }
        jsonData, _ := json.Marshal(msg)
        http.Post(slackURL, "application/json", bytes.NewBuffer(jsonData))
    }

    return nil
}
```

#### Step 3: Build and Run

```bash
go build -o grade-notifier

# Run with environment variables
WEBHOOK_SECRET="your-canvas-secret" \
SLACK_WEBHOOK_URL="https://hooks.slack.com/services/xxx" \
./grade-notifier
```

## Testing Webhooks

### Simulate Events Locally

Test without Canvas by sending mock events:

```bash
# Without signature verification
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event_type": "grade_change",
    "body": {
      "user_id": 12345,
      "assignment_id": 67890,
      "course_id": 111,
      "score": 95,
      "grade": "A"
    }
  }'
```

### Test with Signature Verification

```bash
# Set your secret
SECRET="your-secret"

# Create the payload
BODY='{"event_type":"grade_change","body":{"user_id":123,"score":95}}'

# Generate HMAC-SHA256 signature
SIGNATURE=$(echo -n "$BODY" | openssl dgst -sha256 -hmac "$SECRET" | cut -d' ' -f2)

# Send signed request
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -H "X-Canvas-Signature: $SIGNATURE" \
  -d "$BODY"
```

### Health Check

Verify the listener is running:

```bash
curl http://localhost:8080/health
# Returns: {"status":"ok"}
```

## Architecture

```
┌─────────────┐         ┌──────────────────┐         ┌─────────────────┐
│  Canvas LMS │  POST   │  Webhook Listener │         │  Your Handler   │
│  (or Data   │────────▶│  (canvas webhook  │────────▶│  (script/code)  │
│   Services) │         │   listen)         │         │                 │
└─────────────┘         └──────────────────┘         └─────────────────┘
                               │
                               ▼
                   ┌───────────────────────┐
                   │   Verification Layer  │
                   ├───────────────────────┤
                   │ 1. JWT (Data Services)│
                   │    ↓ fallback         │
                   │ 2. HMAC-SHA256        │
                   └───────────────────────┘
                               │
                               ▼
                   ┌───────────────────────┐
                   │   JWK Cache (1hr TTL) │
                   │   (Canvas JWKs)       │
                   └───────────────────────┘
```

**Flow:**

1. Canvas sends HTTP POST to `/webhook` endpoint
2. Listener verifies signature:
   - If `--canvas-data-services`: Verify JWT using cached JWKs
   - If `--secret`: Verify HMAC-SHA256 signature
   - If both: Try JWT first, then HMAC as fallback
3. Event is parsed and routed to registered handlers
4. Handler processes the event (your custom logic)
5. Returns 200 OK to Canvas (or error to trigger retry)

## Best Practices

!!! tip "Use Canvas Data Services Mode"
    For Instructure-hosted Canvas, use `--canvas-data-services` for JWT verification. It's the recommended approach for production.

!!! tip "Handle Failures Gracefully"
    If your handler returns an error, Canvas will retry the webhook. Make handlers idempotent.

!!! tip "Use Specific Events"
    Only subscribe to events you need. Use `--events` to filter and reduce noise.

!!! tip "Monitor Health"
    Set up monitoring on the `/health` endpoint to ensure your listener is running.

!!! tip "Log for Debugging"
    Use `--log` during development to see full event payloads.

## Troubleshooting

### Webhook Not Receiving Events

1. Check listener is running: `curl http://localhost:8080/health`
2. Verify Canvas can reach your server (use ngrok for local dev)
3. Check Canvas webhook configuration for correct URL
4. Ensure events are enabled in Canvas webhook settings

### 401 Unauthorized Errors

- Signature mismatch - verify the secret matches Canvas configuration
- Check `X-Canvas-Signature` header is being sent
- Ensure the secret has no extra whitespace

### Events Processing But Nothing Happens

- Check your handler code for errors
- Enable `--log` to see event payloads
- Verify the event type matches what you're filtering for

## Next Steps

- [Command Reference: webhook](../commands/canvas_webhook.md) - Full command documentation
- [Scripting Guide](scripting.md) - More automation techniques
- [Canvas API Documentation](https://canvas.instructure.com/doc/api/) - Canvas webhook payload formats
