# Exit Codes

Canvas CLI uses standard exit codes to indicate command status.

## Exit Code Reference

| Code | Name | Description |
|------|------|-------------|
| `0` | Success | Command completed successfully |
| `1` | General Error | Unspecified error occurred |
| `2` | Usage Error | Invalid command syntax or arguments |
| `64` | Usage Error | Command line usage error (EX_USAGE) |
| `65` | Data Error | Input data was incorrect (EX_DATAERR) |
| `66` | No Input | Input file did not exist or was unreadable (EX_NOINPUT) |
| `69` | Unavailable | Service unavailable (EX_UNAVAILABLE) |
| `70` | Software Error | Internal software error (EX_SOFTWARE) |
| `73` | Can't Create | Can't create output file (EX_CANTCREAT) |
| `74` | IO Error | Input/output error (EX_IOERR) |
| `75` | Temp Failure | Temporary failure, retry later (EX_TEMPFAIL) |
| `77` | Permission Denied | Permission denied (EX_NOPERM) |
| `78` | Config Error | Configuration error (EX_CONFIG) |

## Common Scenarios

### Success (0)

```bash
canvas courses list
echo $?  # Output: 0
```

### Authentication Error (77)

```bash
canvas courses list  # Without authentication
echo $?  # Output: 77
```

### Invalid Arguments (64)

```bash
canvas courses list --invalid-flag
echo $?  # Output: 64
```

### Resource Not Found (65)

```bash
canvas courses get 999999  # Non-existent course
echo $?  # Output: 65
```

### API Error (69)

```bash
# Canvas server unavailable
canvas courses list
echo $?  # Output: 69
```

## Using Exit Codes in Scripts

### Basic Check

```bash
#!/bin/bash
if canvas courses list -o json > courses.json; then
  echo "Success!"
else
  echo "Failed with exit code: $?"
  exit 1
fi
```

### Detailed Error Handling

```bash
#!/bin/bash
canvas courses list -o json > courses.json
exit_code=$?

case $exit_code in
  0)
    echo "Success"
    ;;
  77)
    echo "Authentication required. Run: canvas auth login"
    exit 77
    ;;
  69)
    echo "Canvas server unavailable. Retrying in 60 seconds..."
    sleep 60
    exec "$0" "$@"  # Retry
    ;;
  *)
    echo "Failed with exit code: $exit_code"
    exit $exit_code
    ;;
esac
```

### Retry on Temporary Failure

```bash
#!/bin/bash
MAX_RETRIES=3
RETRY_DELAY=10

for i in $(seq 1 $MAX_RETRIES); do
  canvas courses list -o json > courses.json
  exit_code=$?

  if [ $exit_code -eq 0 ]; then
    break
  elif [ $exit_code -eq 75 ] && [ $i -lt $MAX_RETRIES ]; then
    echo "Temporary failure, retry $i of $MAX_RETRIES..."
    sleep $RETRY_DELAY
  else
    echo "Failed after $i attempts"
    exit $exit_code
  fi
done
```

## CI/CD Integration

### GitHub Actions

```yaml
- name: Fetch Canvas Data
  run: |
    canvas courses list -o json > courses.json
  continue-on-error: false  # Fail the job on non-zero exit

- name: Handle Partial Failures
  run: |
    canvas courses list -o json > courses.json || {
      if [ $? -eq 75 ]; then
        echo "::warning::Temporary failure, data may be stale"
      else
        exit 1
      fi
    }
```

### Jenkins Pipeline

```groovy
pipeline {
  stages {
    stage('Fetch Data') {
      steps {
        script {
          def exitCode = sh(
            script: 'canvas courses list -o json > courses.json',
            returnStatus: true
          )
          if (exitCode == 77) {
            error('Authentication required')
          } else if (exitCode != 0) {
            error("Canvas CLI failed with exit code: ${exitCode}")
          }
        }
      }
    }
  }
}
```

## Best Practices

!!! tip "Always Check Exit Codes"
    In production scripts, always check and handle exit codes appropriately.

!!! tip "Log the Exit Code"
    When debugging, log the exit code along with any error messages:
    ```bash
    canvas courses list 2>&1
    echo "Exit code: $?"
    ```

!!! warning "Don't Ignore Errors"
    Use `set -e` in bash scripts to exit on any error:
    ```bash
    #!/bin/bash
    set -e
    canvas courses list -o json > courses.json
    # Script stops here if command fails
    ```
