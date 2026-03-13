// SPDX-License-Identifier: Apache-2.0

package playwright

import (
	"encoding/xml"
	"fmt"
	"io"
	"time"
)

// JUnit XML types for CI output.

type junitTestSuites struct {
	XMLName xml.Name         `xml:"testsuites"`
	Suites  []junitTestSuite `xml:"testsuite"`
}

type junitTestSuite struct {
	XMLName  xml.Name        `xml:"testsuite"`
	Name     string          `xml:"name,attr"`
	Tests    int             `xml:"tests,attr"`
	Failures int             `xml:"failures,attr"`
	Errors   int             `xml:"errors,attr"`
	Time     string          `xml:"time,attr"`
	Cases    []junitTestCase `xml:"testcase"`
}

type junitTestCase struct {
	XMLName   xml.Name      `xml:"testcase"`
	Name      string        `xml:"name,attr"`
	ClassName string        `xml:"classname,attr"`
	Time      string        `xml:"time,attr"`
	Failure   *junitFailure `xml:"failure,omitempty"`
	Error     *junitError   `xml:"error,omitempty"`
	SystemOut string        `xml:"system-out,omitempty"`
}

type junitFailure struct {
	Message string `xml:"message,attr"`
	Content string `xml:",chardata"`
}

type junitError struct {
	Message string `xml:"message,attr"`
	Content string `xml:",chardata"`
}

// WriteJUnitXML writes test results in JUnit XML format.
func WriteJUnitXML(w io.Writer, result *SuiteResult) error {
	suite := junitTestSuite{
		Name:     result.Name,
		Tests:    len(result.Scripts),
		Failures: 0,
		Errors:   0,
		Time:     formatDuration(result.Duration),
	}

	for _, s := range result.Scripts {
		tc := junitTestCase{
			Name:      s.Name,
			ClassName: "playwright-verify",
			Time:      formatDuration(s.Duration),
		}

		switch s.Status {
		case StatusFail:
			suite.Failures++
			tc.Failure = &junitFailure{
				Message: s.Message,
				Content: s.Output,
			}
		case StatusError:
			suite.Errors++
			tc.Error = &junitError{
				Message: s.Message,
				Content: s.Output,
			}
		case StatusPass:
			if s.Output != "" {
				tc.SystemOut = s.Output
			}
		}

		suite.Cases = append(suite.Cases, tc)
	}

	suites := junitTestSuites{
		Suites: []junitTestSuite{suite},
	}

	fmt.Fprintf(w, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	encoder := xml.NewEncoder(w)
	encoder.Indent("", "  ")
	if err := encoder.Encode(suites); err != nil {
		return fmt.Errorf("encoding JUnit XML: %w", err)
	}
	fmt.Fprintln(w)
	return nil
}

func formatDuration(d time.Duration) string {
	return fmt.Sprintf("%.3f", d.Seconds())
}
