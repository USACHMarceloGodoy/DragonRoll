package services

import (
	"context"
	"dragonroll/backend/models"
	"dragonroll/backend/utils"
	"encoding/json"
	"errors"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Estructuras intermedias del archivo JSON: usan "tag" como ID lógico
// (legible) que el seeder resuelve a ObjectID al insertar.
type seedHomebrewFile struct {
	Features []seedFeature `json:"features"`
	Spells   []seedSpell   `json:"spells"`
	Races    []seedRace    `json:"races"`
	Classes  []seedClass   `json:"classes"`
}

type seedFeature struct {
	Tag         string              `json:"tag"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Origin      string              `json:"origin"`
	Cost        string              `json:"cost"`
	Uses        *models.FeatureUses `json:"uses,omitempty"`
	RollHint    string              `json:"roll_hint,omitempty"`
}

type seedSpell struct {
	Tag           string `json:"tag"`
	Name          string `json:"name"`
	Level         int    `json:"level"`
	School        string `json:"school"`
	CastingTime   string `json:"casting_time"`
	Range         string `json:"range"`
	Components    string `json:"components"`
	Duration      string `json:"duration"`
	Description   string `json:"description"`
	HigherLevel   string `json:"higher_level,omitempty"`
	Ritual        bool   `json:"ritual"`
	Concentration bool   `json:"concentration"`
	Damage        string `json:"damage,omitempty"`
	DamageType    string `json:"damage_type,omitempty"`
	SaveAbility   string `json:"save_ability,omitempty"`
}

type seedRace struct {
	Tag            string                `json:"tag"`
	Name           string                `json:"name"`
	ShortDesc      string                `json:"short_desc"`
	Speed          int                   `json:"speed"`
	Size           string                `json:"size"`
	AbilityBonuses []models.AbilityBonus `json:"ability_bonuses"`
	Languages      []string              `json:"languages"`
	TraitTags      []string              `json:"trait_tags"`
}

type seedClassCasting struct {
	Ability       string   `json:"ability"`
	Preparation   string   `json:"preparation"`
	Ritual        bool     `json:"ritual"`
	SpellListTags []string `json:"spell_list_tags"`
}

type seedClassLevel struct {
	Level            int      `json:"level"`
	ProficiencyBonus int      `json:"proficiency_bonus"`
	FeatureTags      []string `json:"feature_tags"`
	SpellSlots       [9]int   `json:"spell_slots"`
	CantripsKnown    int      `json:"cantrips_known"`
	SpellsKnown      int      `json:"spells_known"`
	Notes            string   `json:"notes,omitempty"`
}

type seedClass struct {
	Tag              string            `json:"tag"`
	Name             string            `json:"name"`
	ShortDesc        string            `json:"short_desc"`
	HitDie           string            `json:"hit_die"`
	PrimaryStats     []string          `json:"primary_stats"`
	SaveProficient   []string          `json:"save_proficient"`
	ArmorProfs       []string          `json:"armor_profs"`
	WeaponProfs      []string          `json:"weapon_profs"`
	ToolProfs        []string          `json:"tool_profs"`
	SkillChoiceCount int               `json:"skill_choice_count"`
	SkillChoicePool  []string          `json:"skill_choice_pool"` // codes de skills
	SubclassLevel    int               `json:"subclass_level"`
	Casting          *seedClassCasting `json:"casting,omitempty"`
	Levels           []seedClassLevel  `json:"levels"`
}

// SeedHomebrewExamples carga el archivo si la base no tiene clases.
// Es opcional: el admin puede borrar el archivo si no lo quiere.
func SeedClassesFromFile(path, sourceLabel string, homebrew bool) error {
	count, err := classRepo.Collection.CountDocuments(context.Background(), bson.M{"meta.source_label": sourceLabel})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	// path se recibe como parámetro
	data, err := os.ReadFile(path)
	if err != nil {
		utils.Info("BootstrapSeed: " + path + " no encontrado, se omite")
		return nil
	}
	var file seedHomebrewFile
	if err := json.Unmarshal(data, &file); err != nil {
		return errors.New("seed homebrew parse: " + err.Error())
	}

	now := time.Now()
	meta := models.CatalogMeta{
		Scope:       models.ScopeGlobal,
		IsHomebrew:  homebrew,
		SourceLabel: sourceLabel,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// 1) Insertar features y construir mapa tag → ObjectID
	featureIDs := map[string]primitive.ObjectID{}
	for _, f := range file.Features {
		doc := models.Feature{
			Name:        f.Name,
			Description: f.Description,
			Origin:      f.Origin,
			Cost:        f.Cost,
			Uses:        f.Uses,
			RollHint:    f.RollHint,
			Meta:        meta,
		}
		id, err := featureRepo.InsertOne(&doc)
		if err != nil {
			return err
		}
		featureIDs[f.Tag] = id
	}

	// 2) Spells
	spellIDs := map[string]primitive.ObjectID{}
	for _, s := range file.Spells {
		doc := models.Spell{
			Name: s.Name, Level: s.Level, School: s.School,
			CastingTime: s.CastingTime, Range: s.Range, Components: s.Components,
			Duration: s.Duration, Description: s.Description, HigherLevel: s.HigherLevel,
			Ritual: s.Ritual, Concentration: s.Concentration,
			Damage: s.Damage, DamageType: s.DamageType, SaveAbility: s.SaveAbility,
			Meta: meta,
		}
		id, err := spellRepo.InsertOne(&doc)
		if err != nil {
			return err
		}
		spellIDs[s.Tag] = id
	}

	// 3) Skills (por code) — para resolver el pool de la clase
	skillIDsByCode := map[string]primitive.ObjectID{}
	skillRaws, err := skillRepo.FindManyRaw(bson.M{})
	if err == nil {
		for _, r := range skillRaws {
			var sk models.Skill
			if err := bson.Unmarshal(r, &sk); err == nil {
				skillIDsByCode[sk.Code] = sk.ID
			}
		}
	}

	// 4) Races
	for _, r := range file.Races {
		traitIDs := []primitive.ObjectID{}
		for _, t := range r.TraitTags {
			if id, ok := featureIDs[t]; ok {
				traitIDs = append(traitIDs, id)
			}
		}
		doc := models.Race{
			Name: r.Name, ShortDesc: r.ShortDesc,
			Speed: r.Speed, Size: r.Size,
			AbilityBonuses: r.AbilityBonuses, Languages: r.Languages,
			TraitIDs: traitIDs,
			Meta:     meta,
		}
		if _, err := raceRepo.InsertOne(&doc); err != nil {
			return err
		}
	}

	// 5) Classes
	for _, c := range file.Classes {
		var skillPool []primitive.ObjectID
		for _, code := range c.SkillChoicePool {
			if id, ok := skillIDsByCode[code]; ok {
				skillPool = append(skillPool, id)
			}
		}

		levels := make([]models.LevelStep, 0, len(c.Levels))
		for _, lv := range c.Levels {
			fids := []primitive.ObjectID{}
			for _, t := range lv.FeatureTags {
				if id, ok := featureIDs[t]; ok {
					fids = append(fids, id)
				}
			}
			levels = append(levels, models.LevelStep{
				Level:            lv.Level,
				ProficiencyBonus: lv.ProficiencyBonus,
				FeatureIDs:       fids,
				SpellSlots:       lv.SpellSlots,
				CantripsKnown:    lv.CantripsKnown,
				SpellsKnown:      lv.SpellsKnown,
				Notes:            lv.Notes,
			})
		}

		var casting *models.CastingInfo
		if c.Casting != nil {
			spellList := []primitive.ObjectID{}
			for _, t := range c.Casting.SpellListTags {
				if id, ok := spellIDs[t]; ok {
					spellList = append(spellList, id)
				}
			}
			casting = &models.CastingInfo{
				Ability:      c.Casting.Ability,
				Preparation:  c.Casting.Preparation,
				Ritual:       c.Casting.Ritual,
				SpellListIDs: spellList,
			}
		}

		doc := models.Class{
			Name: c.Name, ShortDesc: c.ShortDesc, HitDie: c.HitDie,
			PrimaryStats: c.PrimaryStats, SaveProficient: c.SaveProficient,
			ArmorProficiencies: c.ArmorProfs, WeaponProficiencies: c.WeaponProfs, ToolProficiencies: c.ToolProfs,
			SkillChoice:   models.SkillChoice{Count: c.SkillChoiceCount, FromSkillIDs: skillPool},
			SubclassLevel: c.SubclassLevel,
			Levels:        levels,
			Casting:       casting,
			Meta:          meta,
		}
		if _, err := classRepo.InsertOne(&doc); err != nil {
			return err
		}
	}

	utils.Info("BootstrapSeed: cargado " + path + " (" + sourceLabel + ")")
	return nil
}

// SeedHomebrewExamples es retrocompat: wrapper sobre SeedClassesFromFile.
func SeedHomebrewExamples() error {
	return SeedClassesFromFile("data/seed/homebrew_examples.json", "homebrew-examples", true)
}
