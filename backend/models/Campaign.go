package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Campaign agrupa lo que un DM gestiona. Se usa para:
//  1. Autorización: el DM tiene "control total" sobre los PJ y NPCs de
//     su(s) campaña(s).
//  2. Scope del catálogo: clases/subclases/conjuros homebrew creados
//     con Scope=ScopeCampaign sólo se ven dentro de esta campaña.
//  3. Estado narrativo: NPCs, ubicaciones, bitácora.
type Campaign struct {
	ID          primitive.ObjectID   `json:"_id"          bson:"_id,omitempty"`
	DMID        primitive.ObjectID   `json:"dm_id"        bson:"dm_id"`
	PlayerIDs   []primitive.ObjectID `json:"player_ids" bson:"player_ids"` // miembros (jugadores)
	Name        string               `json:"name"         bson:"name"`
	Description string               `json:"description"  bson:"description,omitempty"`
	System      string               `json:"system"       bson:"system"` // "D&D 5e" | "Homebrew Castoñegro" | ...

	NPCs      []NPC          `json:"npcs"      bson:"npcs"`
	Locations []Location     `json:"locations" bson:"locations"`
	Sessions  []SessionEntry `json:"sessions"  bson:"sessions"`

	Active    bool      `json:"active"     bson:"active"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type NPC struct {
	ID       string    `json:"id"       bson:"id"`
	Name     string    `json:"name"     bson:"name"`
	Role     string    `json:"role"     bson:"role"`
	Faction  string    `json:"faction"  bson:"faction,omitempty"`
	Location string    `json:"location" bson:"location,omitempty"`
	Alive    bool      `json:"alive"    bson:"alive"`
	Attitude string    `json:"attitude" bson:"attitude"` // "aliado" | "neutral" | "sospechoso" | "enemigo"
	Pic      string    `json:"pic"      bson:"pic,omitempty"`
	Color    string    `json:"color"    bson:"color,omitempty"`
	Summary  string    `json:"summary"  bson:"summary"`
	Details  string    `json:"details"  bson:"details,omitempty"`
	Tags     []string  `json:"tags"     bson:"tags"`
	Comments []NPCNote `json:"comments" bson:"comments"`
}

type NPCNote struct {
	Author string `json:"author" bson:"author"`
	Avatar string `json:"avatar" bson:"avatar,omitempty"`
	Text   string `json:"text"   bson:"text"`
}

type Location struct {
	ID          string   `json:"id"          bson:"id"`
	Name        string   `json:"name"        bson:"name"`
	Type        string   `json:"type"        bson:"type"`
	Description string   `json:"description" bson:"description"`
	Tags        []string `json:"tags"        bson:"tags"`
}

type SessionEntry struct {
	Number  int    `json:"number"  bson:"number"`
	Date    string `json:"date"    bson:"date"`
	Title   string `json:"title"   bson:"title"`
	Summary string `json:"summary" bson:"summary"`
}
