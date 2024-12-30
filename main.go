package main

import (
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
	queryCmd.Flag("config-path", "The path to the config file.").Default("./timelite/config").String()
	queryCmd.Flag("storage-path", "The path to the tsdb data.").Default("./timelite/tsdb").String()

	tsdbCmd := app.Command("tsdb", "Manage the tsdb.")
	tsdbMockCmd := tsdbCmd.Command("mock", "Mock the tsdb data.")

	//tsdbCleanCmd := tsdbCmd.Command("clean", "Clean the tsdb data.")
	runTSDBCmd := tsdbCmd.Command("run", "Run the tsdb")
	runTSDBCmd.Flag("host", "The host to scrape metrics.").Default("0.0.0.0:9090").
		StringVar(&conf.DefaultConfig.Host)

	runTSDBCmd.Flag("storage-path", "The path to store the tsdb data.").Default("./timelite/tsdb").
		StringVar(&conf.DefaultConfig.StoragePath)

	parseCmd := kingpin.MustParse(app.Parse(os.Args[1:]))

	switch parseCmd {
	case queryCmd.FullCommand():
		cmd.Init()
	case tsdbMockCmd.FullCommand():
		util.PromTestServer()
	case runTSDBCmd.FullCommand():
		os.MkdirAll(filepath.Join(conf.DefaultConfig.StoragePath, "log"), 0755)
		conf.DefaultConfig.LogPath = filepath.Join(conf.DefaultConfig.StoragePath, "log")
		logPath := filepath.Join(conf.DefaultConfig.LogPath, "log.info")
		f, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			panic(err)
		}
		logrus.StandardLogger().SetOutput(f)

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
