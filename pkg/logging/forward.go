package logging

import (
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
	bytes, err := f.formatLogsForNewRelic()
	if err != nil {
		return err
	}

	fmt.Println(string(bytes))
	return nil
}

func (f *forwarder) formatLogsForNewRelic() (
	[]byte,
	error,
) {
	lo := &logObject{
		Common: &commonBlock{
			Attributes: make(map[string]string),
		},
		Logs: make([]logBlock, len(f.logs)),
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

	bytes, err := json.Marshal([]logObject{*lo})
	if err != nil {
		return nil, errors.New(LOGS__PAYLOAD_COULD_NOT_BE_CREATED)
	}

	return bytes, nil
}
