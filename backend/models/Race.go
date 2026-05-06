package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Race: raza jugable.
type Race struct {
	ID        primitive.ObjectID `json:"_id"          bson:"_id,omitempty"`
	Name      string             `json:"name"         bson:"name"`
	ShortDesc string             `json:"short_desc"   bson:"short_desc"`

	Speed int    `json:"speed"        bson:"speed"` // pies por turno (típicamente 25/30/35)
	Size  string `json:"size"         bson:"size"`  // "Mediano" | "Pequeño" | ...

	// Bonificadores típicos de stat (libres por raza homebrew).
	AbilityBonuses []AbilityBonus `json:"ability_bonuses" bson:"ability_bonuses"`

	// Lenguajes que la raza otorga por defecto.
	Languages []string `json:"languages" bson:"languages"`

	// Rasgos raciales (referencias a Feature con Origin="race").
	TraitIDs []primitive.ObjectID `json:"trait_ids" bson:"trait_ids"`

	// Subrazas opcionales — se modelan también como Race con ParentRace seteado.
	ParentRace primitive.ObjectID `json:"parent_race" bson:"parent_race,omitempty"`

	Meta CatalogMeta `json:"meta" bson:"meta"`
}

type AbilityBonus struct {
	Stat  string `json:"stat"  bson:"stat"`  // "FUE" | "DES" | ...
	Bonus int    `json:"bonus" bson:"bonus"` // típicamente +1 o +2
}

// Skill: catálogo universal (las 18 mecánicas del sistema d20 + homebrew).
// Es la única entidad de catálogo SIN CatalogMeta.IsHomebrew obligatorio
// porque las 18 base se siembran como fixtures y son inmutables; las
// homebrew sí llevan meta.
type Skill struct {
	ID          primitive.ObjectID `json:"_id"          bson:"_id,omitempty"`
	Name        string             `json:"name"         bson:"name"` // "Atletismo"
	Code        string             `json:"code"         bson:"code"` // identificador estable: "athletics"
	Stat        string             `json:"stat"         bson:"stat"` // "FUE" | "DES" | "INT" | "SAB" | "CAR"
	Description string             `json:"description"  bson:"description,omitempty"`

	Meta CatalogMeta `json:"meta" bson:"meta"`
}
