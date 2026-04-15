package events

import "sync/atomic"

// metrics stores process-local counters for basic event pipeline observability.
var metrics = struct {
	publishSuccessTotal       atomic.Uint64
	publishFailureTotal       atomic.Uint64
	workerProcessedTotal      atomic.Uint64
	workerHandlerFailureTotal atomic.Uint64
	workerAckFailureTotal     atomic.Uint64
}{}

func incrementPublishSuccess() {
	metrics.publishSuccessTotal.Add(1)
}

func incrementPublishFailure() {
	metrics.publishFailureTotal.Add(1)
}

func incrementWorkerProcessed() {
	metrics.workerProcessedTotal.Add(1)
}

func incrementWorkerHandlerFailure() {
	metrics.workerHandlerFailureTotal.Add(1)
}

func incrementWorkerAckFailure() {
	metrics.workerAckFailureTotal.Add(1)
}

type MetricsSnapshot struct {
	PublishSuccessTotal       uint64
	PublishFailureTotal       uint64
	WorkerProcessedTotal      uint64
	WorkerHandlerFailureTotal uint64
	WorkerAckFailureTotal     uint64
}

func SnapshotMetrics() MetricsSnapshot {
	return MetricsSnapshot{
		PublishSuccessTotal:       metrics.publishSuccessTotal.Load(),
		PublishFailureTotal:       metrics.publishFailureTotal.Load(),
		WorkerProcessedTotal:      metrics.workerProcessedTotal.Load(),
		WorkerHandlerFailureTotal: metrics.workerHandlerFailureTotal.Load(),
		WorkerAckFailureTotal:     metrics.workerAckFailureTotal.Load(),
	}
}
