package notify

import (
	"fmt"
	"github.com/hopeio/cherry/utils/sdk/dingding"
	"github.com/urfave/cli/v2"
	"time"
)

type Config struct {
	Repo          string
	CommitAuthor  string
	Commit        string
	CommitTag     string
	CommitRef     string
	CommitMessage string
	CommitBranch  string
	DingToken     string
	DingSecret    string
	BuildLink     string
}

func GetConfig(c *cli.Context) *Config {
	return &Config{
		Repo:          c.String("repo"),
		CommitAuthor:  c.String("commit_author_name"),
		Commit:        c.String("commit"),
		CommitTag:     c.String("commit_tag"),
		CommitRef:     c.String("commit_ref"),
		CommitMessage: c.String("commit_message"),
		CommitBranch:  c.String("commit_branch"),
		DingToken:     c.String("ding_token"),
		DingSecret:    c.String("ding_secret"),
		BuildLink:     c.String("drone_build_link"),
	}
}

func Notify(c *Config) error {

	if c.DingToken == "" {
		return nil
	}

	msg := "\\n # 发布通知 " +
		" \\n ### 项目: " + c.Repo +
		" \\n ### 操作人: " + c.CommitAuthor +
		" \\n ### 参考: " + c.CommitRef +
		" \\n ### 分支: " + c.CommitBranch +
		" \\n ### 标签: " + c.CommitTag +
		" \\n ### 时间: " + fmt.Sprint(time.Now().Format("2006-01-02 15:04:05")) +
		" \\n ### 提交: " + c.Commit +
		" \\n ### 提交信息: " + c.CommitMessage +
		" \\n ### 发布详情: " + c.BuildLink

	var err error
	if c.DingSecret != "" {
		err = dingding.SendRobotMarkDownMessageWithSecret(c.DingToken, c.DingSecret, "发布通知", msg)
	} else {
		err = dingding.SendRobotMarkDownMessage(c.DingToken, "发布通知", msg)
	}

	return err
}
