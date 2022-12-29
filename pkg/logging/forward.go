package logging

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type commonBlock struct {
	Attributes map[string]string `json:"attributes"`
}

type logBlock struct {
	Timestamp  int64             `json:"timestamp"`
	Message    string            `json:"message"`
	Attributes map[string]string `json:"attributes"`
}

type logObject struct {
	Common *commonBlock `json:"common"`
	Logs   []logBlock   `json:"logs"`
}

type forwarder struct {
	levels []logrus.Level
	logs   []logrus.Entry

	client       *http.Client
	licenseKey   string
	logsEndpoint string
}

func newForwarder(
	levels []logrus.Level,
	licenseKey string,
	logsEndpoint string,
) *forwarder {

	// Create HTTP client
	client := http.Client{Timeout: time.Duration(30 * time.Second)}

	return &forwarder{
		levels:       levels,
		logs:         make([]logrus.Entry, 0),
		client:       &client,
		licenseKey:   licenseKey,
		logsEndpoint: logsEndpoint,
	}
}

func (f *forwarder) Levels() []logrus.Level {
	return f.levels
}

func (f *forwarder) Fire(e *logrus.Entry) error {
	copy := *e
	f.logs = append(f.logs, copy)
	return nil
}

func (f *forwarder) flush() error {
	// Create New Relic logs
	nrLogs := f.createNewRelicLogs()

	// Flush data to New Relic
	return f.sendToNewRelic(nrLogs)
}

func (f *forwarder) createNewRelicLogs() []logObject {
	lo := &logObject{
		Common: &commonBlock{
			Attributes: make(map[string]string),
		},
		Logs: make([]logBlock, 0, len(f.logs)),
	}

	// Create common block
	for key, val := range getCommonAttributes() {
		lo.Common.Attributes[key] = val
	}

	// Create logs block
	for _, log := range f.logs {
		logBlock := logBlock{
			Timestamp:  log.Time.UnixMicro(),
			Message:    log.Message,
			Attributes: make(map[string]string),
		}

		for key, val := range log.Data {
			logBlock.Attributes[key] = fmt.Sprintf("%v", val)
		}
		lo.Logs = append(lo.Logs, logBlock)
	}

	return []logObject{*lo}
}

func (f *forwarder) sendToNewRelic(
	nrLogs []logObject,
) error {

	// Create payload
	json, err := json.Marshal(nrLogs)
	if err != nil {
		return errors.New(LOGS__PAYLOAD_COULD_NOT_BE_CREATED)
	}
	payload := bytes.NewReader(json)

	// Create HTTP request
	req, err := http.NewRequest(http.MethodPost, f.logsEndpoint, payload)
	if err != nil {
		return errors.New(LOGS__HTTP_REQUEST_COULD_NOT_BE_CREATED)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Api-Key", f.licenseKey)

	// Perform HTTP request
	res, err := f.client.Do(req)
	if err != nil {
		return errors.New(LOGS__HTTP_REQUEST_HAS_FAILED)
	}
	defer res.Body.Close()

	// Check if call was successful
	if res.StatusCode != http.StatusAccepted {
		return errors.New(LOGS__NEW_RELIC_RETURNED_NOT_OK_STATUS)
	}

	return nil
}
