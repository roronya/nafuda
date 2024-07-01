package main

import (
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/slack-go/slack"
)

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
<style>
    @import url('https://fonts.googleapis.com/css2?family=M+PLUS+1p:wght@400;700&display=swap');
    body {
        font-family: 'M PLUS 1p', sans-serif;
    }
    @page {
        size: A4;
        margin: 1cm;
    }
    body {
      display: grid;
      grid-template-columns: 1fr 1fr;
      grid-auto-rows: 240px;
    }
    .badge {
        display: grid;
        grid-template-columns: 1fr 2fr;
        grid-template-rows: auto auto;
        border: 1px solid black;
        padding: 8px;
        column-gap: 16px;

    }
    .badge img {
        max-width: 100%;
        height: auto;
        grid-row: 1 / span 2;
        justify-self: center;
        align-self: center;
    }
    .badge .name {
        font-size: 48px;
        font-weight: bold;
        align-self: end;
        justify-self: start;
    }
    .badge .title {
        font-size: 16px;
        align-self: start;
        justify-self: start;
    }
</style>
</head>
<body>
{{range .}}
<div class="badge">
    <img src="{{.Image}}" alt="{{.FullName}}">
    <div class="name">{{.FullName}}</div>
    <div class="title">{{.Title}}</div>
</div>
{{end}}
</body>
</html>
`

type Member struct {
	FullName string
	Title    string
	Image    string
}

func main() {
	// 環境変数からSlack APIトークンを取得
	slackAPIToken := os.Getenv("SLACK_TOKEN")
	if slackAPIToken == "" {
		log.Fatalf("SLACK_TOKEN environment variable is not set")
	}

	// コマンドライン引数からチャンネルIDを取得
	if len(os.Args) < 2 {
		log.Fatalf("Channel ID is required as a command-line argument")
	}
	channelID := os.Args[1]

	api := slack.New(slackAPIToken)

	// チャンネルのメンバー情報を取得
	params := &slack.GetUsersInConversationParameters{
		ChannelID: channelID,
	}
	userIDs, _, err := api.GetUsersInConversation(params)
	if err != nil {
		log.Fatalf("Error fetching users from channel: %v", err)
	}

	var members []Member
	for _, userID := range userIDs {
		user, err := api.GetUserInfo(userID)
		if err != nil {
			log.Printf("Error fetching user info for %s: %v", userID, err)
			continue
		}

		member := Member{
			FullName: user.RealName,
			Title:    user.Profile.Title,
			Image:    user.Profile.Image192,
		}
		members = append(members, member)
	}

	// HTMLを生成
	tmpl, err := template.New("badges").Parse(htmlTemplate)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	file, err := os.Create("name_badges.html")
	if err != nil {
		log.Fatalf("Error creating HTML file: %v", err)
	}
	defer file.Close()

	err = tmpl.Execute(file, members)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

	fmt.Println("名札の作成が完了しました。name_badges.html を確認してください。")
}
