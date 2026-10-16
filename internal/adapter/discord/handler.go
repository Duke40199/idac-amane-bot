package discord

import (
	"bytes"
	_ "embed"

	"idac-amane-bot/internal/domain/port"

	"github.com/bwmarrin/discordgo"
)

//go:embed resource/unauthorized.gif
var unauthorizedGif []byte

// Handler is the Discord controller. It holds the session plus any repositories,
// services and clients the command handlers need. Wire additional dependencies in
// as features grow.
type Handler struct {
	OwnerID string
	Session *discordgo.Session
	// db repositories
	AliasRepo          port.AliasRepository
	OBRankingCfgRepo   port.OBRankingCfgRepository
	RankingCfgRepo     port.RankingCfgRepository
	CarRepo            port.CarRepository
	TATimeMetadataRepo port.TATimeMetadataRepository
	CfsStateRepo       port.CfsStateRepository
	AreaRepo           port.AreaRepository
	AreaMetadataRepo   port.IDACAreaMetadataRepository
	StoreLocationRepo  port.AllNetStoreLocationsRepository
}

// NewHandler creates our controller.
func NewHandler(
	s *discordgo.Session,
	ownerID string,
	// db repositories
	aliasRepo port.AliasRepository,
	obRankingCfgRepo port.OBRankingCfgRepository,
	rankingCfgRepo port.RankingCfgRepository,
	carRepo port.CarRepository,
	taTimeMetadataRepo port.TATimeMetadataRepository,
	cfsStateRepo port.CfsStateRepository,
	areaRepo port.AreaRepository,
	areaMetadataRepo port.IDACAreaMetadataRepository,
	storeLocationRepo port.AllNetStoreLocationsRepository,
) *Handler {
	return &Handler{
		Session:            s,
		OwnerID:            ownerID,
		AliasRepo:          aliasRepo,
		OBRankingCfgRepo:   obRankingCfgRepo,
		RankingCfgRepo:     rankingCfgRepo,
		CarRepo:            carRepo,
		TATimeMetadataRepo: taTimeMetadataRepo,
		CfsStateRepo:       cfsStateRepo,
		AreaRepo:           areaRepo,
		AreaMetadataRepo:   areaMetadataRepo,
		StoreLocationRepo:  storeLocationRepo,
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

// IsRequestFromOwner checks if the user is the owner.
// If NOT, it sends the unauthorized GIF/Message and returns false.
// If YES, it returns true.
func (h *Handler) IsRequestFromOwner(s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	// Owner ID comes from config (DISCORD_BOT_OWNER_ID). If it's unset, no one is
	// treated as owner, so owner-gated commands fail closed.
	if user := interactionUser(i); user != nil && h.OwnerID != "" && user.ID == h.OwnerID {
		return true
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "This command is restricted to the bot owner.",
			Files: []*discordgo.File{
				{Name: "unauthorized.gif", Reader: bytes.NewReader(unauthorizedGif)},
			},
		},
	})
	return false
}
