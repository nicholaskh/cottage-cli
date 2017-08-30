package main

import (
	server "github.com/nicholaskh/golib/server"
)

var c *cmd

func init() {
	parseFlags()

	if options.showVersion {
		server.ShowVersionAndExit()
	}

	c = newCmd()
	server.SetupLogging(options.logFile, options.logLevel, options.crashLogFile)
}

func main() {
	c.run()
}
