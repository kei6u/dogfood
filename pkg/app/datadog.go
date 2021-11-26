package app

import (
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

func startAPM() error {
	tracer.Start(tracer.WithAnalytics(true))
	return profiler.Start(
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,
		),
	)
}

func stopAPM() {
	tracer.Stop()
	profiler.Stop()
}
