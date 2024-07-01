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
    .badge {
        width: 48%;
        border: 1px solid black;
        padding: 10px;
        margin: 1%;
        float: left;
        box-sizing: border-box;
    }
    .badge img {
        max-width: 100px;
        height: auto;
        display: block;
        margin-bottom: 10px;
    }
    .badge .name {
        font-size: 18px;
        font-weight: bold;
    }
    .badge .title {
        font-size: 14px;
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

