package query

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/promslog"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/tsdb"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/labels"
)

var globalQuerier *Querier

func GetGlobalQuerier(path string) *Querier {
	if globalQuerier == nil {
		logrus.Info("new quierier")
		querier, err := NewQuerier("./timelite/tsdb")
		if err != nil {
			logrus.Errorf("failed to new querier, err[%s]", err)
			panic(err)
		}
		globalQuerier = querier
	}
	return globalQuerier
}

type Querier struct {
	ctx context.Context

	storagePath string
	engine      *promql.Engine

	logger *slog.Logger
}

func NewQuerier(path string) (*Querier, error) {
	logger := promslog.New(&promslog.Config{})
	// TODO(alkaid): make the config can changeable
	opts := promql.EngineOpts{
		Logger:               logger.With("component", "query_engine"),
		Reg:                  prometheus.DefaultRegisterer,
		MaxSamples:           50000000,
		Timeout:              time.Duration(1 * time.Minute),
		ActiveQueryTracker:   promql.NewActiveQueryTracker(path, 10, logger.With("component", "activeQueryTracker")),
		LookbackDelta:        time.Duration(5 * time.Minute),
		EnableAtModifier:     true,
		EnableNegativeOffset: true,
		EnablePerStepStats:   false,
	}

	return &Querier{
		ctx:         context.Background(),
		storagePath: path,
		engine:      promql.NewEngine(opts),
		logger:      logger,
	}, nil
}

type Series struct {
	Metric labels.Labels `json:"metric"`
	Floats []float64     `json:"values,omitempty"`
}

func (q *Querier) NewRangeQuery(ctx context.Context, qs string, start, end time.Time, interval time.Duration) ([]*Series, []int64, error) {
	db, err := tsdb.OpenDBReadOnly(q.storagePath, "", q.logger.With("component", "tsdb-reader"))
	if err != nil {
		return nil, nil, err
	}
	defer db.Close()

	query, err := q.engine.NewRangeQuery(q.ctx, db, nil, qs, start, end, interval)
	if err != nil {
		logrus.Errorf("failed to get new query, err[%s]", err)
		return nil, nil, err
	}
	defer query.Close()

	result := query.Exec(ctx)
	if result.Err != nil {
		logrus.Errorf("failed to exec query, err[%s]", result.Err.Error())
		return nil, nil, result.Err
	}
	series := make([]*Series, 0)
	timestamp := make([]int64, 0)
	logrus.Info(result.Value.Type())
	switch rs := result.Value.(type) {
	case promql.Matrix:
		var addTS bool
		for _, serie := range rs {
			sample := make([]float64, len(serie.Floats))
			for i, s := range serie.Floats {
				sample[i] = s.F
				if !addTS {
					timestamp = append(timestamp, s.T)
				}
			}
			addTS = true
			series = append(series, &Series{
				Metric: serie.Metric,
				Floats: sample,
			})
		}
	}

	return series, timestamp, nil
}

func (q *Querier) Close() error {
	return q.engine.Close()
}
