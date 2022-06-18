package main

import (
	"fmt"
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

func msg(err error) int {
	if err != nil {
		fmt.Fprintln(color.Error, err.Error())
		return 1
	}
	return 0
}

func run() int {
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
			return fmt.Errorf(
				"%s: 設定ファイルの読み込みに失敗\n"+
					"設定は `%s` を参照してください\n",
				color.RedString("Error"), color.MagentaString(git_repo))
		}
		to := c.String("to")
		from := c.String("from")
		arg := c.Args().First()
		if to == "" {
			return fmt.Errorf(
				"%s: 必要な引数がありません\n"+
					"詳細は `%s` を参照してください。\n",
				color.RedString("Error"), color.CyanString(command+" help"))
		}

		response, err := HttpRequest(urlGen(to, from, arg, conf.Api))
		if err != nil {
			return fmt.Errorf(
				"%s: リクエストに失敗しました\n"+
					"インターネットの接続、APIの設定等を確認してください\n"+
					"[Log]%s\n",
				color.RedString("Error"), err)
		}

		if response.Msg == "unexpected" {
			return fmt.Errorf(
				"%s: 翻訳に失敗しました\n"+
					"翻訳に対応している言語は `%s` を参照してください。\n"+
					"[Log]API Error\n",
				color.RedString("Error"),
				color.MagentaString("https://cloud.google.com/translate/docs/languages"))
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

	return msg(app.Run(os.Args))
}

func main() {
	os.Exit(run())
}
