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
    @page {
        size: A4;
        margin: 1cm;
    }
    body {
        font-family: 'M PLUS 1p', sans-serif;
        display: grid;
        grid-template-columns: repeat(2, 1fr);
        grid-template-rows: repeat(5, 1fr);
        page-break-inside: avoid;
    }
    .nafuda {
        display: grid;
        grid-template-columns: 1fr 2fr;
        grid-template-rows: auto auto auto;
        border: 1px solid black;
        padding: 8px;
        box-sizing: border-box;
        height: 240px;
        page-break-inside: avoid;
        column-gap: 16px;
    }
    .nafuda img {
        max-width: 100%;
        height: auto;
        grid-row: 1 / span 3;
        justify-self: center;
        align-self: center;
    }
    .nafuda .name {
        font-size: 48px;
        font-weight: bold;
        align-self: end;
        justify-self: start;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }
    .nafuda .display_name {
        font-size: 24px;
        align-self: start;
        justify-self: start;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }
    .nafuda .title {
        font-size: 16px;
        align-self: start;
        justify-self: start;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }
</style>
</head>
<body>
{{range .}}
<div class="nafuda">
    <img src="{{.Image}}" alt="{{.FullName}}">
    <div class="name">{{.FullName}}</div>
    <div class="display_name">@{{.DisplayName}}</div>
    <div class="title">{{.Title}}</div>
</div>
{{end}}
</body>
</html>
`

type Member struct {
	FullName    string
	DisplayName string
	Title       string
	Image       string
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
			FullName:    user.RealName,
			DisplayName: user.Profile.DisplayName,
			Title:       user.Profile.Title,
			Image:       user.Profile.Image192,
		}
		members = append(members, member)
	}

	// HTMLを生成
	tmpl, err := template.New("nafuda").Parse(htmlTemplate)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	file, err := os.Create("nafuda.html")
	if err != nil {
		log.Fatalf("Error creating HTML file: %v", err)
	}
	defer file.Close()

	err = tmpl.Execute(file, members)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

	fmt.Println("complete! see nafuda.html")
}
