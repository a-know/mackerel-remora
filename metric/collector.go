package metric

import (
	"context"
	"sync"
)

type collector struct {
	generators []Generator
}

func newCollector(generators []Generator) *collector {
	return &collector{
		generators: generators,
	}
}

func (c *collector) collect(ctx context.Context) (Values, error) {
	var wg sync.WaitGroup
	values := make(Values)
	mu := new(sync.Mutex)
	for _, g := range c.generators {
		wg.Add(1)
		go func(g Generator) {
			defer wg.Done()
			vs, err := g.Generate(ctx)
			if err != nil {
				logger.Errorf("%s", err)
				return
			}
			mu.Lock()
			defer mu.Unlock()
			for key, value := range vs {
				values[key] = value
			}
		}(g)
	}
	wg.Wait()
	return values, nil
}
