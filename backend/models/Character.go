package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Character es la hoja de personaje. Apunta al catálogo (Class, Subclass,
// Race, Spells, Skills, Features) por ObjectID. Lo que vive embebido es
// SOLO lo que pertenece al PJ específico: stats elegidas, estado de
// combate, inventario, ataques compuestos por el jugador, roleplay.
//
// Cuando el frontend pide /characters/:id, el backend hace populate de
// las referencias y devuelve el documento "expandido". Si el admin/DM
// edita una clase del catálogo, los PJ que la usan ven los cambios al
// próximo populate (es la consecuencia esperada del enfoque data-driven).
type Character struct {
	ID         primitive.ObjectID `json:"_id"         bson:"_id,omitempty"`
	OwnerID    primitive.ObjectID `json:"owner_id"    bson:"owner_id"`
	CampaignID primitive.ObjectID `json:"campaign_id" bson:"campaign_id,omitempty"`
	Active     bool               `json:"active"      bson:"active"`

	// Identidad y catálogo
	Name       string             `json:"name"        bson:"name"`
	Pronouns   string             `json:"pronouns"    bson:"pronouns,omitempty"`
	ClassID    primitive.ObjectID `json:"class_id"    bson:"class_id"`
	SubclassID primitive.ObjectID `json:"subclass_id" bson:"subclass_id,omitempty"` // se setea al alcanzar el nivel de subclase
	RaceID     primitive.ObjectID `json:"race_id"     bson:"race_id"`
	Background string             `json:"background"  bson:"background,omitempty"` // texto libre por ahora
	Alignment  string             `json:"alignment"   bson:"alignment,omitempty"`
	Level      int                `json:"level"       bson:"level"`
	XP         int                `json:"xp"          bson:"xp"`
	NextXP     int                `json:"next_xp"     bson:"next_xp"`

	Pic   string `json:"pic"   bson:"pic,omitempty"`
	Color string `json:"color" bson:"color,omitempty"`

	// Combate básico y estado mutable
	HP          int        `json:"hp"           bson:"hp"`
	HPMax       int        `json:"hp_max"       bson:"hp_max"`
	TempHP      int        `json:"temp_hp"      bson:"temp_hp"`
	AC          int        `json:"ac"           bson:"ac"`
	Speed       int        `json:"speed"        bson:"speed"`
	Initiative  int        `json:"initiative"   bson:"initiative"`
	Inspiration bool       `json:"inspiration"  bson:"inspiration"`
	HitDice     HitDice    `json:"hit_dice"     bson:"hit_dice"`
	DeathSaves  DeathSaves `json:"death_saves"  bson:"death_saves"`

	// Stats elegidas por el jugador
	Stats AbilityScores `json:"stats" bson:"stats"`

	// Habilidades elegidas (referencias a la colección skills)
	SkillProficiencyIDs []primitive.ObjectID `json:"skill_proficiency_ids" bson:"skill_proficiency_ids"`
	SkillExpertiseIDs   []primitive.ObjectID `json:"skill_expertise_ids"   bson:"skill_expertise_ids"`

	// Conjuros conocidos / preparados (referencias a spells)
	CantripIDs        []primitive.ObjectID `json:"cantrip_ids"         bson:"cantrip_ids"`
	SpellsKnownIDs    []primitive.ObjectID `json:"spells_known_ids"    bson:"spells_known_ids"`
	SpellsPreparedIDs []primitive.ObjectID `json:"spells_prepared_ids" bson:"spells_prepared_ids"`

	// Slots usados — el total disponible se deriva de Class.Levels[level].SpellSlots
	SlotsUsed [9]int `json:"slots_used" bson:"slots_used"`

	// Recursos consumibles (Imposición de Manos: 22/25, Canalizar Divinidad: 1/1, etc.)
	// Se snapshotean al subir de nivel desde Feature.Uses.
	Resources []ResourceState `json:"resources" bson:"resources"`

	// Ataques activos: combinaciones de arma/feature/spell. Embedded porque
	// el jugador los compone manualmente (ej. "Espada Larga +1 con Estilo Defensa").
	Attacks []Attack `json:"attacks" bson:"attacks"`

	// Estado en combate
	Conditions []string `json:"conditions" bson:"conditions"`

	// Inventario
	Inventory Inventory `json:"inventory" bson:"inventory"`

	// Roleplay
	Personality string `json:"personality" bson:"personality,omitempty"`
	Ideals      string `json:"ideals"      bson:"ideals,omitempty"`
	Bonds       string `json:"bonds"       bson:"bonds,omitempty"`
	Flaws       string `json:"flaws"       bson:"flaws,omitempty"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type HitDice struct {
	Current int    `json:"current" bson:"current"`
	Max     int    `json:"max"     bson:"max"`
	Type    string `json:"type"    bson:"type"`
}

type DeathSaves struct {
	Successes int `json:"successes" bson:"successes"`
	Failures  int `json:"failures"  bson:"failures"`
}

// AbilityScores: las 6 stats. Score base; los modificadores se calculan
// en el frontend con (Score-10)/2 redondeo abajo. Save y ProfSave se
// derivan de Class.SaveProficient + bonus de proficiencia del nivel.
type AbilityScores struct {
	FUE Ability `json:"FUE" bson:"FUE"`
	DES Ability `json:"DES" bson:"DES"`
	CON Ability `json:"CON" bson:"CON"`
	INT Ability `json:"INT" bson:"INT"`
	SAB Ability `json:"SAB" bson:"SAB"`
	CAR Ability `json:"CAR" bson:"CAR"`
}

type Ability struct {
	Score int `json:"score" bson:"score"`
}

// ResourceState es la "instancia" de un Feature.Uses sobre este personaje.
// FeatureID apunta al feature que lo origina (para refrescar al pasar
// descansos cortos/largos según Recharge).
type ResourceState struct {
	FeatureID primitive.ObjectID `json:"feature_id" bson:"feature_id,omitempty"`
	Name      string             `json:"name"       bson:"name"`
	Current   int                `json:"current"    bson:"current"`
	Max       int                `json:"max"        bson:"max"`
	Recharge  string             `json:"recharge"   bson:"recharge"`
}

// Attack: ataque ejecutable desde la hoja. Puede derivarse de un arma,
// una feature o un conjuro. Lo dejamos embebido porque cada PJ lo
// personaliza (mod a daño, descripción, etc.).
type Attack struct {
	Name        string `json:"name"        bson:"name"`
	Type        string `json:"type"        bson:"type"`        // "attack" | "save"
	Bonus       int    `json:"bonus"       bson:"bonus"`       // mod al d20
	Damage      string `json:"damage"      bson:"damage"`      // "1d8+4"
	DamageType  string `json:"damage_type" bson:"damage_type"` // "cortante", "radiante", ...
	Range       string `json:"range"       bson:"range"`
	Cost        string `json:"cost"        bson:"cost"` // "action" | "bonus" | "reaction" | "spell-slot"
	Description string `json:"description" bson:"description"`

	// Origen opcional: si el ataque proviene de un arma/feature/conjuro
	// del catálogo, lo enlazamos para auditoría.
	SourceKind string             `json:"source_kind" bson:"source_kind,omitempty"` // "weapon" | "feature" | "spell"
	SourceRef  primitive.ObjectID `json:"source_ref"  bson:"source_ref,omitempty"`
}

type Inventory struct {
	Armor         []string `json:"armor"         bson:"armor"`
	Weapons       []string `json:"weapons"       bson:"weapons"`
	Tools         []string `json:"tools"         bson:"tools"`
	Languages     []string `json:"languages"     bson:"languages"`
	Proficiencies []string `json:"proficiencies" bson:"proficiencies"`
	Items         []string `json:"items"         bson:"items"`
	Gold          int      `json:"gold"          bson:"gold"`
}

// CharacterExpanded: lo que devolvemos al frontend con el catálogo
// poblado. El backend lo arma a demanda (no se persiste así).
type CharacterExpanded struct {
	*Character
	Class    *Class     `json:"class,omitempty"`
	Subclass *Subclass  `json:"subclass,omitempty"`
	Race     *Race      `json:"race,omitempty"`
	Skills   []*Skill   `json:"skills_resolved,omitempty"` // todas las skills, marcando proficient/expert
	Cantrips []*Spell   `json:"cantrips_resolved,omitempty"`
	Spells   []*Spell   `json:"spells_resolved,omitempty"`
	Features []*Feature `json:"features_resolved,omitempty"` // features actuales según el nivel
}
