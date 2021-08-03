package worker

type workerPoolConfig struct {
	maxWorkers       int
	jobQueueCapacity int
}

type opts func(*workerPoolConfig)

func WithMaxWorkers(maxWorkers int) opts {
	return func(cfg *workerPoolConfig) {
		cfg.maxWorkers = maxWorkers
	}
}

func WithJobQueueCapacity(jobQueueCapacity int) opts {
	return func(cfg *workerPoolConfig) {
		cfg.jobQueueCapacity = jobQueueCapacity
	}
}

func buildWorkerPoolConfig(opts ...opts) *workerPoolConfig {
	cfg := &workerPoolConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.jobQueueCapacity <= 0 {
		cfg.jobQueueCapacity = 100
	}
	return cfg
}
