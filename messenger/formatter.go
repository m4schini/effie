package messenger

import (
	"effie3/broker"
	"effie3/cache"
	"effie3/riot"
	"effie3/util"
	"effie3/values/volume"
	"fmt"
	"github.com/KnutZuidema/golio/riot/lol"
	"strings"
	"time"
)

func formatGameStart(level volume.Level, info *lol.GameInfo, participants []string) string {
	pString := strings.Join(participants, ", ")
	tributeTitle := "tribute"
	if len(participants) > 1 {
		tributeTitle = "tributes"
	}

	switch level {
	case volume.All:
		return fmt.Sprintf("**%s** started survival training!", pString)
	case volume.Ranked:
		return fmt.Sprintf(
			"**%s** volunteered as %s!\n"+
				"*Welcome, welcome! Happy Hunger Games! And may the odds be ever in your favor.*",
			pString,
			tributeTitle)
	default:
		return fmt.Sprintf(
			"**%s** volunteered as %s!\n"+
				"*Welcome, welcome! Happy Promo Games! And may the odds be ever in your favor.*",
			pString,
			tributeTitle,
		)
	}
}

func GetStartedMessage(info *lol.GameInfo, participants []string, level volume.Level) string {
	return formatGameStart(level, info, participants)
}

func GetLoadingMessage(info *lol.GameInfo, participants []string, level volume.Level) string {
	return formatGameStart(level, info, participants) + "\n\nLooking for sponsors..."
}

func GetActiveMessage(info *lol.GameInfo, participants []string, level volume.Level) string {
	length := time.Since(time.UnixMilli(int64(info.GameStartTime)))

	var queue string
	q, err := riot.GetQueueInfo(info.GameQueueConfigID)
	if err != nil {
		queue = string(level)
	} else {
		queue = q.Description
	}

	// message header
	msg := formatGameStart(level, info, participants) +
		"\n\n"

	// participant details
	for _, participant := range participants {
		leagues, err := riot.GetLeagues(participant)
		if err != nil {
			continue
		}

		var l *lol.LeagueItem

		// if current league is not contains in league by summoners, skip
		if !util.SliceContains(leagues, func(lg *lol.LeagueItem) bool {
			if riot.ToQueueConfigId(l.QueueType) == q.ID {
				l = lg
				return true
			} else {
				return false
			}
		}) {
			continue
		}

		msg = msg + fmt.Sprintf("**%s** (Rank: %v %v) (%dLP)", participant, l.Tier, l.Rank, l.LeaguePoints) + "\n"
	}
	msg = msg + "\n"

	// game details
	msg = msg +
		fmt.Sprintf("*Length*: **%02d:%02d**\n", int(length.Minutes()), int(length.Seconds())%60) +
		fmt.Sprintf("*Queue*: **%v**\n", queue)

	if level == volume.Promo {
		msg = msg + fmt.Sprintf("*This Match is part of a promotion series!*\n")
	} else {
		msg = msg + "*This is not a promo game*"
	}

	return msg
}

func GetEndedMessage(info *lol.GameInfo, participants []string, level volume.Level) string {
	go func() {
		summoner, err := riot.GetSummonerByName(participants[0])
		if err != nil {
			return
		}

		for i := 0; i < 12; i++ {
			match, p, err := riot.GetMatch(fmt.Sprintf("%s_%d", info.PlatformID, info.GameID), summoner.ID)
			if err != nil {
				return
			}
			broker.Publish("match.postgame", &broker.MatchPostData{
				GameId:   info.GameID,
				Summoner: p,
				Match:    match,
			})

			time.Sleep(10 * time.Second)
		}
	}()

	return "Game has ended. Waiting for post game data..."
}

func GetPostMessage(info *broker.MatchPostData, participants []string, level volume.Level) string {
	scoreboard := ""
	for _, participant := range info.Match.Info.Participants {
		if participant.Win == info.Summoner.Win {
			scoreboard = scoreboard + fmt.Sprintf("- *%s* *(%s)* KDA: %02d/%02d/%02d\n",
				participant.SummonerName,
				participant.ChampionName,
				participant.Kills,
				participant.Deaths,
				participant.Assists,
			)

		}
	}

	participantString := strings.Join(participants, ", ")
	haveHas := "has"
	if len(participants) > 1 {
		haveHas = "have"
	}

	headline := ""
	if info.Summoner.Win {
		switch level {
		case volume.Promo:
			count := nextPromoGames()

			headline = fmt.Sprintf("**%s** %s __won__ the %d%s Promo Games!",
				participantString, haveHas, count, util.GetNumberPostfix(int(count%10)))
			break
		case volume.Ranked:
			count := nextHungerGames()

			headline = fmt.Sprintf("**%s** %s __won__ the %d%s Hunger Games!",
				participantString, haveHas, count, util.GetNumberPostfix(int(count%10)))
			break
		default:
			headline = "**WIN!**"
			break
		}
	} else {
		switch level {
		case volume.Promo:
			count := nextPromoGames()
			headline = fmt.Sprintf("**%s** %s __lost__ the %d%s Promo Games!",
				participantString, haveHas, count, util.GetNumberPostfix(int(count%10)))
			break
		case volume.Ranked:
			count := nextHungerGames()
			headline = fmt.Sprintf("**%s** %s __lost__ the %d%s Hunger Games!",
				participantString, haveHas, count, util.GetNumberPostfix(int(count%10)))
			break
		default:
			headline = "**LOSE!**"
			break
		}
	}
	headline = headline + "\n"

	return headline + scoreboard
}

func nextHungerGames() int {
	getCount, setCount := cache.GetHungerGamesCount()
	count := getCount()
	if count < 0 {
		count = 1
	}

	setCount(count + 1)
	return count
}

func nextPromoGames() int {
	getCount, setCount := cache.GetPromoGamesCount()
	count := getCount()
	if count < 0 {
		count = 1
	}

	setCount(count + 1)
	return count
}
