# New Relic Kubernetes Endpoint Scraper

This repository is meant to scrape the values which are exposed
by the custom endpoints of your applications running in any
pod.

The scraper is set up to run as a cron job every minute.
It will be triggered automatically by Kubernetes and fetch
all the endpoints that you have defined in the configuration.

## Configuration

In order to let the scraper know which endpoints to check, the
following configuration options are provided within the
[values.yaml](/charts/scraper/values.yaml) file.

**Example:**

```yaml
config:
  # Mount path for the container
  mountPath: /etc/config
  # Configuration data itself
  data:
    newrelic:
      # New Relic account ID
      accountId: "<YOUR_NEWRELIC_ACCOUNT_ID>"
      # New Relic license key
      licenseKey: "<YOUR_NEWRELIC_LICENSE_KEY>"
      # Log level can be: DEBUG, ERROR
      logLevel: ERROR
    # Endpoints which are to be scraped
    # - type
    #   - kvp: key value pair
    #   - json: json
    endpoints:
      - type: "kvp"
        name: "my_endpoint_1"
        url: "http://<SERVICE>.<NAMESPACE>.svc.cluster.local:<PORT>/<ENDPOINT>"
      - type: "kvp"
        name: "my_endpoint_2"
        url: "http://<IP_ADDRESS_OF_POD>:<PORT>/<ENDPOINT>"
```

## Scraping

Currently only endpoints which are exposing key-value pairs (`kvp`)
can be scraped and formatted. The key and the value should be separated
by a semicolon (`:`).

## Building your Docker image

If you would like to make your changes to the code and create your
own image, refer to the [`build.sh`](/scripts/build.sh). You can
build to `amd` or `arm` processors.

## Deploying the Helm chart

In order to deploy the solution, refer to the [`deploy.sh`](/scripts/deploy.sh).
You can build to `amd` or `arm` processors.

## Querying your data in New Relic

The scraped endpoints will be forwarded to New Relic as custom
events which will have the naem `K8sCustomEndpointScrapeSample`.

You can query them as follows:

```
FROM K8sCustomEndpointScrapeSample SELECT *
```
