package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// Job represents a scheduled job
type Job struct {
	Name     string
	Schedule string // cron-like: "daily@02:00", "hourly", "every:5m"
	Handler  func(ctx context.Context) error
	Enabled  bool
}

// Scheduler manages scheduled jobs
type Scheduler struct {
	jobs    []*jobRunner
	logger  zerolog.Logger
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
	running bool
	mu      sync.Mutex
}

type jobRunner struct {
	job      *Job
	ticker   *time.Ticker
	stopChan chan struct{}
}

// New creates a new Scheduler
func New(logger zerolog.Logger) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		jobs:   make([]*jobRunner, 0),
		logger: logger,
		ctx:    ctx,
		cancel: cancel,
	}
}

// AddJob adds a job to the scheduler
func (s *Scheduler) AddJob(job *Job) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !job.Enabled {
		s.logger.Info().Str("job", job.Name).Msg("Job is disabled, skipping")
		return
	}

	interval, err := parseSchedule(job.Schedule)
	if err != nil {
		s.logger.Error().Err(err).Str("job", job.Name).Msg("Invalid schedule")
		return
	}

	runner := &jobRunner{
		job:      job,
		ticker:   time.NewTicker(interval),
		stopChan: make(chan struct{}),
	}

	s.jobs = append(s.jobs, runner)
	s.logger.Info().
		Str("job", job.Name).
		Str("schedule", job.Schedule).
		Dur("interval", interval).
		Msg("Job registered")
}

// Start starts the scheduler
func (s *Scheduler) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	s.logger.Info().Int("jobs", len(s.jobs)).Msg("Starting scheduler")

	for _, runner := range s.jobs {
		s.wg.Add(1)
		go s.runJob(runner)
	}
}

// Stop stops the scheduler and waits for all jobs to complete
func (s *Scheduler) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	s.running = false
	s.mu.Unlock()

	s.logger.Info().Msg("Stopping scheduler...")
	s.cancel()

	for _, runner := range s.jobs {
		runner.ticker.Stop()
		close(runner.stopChan)
	}

	s.wg.Wait()
	s.logger.Info().Msg("Scheduler stopped")
}

func (s *Scheduler) runJob(runner *jobRunner) {
	defer s.wg.Done()

	// Run immediately on start
	s.executeJob(runner.job)

	for {
		select {
		case <-runner.ticker.C:
			s.executeJob(runner.job)
		case <-runner.stopChan:
			return
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *Scheduler) executeJob(job *Job) {
	start := time.Now()
	s.logger.Info().Str("job", job.Name).Msg("Starting job execution")

	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Minute)
	defer cancel()

	if err := job.Handler(ctx); err != nil {
		s.logger.Error().
			Err(err).
			Str("job", job.Name).
			Dur("duration", time.Since(start)).
			Msg("Job execution failed")
		return
	}

	s.logger.Info().
		Str("job", job.Name).
		Dur("duration", time.Since(start)).
		Msg("Job execution completed")
}

// parseSchedule parses schedule strings like "daily@02:00", "hourly", "every:5m"
func parseSchedule(schedule string) (time.Duration, error) {
	switch schedule {
	case "hourly":
		return time.Hour, nil
	case "daily":
		return 24 * time.Hour, nil
	default:
		// Parse "every:5m", "every:1h", etc.
		if len(schedule) > 6 && schedule[:6] == "every:" {
			return time.ParseDuration(schedule[6:])
		}
		// Default to daily if parsing fails
		return 24 * time.Hour, nil
	}
}
