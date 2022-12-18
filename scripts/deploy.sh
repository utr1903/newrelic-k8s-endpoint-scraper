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
  "../charts/scraper"
fi
