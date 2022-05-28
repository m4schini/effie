package messenger

import (
	"effie3/bot"
	"effie3/broker"
	"effie3/conf"
	"effie3/logger"
	"effie3/riot"
	"time"
)

var log = logger.Get("messenger").Sugar()

var started = false

func Start(channel bot.Api, channelId string) {
	if started {
		log.Error("messenger already started")
		return
	} else {
		started = true
	}

	broker.Subscribe("match.started", onMatchStart(channel, channelId))
	broker.Subscribe("match.ranked.solo.started", onMatchStart(channel, channelId))
	broker.Subscribe("match.ranked.flex.started", onMatchStart(channel, channelId))

	broker.Subscribe("match.updated", onMatchUpdate(channel, channelId))
	broker.Subscribe("match.ended", onMatchEnd(channel, channelId))
	broker.Subscribe("match.postgame", onPostGame(channel, channelId))

}

func onMatchStart(channel bot.Api, channelId string) func(topic string, message interface{}) {
	return func(topic string, message interface{}) {
		msg, ok := message.(*broker.MatchUpdate)
		if !ok {
			log.Warnf("unexpected message type %s: %T", topic, msg)
			return
		}
		gameId := msg.Game.GameID

		level, _ := riot.GetGameLevel(msg.SummonerId, msg.Game)
		if !conf.VolumeOk(level) {
			log.Infow("game doesn't reach output level", "gameId", gameId, "summonerId", msg.SummonerId)
			return
		}

		err := AddGame(msg, level)
		if err != nil {
			log.Errorw(err.Error(), "gameId", gameId, "summonerId", msg.SummonerId)

			game, ok := GetGame(msg.Game.GameID)
			if !ok {
				return
			}

			err := game.Message.Edit(GetMessageString(gameId))
			if err != nil {
				log.Errorw(err.Error(), "gameId", gameId, "summonerId", msg.SummonerId)
			}

			return
		}
		SetState(gameId, Started, msg, nil)

		dmsg, err := channel.SendMessage(channelId, GetMessageString(gameId))
		if err != nil {
			log.Errorw(err.Error(), "gameId", gameId, "summonerId", msg.SummonerId)
			return
		}

		SetDMessage(gameId, dmsg)
	}
}

func onMatchUpdate(channel bot.Api, channelId string) func(topic string, message interface{}) {
	return func(topic string, message interface{}) {
		msg, ok := message.(*broker.MatchUpdate)
		if !ok {
			log.Warnf("unexpected message type %s: %T", topic, msg)
			return
		}
		gameId := msg.Game.GameID

		start := time.UnixMilli(int64(msg.Game.GameStartTime))
		if time.Since(start) > 2*time.Hour {
			SetState(gameId, Loading, msg, nil)
		} else {
			SetState(gameId, Active, msg, nil)
		}

		updateMessage(channel, channelId, gameId)
	}
}

func onMatchEnd(channel bot.Api, channelId string) func(topic string, message interface{}) {
	return func(topic string, message interface{}) {
		msg, ok := message.(*broker.MatchUpdate)
		if !ok {
			log.Warnf("unexpected message type %s: %T", topic, msg)
			return
		}
		gameId := msg.Game.GameID

		SetState(gameId, Ended, msg, nil)

		updateMessage(channel, channelId, gameId)
	}
}

func onPostGame(channel bot.Api, channelId string) func(topic string, message interface{}) {
	return func(topic string, message interface{}) {
		msg, ok := message.(*broker.MatchPostData)
		if !ok {
			log.Warnf("unexpected message type %s: %T", topic, msg)
			return
		}
		gameId := msg.GameId

		SetState(gameId, PostGame, nil, msg)

		updateMessage(channel, channelId, gameId)

		CloseGame(gameId)
	}
}

func updateMessage(channel bot.Api, channelId string, gameId int) {
	_, ok := GetGame(gameId)
	if !ok {
		log.Warnw("game is not being tracked", "gameId", gameId)
		return
	}

	dmsg := GetDMessage(gameId)
	if dmsg == nil {
		log.Warnw("no discord message associated", "gameId", gameId)
		return
	}

	err := dmsg.Edit(GetMessageString(gameId))
	if err != nil {
		log.Errorw(err.Error(), "gameId", gameId)
		return
	}
}
