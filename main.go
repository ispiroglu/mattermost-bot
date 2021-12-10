package main

import (
	"fmt"
	"github.com/mattermost/mattermost-server/v6/model"
	"os"
	"os/signal"
)

const (
	NAME = "bot_Mattermost"

	USER_EMAIL    = "bot@example.com"
	USER_PASSWORD = "P*ssw0rd02"
	USER_NAME     = "bot_Mattermost"

	USER_FIRST = "bot"
	USER_LAST  = "Mattermost"

	TEAM_NAME        = "evren"
	CHANNEL_LOG_NAME = "evren"

	URL = "http://localhost:8065"
)

var client *model.Client4
var webSocketClient *model.WebSocketClient

var botUser *model.User
var botTeam *model.Team
var debuggingChannel *model.Channel

type MyStruct struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	UserId    string `json:"user_id"`
	CreatedAt int64  `json:"created_at"`
}

func main() {
	println("Bot is trying to wake up.")

	SetupGracefulShutdown()

	client = model.NewAPIv4Client("http://localhost:8065")

	start()

	//SendMsg("Hello Mattermost", debuggingChannel)
	sendCoolMsg("Evren", "10.12.2021", "20.12.2021")
	//sendPrivate("eispiroglu")

}
func sendCoolMsg(name string, startD string, finishD string) {

	str := fmt.Sprintf(""+
		"| Izine ayirilan kisi  | Izin baslangic  | Izin bitis |\n"+
		"| :----------------- | ----------------- | -----------------: |\n"+
		"|         %s         |        %s         |         %s         |\n", name, startD, finishD)
	SendMsg(str, debuggingChannel)

}
func sendPrivate(userName string) {

	tempUser, _, err := client.GetUserByUsername(userName, "")
	if err != nil {
		println("Couldnt get user by username")
		panic(err)
	}

	direct, _, err := client.CreateDirectChannel(tempUser.Id, botUser.Id)
	if err != nil {
		println("Couldnt get direct")
		return
	}
	SendMsg(tempUser.Username+" Izin Aldi", direct)

}
func SendMsg(msg string, channel *model.Channel) {
	post := &model.Post{}
	post.ChannelId = channel.Id
	post.Message = msg

	_, _, err := client.CreatePost(post)
	if err != nil {
		println("Couldnt create post")
		panic(err)
	}
}

func start() {
	MakeSureServerIsRunning()
	LoginAsUser()
	FindBotTeam()
	CreateBotDebuggingChannelIfNeeded()
}
func LoginAsUser() {
	user, _, err := client.Login(USER_EMAIL, USER_PASSWORD)
	if err != nil {
		println("Couldnt login")
		panic(err)
	}
	botUser = user
}
func CreateBotDebuggingChannelIfNeeded() {
	if rchannel, _, err := client.GetChannelByName(CHANNEL_LOG_NAME, botTeam.Id, ""); err != nil {
		println("We failed to get the channels")
		panic(err)
	} else {
		debuggingChannel = rchannel
		return
	}

	// Looks like we need to create the logging channel
	channel := &model.Channel{}
	channel.Name = CHANNEL_LOG_NAME
	channel.DisplayName = "Debugging For Sample Bot"
	channel.Purpose = "This is used as a test channel for logging bot debug messages"
	channel.Type = model.ChannelTypeOpen
	channel.TeamId = botTeam.Id
	if rchannel, _, err := client.CreateChannel(channel); err != nil {
		println("We failed to create the channel " + CHANNEL_LOG_NAME)
		panic(err)
	} else {
		debuggingChannel = rchannel
		println("Looks like this might be the first run so we've created the channel " + CHANNEL_LOG_NAME)
	}
}
func FindBotTeam() {
	if team, _, err := client.GetTeamByName(TEAM_NAME, ""); err != nil {
		println("We failed to get the initial load")
		println("or we do not appear to be a member of the team '" + TEAM_NAME + "'")
		panic(err)
		//os.Exit(1)
	} else {
		botTeam = team
	}
}
func SetupGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			if webSocketClient != nil {
				webSocketClient.Close()
			}

			SendMsg("_"+NAME+" has **stopped** running_", debuggingChannel)
			os.Exit(0)
		}
	}()
}
func MakeSureServerIsRunning() {
	if props, _, err := client.GetOldClientConfig(""); err != nil {
		println("There was a problem pinging the Mattermost server.  Are you sure it's running?")
		panic(err)
	} else {
		println("Server detected and is running version " + props["Version"])
	}
}
