package actions

import (
	tdlib "github.com/Arman92/go-tdlib"
	"github.com/sirupsen/logrus"
)

var (
// GroupContact = make(map[string]tdlib.User)
// allMembers   = []*tdlib.ChatMember{}
// isFull       = false
)

func GetUsers(client *tdlib.Client, chatid int64, limit int32) (chats []tdlib.ChatMember, err error) {
	// if !isFull && int(limit) > len(allMembers) {

	// offsetOrder := int64(math.MaxInt64)
	offset := int32(0)
	var chatsMems *tdlib.ChatMembers
	// suoergroupid := client.SearchChatMembers
	groupid := GetGroupId(chatid)
	// logrus.Info("g:", groupid)
	//
	chat, cerr := client.GetChat(chatid)
	if cerr != nil {
		logrus.Error("Makesure chat exists error:", err)
		return nil, cerr
	}
	for {
		switch chat.Type.GetChatTypeEnum() {
		case tdlib.ChatTypeSupergroupType:
			chatsMems, err = client.GetSupergroupMembers(int32(groupid), nil, offset, limit)
		case tdlib.ChatTypeBasicGroupType:
			// chatsMems, err = client.GetSupergroupMembers(int32(groupid), nil, offset, limit)
			if bgFullInfo, gerr := client.GetBasicGroupFullInfo(groupid); gerr != nil {
				logrus.Error("Get Basic group err:", gerr)
				return nil, gerr
			} else {
				return bgFullInfo.Members, nil
			}
		}
		if err != nil {
			logrus.Error(err)
			return

		}
		for _, chat := range chatsMems.Members {
			chats = append(chats, chat)
		}
		if len(chats) < int(offset+200) {
			return
		}
		logrus.Info("Get Members: ", len(chats))
		offset += 200
	}

	// chat,errs := client.GetChat(chatid)
	// client.Member
	// logrus.Info("group id:", chatid)
	// groupid := int32(0 - chatid)
	// chatsMems, err = client.SearchChatMembers(chatid, "", limit, nil)
	// client.GetBasicGroupFullInfo()
	// chats = chatsMems.Members
	// return
}

//GetGrouId : get groupid from chatid(int64) -> groupid(int32)
func GetGroupId(chatid int64) int32 {
	return int32(0 - 1000000000000 - chatid)
}
