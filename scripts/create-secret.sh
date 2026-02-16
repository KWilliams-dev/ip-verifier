#!/bin/bash

# Script to create MaxMind credentials secret from .env file
# Usage: ./scripts/create-secret.sh

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Creating MaxMind Credentials Secret...${NC}"

# Check if .env exists
if [ ! -f .env ]; then
    echo -e "${RED}Error: .env file not found${NC}"
    echo "Please create a .env file with:"
    echo "  GEOIPUPDATE_ACCOUNT_ID=your_account_id"
    echo "  GEOIPUPDATE_LICENSE_KEY=your_license_key"
    exit 1
fi

# Source .env file
source .env

# Validate variables exist
if [ -z "$GEOIPUPDATE_ACCOUNT_ID" ] || [ -z "$GEOIPUPDATE_LICENSE_KEY" ]; then
    echo -e "${RED}Error: Missing credentials in .env${NC}"
    echo "Required variables:"
    echo "  GEOIPUPDATE_ACCOUNT_ID"
    echo "  GEOIPUPDATE_LICENSE_KEY"
    exit 1
fi

# Create namespace if it doesn't exist
echo "Creating namespace..."
kubectl apply -f k8s/namespace.yaml

# Delete secret if it exists (to update)
kubectl delete secret maxmind-credentials -n ip-verifier 2>/dev/null || true

# Create secret
echo "Creating secret..."
kubectl create secret generic maxmind-credentials \
  --namespace=ip-verifier \
  --from-literal=GEOIPUPDATE_ACCOUNT_ID="${GEOIPUPDATE_ACCOUNT_ID}" \
  --from-literal=GEOIPUPDATE_LICENSE_KEY="${GEOIPUPDATE_LICENSE_KEY}" \
  --from-literal=GEOIPUPDATE_EDITION_IDS=GeoLite2-Country

# Verify
echo -e "\n${GREEN}✓ Secret created successfully${NC}"
echo -e "\nVerifying secret..."
kubectl get secret maxmind-credentials -n ip-verifier

echo -e "\n${GREEN}✓ Done!${NC}"
echo -e "\nNext steps:"
echo "  1. kubectl apply -f k8s/pvc-geoip.yaml"
echo "  2. kubectl apply -f k8s/configmap-geoip.yaml"
echo "  3. kubectl apply -f k8s/deployment.yaml"
echo "  4. kubectl apply -f k8s/service.yaml"
echo "  5. kubectl apply -f k8s/geoip-cronjob.yaml"
