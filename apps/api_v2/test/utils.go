package test_utils

import (
	"sync"
	"testing"
	"time"

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

type TestConcurrentProcess struct {
	count 			int
	name  			string
	t	   			*testing.T
	steps 			[]*TestConcurrentStepResult
	Application 	*container.Application
}

type TestConcurrentStepResult struct {
	t	   			*testing.T
	start  			time.Time
	duration    	*time.Duration
	time_percentage float64
	count  			int
	info   			string
	function        func() (int, []byte, error)

	// concurrent test results
	tries 			[]*TryResult
	average_time 	int64
	success_rate 	float64
}

type TryResult struct {
	Status int    `json:"status"`
	Body   []byte `json:"body"`
	Time   int64  `json:"time"`
	Err    error  `json:"-"`
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

func StartConcurrentTest (t *testing.T, name string) TestConcurrentProcess {
	t.Log(name)
	return TestConcurrentProcess{
		count: 0,
		name: name,
		t: t,
		steps: []*TestConcurrentStepResult{},
	}
}

func StartConcurrentTestWithApp (t *testing.T, name string) TestConcurrentProcess {
	appContainer, err := container.BuildApplicationContainer()
	if err != nil {
		t.Fatalf("initialize application container: %s", err)
	}

	t.Log(name)
	return TestConcurrentProcess{
		count: 0,
		name: name,
		t: t,
		steps: []*TestConcurrentStepResult{},
		Application: appContainer,
	}
}

func (test *TestConcurrentProcess) StartConcurrentStep (count int, info string, function func() (int, []byte, error)) {
	test.count = count
	test.t.Logf("[STEP %d] %s", test.count, info)

	test.steps = append(test.steps, &TestConcurrentStepResult{
		t: test.t,
		start: time.Now(),
		time_percentage: 0,
		count: test.count,
		info: info,
		average_time: 0,
		success_rate: 0,
		function: function,
	})
}

func (test *TestConcurrentProcess) RunStep (concurrency int) {
	if len(test.steps) == 0 {
		return
	}
	step := test.steps[len(test.steps)-1]
	if step.function == nil || concurrency <= 0 {
		return
	}

	step.tries = make([]*TryResult, concurrency)

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			start := time.Now()
			status, body, err := step.function()
			elapsed := time.Since(start).Milliseconds()
			step.tries[idx] = &TryResult{
				Status: status,
				Body:   body,
				Time:   elapsed,
				Err:    err,
			}
		}(i)
	}
	wg.Wait()

	duration := time.Since(step.start)
	step.duration = &duration

	var totalTime int64
	successes := 0
	for _, try := range step.tries {
		totalTime += try.Time
		if try.Err == nil && try.Status >= 200 && try.Status < 300 {
			successes++
		}
	}
	step.average_time = totalTime / int64(concurrency)
	step.success_rate = float64(successes) / float64(concurrency) * 100

	test.t.Logf("  |-> [SUCCESS] Finished in %v | runs=%d avg=%dms success=%.2f%%",
		duration, concurrency, step.average_time, step.success_rate)
}

func (test *TestConcurrentProcess) Log (info string) {
	if len(test.steps) > 0 {
		step := test.steps[len(test.steps)-1]
		duration := time.Since(step.start)
		test.t.Logf("  |-> [MESSAGE] %s | from start: %s", info, duration)
	} else {
		test.t.Logf("  |-> [MESSAGE] %s", info)
	}
}

func (test *TestConcurrentProcess) Fail (section string, err error) {
	test.t.Fatalf("  |-> [ERROR] %s: %s", section, err)
}

func (test *TestConcurrentProcess) End () {
	test.t.Log("----------------------------|Critical Process Resume|----------------------------")
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
		test.t.Logf("  - Total: %v ms (%.2f%%) | runs=%d avg=%d ms | success=%.2f%%",
			step.duration.Milliseconds(), step.time_percentage,
			len(step.tries), step.average_time, step.success_rate)
	}
	test.t.Log("------------------------------------------------------------------------")
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