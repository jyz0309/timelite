package conf

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
	StoragePath   string
	LogPath       string
	DashboardPath string
	Host          string `default:"0.0.0.0:9090"`

	PromConfig *config.Config `json:"-"`
}

func GetConfig(configPath string) error {
	config := &Config{}
	conf, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file(%s): %v", configPath, err)
	}
	if err = json.Unmarshal(conf, config); err != nil {
		return fmt.Errorf("failed to unmarshal config file(%s): %v", configPath, err)
	}
	DefaultConfig = config
	DefaultConfig.PromConfig = GetPromConfig(DefaultConfig.Host)
	return nil
}

func GetPromConfig(host string) *config.Config {
	staticConfig := &discovery.StaticConfig{
		{
			Targets: []model.LabelSet{
				{
					model.AddressLabel: model.LabelValue(host),
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
	return promConf
}

func InitConfig(configFile, storagePath, host string) error {
	if configFile != "" {
		if _, err := os.Stat(configFile); !os.IsNotExist(err) {
			// read config from file
			err := GetConfig(configFile)
			if err != nil {
				// if file not exist, create it
				return err
			}
		}
	}

	dashboardPath := filepath.Join(filepath.Dir(configFile), "dashboards")
	logPath := filepath.Join(storagePath, "log")

	promConf := GetPromConfig(host)

	DefaultConfig = &Config{
		StoragePath:   storagePath,
		DashboardPath: dashboardPath,
		LogPath:       logPath,
		PromConfig:    promConf,
		Host:          host,
	}
	return SaveConfig(configFile)

}

func SaveConfig(configFile string) error {
	if _, err := os.Stat(filepath.Dir(configFile)); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(configFile), 0755)
		if err != nil {
			return fmt.Errorf("failed to mkdir config dir: %v", err)
		}
	}
	bytes, err := json.Marshal(DefaultConfig)
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, bytes, 0644)
}
