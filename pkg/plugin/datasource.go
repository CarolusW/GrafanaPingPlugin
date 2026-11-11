package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

const (
	defaultInterval   = 10 * time.Second
	defaultTimeout    = 3 * time.Second
	defaultHistoryLen = 5000
)

type point struct {
	ts    time.Time
	value float64
}

type target struct {
	host     string
	interval time.Duration
	points   []point
}

// Datasource implements backend.QueryDataHandler and backend.CheckHealthHandler.
// It keeps a background scheduler that pings every registered host and stores
// the results in memory, so panels get a proper time series (x = check time,
// y = round-trip time in ms).
type Datasource struct {
	mu                      sync.Mutex
	targets                 map[string]*target
	timeout                 time.Duration
	defaultIntervalDuration time.Duration
	history                 int
	stopOnce                sync.Once
	stopCh                  chan struct{}
}

type queryModel struct {
	Host     string `json:"host"`
	Interval string `json:"interval"`
}

type configSettings struct {
	Interval string `json:"defaultInterval"`
	Timeout  string `json:"timeout"`
}

// NewDatasource creates the datasource and starts the ping scheduler.
func NewDatasource() *Datasource {
	d := &Datasource{
		targets:                 map[string]*target{},
		timeout:                 defaultTimeout,
		defaultIntervalDuration: defaultInterval,
		history:                 defaultHistoryLen,
		stopCh:                  make(chan struct{}),
	}
	go d.scheduler()
	return d
}

// Dispose stops the scheduler.
func (d *Datasource) Dispose() {
	d.stopOnce.Do(func() { close(d.stopCh) })
}

func (d *Datasource) scheduler() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-d.stopCh:
			return
		case <-ticker.C:
			d.runDuePings(time.Now())
		}
	}
}

func (d *Datasource) runDuePings(now time.Time) {
	type job struct {
		t *target
	}

	d.mu.Lock()
	jobs := make([]job, 0, len(d.targets))
	for _, t := range d.targets {
		if len(t.points) == 0 || now.Sub(t.points[len(t.points)-1].ts) >= t.interval {
			jobs = append(jobs, job{t})
		}
	}
	timeout := d.timeout
	d.mu.Unlock()

	var wg sync.WaitGroup
	for _, j := range jobs {
		wg.Add(1)
		go func(t *target) {
			defer wg.Done()
			lat, err := Ping(t.host, timeout)
			if err != nil {
				log.DefaultLogger.Debug("ping failed", "host", t.host, "error", err.Error())
				d.mu.Lock()
				t.points = append(t.points, point{ts: now, value: -1})
				d.trimLocked(t)
				d.mu.Unlock()
				return
			}
			d.mu.Lock()
			t.points = append(t.points, point{ts: now, value: lat})
			d.trimLocked(t)
			d.mu.Unlock()
		}(j.t)
	}
	wg.Wait()
}

func (d *Datasource) trimLocked(t *target) {
	if len(t.points) > d.history {
		t.points = t.points[len(t.points)-d.history:]
	}
}

// register ensures a host is being pinged on its interval. Returns the key.
func (d *Datasource) register(host string, interval time.Duration) {
	if host == "" {
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	key := host + "\x00" + interval.String()
	t, ok := d.targets[key]
	if !ok {
		t = &target{host: host, interval: interval}
		d.targets[key] = t
		log.DefaultLogger.Info("ping target registered", "host", host, "interval", interval.String())
	}
}

// QueryData returns all stored points for the queried host within the frame
// time range.
func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	resp := backend.NewQueryDataResponse()

	for _, q := range req.Queries {
		if settings := req.PluginContext.DataSourceInstanceSettings; settings != nil {
			d.applySettings(settings.JSONData)
		}
		resp.Responses[q.RefID] = d.query(ctx, q)
	}
	return resp, nil
}

func (d *Datasource) query(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	var model queryModel
	if err := json.Unmarshal(q.JSON, &model); err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("invalid query JSON: %v", err))
	}
	if model.Host == "" {
		return backend.ErrDataResponse(backend.StatusBadRequest, "host is required")
	}

	interval := d.defaultInterval()
	if model.Interval != "" {
		parsed, err := time.ParseDuration(model.Interval)
		if err != nil || parsed <= 0 {
			return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("invalid interval %q", model.Interval))
		}
		interval = parsed
	}

	d.register(model.Host, interval)

	d.mu.Lock()
	key := model.Host + "\x00" + interval.String()
	t := d.targets[key]
	var pts []point
	if t != nil {
		from := q.TimeRange.From
		to := q.TimeRange.To
		for _, p := range t.points {
			if !p.ts.Before(from) && !p.ts.After(to) && p.value >= 0 {
				pts = append(pts, p)
			}
		}
	}
	d.mu.Unlock()

	sort.Slice(pts, func(i, j int) bool { return pts[i].ts.Before(pts[j].ts) })

	frame := data.NewFrame(
		model.Host,
		data.NewField("time", nil, make([]time.Time, 0, len(pts))),
		data.NewField("rtt_ms", nil, make([]*float64, 0, len(pts))),
	)
	for _, p := range pts {
		v := p.value
		frame.AppendRow(p.ts, &v)
	}

	resp := backend.DataResponse{}
	resp.Frames = append(resp.Frames, frame)
	return resp
}

// applySettings reads datasource configuration (default interval, timeout).
func (d *Datasource) applySettings(raw json.RawMessage) {
	var s configSettings
	if err := json.Unmarshal(raw, &s); err != nil {
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if v, err := time.ParseDuration(s.Timeout); err == nil && v > 0 {
		d.timeout = v
	}
	if v, err := time.ParseDuration(s.Interval); err == nil && v > 0 {
		d.defaultIntervalDuration = v
	}
}

func (d *Datasource) defaultInterval() time.Duration {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.defaultIntervalDuration
}

// CheckHealth pings localhost to verify the system ping binary works.
func (d *Datasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	if settings := req.PluginContext.DataSourceInstanceSettings; settings != nil {
		d.applySettings(settings.JSONData)
	}

	d.mu.Lock()
	timeout := d.timeout
	d.mu.Unlock()

	if _, err := Ping("127.0.0.1", timeout); err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: fmt.Sprintf("ping check failed: %v", err),
		}, nil
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Ping datasource is running",
	}, nil
}
