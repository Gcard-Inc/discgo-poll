package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        "poll",
			Description: "basic command route for starting a poll",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "question",
					Description: "question for the poll",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "multiple-options",
					Description: "able to cast multiple votes",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "answer1",
					Description: "first answer",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "answer2",
					Description: "second answer",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "answer-3",
					Description: "third answer",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "answer-4",
					Description: "fourth answer",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "answer-5",
					Description: "fifth answer",
					Required:    false,
				},
			},
		},
		{
			Name:        "pollist",
			Description: "List all open polls",
		},
		{
			Name:        "pollhelp",
			Description: "get help on all commands",
		},
		{
			Name:        "closepoll",
			Description: "Close a poll by id",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "poll-id",
					Description: "Poll id",
					Required:    true,
				},
			},
		},
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"poll": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			margs := []interface{}{}
			msgformat := ` New poll: \n`
			if len(i.ApplicationCommandData().Options) >= 3 {
				for j, opt := range i.ApplicationCommandData().Options {
					if opt.Name == "question" {
						msgformat += "question: %s \n"
						margs = append(margs, opt.StringValue())
					} else if opt.Name == "multipleOptions" {
						msgformat += "> multipleOptions: %v\n"
						margs = append(margs, opt.BoolValue())
					} else {
						msgformat += fmt.Sprintf("answer %d", j)
						msgformat += ": %v\n"
						margs = append(margs, opt.StringValue())
					}
				}
				margs = append(margs, i.ApplicationCommandData().Options[0].StringValue())
				msgformat += "> poll-id: <#%s>\n"
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(
						msgformat,
						margs...,
					),
				},
			})
		},
		"pollist": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "List all polls slash command",
				},
			})
		},
		"pollhelp": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "List all commands and how to use them",
				},
			})
		},
		"closepoll": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			margs := []interface{}{
				// Here we need to convert raw interface{} value to wanted type.
				// Also, as you can see, here is used utility functions to convert the value
				// to particular type. Yeah, you can use just switch type,
				// but this is much simpler
				i.ApplicationCommandData().Options[0].StringValue(),
			}
			msgformat :=
				` Attempting to close:
				> poll-id: %s
`
			if len(i.ApplicationCommandData().Options) >= 2 {
				margs = append(margs, i.ApplicationCommandData().Options[0].StringValue())
				msgformat += "> poll-id: <#%s>\n"
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				// Ignore type for now, we'll discuss them in "responses" part
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(
						msgformat,
						margs...,
					),
				},
			})
		},
	}
)
