package main

import (
	"fmt"
	"os"
	"path/filepath"
	"timelite/cmd"
	"timelite/conf"
	"timelite/engine"
	"timelite/util"

	"github.com/alecthomas/kingpin/v2"
	"github.com/sirupsen/logrus"
)

func main() {
	app := kingpin.New(filepath.Base(os.Args[0]), "Tooling for timelite.").UsageWriter(os.Stdout)
	app.Version(conf.Version)
	app.HelpFlag.Short('h')

	queryCmd := app.Command("query", "Query the data in the timeline.")
	queryConfigFile := queryCmd.Flag("config-file", "The path to the config file.").String()
	queryCmd.Flag("storage-path", "The path to the tsdb data.").Default("./storage/tsdb").String()

	tsdbCmd := app.Command("tsdb", "Manage the tsdb.")
	tsdbMockCmd := tsdbCmd.Command("mock", "Mock the tsdb data.")

	//tsdbCleanCmd := tsdbCmd.Command("clean", "Clean the tsdb data.")
	runTSDBCmd := tsdbCmd.Command("run", "Run the tsdb")
	host := runTSDBCmd.Flag("host", "The host to scrape metrics.").Default("0.0.0.0:9090").String()

	storagePath := runTSDBCmd.Flag("storage-path", "The path to store the tsdb data.").Default("./storage/tsdb").
		String()

	configPath := runTSDBCmd.Flag("config-path", "The path to store the tsdb data.").Default("./storage/conf").String()

	parseCmd := kingpin.MustParse(app.Parse(os.Args[1:]))

	switch parseCmd {
	case queryCmd.FullCommand():
		if queryConfigFile == nil {
			panic("config path is required")
		}
		err := conf.GetConfig(*queryConfigFile)
		if err != nil {
			panic(err)
		}
		conf.InitDashboardConf(conf.DefaultConfig.DashboardPath)
		cmd.Init()
	case tsdbMockCmd.FullCommand():
		util.PromTestServer()
	case runTSDBCmd.FullCommand():
		prepareTSDBEnv(*configPath, *storagePath, *host)
		// save config
		conf.SaveConfig(filepath.Join(*configPath, "config.json"))

		engine, err := engine.NewEngine(conf.DefaultConfig.StoragePath, conf.DefaultConfig.LogPath)
		if err != nil {
			logrus.Errorf("Failed to create engine: %v", err)
			os.Exit(1)
		}
		engine.Run()
		stop := make(chan struct{})
		<-stop
	}
}

func prepareTSDBEnv(configPath, storagePath, host string) {
	err := conf.InitConfig(configPath, storagePath, host)
	if err != nil {
		panic(err)
	}
	// mkdir log dir
	err = os.MkdirAll(conf.DefaultConfig.LogPath, 0755)
	if err != nil {
		panic(fmt.Sprintf("failed to mkdir log(%s) dir: %v", conf.DefaultConfig.LogPath, err))
	}
	// mkdir dashboard dir
	err = os.MkdirAll(conf.DefaultConfig.DashboardPath, 0755)
	if err != nil {
		panic(fmt.Sprintf("failed to mkdir dashboard(%s) dir: %v", conf.DefaultConfig.DashboardPath, err))
	}
	logFile := filepath.Join(conf.DefaultConfig.LogPath, "log.info")
	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(fmt.Sprintf("failed to open log file(%s): %v", logFile, err))
	}
	logrus.StandardLogger().SetOutput(f)
}
