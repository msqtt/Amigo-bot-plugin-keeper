package keeper

import (
	"strings"

	"github.com/mosqu1t0/Amigo-bot/bot"
	"github.com/mosqu1t0/Amigo-bot/utils/logcat"
)

type Keeper struct {
}

func init() {
	keeper := new(Keeper)
	bot.PluginMgr.AddPlugin(keeper)
}

func (keeper *Keeper) GetType() string {
	return bot.MsgPostType
}

func (keeper *Keeper) Init() {
	logcat.Good("[管事的] 简单指令管理机器人的插件已加载! <3")
}

func (keeper *Keeper) Action(b *bot.Bot, v interface{}) {
	msg := v.(*bot.RecvMessage)
	slice := strings.Split(msg.Message, " ")

	switch slice[0] {
	case "！改名":
		handleRename(b, msg, slice[1:])
	case "！好友":
		handleFriend(b, msg)
	case "！拉黑":
		handleBlack(b, msg, slice[1:])
	case "！群组":
		handleGroup(b, msg)
	case "！退群":
		handleLeave(b, msg, slice[1:])
	default:
	}
}
