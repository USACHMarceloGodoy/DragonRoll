package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Subclass extiende una clase base. Tiene su propia progresión (qué
// features desbloquea en qué niveles) y, opcionalmente, su propia lista
// de conjuros añadidos.
type Subclass struct {
	ID          primitive.ObjectID `json:"_id"          bson:"_id,omitempty"`
	ParentClass primitive.ObjectID `json:"parent_class" bson:"parent_class"`

	Name      string `json:"name"       bson:"name"`
	ShortDesc string `json:"short_desc" bson:"short_desc"`
	Flavor    string `json:"flavor"     bson:"flavor,omitempty"`

	// Niveles donde la subclase aporta features (suelen ser 3,7,15,18 o
	// el subset que defina la subclase). El admin/DM define exactamente.
	Levels []LevelStep `json:"levels" bson:"levels"`

	// Conjuros adicionales que esta subclase concede automáticamente
	// (p. ej. el Juramento de Devoción otorga ciertos conjuros).
	GrantedSpells []GrantedSpell `json:"granted_spells" bson:"granted_spells"`

	Meta CatalogMeta `json:"meta" bson:"meta"`
}

type GrantedSpell struct {
	Level   int                `json:"level"    bson:"level"` // a qué nivel se otorga
	SpellID primitive.ObjectID `json:"spell_id" bson:"spell_id"`
}
