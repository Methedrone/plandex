package performance

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"runtime"
	"strings"
	"sync"
	"time"
)

// WorkerPool provides a pool of goroutines for handling HTTP requests
type WorkerPool struct {
	workerCount   int
	jobQueue      chan Job
	workers       []*Worker
	quit          chan bool
	wg            sync.WaitGroup
	metrics       *PoolMetrics
	maxQueueSize  int
}

// Job represents a unit of work to be processed by the worker pool
type Job struct {
	ID       string
	Handler  func() error
	Priority int
	Created  time.Time
	ctx      context.Context
}

// Worker represents a single worker goroutine
type Worker struct {
	id         int
	jobQueue   chan Job
	quit       chan bool
	wg         *sync.WaitGroup
	metrics    *PoolMetrics
}

// PoolMetrics tracks performance metrics for the worker pool
type PoolMetrics struct {
	mu                 sync.RWMutex
	JobsProcessed      int64
	JobsQueued         int64
	JobsFailed         int64
	AverageProcessTime time.Duration
	ActiveWorkers      int32
	QueueLength        int32
	TotalProcessTime   time.Duration
}

// ObjectPools provides efficient object reuse for common allocations
type ObjectPools struct {
	JSONEncoderPool   sync.Pool
	JSONDecoderPool   sync.Pool
	ByteBufferPool    sync.Pool
	StringBuilderPool sync.Pool
	RequestPool       sync.Pool
	ResponsePool      sync.Pool
}

// Global pool instances optimized for MacBook 2012 constraints
var (
	GlobalWorkerPool *WorkerPool
	GlobalObjectPools *ObjectPools
	poolOnce sync.Once
)

// InitializePools sets up optimized pools for MacBook 2012 environment
func InitializePools() {
	poolOnce.Do(func() {
		// Conservative settings for MacBook 2012 (4-8GB RAM, dual/quad-core)
		workerCount := runtime.GOMAXPROCS(0)
		if workerCount > 4 {
			workerCount = 4 // Limit for MacBook 2012
		}
		
		maxQueueSize := workerCount * 10 // Conservative queue size
		
		GlobalWorkerPool = NewWorkerPool(workerCount, maxQueueSize)
		GlobalObjectPools = NewObjectPools()
		
		// Start the worker pool
		GlobalWorkerPool.Start()
	})
}

// NewWorkerPool creates a new worker pool with specified parameters
func NewWorkerPool(workerCount, maxQueueSize int) *WorkerPool {
	return &WorkerPool{
		workerCount:  workerCount,
		jobQueue:     make(chan Job, maxQueueSize),
		workers:      make([]*Worker, workerCount),
		quit:         make(chan bool),
		maxQueueSize: maxQueueSize,
		metrics: &PoolMetrics{
			AverageProcessTime: 0,
			ActiveWorkers:      0,
			QueueLength:        0,
		},
	}
}

// Start initializes and starts all workers in the pool
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workerCount; i++ {
		worker := &Worker{
			id:       i,
			jobQueue: wp.jobQueue,
			quit:     make(chan bool),
			wg:       &wp.wg,
			metrics:  wp.metrics,
		}
		wp.workers[i] = worker
		wp.wg.Add(1)
		go worker.Start()
	}
}

// Submit adds a job to the worker pool queue
func (wp *WorkerPool) Submit(job Job) error {
	wp.metrics.mu.Lock()
	wp.metrics.JobsQueued++
	wp.metrics.QueueLength = int32(len(wp.jobQueue))
	wp.metrics.mu.Unlock()

	select {
	case wp.jobQueue <- job:
		return nil
	case <-time.After(5 * time.Second):
		return ErrPoolTimeout
	}
}

// Stop gracefully shuts down the worker pool
func (wp *WorkerPool) Stop() {
	close(wp.quit)
	close(wp.jobQueue)
	wp.wg.Wait()
}

// GetMetrics returns current pool performance metrics
func (wp *WorkerPool) GetMetrics() PoolMetrics {
	wp.metrics.mu.RLock()
	defer wp.metrics.mu.RUnlock()
	
	// Copy metrics safely without copying the mutex
	return PoolMetrics{
		JobsProcessed:      wp.metrics.JobsProcessed,
		JobsQueued:         wp.metrics.JobsQueued,
		JobsFailed:         wp.metrics.JobsFailed,
		AverageProcessTime: wp.metrics.AverageProcessTime,
		ActiveWorkers:      wp.metrics.ActiveWorkers,
		QueueLength:        int32(len(wp.jobQueue)),
		TotalProcessTime:   wp.metrics.TotalProcessTime,
	}
}

// Start begins the worker's job processing loop
func (w *Worker) Start() {
	defer w.wg.Done()
	
	for {
		select {
		case job := <-w.jobQueue:
			w.processJob(job)
		case <-w.quit:
			return
		}
	}
}

// processJob handles a single job and updates metrics
func (w *Worker) processJob(job Job) {
	start := time.Now()
	
	w.metrics.mu.Lock()
	w.metrics.ActiveWorkers++
	w.metrics.mu.Unlock()
	
	defer func() {
		duration := time.Since(start)
		
		w.metrics.mu.Lock()
		w.metrics.ActiveWorkers--
		w.metrics.JobsProcessed++
		w.metrics.TotalProcessTime += duration
		w.metrics.AverageProcessTime = w.metrics.TotalProcessTime / time.Duration(w.metrics.JobsProcessed)
		w.metrics.mu.Unlock()
		
		// Recover from any panics in job processing
		if r := recover(); r != nil {
			w.metrics.mu.Lock()
			w.metrics.JobsFailed++
			w.metrics.mu.Unlock()
		}
	}()
	
	// Execute the job with timeout context
	ctx, cancel := context.WithTimeout(job.ctx, 30*time.Second)
	defer cancel()
	
	done := make(chan error, 1)
	go func() {
		done <- job.Handler()
	}()
	
	select {
	case err := <-done:
		if err != nil {
			w.metrics.mu.Lock()
			w.metrics.JobsFailed++
			w.metrics.mu.Unlock()
		}
	case <-ctx.Done():
		w.metrics.mu.Lock()
		w.metrics.JobsFailed++
		w.metrics.mu.Unlock()
	}
}

