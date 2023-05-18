package main

import (
	"github.com/scnon/go-utils/logger"
	"github.com/scnon/md-doc/cmd"
)

func main() {
	logger.Config("./logs/md-doc.log", logger.DebugLevel, true)
	cmd.Execute()

}
