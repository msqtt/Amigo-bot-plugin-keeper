package keeper

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/mosqu1t0/Amigo-bot/bot"
	"github.com/mosqu1t0/Amigo-bot/utils/logcat"
)

const (
	goodReply  = "收到辣！"
	wrongReply = "指令格式有误..."
	whereReply = "这个指令在这里输入才有用嗷..."
)

func handleRename(b *bot.Bot, msg *bot.RecvMessage, slice []string) {
	if !isRoot(msg) {
		return
	}

	switch msg.MessageType {
	case bot.PriMsgType:
		if len(slice) >= 2 {
			var group_id int64
			var err error
			card := slice[0]
			group_id, err = strconv.ParseInt(slice[1], 10, 64)
			if err != nil {
				logcat.Error("Can't ParseInt when handleRename: ", err)
				return
			}
			err = b.Send(
				bot.GroupCardApi,
				struct {
					GroupId int64  `json:"group_id"`
					UserId  int64  `json:"user_id"`
					Card    string `json:"card"`
				}{group_id, b.Info.UserId, card},
			)
			if err != nil {
				logcat.Error("插件执行失败: ", err)
				return
			}
			replyMessageFrom(b, msg, goodReply)
		} else {
			replyMessageFrom(b, msg, wrongReply)
		}

		// 回复信息
	case bot.GruMsgType:
		if len(slice) < 1 {
			replyMessageFrom(b, msg, wrongReply)
			return
		}
		card := slice[0]
		group_id := msg.GroupId
		var err error
		err = b.Send(
			bot.GroupCardApi,
			struct {
				GroupId int64  `json:"group_id"`
				UserId  int64  `json:"user_id"`
				Card    string `json:"card"`
			}{group_id, b.Info.UserId, card},
		)
		if err != nil {
			logcat.Error("插件执行失败: ", err)
			return
		}
		replyMessageFrom(b, msg, goodReply)
	}
}

func handleLeave(b *bot.Bot, msg *bot.RecvMessage, slice []string) {
	if !isRoot(msg) {
		return
	}

	// 不能在群组中使用
	if msg.MessageType == bot.GruMsgType {
		replyPrivate(b, msg, whereReply)
		return
	}

	if len(slice) < 1 {
		replyMessageFrom(b, msg, wrongReply)
		return
	}
	group_id, err := strconv.ParseInt(slice[0], 10, 64)
	if err != nil {
		logcat.Error("Can't ParseInt when handleLeave: ", err)
		replyMessageFrom(b, msg, wrongReply)
		return
	}
	err = b.Send(
		bot.GroupLeaveApi,
		struct {
			GroupId int64 `json:"group_id"`
		}{group_id},
	)
	if err != nil {
		logcat.Error("插件执行失败: ", err)
	}
	replyMessageFrom(b, msg, goodReply)
}

func handleBlack(b *bot.Bot, msg *bot.RecvMessage, slice []string) {
	if !isRoot(msg) {
		return
	}
	// 不能在群组中使用
	if msg.MessageType == bot.GruMsgType {
		replyPrivate(b, msg, whereReply)
		return
	}
	if len(slice) < 1 {
		replyMessageFrom(b, msg, wrongReply)
		return
	}
	user_id, err := strconv.ParseInt(slice[0], 10, 64)

	if err != nil {
		logcat.Error("Can't ParseInt when handleLeave: ", err)
		replyMessageFrom(b, msg, wrongReply)
		return
	}

	err = b.Send(
		bot.FriDelApi,
		struct {
			UserId int64 `json:"user_id"`
		}{user_id},
	)
	if err != nil {
		logcat.Error("插件执行失败: ", err)
	}
	replyMessageFrom(b, msg, goodReply)
}

func handleFriend(b *bot.Bot, msg *bot.RecvMessage) {
	if !isRoot(msg) {
		return
	}
	// 不能在群组中使用
	if msg.MessageType == bot.GruMsgType {
		replyPrivate(b, msg, whereReply)
		return
	}

	bytes, err := b.QuickTalk(bot.FriGetListApi, nil)
	if err != nil {
		logcat.Error("获取好友列表失败: ", err)
	}
	var friends struct {
		Data []struct {
			UserId   int64  `json:"user_id"`
			Nickname string `json:"nickname"`
			Remark   string `json:"remark"`
		} `json:"data"`
	}

	json.Unmarshal(bytes, &friends)
	friendsList := fmt.Sprintf("%-12s\t%s(%s)\n", "QQ", "昵称", "备注")
	for _, f := range friends.Data {
		if f.Nickname == f.Remark {
			friendsList += fmt.Sprintf("%-10d %s\n", f.UserId, f.Nickname)
		} else {
			friendsList += fmt.Sprintf("%-10d %s(%s)\n",
				f.UserId, f.Nickname, f.Remark,
			)
		}
	}
	replyPrivate(b, msg, friendsList)
}

func handleGroup(b *bot.Bot, msg *bot.RecvMessage) {
	if !isRoot(msg) {
		return
	}
	// 不能在群组中使用
	if msg.MessageType == bot.GruMsgType {
		replyPrivate(b, msg, whereReply)
		return
	}

	bytes, err := b.QuickTalk(bot.GruGetListApi, nil)
	if err != nil {
		logcat.Error("获取群组列表失败: ", err)
	}
	var groups struct {
		Data []struct {
			GroupId   int64  `json:"group_id"`
			GroupName string `json:"group_name"`
		} `json:"data"`
	}
	json.Unmarshal(bytes, &groups)
	groupsList := fmt.Sprintf("%-10s\t%-12s\n", "群号", "群名")
	for _, g := range groups.Data {
		groupsList += fmt.Sprintf("%-10d %s\n", g.GroupId, g.GroupName)
	}
	replyPrivate(b, msg, groupsList)
}

func replyPrivate(b *bot.Bot, msg *bot.RecvMessage, message string) {
	err := b.Send(
		bot.MsgSendApi,
		bot.MsgSend{
			MessageType: bot.PriMsgType,
			UserId:      msg.Sender.UserId,
			Message:     message,
		},
	)
	if err != nil {
		logcat.Error("插件执行失败: ", err)
	}
}

func replyMessageFrom(b *bot.Bot, msg *bot.RecvMessage, message string) {
	err := b.Send(
		bot.MsgSendApi,
		bot.MsgSend{
			MessageType: msg.MessageType,
			UserId:      msg.Sender.UserId,
			GroupId:     msg.GroupId,
			Message:     message,
		},
	)
	if err != nil {
		logcat.Error("插件执行失败: ", err)
	}
}

func isRoot(msg *bot.RecvMessage) bool {
	for _, root := range bot.DefaultBotConfig.Root {
		if root == msg.Sender.UserId {
			return true
		}
	}
	return false
}
