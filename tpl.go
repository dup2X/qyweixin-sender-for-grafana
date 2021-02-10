package main

import (
	"bytes"
	"context"
	"fmt"
	"text/template"
	"time"

	"github.com/dup2X/gopkg/logger"
)

func genContent(ctx context.Context, msg GrafanaMessage, level string) string {
	t, err := template.New("send").Parse(tpl)
	if err != nil {
		logger.Errorf(ctx, logger.DLTagUndefined, "InternalServerError: %v", err)

		return tpl
	}

	var body bytes.Buffer
	err = t.Execute(&body, map[string]interface{}{
		"Status":   ET[msg.State],
		"Priority": level,
		"Sname":    msg.RuleName,
		"Items":    msg.MatchInfos,
		"Elink":    msg.RuleLink,
		"Etime":    time.Now().Format("2006-01-02 15:04:05"),
	})

	if err != nil {
		logger.Errorf(ctx, logger.DLTagUndefined, "InternalServerError: %v", err)
		return fmt.Sprintf("InternalServerError: %v", err)
	}

	return body.String()
}

var tpl = `
**事件状态**：P{{.Priority}} {{.Status}}
**策略名称**：{{.Sname}}
{{range .Items}}
metric：{{.Metric}}
tags：
> {{range $key,$val := .Tags}}
> - {{$key}}={{$val}}
{{end}}

**当前值**：{{.Value}}
{{end}}
**触发时间**：{{.Etime}}
**报警详情**：{{.Elink}}
{{if .IsUpgrade}}
---
报警已升级!!!
{{end}}
`

var ET = map[string]string{
	"alerting": "告警",
	"recovery": "恢复",
}
