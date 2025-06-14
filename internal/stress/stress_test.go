package stress

import (
	"bytes"
	"net/http"
	"os"
	"runtime/pprof"
	"strconv"
	"sync"
	"testing"
	"time"
)

type metrics struct {
	total    int
	success  int
	fail     int
	duration time.Duration
	mu       sync.Mutex
}

var m metrics

func TestStressPublishEndpoint(t *testing.T) {
	cpu, _ := os.Create("cpu.prof")
	pprof.StartCPUProfile(cpu)
	defer func() {
		pprof.StopCPUProfile()
		cpu.Close()
	}()

	var wg sync.WaitGroup
	concurrency := 100
	requestsPerWorker := 10

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < requestsPerWorker; j++ {
				payload := `{"worker":` + strconv.Itoa(workerID) + `,"seq":` + strconv.Itoa(j) + `}`
				start := time.Now()
				resp, err := http.Post("http://localhost:8080/publish", "application/json", bytes.NewBufferString(payload))
				elapsed := time.Since(start)

				m.mu.Lock()
				m.total++
				m.duration += elapsed
				if err != nil {
					m.fail++
					m.mu.Unlock()
					t.Errorf("worker %d request %d failed: %v", workerID, j, err)
					continue
				}
				defer resp.Body.Close()
				if resp.StatusCode != http.StatusAccepted {
					m.fail++
					m.mu.Unlock()
					t.Errorf("worker %d request %d got status: %d", workerID, j, resp.StatusCode)
					continue
				}
				m.success++
				m.mu.Unlock()
			}
		}(i)
	}

	wg.Wait()

	t.Logf("Total: %d | Success: %d | Fail: %d", m.total, m.success, m.fail)
	if m.total > 0 {
		t.Logf("Avg response time: %s", m.duration/time.Duration(m.total))
	}

	mem, _ := os.Create("mem.prof")
	pprof.WriteHeapProfile(mem)
	mem.Close()
}
