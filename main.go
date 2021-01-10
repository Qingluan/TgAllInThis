package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Arman92/go-tdlib"
	"github.com/Qingluan/TgAllInThis/tools"

	"github.com/Qingluan/TgAllInThis/actions"
	"github.com/sirupsen/logrus"

	"github.com/alyu/configparser"
)

var (
	action string
	conf   string
	cli    bool
	ini    bool
)

func main() {
	flag.StringVar(&action, "action", "", "set action to do")
	flag.StringVar(&conf, "conf", "conf.ini", "set config ini file path")
	flag.BoolVar(&cli, "cli", false, "true to cli mode .")
	flag.BoolVar(&ini, "init", false, "genreate config ini  .")
	flag.Parse()
	if ini {
		fmt.Println(tools.GenerateConfIni())
		return
	}
	config := readConfig(conf)
	authConfig, _ := config.Section("auth")
	client, err := actions.AuthClient(authConfig)
	if err != nil {
		log.Fatal(err)
	}
	if cli {
		actions.WaitInterrupt(client)
		for {
			RunAction(client, config)
			action = ""
		}
		return
	} else {
		RunAction(client, config)
	}

}

func readConfig(path string) *configparser.Configuration {
	configparser.Delimiter = "="

	config, err := configparser.Read(path)
	if err != nil {
		log.Fatal(err)
	}

	return config
}

func RunAction(client *tdlib.Client, config *configparser.Configuration) {
	if action == "" {
		action = tools.Tui.Select("How to Do?", "getChats", "getContacts", "exit")
	}
	switch action {
	case "getChats":
		limitstr, _ := config.StringValue("getchats", "limit")
		limit, _ := strconv.Atoi(limitstr)
		chats, _ := actions.GetChatList(client, limit)
		for _, chat := range chats {
			ctp := actions.GetChatType(chat)
			tools.Log("All-Chat", "GroupTitle:%s | chat id: %d | groupType: %s", chat.Title, chat.ID, ctp)
		}

	case "getContacts":
		limitstr, _ := config.StringValue("getcontacts", "limit")
		limit, _ := strconv.Atoi(limitstr)
		chats, _ := actions.GetChatList(client, limit, "Group", "GroupPrivate", "GroupOnlyReading")
		datas := make(tools.Datas)
		for _, chat := range chats {
			datas[chat.Title] = actions.GetChatType(chat)
		}
		name := tools.Tui.Input("Search group name>", datas)
		chatid := actions.ChatsMap[name].ID

		contacts, _ := actions.GetUsers(client, chatid, int32(limit))
		for _, contact := range contacts {
			// user ,err := client.GetChatMember(chantid, contact.UserID)
			user, uerr := client.GetUser(contact.UserID)
			if uerr != nil {
				logrus.Error(uerr)
				continue
			}
			// client.GetGroupsInCommon
			// logrus.Info(contact.Extra)
			// client.GetGroups
			tools.Log(name+".group.contact", "id:%d | name : %s |username: %s | phone: %s", user.ID, user.FirstName+user.LastName, user.Username, user.PhoneNumber)
		}
	case "exit":
		basedr, _ := config.StringValue("auth", "tddir")
		exitaction := tools.Tui.Select("If exists save log info as csv?", "yes", "no")
		if exitaction == "yes" {
			tools.SaveAsCsv(basedr)
		}
		os.Exit(0)
	}
}
