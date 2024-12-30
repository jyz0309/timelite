package engine

import (
	"os"
	"path/filepath"
	"time"
	"timelite/conf"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
	"github.com/prometheus/common/promslog"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"github.com/prometheus/prometheus/scrape"
	"github.com/prometheus/prometheus/tsdb"
)

type Engine struct {
	scraper *scraper

	db *tsdb.DB
}

func NewEngine(storagePath, logPath string) (*Engine, error) {
	f, err := os.OpenFile(filepath.Join(storagePath, "tsdb.info"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	logger := promslog.New(&promslog.Config{
		Writer: f,
	})
	db, err := tsdb.Open(storagePath, logger, nil, tsdb.DefaultOptions(), nil)
	if err != nil {
		return nil, err
	}
	scraper, err := newScraper(db, logPath)
	if err != nil {
		return nil, err
	}
	return &Engine{
		scraper: scraper,
		db:      db,
	}, err
}

func (e *Engine) Run() {
	e.scraper.Run()
}

type scraper struct {
	manager *scrape.Manager
}

func newScraper(app *tsdb.DB, logPath string) (*scraper, error) {
	opts := scrape.Options{
		DiscoveryReloadInterval: model.Duration(time.Second * 15),
		AppendMetadata:          true,
	}

	f, err := os.OpenFile(filepath.Join(logPath, "scrape.info"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	manager, _ := scrape.NewManager(&opts, promslog.New(&promslog.Config{
		Writer: f,
	}), nil, app, prometheus.DefaultRegisterer)

	err = manager.ApplyConfig(conf.DefaultConfig.PromConfig)
	if err != nil {
		return nil, err
	}
	return &scraper{
		manager: manager,
	}, nil
}

func (s *scraper) Run() {
	ts := make(chan map[string][]*targetgroup.Group)
	go s.manager.Run(ts)
	ts <- map[string][]*targetgroup.Group{
		conf.JobName: {
			{
				Targets: []model.LabelSet{
					{
						"__address__": model.LabelValue(conf.DefaultConfig.Host), // TODO
					},
				},
			},
		},
	}
}
