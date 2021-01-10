package actions

import (
	"math"

	"github.com/sirupsen/logrus"

	"github.com/Arman92/go-tdlib"
)

var AllChats []*tdlib.Chat
var ChatsMap = make(map[string]*tdlib.Chat)
var haveFullChatList bool

// see https://stackoverflow.com/questions/37782348/how-to-use-getchats-in-tdlib
func getChatList(client *tdlib.Client, limit int) error {

	if !haveFullChatList && limit > len(AllChats) {
		offsetOrder := int64(math.MaxInt64)
		offsetChatID := int64(0)
		var chatList = tdlib.NewChatListMain()
		var lastChat *tdlib.Chat

		if len(AllChats) > 0 {
			lastChat = AllChats[len(AllChats)-1]
			for i := 0; i < len(lastChat.Positions); i++ {
				//Find the main chat list
				if lastChat.Positions[i].List.GetChatListEnum() == tdlib.ChatListMainType {
					offsetOrder = int64(lastChat.Positions[i].Order)
				}
			}
			offsetChatID = lastChat.ID
		}

		// get chats (ids) from tdlib
		logrus.Info("To get chats ...")
		chats, err := client.GetChats(chatList, tdlib.JSONInt64(offsetOrder),
			offsetChatID, int32(limit-len(AllChats)))
		if err != nil {
			return err
		}
		if len(chats.ChatIDs) == 0 {
			haveFullChatList = true
			return nil
		}

		for _, chatID := range chats.ChatIDs {
			// get chat info from tdlib
			chat, err := client.GetChat(chatID)
			if err == nil {
				AllChats = append(AllChats, chat)
				ChatsMap[chat.Title] = chat
			} else {
				return err
			}
		}
		return getChatList(client, limit)
	}
	return nil
}

/*GetChatList get chat list
@param typeFilter ...stirng :
	ctp := "Contact"
	switch chat.Type.GetChatTypeEnum() {
	case tdlib.ChatTypeSecretType:
		ctp = "GroupPrivate"
	case tdlib.ChatTypeBasicGroupType:
		ctp = "Group"
	case tdlib.ChatTypeSupergroupType:
		ctp = "GroupOnlyReading"
	}
*/
func GetChatList(client *tdlib.Client, limit int, typeFilter ...string) (chats []*tdlib.Chat, err error) {
	err = getChatList(client, limit)
	if err != nil {
		return
	}

	for _, chat := range AllChats {
		ctp := GetChatType(chat)

		if typeFilter != nil {
			for _, tp := range typeFilter {
				if tp == ctp {
					chats = append(chats, chat)
					break
				}
			}
		} else {
			chats = append(chats, chat)
		}
		// logrus.Infof("[%d] : %s | %d Type: %s", no, chat.Title, chat.ID, ctp)
	}
	return
}

func GetChatType(chat *tdlib.Chat) string {
	ctp := "Contact"
	switch chat.Type.GetChatTypeEnum() {
	case tdlib.ChatTypeSecretType:
		ctp = "GroupPrivate"
	case tdlib.ChatTypeBasicGroupType:
		ctp = "Group"
	case tdlib.ChatTypeSupergroupType:
		ctp = "GroupOnlyReading"
	}
	return ctp
}
