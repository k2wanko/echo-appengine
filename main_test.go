package appengine

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"

	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
)

var (
	aetInstWg sync.WaitGroup
	aetInstMu sync.Mutex
	aetInst   = make(map[*testing.T]aetest.Instance)

	aetestInstance aetest.Instance
)

func TestMain(m *testing.M) {
	aetestInstance, err := aetest.NewInstance(nil)
	if err != nil {
		panic(fmt.Sprintf("aetestInstace: %v", err))
	}

	code := m.Run()

	// cleanup
	aetestInstance.Close()
	cleanupTests()

	os.Exit(code)
}

func newTestContext(r *http.Request) echo.Context {
	e := echo.New()
	if r == nil {
		r, _ = aetestInstance.NewRequest("GET", "/", nil)
	}
	w := httptest.NewRecorder()
	c := e.NewContext(standard.NewRequest(r, e.Logger()), standard.NewResponse(w, e.Logger()))
	c.SetContext(appengine.WithContext(c.Context(), r))
	return c
}

func newTestRequest(t *testing.T, method, url string, body io.Reader) *http.Request {
	req, err := aetInstance(t).NewRequest(method, url, body)
	if err != nil {
		t.Fatalf("newTestRequest(%q, %q): %v", method, url, err)
	}
	return req
}

func resetTestState(t *testing.T) {
	aetInstMu.Lock()
	defer aetInstMu.Unlock()
	inst, ok := aetInst[t]
	if !ok {
		return
	}
	aetInstWg.Add(1)
	go func() {
		if err := inst.Close(); err != nil {
			t.Logf("resetTestState: %v", err)
		}
		aetInstWg.Done()
	}()
	delete(aetInst, t)
}

func cleanupTests() {
	aetInstMu.Lock()
	tts := make([]*testing.T, 0, len(aetInst))
	for t := range aetInst {
		tts = append(tts, t)
	}
	aetInstMu.Unlock()
	for _, t := range tts {
		resetTestState(t)
	}
	aetInstWg.Wait()
}

func aetInstance(t *testing.T) aetest.Instance {
	aetInstMu.Lock()
	defer aetInstMu.Unlock()
	if inst, ok := aetInst[t]; ok {
		return inst
	}
	inst, err := aetest.NewInstance(nil)
	if err != nil {
		t.Fatalf("aetest.NewInstance: %v", err)
	}
	aetInst[t] = inst
	return inst
}
