package app

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type job struct {
	source  *Source
	created time.Time
}

func Run(ctx context.Context, sources []*Source, workers, samples int) (int, error) {
	if workers < 1 || workers > 8 {
		return 0, fmt.Errorf("workers must be between 1 and 8")
	}
	if samples < 1 || len(sources) == 0 {
		return 0, fmt.Errorf("samples and sources must be positive")
	}
	jobs := make(chan job, 64)
	var group sync.WaitGroup
	for worker := 0; worker < workers; worker++ {
		group.Add(1)
		go func() {
			defer group.Done()
			for task := range jobs {
				task.source.Produce(task.created, len(jobs))
			}
		}()
	}
	accepted := 0
	for accepted < samples && ctx.Err() == nil {
		task := job{sources[accepted%len(sources)], time.Now()}
		select {
		case jobs <- task:
			accepted++
		case <-ctx.Done():
		}
	}
	close(jobs)
	group.Wait()
	return accepted, nil
}
