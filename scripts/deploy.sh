#!/bin/bash

helm template "scraper" \
  --create-namespace \
  --namespace "newrelic" \
  "../charts/scraper"

# helm upgrade "scraper" \
#   --install \
#   --wait \
#   --debug \
#   --create-namespace \
#   --namespace "test" \
#   "../charts/scraper"