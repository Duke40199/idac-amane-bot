package discord

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"

	"github.com/bwmarrin/discordgo"

	"idac-amane-bot/pkg/logger"
)

func interactionUser(i *discordgo.InteractionCreate) *discordgo.User {
	if i == nil {
		return nil
	}
	if i.Member != nil && i.Member.User != nil {
		return i.Member.User
	}
	return i.User
}

func recoverInteraction(label string) {
	if r := recover(); r != nil {
		logger.Error(context.Background(), "recovered from panic in "+label, fmt.Errorf("%v", r), slog.String("panic_stack", string(debug.Stack())))
	}
}

type CommandName string

const (
	CommandNameStatus CommandName = "status"

	CommandNameIDAC            CommandName = "idac"
	CommandNameIDACPlayerAlias CommandName = "player-alias"
)

func (h *Handler) RegisterCommands() error {
	// 1. Get Bot ID explicitly
	// We use User("@me") because h.Session.State.User is nil until Open() is called.
	u, err := h.Session.User("@me")
	if err != nil {
		return fmt.Errorf("failed to fetch bot ID: %w", err)
	}
	botID := u.ID

	// 2. Define All Commands
	commands := []*discordgo.ApplicationCommand{
		// --- IDAC Complex Command ---
		{
			Name:        string(CommandNameIDAC),
			Description: "Commands related to the game Initial D The Arcade",
			Options: []*discordgo.ApplicationCommandOption{
				// Player Alias
				{
					Name:        string(CommandNameIDACPlayerAlias),
					Description: "Manage Player Aliases",
					Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
					Options: []*discordgo.ApplicationCommandOption{
						// 1. Link Discord User
						{
							Name:        "set",
							Description: "Link a Discord User to an IGN and Area",
							Type:        discordgo.ApplicationCommandOptionSubCommand,
							Options: []*discordgo.ApplicationCommandOption{
								{
									Type:        discordgo.ApplicationCommandOptionUser,
									Name:        "user",
									Description: "The Discord User",
									Required:    true,
								},
								{
									Type:        discordgo.ApplicationCommandOptionString,
									Name:        "ign",
									Description: "In-Game Name",
									Required:    true,
								},
								{
									Type:        discordgo.ApplicationCommandOptionString,
									Name:        "area",
									Description: "Area Code",
									Required:    true,
								},
							},
						},
						// 2. Link Custom Tag (For non-Discord players)
						{
							Name:        "set-custom",
							Description: "Create a custom alias tag (e.g. 'rival', 'alt')",
							Type:        discordgo.ApplicationCommandOptionSubCommand,
							Options: []*discordgo.ApplicationCommandOption{
								{
									Type:        discordgo.ApplicationCommandOptionString,
									Name:        "tag",
									Description: "Unique Tag Name (e.g. 'takumi', 'rival1')",
									Required:    true,
								},
								{
									Type:        discordgo.ApplicationCommandOptionString,
									Name:        "ign",
									Description: "In-Game Name",
									Required:    true,
								},
								{
									Type:        discordgo.ApplicationCommandOptionString,
									Name:        "area",
									Description: "Area Code",
									Required:    true,
								},
							},
						},
					},
				},
			},
		},
		// --- Simple Commands ---
		{Name: string(CommandNameStatus), Description: "Show current bot status"},
	}

	// 3. Register Event Router
	h.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		defer recoverInteraction("interaction router")
		switch i.Type {
		// A. Execute Commands
		case discordgo.InteractionApplicationCommand:
			name := i.ApplicationCommandData().Name
			switch CommandName(name) {
			case CommandNameIDAC:
				h.routeIdacCommand(i)
			case CommandNameStatus:
				h.HandleStatus(i)
			}
		}
	})

	// 4. Register with Discord
	logger.Info(context.Background(), "registering commands")
	_, err = h.Session.ApplicationCommandBulkOverwrite(botID, "", commands)
	return err
}

// routeIdacCommand routes subcommands for /idac
func (h *Handler) routeIdacCommand(i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		return
	}
	rootOption := options[0]

	// --- HANDLE SUBCOMMAND GROUPS ("player-alias") ---
	if rootOption.Type == discordgo.ApplicationCommandOptionSubCommandGroup {
		// Owner Check
		if !h.IsRequestFromOwner(h.Session, i) {
			return
		}
		subCmd := rootOption.Options[0]
		optMap := make(map[string]string)

		// Extract options for alias
		for _, opt := range subCmd.Options {
			if opt.Type == discordgo.ApplicationCommandOptionUser {
				optMap["user"] = opt.UserValue(h.Session).ID
			} else {
				optMap[opt.Name] = opt.StringValue()
			}
		}
		if CommandName(rootOption.Name) == CommandNameIDACPlayerAlias {
			switch subCmd.Name {
			case "set":
				h.HandleSetPlayerAlias(i, optMap["user"], optMap)
			case "set-custom":
				h.HandleSetPlayerAlias(i, optMap["tag"], optMap)
			}
		}
		return
	}
}
