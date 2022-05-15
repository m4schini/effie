package bot

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
)

type Option[V any] struct {
	Type  OptionType
	Name  string
	Value V
}

func GetSubCommandOptions(data discordgo.InteractionData) ([]*Option[string], error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}

	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}

	interationOptions := result["options"].([]interface{})

	subCommand := interationOptions[0].(map[string]interface{})

	scOptions := subCommand["options"].([]interface{})

	options := make([]*Option[string], 0, len(scOptions))
	for _, option := range scOptions {
		optionMap := option.(map[string]interface{})

		options = append(options, &Option[string]{Summoner, optionMap["name"].(string), optionMap["value"].(string)})
	}

	return options, nil
}
