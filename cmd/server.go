package cmd

import (
	"log"
	"os/exec"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/scnon/md-doc/internal"
	"github.com/scnon/md-doc/model"
	"github.com/scnon/md-doc/utils"
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

	e.GET("/", listHandler)
	e.Static("/static", "./static")
	e.Any("/git/:repo/:action", gitHandler)
	e.Any("/doc/:repo/*", docHandler)

	return e.Start(":80")
}

func gitHandler(c echo.Context) error {
	repo := c.Param("repo")
	action := c.Param("action")

	log.Println(repo, action, c.ParamNames(), c.ParamValues())

	switch action {
	case "create":
		err := utils.CreateRepo(repo)
		if err != nil {
			return utils.Resp500(c, err)
		}

		return c.JSON(200, model.Response{
			Code: 200,
			Msg:  "create repo success",
		})
	default:
		internal.Handler()(c.Response(), c.Request())
		return nil
	}
}

func listHandler(c echo.Context) error {
	cmd := exec.Command("git", "ls-tree", "HEAD", "-r")
	cmd.Dir = utils.GetDataPath("test")
	out, err := cmd.Output()
	log.Println(cmd.Dir)
	if err != nil {
		return c.String(200, err.Error())
	}
	res := string(out)
	log.Println(res)
	sp := strings.Split(res, "\n")

	rr := ""
	for _, v := range sp {
		spv := strings.Split(v, "\t")
		if len(spv) < 2 {
			continue
		}
		rr += spv[1] + "<br>"
	}

	return c.HTML(200, rr)
}

func docHandler(c echo.Context) error {
	repo := c.Param("repo")
	path := c.Param("*")
	log.Println(repo, path)

	cmd := exec.Command("git", "show", "HEAD:"+path)
	cmd.Dir = utils.GetDataPath(repo)
	out, err := cmd.Output()

	if err != nil {
		return utils.Resp404(c)
	}

	author, created, updated := utils.GetFileInfo(repo, path)
	res, err := utils.ReaderDoc(repo, path, author, created, updated, out)
	if err != nil {
		return utils.Resp500(c, err)
	}

	return c.HTML(200, res)
}
