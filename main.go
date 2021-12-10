package main

import (
	"encoding/json"
	"fmt"
	"github.com/mattermost/mattermost-server/v6/model"
	"io/ioutil"
)

const URL = "http://localhost:8065"

type Bot struct {
	UserEmail    string `json:"UserEmail"`
	UserPassword string `json:"UserPassword"`
	TeamName     string `json:"TeamName"`
	ChannelName  string `json:"ChannelName"`
}

var bot Bot
var client *model.Client4
var botTeam *model.Team
var debuggingChannel *model.Channel

func init() {
	file, _ := ioutil.ReadFile("config.json")
	err := json.Unmarshal(file, &bot)
	if err != nil {
		println("Error at init")
		panic(err)
	}
}

func main() {
	println("Bot is trying to wake up.")

	client = model.NewAPIv4Client(URL)
	start()
	SendMsg(debuggingChannel)

}
func createTable(name string, department string, startD string, finishD string) string {
	str := fmt.Sprintf(""+
		"| Departman  | Izine ayirilan kisi  | Izin baslangic  | Izin bitis |\n"+
		"| :-----------------| ----------------- | ----------------- | -----------------: |\n"+
		"|         %s        |         %s        |        %s         |         %s         |\n", department, name, startD, finishD)
	return str
}
func SendMsg(channel *model.Channel) {
	msg := createTable("Evren Ispiroglu", "bilsiimHR", "10.12.2021", "20.12.2021")
	slackAttachment := model.SlackAttachment{
		Text:      msg,
		Title:     "Evren Ispiroglu",
		TitleLink: fmt.Sprintf("[HR linki]%s", "meetingUrl"),
	}

	post := &model.Post{
		ChannelId: channel.Id,
		Message:   "",
		Props: map[string]interface{}{
			"attachments":              []*model.SlackAttachment{&slackAttachment},
			"meetingID":                "meetingID",
			"meeting_link":             "meetingURL",
			"meeting_status":           "zoom.WebhookStatusStarted",
			"meeting_personal":         true,
			"meeting_topic":            "Izin'e Ayrılma !!",
			"meeting_creator_username": "creator.Username",
			"meeting_provider":         "zoomProviderName",
		},
	}

	_, _, err := client.CreatePost(post)
	if err != nil {
		println("Couldnt create post")
		panic(err)
	}
}
func start() {
	LoginAsUser()
	FindBotTeam()
	CreateBotDebuggingChannelIfNeeded()
}
func LoginAsUser() {
	_, _, err := client.Login(bot.UserEmail, bot.UserPassword)
	if err != nil {
		println("Couldnt login")
		panic(err)
	}
}
func CreateBotDebuggingChannelIfNeeded() {
	if rchannel, _, err := client.GetChannelByName(bot.ChannelName, botTeam.Id, ""); err != nil {
		println("We failed to get the channels")
		panic(err)
	} else {
		debuggingChannel = rchannel
		return
	}
}
func FindBotTeam() {
	if team, _, err := client.GetTeamByName(bot.TeamName, ""); err != nil {
		println("We failed to get the initial load")
		println("or we do not appear to be a member of the team '" + bot.TeamName + "'")
		panic(err)
		//os.Exit(1)
	} else {
		botTeam = team
	}
}
