package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

const (
	// DEV MODE
	dev      = false
	app_name = "ConsoleTranslate"
	version  = "1.0.0"
	command  = "translate"
	git_repo = "https://github.com/Ablaze-MIRAI/ConsoleTranslate"
)

var commands = []*cli.Command{
	{
		Name:  "help",
		Usage: "ヘルプ",
		Action: func(c *cli.Context) error {
			topic := c.Args().First()
			if topic == "api" {
				fmt.Fprintf(color.Output,
					"APIの設定は %s を参照してください\n",
					color.MagentaString(git_repo))
			} else {
				fmt.Fprintf(color.Output,
					"Example\n\n"+
						"%s <テキスト> -t [翻訳先] (-f [翻訳元]:任意)\n\n"+
						"%s -t, --to : 翻訳先の言語コードを指定\n"+
						"%s -f, --from : 翻訳元の言語コードを指定\n\n"+
						"対応している言語の言語コード一覧は `%s` を参照\n",
					command, color.RedString("(必須)"),
					color.CyanString("(任意)"),
					color.MagentaString("https://cloud.google.com/translate/docs/languages"))
			}
			return nil
		},
	},
	{
		Name:  "version",
		Usage: "バージョン",
		Action: func(c *cli.Context) error {
			fmt.Fprintf(color.Output,
				"\n%s v%s\n"+
					"Github: %s\n"+
					"Help: `%s`\n\n",
				app_name, version, git_repo, color.CyanString(command+" help"))
			return nil
		},
	},
}

func main() {
	app := cli.NewApp()
	app.Name = app_name
	app.Version = version
	app.Commands = commands
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "to",
			Usage:   "翻訳先",
			Aliases: []string{"t"},
		},
		&cli.StringFlag{
			Name:    "from",
			Usage:   "翻訳元",
			Aliases: []string{"f"},
		},
	}
	app.Action = func(c *cli.Context) error {
		conf, err := loadConfig(dev)
		if err != nil {
			fmt.Fprintf(color.Output,
				"%s: 設定ファイルの読み込みに失敗\n"+
					"設定は `%s` を参照してください\n",
				color.RedString("Error"), color.MagentaString(git_repo))
			os.Exit(0)
		}
		to := c.String("to")
		from := c.String("from")
		arg := c.Args().First()
		if to == "" {
			fmt.Fprintf(color.Output,
				"%s: 必要な引数がありません\n"+
					"詳細は `%s` を参照してください。\n",
				color.RedString("Error"), color.CyanString(command+" help"))
			os.Exit(0)
		}

		response := &struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
			Text string `json:"text"`
		}{}

		resp, err := http.Get(urlGen(to, from, arg, conf.Api))
		if err != nil {
			fmt.Fprintf(color.Output,
				"%s: リクエストに失敗しました\n"+
					"インターネットの接続、APIの設定等を確認してください\n"+
					"[Log]%s\n",
				color.RedString("Error"), err)
			os.Exit(0)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			fmt.Fprintf(color.Output,
				"%s: リクエストに失敗しました\n"+
					"[Log]HTTP Status: `%s`\n",
				color.RedString("Error"), resp.Status)
			os.Exit(0)
		}

		body, _ := io.ReadAll(resp.Body)
		if err := json.Unmarshal(body, response); err != nil {
			fmt.Fprintf(color.Output,
				"%s: リクエストの解析に失敗しました\n"+
					"[Log]%s\n",
				color.RedString("Error"), err)
			os.Exit(0)
		}
		if response.Msg == "unexpected" {
			fmt.Fprintf(color.Output,
				"%s: 翻訳に失敗しました\n"+
					"翻訳に対応している言語は `%s` を参照してください。\n"+
					"[Log]API Error\n",
				color.RedString("Error"),
				color.MagentaString("https://cloud.google.com/translate/docs/languages"))
			os.Exit(0)
		}

		var lang_info string
		if from == "" {
			lang_info = "auto"
		} else {
			lang_info = from
		}
		fmt.Fprintf(color.Output,
			"%s\n %s\n",
			color.MagentaString("[Before: "+lang_info+"]"), arg)
		fmt.Fprint(color.Output, "  ↓\n")
		fmt.Fprintf(color.Output,
			"%s\n %s\n",
			color.GreenString("[After: "+to+"]"), response.Text)
		return nil
	}
	app.Run(os.Args)
}
