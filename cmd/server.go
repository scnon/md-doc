package cmd

import (
	"bytes"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"text/template"

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
	e.Any("/git/:repo/:action", gitActionHandler)
	e.Any("/doc/:repo/*", docHandler)

	e.GET("/demo", func(c echo.Context) error {
		byte, err := ioutil.ReadFile("./data/demo.md")
		if err != nil {
			return err
		}
		tmpl, err := template.New("doc").Parse(model.DocTmpl)
		if err != nil {
			return err
		}

		var reader bytes.Buffer
		err = tmpl.Execute(&reader, map[string]string{
			"Title":   "demo doc",
			"Content": string(internal.Render2Html(byte)),
		})
		if err != nil {
			return err
		}

		return c.HTML(200, reader.String())
	})

	return e.Start(":80")
}

func gitActionHandler(c echo.Context) error {
	repo := c.Param("repo")
	action := c.Param("action")

	log.Println(repo, action, c.ParamNames(), c.ParamValues())

	switch action {
	case "create":
		if utils.CheckRepoExist(repo) {
			log.Println("create failed: repo ", repo, " exist")
			return c.JSON(200, model.Response{
				Code: 201,
				Msg:  "repo exist",
			})
		}

		log.Println("begin create repo: ", repo)
		code := 200
		msg := "create repo success"
		err := utils.CreateRepo(repo)
		if err != nil {
			code = 202
			msg = err.Error()
		}

		return c.JSON(200, model.Response{
			Code: code,
			Msg:  msg,
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
	return c.String(200, repo+" - "+path)
}
