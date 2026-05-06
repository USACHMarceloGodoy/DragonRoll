package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Spell: conjuro como entrada de catálogo, referenciable desde Class,
// Subclass y Character. El campo SaveDC y AttackBonus se calculan
// dinámicamente desde el lanzador en runtime; aquí sólo va lo estático.
type Spell struct {
	ID            primitive.ObjectID `json:"_id"          bson:"_id,omitempty"`
	Name          string             `json:"name"         bson:"name"`
	Level         int                `json:"level"        bson:"level"`        // 0..9 (0 = truco)
	School        string             `json:"school"       bson:"school"`       // "evocación", "abjuración", ...
	CastingTime   string             `json:"casting_time" bson:"casting_time"` // "action" | "bonus" | "reaction" | "1 minute" | ...
	Range         string             `json:"range"        bson:"range"`        // "60 ft", "Toque", "Personal"
	Components    string             `json:"components"   bson:"components"`   // "V/S/M (un símbolo sagrado)"
	Duration      string             `json:"duration"     bson:"duration"`     // "Instantáneo" | "Concentración, 1 minuto"
	Description   string             `json:"description"  bson:"description"`
	HigherLevel   string             `json:"higher_level" bson:"higher_level,omitempty"` // efecto al subir el slot
	Ritual        bool               `json:"ritual"       bson:"ritual"`
	Concentration bool               `json:"concentration" bson:"concentration"`

	// Damage opcional para conjuros de daño directo.
	Damage     string `json:"damage"      bson:"damage,omitempty"`      // "8d6"
	DamageType string `json:"damage_type" bson:"damage_type,omitempty"` // "fuego", "radiante", ...

	// SaveAbility opcional: si se permite TS contra el conjuro.
	SaveAbility string `json:"save_ability" bson:"save_ability,omitempty"` // "DES", "CON", ...

	Meta CatalogMeta `json:"meta" bson:"meta"`
}
