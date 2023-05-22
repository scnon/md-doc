package cmd

import (
	"os/exec"

	"github.com/labstack/echo/v4"
	"github.com/scnon/md-doc/internal"
	"github.com/scnon/md-doc/logic"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:     "server",
	Long:    "Run server",
	Aliases: []string{"s"},
	RunE:    runServer,
}

func runServer(cmd *cobra.Command, args []string) error {
	e := echo.New()
	e.Debug = false

	internal.InitConfig(internal.Config{
		AuthPassEnvVar: "",
		AuthUserEnvVar: "",
		DefaultEnv:     "",
		ProjectRoot:    "/Users/x/go/md-doc/data",
		GitBinPath:     "/usr/bin/git",
		UploadPack:     true,
		ReceivePack:    true,
		RoutePrefix:    "",
		CommandFunc:    func(*exec.Cmd) {},
	})

	e.GET("/", logic.ListHandler)
	e.Static("/static", "./static")

	e.Any("/repo/:repo/:action", echo.WrapHandler(internal.Handler()))

	e.Any("/doc/:repo/*", logic.DocHandler)
	e.POST("/api/doc/search", logic.SearchHander)

	return e.Start(":80")
}
