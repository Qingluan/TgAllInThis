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

func MakeSureUser(client *tdlib.Client, each func(contact tdlib.Contact), phone ...string) (existsphone []tdlib.Contact) {

	testContacts := []tdlib.Contact{}

	for _, ph := range phone {
		testcontact := *tdlib.NewContact(ph, "", "", "", 0)
		testContacts = append(testContacts, testcontact)
	}
	imted, err := client.ImportContacts(testContacts)
	if err != nil {
		logrus.Error(err)
		return
	}
	for _, uid := range imted.UserIDs {
		// for _, uid := range users.UserIDs {
		user, cerr := client.GetUser(uid)
		if cerr != nil {
			logrus.Error(cerr)
			if uid != 0 && uid != -1 {
				contact := *tdlib.NewContact("", "", "", "", uid)
				// contact.Extra = user.Username
				each(contact)
				existsphone = append(existsphone, contact)

			}
			break
		} else {

			contact := *tdlib.NewContact(user.PhoneNumber, user.FirstName, user.LastName, "", uid)
			contact.Extra = user.Username
			each(contact)
			existsphone = append(existsphone, contact)
		}
		// }
	}
	return

}

func GetUsers(client *tdlib.Client, chatid int64, limit int32) (chats []tdlib.ChatMember, err error) {
	// if !isFull && int(limit) > len(allMembers) {

	// offsetOrder := int64(math.MaxInt64)
	offset := int32(0)
	var chatsMems *tdlib.ChatMembers
	// suoergroupid := client.SearchChatMembers
	groupid := GetGroupId(chatid)
	logrus.Info("group id: ", groupid, "| chat id:", chatid)
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
			if bgFullInfo, gerr := client.GetBasicGroupFullInfo(int32(groupid)); gerr != nil {
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
	// return349157802
	//1349157802
	//1068773197
}

//GetGrouId : get groupid from chatid(int64) -> groupid(int32)
func GetGroupId(chatid int64) int32 {
	if chatid < -1000000000000 {
		return int32(0 - 1000000000000 - chatid)
	} else {
		return int32(0 - chatid)
	}
}
