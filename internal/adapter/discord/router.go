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
