package messenger

import (
	"effie/bot"
	"effie/broker"
	"effie/logger"
	"effie/output"
	"effie/riot"
	"time"
)

var log = logger.Get("messenger").Sugar()

var started = false

func Start(channel bot.Messenger) {
	if started {
		log.Error("messenger already started")
		return
	} else {
		started = true
	}

	broker.Subscribe("match.started", onMatchStart(channel))
	broker.Subscribe("match.ranked.solo.started", onMatchStart(channel))
	broker.Subscribe("match.ranked.flex.started", onMatchStart(channel))

	broker.Subscribe("match.updated", onMatchUpdate(channel))
	broker.Subscribe("match.ended", onMatchEnd(channel))
	broker.Subscribe("match.postgame", onPostGame(channel))

}

func onMatchStart(channel bot.Messenger) func(topic string, message interface{}) {
	return func(topic string, message interface{}) {
		msg, ok := message.(*broker.MatchUpdate)
		if !ok {
			log.Warnf("unexpected message type %s: %T", topic, msg)
			return
		}
		gameId := msg.Game.GameID

		level, _ := riot.GetGameLevel(msg.SummonerId, msg.Game)
		if !output.Ok(level) {
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

			_, err := channel.UpdateMessage(game.DMessage.ID, GetMessageString(gameId))
			if err != nil {
				log.Errorw(err.Error(), "gameId", gameId, "summonerId", msg.SummonerId)
			}

			return
		}
		SetState(gameId, Started, msg, nil)

		dmsg, err := channel.SendMessage(GetMessageString(gameId))
		if err != nil {
			log.Errorw(err.Error(), "gameId", gameId, "summonerId", msg.SummonerId)
			return
		}

		SetDMessage(gameId, dmsg)
	}
}

func onMatchUpdate(channel bot.Messenger) func(topic string, message interface{}) {
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

		updateMessage(channel, gameId)
	}
}

func onMatchEnd(channel bot.Messenger) func(topic string, message interface{}) {
	return func(topic string, message interface{}) {
		msg, ok := message.(*broker.MatchUpdate)
		if !ok {
			log.Warnf("unexpected message type %s: %T", topic, msg)
			return
		}
		gameId := msg.Game.GameID

		SetState(gameId, Ended, msg, nil)

		updateMessage(channel, gameId)
	}
}

func onPostGame(channel bot.Messenger) func(topic string, message interface{}) {
	return func(topic string, message interface{}) {
		msg, ok := message.(*broker.MatchPostData)
		if !ok {
			log.Warnf("unexpected message type %s: %T", topic, msg)
			return
		}
		gameId := msg.GameId

		SetState(gameId, PostGame, nil, msg)

		updateMessage(channel, gameId)

		CloseGame(gameId)
	}
}

func updateMessage(channel bot.Messenger, gameId int) {
	_, ok := GetGame(gameId)
	if !ok {
		return
	}

	dmsg := GetDMessage(gameId)
	if dmsg == nil {
		log.Warnw("no discord message associated", "gameId", gameId)
		return
	}

	_, err := channel.UpdateMessage(dmsg.ID, GetMessageString(gameId))
	if err != nil {
		log.Errorw(err.Error(), "gameId", gameId)
		return
	}
}
