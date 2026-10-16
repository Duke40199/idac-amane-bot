package discord

import (
	"github.com/bwmarrin/discordgo"
)

// Handler is the Discord controller. It holds the session plus any repositories,
// services and clients the command handlers need. In this wireframe it carries
// only the session and owner ID; wire additional dependencies in as features grow.
type Handler struct {
	OwnerID string
	Session *discordgo.Session
}

// NewHandler creates our controller.
func NewHandler(s *discordgo.Session, ownerID string) *Handler {
	return &Handler{
		Session: s,
		OwnerID: ownerID,
	}
}

func (h *Handler) HandleStatus(i *discordgo.InteractionCreate) {
	h.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Up and running, desuwa~",
		},
	})
}
