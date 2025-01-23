package sealights_ginkgo_agent

import (
	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/ginkgo/v2/types"
	"log"
	"sync"
	"time"
)

type testResult struct {
	name   string
	start  time.Time
	end    time.Time
	result string
}

func newTestResult() *testResult {
	return &testResult{}
}

func (t *testResult) setName(value string) *testResult {
	t.name = value
	return t
}

func (t *testResult) setStart(value time.Time) *testResult {
	t.start = time.Now()
	return t
}

func (t *testResult) setEnd(value time.Time) *testResult {
	t.end = time.Now()
	return t
}

func (t *testResult) setResult(value string) *testResult {
	t.result = value
	return t
}

type sealightsAgent struct {
	mu          sync.RWMutex
	testsToSkip map[string]bool
	testsReport []*testResult
	isDisabled  bool
}

func newSealightsAgent() *sealightsAgent {
	return &sealightsAgent{
		mu:          sync.RWMutex{},
		testsToSkip: make(map[string]bool),
	}
}

func (s *sealightsAgent) init() error {
	log.Println("Calling Sealights backend to retrieve the test data")
	time.Sleep(2 * time.Second)
	return nil
}

func (s *sealightsAgent) isSkipped(testName string) bool {
	return s.testsToSkip[testName]
}

func (s *sealightsAgent) addTestResult(test *testResult) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.testsReport = append(s.testsReport, test)
}

var sealights *sealightsAgent

var _ = BeforeEach(
	func() {
		if sealights.isDisabled {
			return
		}
		log.Println("Sealights hook BeforeEach checking if the test should be skipped")
		currentTest := CurrentSpecReport().LeafNodeText
		if sealights.isSkipped(currentTest) {
			log.Printf("Sealights skipping test %s\n", currentTest)
			Skip("Skipping test")
		}
	})

var _ = AfterEach(
	func() {
		if sealights.isDisabled {
			return
		}
		log.Println("Sealights hook AfterEach adding the test result to the report")
		testReport := CurrentSpecReport()
		testResult := newTestResult().
			setName(testReport.LeafNodeText).
			setStart(testReport.StartTime).
			setEnd(testReport.EndTime)
		if testReport.State.Is(types.SpecStatePassed) {

			testResult.setResult("Passed")
		}
		if testReport.State.Is(types.SpecStateFailed) {

			testResult.setResult("Failed")
		}
		if testReport.State.Is(types.SpecStateSkipped) {

			testResult.setResult("Skipped")
		}
		sealights.addTestResult(testResult)
	})

var _ = ReportAfterSuite("Sealights Reporting", func(ctx SpecContext, r Report) {
	if sealights.isDisabled {
		return
	}
	log.Println("Sealights hook ReportAfterSuite reporting the test results")

}, NodeTimeout(1*time.Minute))

var _ = ReportAfterSuite("Sealights Reporting", func(ctx SpecContext, r Report) {
	if sealights.isDisabled {
		return
	}
	log.Println("Sealights hook ReportAfterSuite reporting the test results")

}, NodeTimeout(1*time.Minute))

func RegisterSealightAgent() {
	log.Println("This is the start of Sealights hook init, we are adding a place holder for calling the Sealights backend to retrieve the test data")
	sealights = newSealightsAgent()
	if err := sealights.init(); err != nil {
		sealights.isDisabled = true
		log.Println("Failed to init sealights agent")
	}
}
