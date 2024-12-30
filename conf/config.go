package conf

import (
	"time"

	common_config "github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/discovery"
)

const (
	scrapeInterval = model.Duration(15 * time.Second)
	JobName        = "timelite"

	Version = "0.0.1"
)

var DefaultConfig *Config

type Config struct {
	StoragePath string
	LogPath     string
	Host        string

	PromConfig *config.Config
}

func init() {
	staticConfig := &discovery.StaticConfig{
		{
			Targets: []model.LabelSet{
				{
					model.AddressLabel: model.LabelValue(":9090"),
				},
			},
		},
	}

	promConf := &config.DefaultConfig
	promConf.GlobalConfig.MetricNameValidationScheme = config.LegacyValidationConfig
	promConf.ScrapeConfigs = []*config.ScrapeConfig{
		{
			JobName:           JobName,
			HonorLabels:       true,
			HonorTimestamps:   true,
			ScrapeInterval:    scrapeInterval,
			ScrapeTimeout:     config.DefaultGlobalConfig.ScrapeTimeout,
			EnableCompression: true,
			MetricsPath:       config.DefaultScrapeConfig.MetricsPath,
			Scheme:            config.DefaultScrapeConfig.Scheme,
			HTTPClientConfig: common_config.HTTPClientConfig{
				FollowRedirects: true,
				EnableHTTP2:     true,
			},

			ServiceDiscoveryConfigs: discovery.Configs{
				staticConfig,
			},
		},
	}
	DefaultConfig = &Config{
		PromConfig: promConf,
	}
}
