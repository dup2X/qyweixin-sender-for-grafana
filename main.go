package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dup2X/gopkg/logger"
	"github.com/dup2X/qyweixin-sender-for-grafana/utils"
)

type grafanaAlertProxyHandler struct {
}

func (gh *grafanaAlertProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()
	logger.Debugf(ctx, logger.DLTagUndefined, "recv path:%s", r.URL.Path)
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	logger.Infof(ctx, logger.DLTagUndefined, "recv body:%s", string(data))
	var msg = GrafanaMessage{}
	err = json.Unmarshal(data, &msg)
	if err != nil {
		logger.Warnf(ctx, logger.DLTagUndefined, "parse message failed,err:%s", err)
		w.Write([]byte(err.Error()))
	}
	sendToWeiXin(ctx, msg)
}

func main() {
	http.ListenAndServe("127.0.0.1:20009", &grafanaAlertProxyHandler{})
}

type GrafanaMessage struct {
	DashboardID int64             `json:"dashboardId"`
	MatchInfos  []*MatchInfo      `json:"evalMatches"`
	Message     string            `json:"message"`
	RuleID      int64             `json:"ruleId"`
	RuleLink    string            `json:"ruleUrl"`
	RuleName    string            `json:"ruleName"`
	State       string            `json:"state"`
	Tags        map[string]string `json:"tags"`
	Title       string            `json:"title"`
}

type MatchInfo struct {
	Value        int64                  `json:"value"`
	Metric       string                 `json:"metric"`
	Tags         map[string]interface{} `json:"tags"`
	ReadableTags string                 `json:"-"`
}

func sendToWeiXin(ctx context.Context, msg GrafanaMessage) {
	token := msg.Tags["token"]
	if token == "" {
		token = "b10acb50-feda-4808-9d9a-b919e53f6535"
	}
	level := msg.Tags["level"]
	if level == "" {
		level = "4"
	}
	ms := msg.Tags["at"]
	mobiles := strings.Split(ms, ",")
	isAtAll := false
	if _, ok := msg.Tags["is_at_all"]; ok {
		isAtAll = true
	}
	wxMsg := genContent(ctx, msg, level)
	c := utils.New(token, mobiles, isAtAll)
	_ = wxMsg
	_ = c
	println(wxMsg)
	c.Send(token, mobiles, wxMsg)
}
