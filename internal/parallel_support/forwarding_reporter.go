package parallel_support

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/onsi/ginkgo/internal"

	"github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/types"
)

type ConfigAndSummary struct {
	Config  config.GinkgoConfigType `json:"config"`
	Summary types.SuiteSummary      `json:"suite-summary"`
}

/*
The ForwardingReporter is a Ginkgo reporter that forwards information to
a Ginkgo remote server.

When streaming parallel test output, this repoter is automatically installed by Ginkgo.

This is accomplished by passing in the GINKGO_REMOTE_REPORTING_SERVER environment variable to `go test`, the Ginkgo test runner
detects this environment variable (which should contain the host of the server) and automatically installs a ForwardingReporter
in place of Ginkgo's DefaultReporter.
*/

type ForwardingReporter struct {
	serverHost string
}

func NewForwardingReporter(config config.DefaultReporterConfigType, serverHost string, ginkgoWriter *internal.Writer) *ForwardingReporter {
	reporter := &ForwardingReporter{
		serverHost: serverHost,
	}

	return reporter
}

func (reporter *ForwardingReporter) post(path string, data interface{}) {
	encoded, _ := json.Marshal(data)
	buffer := bytes.NewBuffer(encoded)
	http.Post(reporter.serverHost+path, "application/json", buffer)
}

func (reporter *ForwardingReporter) SpecSuiteWillBegin(conf config.GinkgoConfigType, summary types.SuiteSummary) {
	reporter.post("/SpecSuiteWillBegin", ConfigAndSummary{Config: conf, Summary: summary})
}

func (reporter *ForwardingReporter) WillRun(report types.SpecReport) {
}

func (reporter *ForwardingReporter) DidRun(report types.SpecReport) {
	reporter.post("/DidRun", report)
}

func (reporter *ForwardingReporter) SpecSuiteDidEnd(summary types.SuiteSummary) {
	reporter.post("/SpecSuiteDidEnd", summary)
}