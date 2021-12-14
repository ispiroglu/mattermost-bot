package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const URL = "http://localhost:8065"

/*
type Bot struct {
	UserEmail    string `json:"UserEmail"`
	UserPassword string `json:"UserPassword"`
	TeamName     string `json:"TeamName"`
	ChannelName  string `json:"ChannelName"`
}
*/

type Post struct {
	Channel     string `json:"channel"`
	Message     string `json:"message"`
	Attachments []struct {
		Fallback   string `json:"fallback"`
		Color      string `json:"color"`
		Pretext    string `json:"pretext"`
		AuthorName string `json:"author_name"`
		AuthorIcon string `json:"author_icon"`
		AuthorLink string `json:"author_link"`
		Title      string `json:"title"`
		TitleLink  string `json:"title_link"`
		Fields     []struct {
			Short bool   `json:"short"`
			Title string `json:"title"`
			Value string `json:"value"`
		} `json:"fields"`
		ImageUrl string `json:"image_url"`
	} `json:"attachments"`
}

var tmp string
var post Post

/*
var bot Bot
var post Post
var client *model.Client4
var botTeam *model.Team
var debuggingChannel *model.Channel*/

func init() {
	file2, err := ioutil.ReadFile("msg.json")
	err = json.Unmarshal(file2, &post)
	if err != nil {
		println("Error at init")
		panic(err)
	}
}

func main() {
	println("Bot is trying to wake up.")
	run()
}
func run() {
	whoGotPermit()
	getStartDate()
	getEndDate()
	getDepartmanVal()
	//saveFile()
	sendWebHook()
}
func sendWebHook() {
	resp, err := http.Post(URL+"/hooks/yt44wgr6tbfc88pfnwy7anj8ec", "application/json", readFile())
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	log.Printf(sb)
}

func readFile() *bytes.Buffer {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(post)
	if err != nil {
		println("Coulnt encode")
		panic(err)
	}
	return &buf
}
func getDepartmanVal() {
	println("Departman")
	_, _ = fmt.Scanln(&tmp)
	post.Attachments[0].Fields[3].Value = tmp
}

func getStartDate() {
	println("Izin baslangic tarihi")
	_, _ = fmt.Scanln(&tmp)
	post.Attachments[0].Fields[1].Value = tmp
}

func getEndDate() {
	println("Izin bitis tarihi")
	_, _ = fmt.Scanln(&tmp)
	post.Attachments[0].Fields[2].Value = tmp
}

func whoGotPermit() {
	println("Izin alan kisi")
	_, _ = fmt.Scanln(&tmp)
	post.Attachments[0].Fields[0].Value = tmp
}

/*
func createTable(name string, department string, startD string, finishD string) string {
	str := fmt.Sprintf(""+
		"| Departman  | Izine ayirilan kisi  | Izin baslangic  | Izin bitis |\n"+
		"| :-----------------| ----------------- | ----------------- | -----------------: |\n"+
		"|         %s        |         %s        |        %s         |         %s         |\n", department, name, startD, finishD)
	return str
}
func SendMsg(channel *model.Channel) model.SlackAttachment {
	msg := createTable("Evren Ispiroglu", "bilsiimHR", "10.12.2021", "20.12.2021")

	post.Attachments{
		Title,
	}
}
func EncodeToBytes(p interface{}) []byte {

	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("uncompressed size (bytes): ", len(buf.Bytes()))
	return buf.Bytes()
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
		os.Exit(1)
	} else {
		botTeam = team
	}
}
*/
