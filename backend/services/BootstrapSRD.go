package services

import (
	"context"
	"dragonroll/backend/models"
	"dragonroll/backend/utils"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// SeedSRDSpellsIfEmpty carga `data/seed/srd_spells.json` (generado por
// scripts/parse_srd.py) si la colección de spells está vacía o sólo
// tiene los homebrew de ejemplo. Es idempotente por nombre.
//
// Atribución: el SRD 5.1 está bajo CC-BY-4.0. Marcamos meta.source_label
// = "SRD 5.1" en cada documento para que el frontend pueda mostrar la
// nota de atribución requerida por la licencia.
func SeedSRDSpellsIfEmpty() error {
	path := utils.GetEnvWithDefault("SEED_SRD_SPELLS_PATH", "data/seed/srd_spells.json")
	data, err := os.ReadFile(path)
	if err != nil {
		// Archivo opcional: no es error si no existe.
		return nil
	}

	// Si ya hay >50 hechizos asumimos que el SRD ya se cargó.
	count, err := spellRepo.Collection.CountDocuments(context.Background(), bson.M{"meta.source_label": "SRD 5.1"})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	var file struct {
		Spells []seedSpell `json:"spells"`
	}
	if err := json.Unmarshal(data, &file); err != nil {
		return errors.New("seed SRD parse: " + err.Error())
	}

	now := time.Now()
	meta := models.CatalogMeta{
		Scope:       models.ScopeGlobal,
		IsHomebrew:  false,
		SourceLabel: "SRD 5.1",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	for _, s := range file.Spells {
		doc := models.Spell{
			Name: s.Name, Level: s.Level, School: s.School,
			CastingTime: s.CastingTime, Range: s.Range, Components: s.Components,
			Duration: s.Duration, Description: s.Description, HigherLevel: s.HigherLevel,
			Ritual: s.Ritual, Concentration: s.Concentration,
			Damage: s.Damage, DamageType: s.DamageType, SaveAbility: s.SaveAbility,
			Meta: meta,
		}
		if _, err := spellRepo.InsertOne(&doc); err != nil {
			return err
		}
	}
	utils.Info("BootstrapSeed: cargados " + strconv.Itoa(len(file.Spells)) + " hechizos del SRD 5.1")
	return nil
}
