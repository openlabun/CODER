package test_utils

import (
	"time"
	"testing"

	container "github.com/openlabun/CODER/apps/api_v2/internal/application/container"
)

type HTTPTestResult struct {
	Status string `json:"status"`
	Body  []byte  `json:"body"`
	Time  int64   `json:"time"`
}

type TestProcess struct {
	count int
	name  string
	t	   *testing.T
	steps []*TestStepResult
	Application *container.Application
}

type TestStepResult struct {
	t	   *testing.T
	start  time.Time
	duration    *time.Duration
	time_percentage float64
	count  int
	info   string
}

func StartTest (t *testing.T, name string) TestProcess {
	t.Log(name)
	return TestProcess{
		count: 0,
		name: name,
		t: t,
		steps: []*TestStepResult{},
	}
}

func StartTestWithApp (t *testing.T, name string) TestProcess {
	appContainer, err := container.BuildApplicationContainer()
	if err != nil {
		t.Fatalf("initialize application container: %s", err)
	}

	t.Log(name)
	return TestProcess{
		count: 0,
		name: name,
		t: t,
		steps: []*TestStepResult{},
		Application: appContainer,
	}
}

func (test *TestProcess) StartStep (info string) {
	test.count = test.count + 1
	test.t.Logf("[STEP %d] %s", test.count, info)

	test.steps = append(test.steps, &TestStepResult{
		t: test.t,
		start: time.Now(),
		count: test.count,
		info: info,
	})
}

func (test *TestProcess) Log (info string) {
	if len(test.steps) > 0 {
		step := test.steps[len(test.steps)-1]
		duration := time.Since(step.start)
		test.t.Logf("  |-> [MESSAGE] %s | from start: %s", info, duration)
	} else {
		test.t.Logf("  |-> [MESSAGE] %s", info)
	}
}

func (test *TestProcess) Fail (section string, err error) {
	test.t.Fatalf("  |-> [ERROR] %s: %s", section, err)
}

func (test *TestProcess) EndStep () {
	if len(test.steps) > 0 {
		step := test.steps[len(test.steps)-1]
		duration := time.Since(step.start)

		step.duration = &duration
		test.t.Logf("  |-> [SUCCESS] Finished in %v", step.duration)
	}
}

func (test *TestProcess) End () {
	test.t.Log("----------------------------|Process Resume|----------------------------")
	total_time := 0
	for _, step := range test.steps {
		if step.duration != nil {
			total_time = total_time + int(step.duration.Milliseconds())
		}
	}

	for _, step := range test.steps {
		if step == nil || step.duration == nil {
			continue
		}

		step.time_percentage = (float64(step.duration.Milliseconds()) / float64(total_time)) * 100
		test.t.Logf("[STEP %d] %s", step.count, step.info)
		test.t.Logf("  - Time taken: %v ms (%.2f%%)", step.duration.Milliseconds(), step.time_percentage)
	}
	test.t.Log("------------------------------------------------------------------------")
}