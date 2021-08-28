package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const (
	CmdPoll      = "poll"
	CmdPolList   = "pollist"
	CmdPollHelp  = "pollhelp"
	CmdClosePoll = "closepoll"
)

var numMap = map[int]string{
	1: "1️⃣",
	2: "2️⃣",
	3: "3️⃣",
	4: "4️⃣",
	5: "5️⃣",
}

var commandList = []string{CmdPoll, CmdPolList, CmdPollHelp, CmdClosePoll}

var (
	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        CmdPoll,
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
					Name:        "answer3",
					Description: "third answer",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "answer4",
					Description: "fourth answer",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "answer5",
					Description: "fifth answer",
					Required:    false,
				},
			},
		},
		{
			Name:        CmdPolList,
			Description: "List all open polls",
		},
		{
			Name:        CmdPollHelp,
			Description: "get help on all commands",
		},
		{
			Name:        CmdClosePoll,
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
		CmdPoll: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			msgformatted := formMessageContentPollID(i.ApplicationCommandData())

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: msgformatted,
					Flags:   1 << 6,
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: emojiNumbers(i.ApplicationCommandData()),
						},
					},
				},
			})
		},
		CmdPolList: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "List of all open polls",
				},
			})
		},
		CmdPollHelp: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			msgFormat := "All available commands: \n"
			cmdFormat := "/%s \n"
			for _, c := range commandList {
				msgFormat += fmt.Sprintf(cmdFormat, c)
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: msgFormat,
				},
			})
		},
		CmdClosePoll: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

func formMessageContentPollID(data discordgo.ApplicationCommandInteractionData) string {
	margs := []interface{}{}
	msgformat := "New poll: \n"
	if len(data.Options) >= 3 {
		for j, opt := range data.Options {
			if opt.Name == "question" {
				msgformat += "Question:\n%s \n"
				margs = append(margs, opt.StringValue())
			} else if opt.Name == "multipleOptions" {
				if opt.BoolValue() {
					msgformat += "Allows multiple options\n"
				} else {
					msgformat += "Multiple options not allowed\n"
				}
			} else {
				msgformat += fmt.Sprintf("Answer %d", j)
				msgformat += ": %v\n"
				margs = append(margs, opt.StringValue())
			}
		}
		margs = append(margs, data.Options[0].StringValue())
		msgformat += "PollID: <#%s>\n"
	}

	return fmt.Sprint(msgformat, margs)
}

func emojiNumbers(data discordgo.ApplicationCommandInteractionData) []discordgo.MessageComponent {
	messageComponent := []discordgo.MessageComponent{}
	index := 0
	for _, opt := range data.Options {
		if opt.Name == "question" {
			index++
			messageComponent = append(messageComponent, discordgo.Button{
				Label: opt.StringValue(),
				Style: discordgo.PrimaryButton,
				Emoji: discordgo.ButtonEmoji{
					Name: numMap[index],
				},
			})
		}
	}
	return messageComponent
}

/*func formEmbeddedArray(data discordgo.ApplicationCommandInteractionData) *discordgo.MessageEmbed {
	qEmbed := &discordgo.MessageEmbed{
		Type:  discordgo.EmbedTypeRich,
		Title: "Question",
		Color: 3447003,
	}
	messageFields := []*discordgo.MessageEmbedField{}
	if len(data.Options) >= 3 {
		for _, opt := range data.Options {
			if opt.Name == "question" {
				qEmbed.Description = opt.StringValue()

			} else if opt.Name == "multipleOptions" {
				field := &discordgo.MessageEmbedField{}
				if opt.BoolValue() {
					field.Name = "Allows multiple options"
				} else {
					field.Name = "Multiple options not allowed"
				}
				messageFields = append(messageFields, field)
			} else {

				field := &discordgo.MessageEmbedField{
					Name:   opt.Name,
					Value:  opt.StringValue(),
					Inline: true,
				}
				messageFields = append(messageFields, field)
			}
		}
	}
	qEmbed.Fields = messageFields
	return qEmbed
}*/
