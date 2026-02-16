# Testing Guide

Comprehensive guide for testing the IP Verifier service, including manual API testing and GeoIP database update verification.

---

## 📋 Table of Contents

1. [Manual API Testing](#manual-api-testing)
2. [Database Update Testing](#database-update-testing)
3. [Automated Testing](#automated-testing)
4. [Troubleshooting](#troubleshooting)

---

## 🧪 Manual API Testing

### Health Check

Basic health endpoint to verify the service is running:

```bash
curl -s http://localhost/api/v1/health | jq .
```

**Expected Response:**
```json
{
  "status": "healthy"
}
```

### IP Verification Tests

#### Test 1: Allowed Country (US IP)
```bash
curl -s -X POST http://localhost/api/v1/ip-verifier \
  -H "Content-Type: application/json" \
  -d '{
    "ip": "8.8.8.8",
    "allowed_countries": ["US"]
  }' | jq .
```

**Expected Response:**
```json
{
  "ip": "8.8.8.8",
  "country_code": "US",
  "is_allowed": true
}
```

#### Test 2: Not Allowed Country (Russia IP)
```bash
curl -s -X POST http://localhost/api/v1/ip-verifier \
  -H "Content-Type: application/json" \
  -d '{
    "ip": "77.88.8.8",
    "allowed_countries": ["US", "CA"]
  }' | jq .
```

**Expected Response:**
```json
{
  "ip": "77.88.8.8",
  "country_code": "RU",
  "is_allowed": false
}
```

#### Test 3: Multiple Allowed Countries
```bash
curl -s -X POST http://localhost/api/v1/ip-verifier \
  -H "Content-Type: application/json" \
  -d '{
    "ip": "1.1.1.1",
    "allowed_countries": ["US", "AU"]
  }' | jq .
```

**Expected Response:**
```json
{
  "ip": "1.1.1.1",
  "country_code": "AU",
  "is_allowed": true
}
```

#### Test 4: Invalid IP Address
```bash
curl -s -X POST http://localhost/api/v1/ip-verifier \
  -H "Content-Type: application/json" \
  -d '{
    "ip": "invalid-ip",
    "allowed_countries": ["US"]
  }' | jq .
```

**Expected Response:**
```json
{
  "error": "invalid IP address format"
}
```

#### Test 5: Empty Allowed Countries
```bash
curl -s -X POST http://localhost/api/v1/ip-verifier \
  -H "Content-Type: application/json" \
  -d '{
    "ip": "8.8.8.8",
    "allowed_countries": []
  }' | jq .
```

**Expected Response:**
```json
{
  "ip": "8.8.8.8",
  "country_code": "US",
  "is_allowed": false
}
```

### Test Local Development Server

If running locally on port 8080:

```bash
# Health check
curl -s http://localhost:8080/api/v1/health | jq .

# IP verification
curl -s -X POST http://localhost:8080/api/v1/ip-verifier \
  -H "Content-Type: application/json" \
  -d '{"ip":"8.8.8.8","allowed_countries":["US"]}' | jq .
```

### Common Test IPs

| IP Address | Country | Description |
|------------|---------|-------------|
| 8.8.8.8 | US | Google DNS |
| 1.1.1.1 | AU | Cloudflare DNS |
| 77.88.8.8 | RU | Yandex DNS |
| 208.67.222.222 | US | OpenDNS |
| 9.9.9.9 | US | Quad9 DNS |
| 185.228.168.9 | NL | CleanBrowsing DNS |

---

## 🗄️ Database Update Testing

The GeoIP database is automatically updated weekly via a Kubernetes CronJob. Here's how to verify it's working correctly.

### Check Update Schedule

View the CronJob configuration:

```bash
kubectl get cronjob geoip-updater -n ip-verifier
```

**Expected Output:**
```
NAME             SCHEDULE      SUSPEND   ACTIVE   LAST SCHEDULE   AGE
geoip-updater    0 3 * * 3     False     0        6d              7d
```

- **Schedule**: `0 3 * * 3` = Every Wednesday at 3:00 AM UTC
- **LAST SCHEDULE**: Time since last run

### View Update History

Check past update jobs:

```bash
kubectl get jobs -n ip-verifier -l cronjob-name=geoip-updater --sort-by=.status.startTime
```

**Example Output:**
```
NAME                      COMPLETIONS   DURATION   AGE
geoip-updater-28567890    1/1           12s        6d
geoip-updater-28653210    1/1           14s        13d
```

### Check Update Job Logs

View logs from the most recent update:

```bash
# Get the most recent job name
JOB=$(kubectl get jobs -n ip-verifier -l cronjob-name=geoip-updater --sort-by=.status.startTime -o jsonpath='{.items[-1].metadata.name}')

# View logs
kubectl logs -n ip-verifier job/$JOB
```

**Expected Log Output:**
```
time="2026-02-12T03:00:15Z" level=info msg="GeoIP database update starting"
time="2026-02-12T03:00:16Z" level=info msg="Downloading GeoLite2-Country.mmdb"
time="2026-02-12T03:00:25Z" level=info msg="Database updated successfully"
time="2026-02-12T03:00:25Z" level=info msg="GeoIP update complete"
```

### Manually Trigger Database Update

Force an immediate update without waiting for the scheduled time:

```bash
# Trigger update
make k8s-update-db

# Or manually:
kubectl create job -n ip-verifier geoip-manual-$(date +%s) --from=cronjob/geoip-updater

# Wait 10-15 seconds, then check status
kubectl get jobs -n ip-verifier | grep geoip-manual
```

### Verify Database File in Pod

Check if the MMDB file exists and its size:

```bash
# Get pod name
POD=$(kubectl get pod -n ip-verifier -l app=ip-verifier -o jsonpath='{.items[0].metadata.name}')

# Check file (limited in distroless container)
kubectl exec -n ip-verifier $POD -- ls -lh /var/lib/geoip/ 2>/dev/null || echo "Limited shell in distroless"
```

**Note**: The distroless container has no shell utilities, so direct file inspection is limited. The best verification is through functional testing.

### Verify Database is Working

The most reliable test is to verify IP lookups work correctly:

```bash
# Test multiple IPs from different countries
make test-api

# Or manually:
curl -s -X POST http://localhost/api/v1/ip-verifier \
  -H "Content-Type: application/json" \
  -d '{"ip":"8.8.8.8","allowed_countries":["US"]}' | jq .
```

If IP lookups return correct country codes, the database is loaded and working.

### Check Init Container Logs

The init container downloads the database on pod startup:

```bash
# Get pod name
POD=$(kubectl get pod -n ip-verifier -l app=ip-verifier -o jsonpath='{.items[0].metadata.name}')

# View init container logs
kubectl logs -n ip-verifier $POD -c geoip-init
```

**Expected Output:**
```
time="2026-02-16T10:30:05Z" level=info msg="GeoIP database initialization"
time="2026-02-16T10:30:06Z" level=info msg="Downloading GeoLite2-Country.mmdb"
time="2026-02-16T10:30:18Z" level=info msg="Database downloaded successfully"
```

### Quick Database Check Command

Use the Makefile for a quick overview:

```bash
make k8s-check-updates
```

This displays:
- Last 5 update jobs with completion status
- CronJob schedule and configuration
- Database file information (if accessible)

---

## 🤖 Automated Testing

### Run Unit Tests

```bash
# All tests
make test

# With coverage
make test-coverage
open coverage.html
```

**Expected Output:**
```
ok      github.com/KWilliams-dev/ip-verifier/cmd/ip-verifier-api       0.234s
ok      github.com/KWilliams-dev/ip-verifier/internal/handler           0.156s
ok      github.com/KWilliams-dev/ip-verifier/internal/service           0.189s
...
```

### Run E2E Tests

End-to-end tests require a running service:

```bash
# Ensure service is running
make k8s-status

# Run E2E tests
make test-e2e
```

### Automated API Tests

Quick smoke test of all endpoints:

```bash
make test-api
```

This tests:
1. Health endpoint
2. Valid IP with allowed country
3. Valid IP with non-allowed country (optional, check Makefile)

---

## 🔍 Troubleshooting

### Database Not Updating

**Problem**: CronJob not running on schedule

```bash
# Check CronJob exists
kubectl get cronjob -n ip-verifier

# Check for suspended status
kubectl get cronjob geoip-updater -n ip-verifier -o yaml | grep suspend

# Describe CronJob for events
kubectl describe cronjob geoip-updater -n ip-verifier
```

**Solution**: Ensure CronJob is not suspended and check events for errors.

---

**Problem**: Update jobs failing

```bash
# Check failed jobs
kubectl get jobs -n ip-verifier -l cronjob-name=geoip-updater

# View logs from failed job
kubectl logs -n ip-verifier job/<failed-job-name>
```

**Common Issues**:
- Invalid MaxMind credentials
- Network connectivity issues
- PVC mount problems

**Solution**: Verify secret exists and has correct credentials:
```bash
kubectl get secret maxmind-credentials -n ip-verifier
```

---

**Problem**: Database file not accessible in pod

```bash
# Check PVC exists and is bound
kubectl get pvc -n ip-verifier

# Check pod volume mounts
kubectl describe pod -n ip-verifier -l app=ip-verifier | grep -A 5 Mounts
```

**Solution**: Ensure PVC is bound and mounted at `/var/lib/geoip`

---

### IP Lookups Returning Wrong Countries

**Problem**: Country codes are incorrect or outdated

```bash
# Trigger manual database update
make k8s-update-db

# Wait for completion (15-20 seconds)
sleep 20

# Restart pods to reload database
make k8s-restart

# Test again
make test-api
```

---

### Service Not Responding

**Problem**: API requests timing out or failing

```bash
# Check pod status
make k8s-status

# Check logs for errors
make k8s-logs

# Check service endpoints
kubectl get endpoints -n ip-verifier
```

**Solution**: 
1. Ensure pods are Running (2/2 ready)
2. Check logs for application errors
3. Verify service has endpoints

---

### Testing Checklist

Use this checklist to verify full functionality:

- [ ] Health endpoint returns `{"status":"healthy"}`
- [ ] US IP (8.8.8.8) correctly identified as US
- [ ] RU IP (77.88.8.8) correctly identified as RU
- [ ] AU IP (1.1.1.1) correctly identified as AU
- [ ] Invalid IP returns error
- [ ] CronJob exists and is not suspended
- [ ] At least one successful update job exists
- [ ] Init container logs show successful database download
- [ ] All unit tests pass (`make test`)
- [ ] Pods are running (2/2 ready)
- [ ] Service has valid endpoints

---

## 📊 Monitoring Database Updates

### Set Up Recurring Checks

Add to crontab for daily verification:

```bash
# Check update status daily at 9 AM
0 9 * * * kubectl get jobs -n ip-verifier -l cronjob-name=geoip-updater --sort-by=.status.startTime | tail -1
```

### Alert on Failed Updates

Monitor for failed jobs:

```bash
# Check for any failed jobs
kubectl get jobs -n ip-verifier -l cronjob-name=geoip-updater -o jsonpath='{.items[?(@.status.failed>0)].metadata.name}'
```

If output is not empty, an update job has failed and needs investigation.

---

## 🎯 Best Practices

1. **Test after deployment**: Always run `make test-api` after deploying changes
2. **Monitor weekly updates**: Check logs the day after Wednesday to ensure updates succeeded
3. **Manual updates before major releases**: Trigger `make k8s-update-db` before important deployments
4. **Keep credentials current**: MaxMind license keys should be rotated periodically
5. **Test diverse IPs**: Use IPs from various countries to ensure database coverage
6. **Check pod restarts**: Frequent restarts may indicate database loading issues

---

**Quick Reference**: `make help` for all available commands
