package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Scope describe quién ve un objeto del catálogo.
//   - global:   visible para toda la plataforma; sólo el admin lo crea/edita.
//   - campaign: visible sólo dentro de la campaña X; el DM lo crea/edita.
//   - user:     visible sólo para el usuario que lo creó (sandbox personal).
type Scope string

const (
	ScopeGlobal   Scope = "global"
	ScopeCampaign Scope = "campaign"
	ScopeUser     Scope = "user"
)

// CatalogMeta es lo que toda entidad del catálogo (Class, Subclass, Spell,
// Feature, Race) tiene en común: scope, autoría y timestamps.
// Se embebe en cada modelo concreto.
type CatalogMeta struct {
	Scope       Scope              `json:"scope"          bson:"scope"`
	CampaignID  primitive.ObjectID `json:"campaign_id"    bson:"campaign_id,omitempty"`
	OwnerID     primitive.ObjectID `json:"owner_id"       bson:"owner_id,omitempty"`
	IsHomebrew  bool               `json:"is_homebrew"    bson:"is_homebrew"`
	SourceLabel string             `json:"source_label"   bson:"source_label,omitempty"` // texto libre: "Homebrew Marcelo", "Castoñegro v2"
	CreatedAt   time.Time          `json:"created_at"     bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"     bson:"updated_at"`
}

// Class: clase base (Bárbaro, Mago, etc., u homebrew).
//
// La progresión por nivel se modela como una lista de LevelStep para los
// niveles 1..20. Cada step indica el bonus de proficiencia, los slots de
// conjuro, y qué Features se desbloquean (referencias a la colección
// features). Los slots de subclase aparecen reflejados en Subclass.
type Class struct {
	ID        primitive.ObjectID `json:"_id"               bson:"_id,omitempty"`
	Name      string             `json:"name"              bson:"name"`
	ShortDesc string             `json:"short_desc"        bson:"short_desc"`
	HitDie    string             `json:"hit_die"           bson:"hit_die"` // "d6"|"d8"|"d10"|"d12"

	// Stats primarias (para sugerir asignación) y TS competentes.
	PrimaryStats   []string `json:"primary_stats"    bson:"primary_stats"`   // p.ej. ["FUE","CAR"]
	SaveProficient []string `json:"save_proficient"  bson:"save_proficient"` // p.ej. ["FUE","CAR"]

	// Competencias iniciales: armaduras/armas/tools (texto libre, son labels)
	ArmorProficiencies  []string `json:"armor_profs"   bson:"armor_profs"`
	WeaponProficiencies []string `json:"weapon_profs"  bson:"weapon_profs"`
	ToolProficiencies   []string `json:"tool_profs"    bson:"tool_profs"`

	// Skill picks: el jugador elige N habilidades de FromSkillIDs.
	SkillChoice SkillChoice `json:"skill_choice" bson:"skill_choice"`

	// A qué nivel se elige subclase (3 para casi todas, 1 para Clérigo/Hechicero/Brujo, 2 para Mago/Druida).
	SubclassLevel int                  `json:"subclass_level" bson:"subclass_level"`
	Subclasses    []primitive.ObjectID `json:"subclasses"     bson:"subclasses"`

	// Tabla de progresión 1..20.
	Levels []LevelStep `json:"levels" bson:"levels"`

	// Casting (opcional). Si la clase no lanza conjuros, se omite.
	Casting *CastingInfo `json:"casting,omitempty" bson:"casting,omitempty"`

	Meta CatalogMeta `json:"meta" bson:"meta"`
}

type SkillChoice struct {
	Count        int                  `json:"count"           bson:"count"`          // p.ej. 2
	FromSkillIDs []primitive.ObjectID `json:"from_skill_ids"  bson:"from_skill_ids"` // pool donde elegir
}

// LevelStep describe los recursos que la clase otorga al alcanzar Level.
type LevelStep struct {
	Level            int                  `json:"level"             bson:"level"`
	ProficiencyBonus int                  `json:"proficiency_bonus" bson:"proficiency_bonus"`
	FeatureIDs       []primitive.ObjectID `json:"feature_ids"       bson:"feature_ids"`
	// Slots de conjuro disponibles a este nivel (índice = nivel del conjuro 1..9).
	// Por ejemplo: SpellSlots = [4,2,0,0,0,0,0,0,0] significa 4 de nv1 y 2 de nv2.
	SpellSlots [9]int `json:"spell_slots" bson:"spell_slots"`
	// Cantidad total de cantrips conocidos a este nivel (0 si no aplica).
	CantripsKnown int `json:"cantrips_known" bson:"cantrips_known"`
	// Cantidad de conjuros de nivel 1+ conocidos o preparables. -1 = "todos los del libro".
	SpellsKnown int `json:"spells_known" bson:"spells_known"`
	// Notas libres (p.ej. "Ataque extra", "Mejora de ASI/Don").
	Notes string `json:"notes" bson:"notes,omitempty"`
}

// CastingInfo describe el motor de conjuros de la clase.
type CastingInfo struct {
	Ability     string `json:"ability"      bson:"ability"`     // "INT" | "SAB" | "CAR"
	Preparation string `json:"preparation"  bson:"preparation"` // "prepared" | "known" | "spellbook" | "pact"
	Ritual      bool   `json:"ritual"       bson:"ritual"`
	// Lista de conjuros disponibles para esta clase (referencia a Spells).
	// Subclases pueden ampliarla con sus propios SubclassSpells.
	SpellListIDs []primitive.ObjectID `json:"spell_list_ids" bson:"spell_list_ids"`
}
