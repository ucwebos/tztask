package command

import (
	"bytes"
	"context"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"

	"tztask/domain/entity"
	"tztask/domain/repo"
	"tztask/utils"
)

const (
	HttpGET = "GET"
)

var verifyMethodMap = map[string]struct{}{
	HttpGET: {},
}

type HttpCommand struct {
	Expr    string          `json:"expr" yaml:"expr"`
	Name    string          `json:"name" yaml:"name"`
	Command *entity.Command `json:"command" yaml:"command"`
	// 存储实现
	HttpResultRepo repo.HttpResultRepo `di:"repo_impl.HttpResultRepoSQLite"`
}

func (h *HttpCommand) ID() string {
	return h.Name
}

func (h *HttpCommand) Spec() string {
	return h.Expr
}

func (h *HttpCommand) Verify() error {
	if h.Expr == "" || h.Name == "" {
		return errors.New("Expr or Name error")
	}
	if h.Command.Method == "" {
		return errors.New("not found HttpCommand args: method")
	}
	if _, ok := verifyMethodMap[h.Command.Method]; !ok {
		return errors.Errorf("HttpCommand method[%s]: not allowd", h.Command.Method)
	}
	if h.Command.Target == "" {
		return errors.New("not found HttpCommand args: target")
	}
	return nil
}

func (h *HttpCommand) TryNum() int {
	return 3
}

func (h *HttpCommand) Run() error {
	switch h.Command.Method {
	case HttpGET:
		return h.methodGet(h.Command.Target)
	}
	return nil
}

func (h *HttpCommand) methodGet(target string) error {
	raw, resp, err := utils.GetRaw(nil, target, nil, nil, 5000) // 5000ms timeout
	if err != nil || resp == nil {
		return err
	}
	rs := &entity.HttpResult{
		Time:          time.Now().Unix(),
		TaskName:      h.Name,
		Target:        h.Command.Target,
		TargetTo:      target,
		StatusCode:    resp.StatusCode,
		ContentLength: resp.ContentLength,
		ContentType:   resp.Header.Get("Content-Type"),
		Raw:           string(raw),
	}
	if strings.Contains(rs.ContentType, "text/html") {
		if targetTo := h.ParseHTML(raw, rs); targetTo != "" {
			return h.methodGet(targetTo)
		}
	}
	err = h.HttpResultRepo.Save(context.Background(), rs)
	if err != nil {
		return errors.Wrap(err, "HttpResultRepo.Save ")
	}
	return nil
}

func (h *HttpCommand) ParseHTML(raw []byte, rs *entity.HttpResult) string {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(raw))
	if err != nil {
		log.Printf("ParseHTML err:%v", err)
		return ""
	}
	if v, ok := doc.Find("meta[http-equiv='refresh']").Attr("content"); ok {
		targetTo := v[strings.IndexAny(v, "url=")+4:]
		if targetTo != "" || len(targetTo) > 6 && rs.TargetTo != targetTo {
			return targetTo
		}
	}
	rs.Title = doc.Find("title").Text()
	if v, ok := doc.Find("meta[name='description']").Attr("content"); ok {
		rs.Description = v
	}
	return ""
}
