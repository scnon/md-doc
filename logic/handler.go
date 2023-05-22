package logic

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/scnon/md-doc/model"
	"github.com/scnon/md-doc/utils"
)

func ListHandler(c echo.Context) error {
	return c.HTML(200, "")
}

func DocHandler(c echo.Context) error {
	repo := c.Param("repo")
	path := c.Param("*")
	log.Println(repo, path)
	if strings.HasSuffix(path, ".png") || strings.HasSuffix(path, ".webp") ||
		strings.HasSuffix(path, ".jpg") || strings.HasSuffix(path, ".jpeg") {
		file := fmt.Sprint(utils.GetGitPath(repo), path)
		return c.File(file)
	}

	out, err := utils.GetFile(repo, path)
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

func SearchHander(c echo.Context) error {
	var req model.SearchReq
	err := json.NewDecoder(c.Request().Body).Decode(&req)
	if err != nil {
		return utils.Resp500(c, err)
	}

	res := utils.RenderSearchItem([]model.SearchItem{
		{
			Title:   req.Key,
			Content: req.Key,
			Class:   "search_item",
		},
		{
			Title:   req.Key,
			Content: req.Key,
			Class:   "search_item",
		},
		{
			Title:   req.Key,
			Content: req.Key,
			Class:   "search_item",
		},
		{
			Title:   req.Key,
			Content: req.Key,
			Class:   "search_item_last",
		},
	})

	return c.JSON(200, model.Response{
		Code: 200,
		Msg:  "success",
		Data: res,
	})
}
