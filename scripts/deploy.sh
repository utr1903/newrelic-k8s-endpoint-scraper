#!/bin/bash

# Get commandline arguments
while (( "$#" )); do
  case "$1" in
    --arm)
      arm="true"
      shift
      ;;
    *)
      shift
      ;;
  esac
done

# ARM deployment
if [[ $arm == "true" ]]; then
  helm upgrade "scraper" \
    --install \
    --wait \
    --debug \
    --create-namespace \
    --namespace "newrelic" \
    --set scraper.image.repository="uturkarslan/newrelic-kubernetes-endpoint-scraper-arm" \
    --set scraper.image.tag="1.1.0" \
    --set scraper.config.endpoints[0].type="kvp" \
    --set scraper.config.endpoints[0].name="TestEndpointSample" \
    --set scraper.config.endpoints[0].url="http://server.test.svc.cluster.local:8080/kvp" \
    --set scraper.config.newrelic.logLevel="DEBUG" \
    --set scraper.config.newrelic.accountId=$NEWRELIC_ACCOUNT_ID \
    --set scraper.config.newrelic.licenseKey=$NEWRELIC_LICENSE_KEY \
    "../charts/scraper"

# AMD deployment
else
helm upgrade "scraper" \
  --install \
  --wait \
  --debug \
  --create-namespace \
  --namespace "newrelic" \
  --set scraper.image.repository="uturkarslan/newrelic-kubernetes-endpoint-scraper-amd" \
  --set scraper.image.tag="1.1.0" \
  --set scraper.config.endpoints[0].type="kvp" \
  --set scraper.config.endpoints[0].name="TestEndpointSample" \
  --set scraper.config.endpoints[0].url="http://server.test.svc.cluster.local:8080/kvp" \
  --set scraper.config.newrelic.logLevel="DEBUG" \
  --set scraper.config.newrelic.accountId=$NEWRELIC_ACCOUNT_ID \
  --set scraper.config.newrelic.licenseKey=$NEWRELIC_LICENSE_KEY \
  "../charts/scraper"
fi