// NewObjectPools creates optimized object pools for memory efficiency
func NewObjectPools() *ObjectPools {
	return &ObjectPools{
		JSONEncoderPool: sync.Pool{
			New: func() interface{} {
				return json.NewEncoder(&bytes.Buffer{})
			},
		},
		JSONDecoderPool: sync.Pool{
			New: func() interface{} {
				return json.NewDecoder(strings.NewReader(""))
			},
		},
		ByteBufferPool: sync.Pool{
			New: func() interface{} {
				// Pre-allocate with reasonable capacity for API responses
				return bytes.NewBuffer(make([]byte, 0, 1024))
			},
		},
		StringBuilderPool: sync.Pool{
			New: func() interface{} {
				var sb strings.Builder
				sb.Grow(512) // Pre-allocate capacity
				return &sb
			},
		},
		RequestPool: sync.Pool{
			New: func() interface{} {
				return make(map[string]interface{})
			},
		},
		ResponsePool: sync.Pool{
			New: func() interface{} {
				return make(map[string]interface{})
			},
		},
	}
}

// GetJSONEncoder retrieves a JSON encoder from the pool
func (op *ObjectPools) GetJSONEncoder(buf *bytes.Buffer) *json.Encoder {
	// Create new encoder each time since Reset was added in Go 1.21+
	return json.NewEncoder(buf)
}

// PutJSONEncoder returns a JSON encoder to the pool (no-op for now)
func (op *ObjectPools) PutJSONEncoder(encoder *json.Encoder) {
	// No-op since we create new encoders
}

// GetJSONDecoder retrieves a JSON decoder from the pool
func (op *ObjectPools) GetJSONDecoder(reader *strings.Reader) *json.Decoder {
	// Create new decoder each time since Reset was added in Go 1.21+
	return json.NewDecoder(reader)
}

// PutJSONDecoder returns a JSON decoder to the pool (no-op for now)
func (op *ObjectPools) PutJSONDecoder(decoder *json.Decoder) {
	// No-op since we create new decoders
}

// GetBuffer retrieves a byte buffer from the pool
func (op *ObjectPools) GetBuffer() *bytes.Buffer {
	buf := op.ByteBufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

// PutBuffer returns a byte buffer to the pool
func (op *ObjectPools) PutBuffer(buf *bytes.Buffer) {
	// Only reuse buffers under 64KB to prevent memory bloat
	if buf.Cap() < 64*1024 {
		op.ByteBufferPool.Put(buf)
	}
}

// GetStringBuilder retrieves a string builder from the pool
func (op *ObjectPools) GetStringBuilder() *strings.Builder {
	sb := op.StringBuilderPool.Get().(*strings.Builder)
	sb.Reset()
	return sb
}

// PutStringBuilder returns a string builder to the pool
func (op *ObjectPools) PutStringBuilder(sb *strings.Builder) {
	// Only reuse builders under 16KB to prevent memory bloat
	if sb.Cap() < 16*1024 {
		op.StringBuilderPool.Put(sb)
	}
}

// GetRequestMap retrieves a request map from the pool
func (op *ObjectPools) GetRequestMap() map[string]interface{} {
	reqMap := op.RequestPool.Get().(map[string]interface{})
	// Clear the map
	for k := range reqMap {
		delete(reqMap, k)
	}
	return reqMap
}

// PutRequestMap returns a request map to the pool
func (op *ObjectPools) PutRequestMap(reqMap map[string]interface{}) {
	// Only reuse maps with reasonable size
	if len(reqMap) < 100 {
		op.RequestPool.Put(reqMap)
	}
}

// GetResponseMap retrieves a response map from the pool
func (op *ObjectPools) GetResponseMap() map[string]interface{} {
	respMap := op.ResponsePool.Get().(map[string]interface{})
	// Clear the map
	for k := range respMap {
		delete(respMap, k)
	}
	return respMap
}

// PutResponseMap returns a response map to the pool
func (op *ObjectPools) PutResponseMap(respMap map[string]interface{}) {
	// Only reuse maps with reasonable size
	if len(respMap) < 100 {
		op.ResponsePool.Put(respMap)
	}
}

// Cleanup performs memory cleanup and pool maintenance
func (op *ObjectPools) Cleanup() {
	// Force garbage collection to free unused objects
	runtime.GC()
}

// Performance error types
var (
	ErrPoolTimeout = errors.New("worker pool timeout")
	ErrPoolFull    = errors.New("worker pool queue full")
)

// SubmitHTTPJob submits an HTTP request handling job to the worker pool
func SubmitHTTPJob(ctx context.Context, id string, handler func() error) error {
	if GlobalWorkerPool == nil {
		InitializePools()
	}
	
	job := Job{
		ID:       id,
		Handler:  handler,
		Priority: 1,
		Created:  time.Now(),
		ctx:      ctx,
	}
	
	return GlobalWorkerPool.Submit(job)
}