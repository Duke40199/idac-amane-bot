package discord

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"idac-amane-bot/internal/domain/entity"
	"idac-amane-bot/pkg/logger"
	"idac-amane-bot/pkg/tracer"

	"github.com/bwmarrin/discordgo"
)

func (h *Handler) HandleSetPlayerAlias(i *discordgo.InteractionCreate, key string, optMap map[string]string) {
	ctx := tracer.NewContext(context.Background())
	if u := interactionUser(i); u != nil {
		ctx = tracer.WithUserID(ctx, u.ID)
	}

	// Defer first: the DB write and the blocking h.Session.User() lookup below can
	// exceed Discord's 3-second response window, which would expire the token and
	// make the final reply fail.
	if err := h.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}); err != nil {
		logger.Error(ctx, "failed to defer interaction", err)
		return
	}

	ign := optMap["ign"]
	area := optMap["area"]

	// Normalize Area
	if val, ok := entity.AreaAliases[strings.ToLower(area)]; ok {
		area = val
	}

	// Normalize Key (custom tags)
	isNumeric := regexp.MustCompile(`^\d+$`).MatchString(key)
	if !isNumeric {
		key = strings.ToLower(key)
	}

	// Save
	err := h.AliasRepo.SetPlayerAlias(key, ign, area)
	msg := ""
	if err != nil {
		msg = "❌ Failed to save alias: " + err.Error()
	} else {
		areaName := entity.AreaDisplayNameByCode[area]
		if areaName == "" {
			areaName = area
		}

		if isNumeric && len(key) > 15 {
			// --- Fetch Username to avoid tagging ---
			displayName := key // Fallback to ID if fetch fails

			// Try fetching user from Discord
			if u, err := h.Session.User(key); err == nil {
				displayName = u.Username
				// Optional: Use u.GlobalName if you prefer the display name
				if u.GlobalName != "" {
					displayName = u.GlobalName
				}
			}
			msg = fmt.Sprintf("✅ **User Linked:** %s → **%s** (%s)", displayName, ign, areaName)
		} else {
			msg = fmt.Sprintf("✅ **Tag Registered!**\nTag `%s` → **%s** (%s)\nUse it like: `/idac player-info player:%s`", key, ign, areaName, key)
		}
	}

	if _, err := h.Session.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content:         &msg,
		AllowedMentions: &discordgo.MessageAllowedMentions{},
	}); err != nil {
		logger.Error(ctx, "failed to edit response", err)
	}
}
