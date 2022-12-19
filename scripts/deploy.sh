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
    --set image.repository="uturkarslan/newrelic-kubernetes-endpoint-scraper-arm" \
    --set image.tag="1.0.0" \
    --set config.data.endpoints[0].type="kvp" \
    --set config.data.endpoints[0].name="test" \
    --set config.data.endpoints[0].url="http://server.test.svc.cluster.local:8080/kvp" \
    --set config.data.newrelic.logLevel="DEBUG" \
    --set config.data.newrelic.accountId=$NEWRELIC_ACCOUNT_ID \
    --set config.data.newrelic.licenseKey=$NEWRELIC_LICENSE_KEY \
    "../charts/scraper"

# AMD deployment
else
helm upgrade "scraper" \
  --install \
  --wait \
  --debug \
  --create-namespace \
  --namespace "newrelic" \
  --set image.repository="uturkarslan/newrelic-kubernetes-endpoint-scraper-amd" \
  --set image.tag="1.0.0" \
  --set config.data.endpoints[0].type="kvp" \
  --set config.data.endpoints[0].name="test" \
  --set config.data.endpoints[0].url="http://server.test.svc.cluster.local:8080/kvp" \
  --set config.data.newrelic.logLevel="DEBUG" \
  --set config.data.newrelic.accountId=$NEWRELIC_ACCOUNT_ID \
  --set config.data.newrelic.licenseKey=$NEWRELIC_LICENSE_KEY \
  "../charts/scraper"
fi
