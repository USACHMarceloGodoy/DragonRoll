package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Feature representa un rasgo de clase/subclase/raza/fondo. Es el átomo
// que aparece en LevelStep.FeatureIDs y en Race.TraitIDs.
//
// Ejemplos:
//   - "Imposición de Manos" (clase): cost=action, uses.max = 5*level, recharge=long_rest
//   - "Aura de Protección" (clase, pasivo): cost=passive
//   - "Visión en la Oscuridad" (raza): cost=passive
//   - "Castigo Divino" (clase, action mod): cost=spell-slot
type Feature struct {
	ID          primitive.ObjectID `json:"_id"          bson:"_id,omitempty"`
	Name        string             `json:"name"         bson:"name"`
	Description string             `json:"description"  bson:"description"`

	// Origen narrativo: clase X, subclase Y, raza Z, fondo, custom.
	Origin    string             `json:"origin"        bson:"origin"`               // "class" | "subclass" | "race" | "background" | "feat" | "item"
	OriginRef primitive.ObjectID `json:"origin_ref"    bson:"origin_ref,omitempty"` // _id de Class/Subclass/Race según corresponda

	// Costo en economía de acción.
	Cost string `json:"cost" bson:"cost"` // "action" | "bonus" | "reaction" | "passive" | "spell-slot" | "free"

	// Recurso opcional. Si no se setea, el feature es ilimitado (típico de pasivos).
	Uses *FeatureUses `json:"uses,omitempty" bson:"uses,omitempty"`

	// Tirada asociada al uso (informativo, el motor de combate lo usa
	// para sugerir botones y cálculos).
	RollHint string `json:"roll_hint" bson:"roll_hint,omitempty"` // "1d8+CAR", "2d6 radiante", "cura X"

	Meta CatalogMeta `json:"meta" bson:"meta"`
}

type FeatureUses struct {
	// Max puede ser un número fijo o una fórmula textual ("CAR mod", "1+nivel/3").
	// Para v1 dejamos un número; las fórmulas las evalúa el frontend al renderizar.
	Max      int    `json:"max"      bson:"max"`
	Formula  string `json:"formula"  bson:"formula,omitempty"` // documentación: "5 × nivel"
	Recharge string `json:"recharge" bson:"recharge"`          // "short_rest" | "long_rest" | "per_day" | "none"
}
