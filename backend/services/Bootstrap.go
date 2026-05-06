package services

import (
	"context"
	"dragonroll/backend/config"
	"dragonroll/backend/models"
	"dragonroll/backend/utils"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// BootstrapSeed se ejecuta en el arranque. Es idempotente.
//
//   - Si users está vacío, crea un admin inicial leyendo ADMIN_USERNAME
//     y ADMIN_PASSWORD del entorno (con defaults para dev).
//   - Si skills está vacío, carga las 18 mecánicas universales del sistema
//     d20 desde data/seed/skills.json (ideas, no copyrightables).
//   - Si features-pasivos universales (condiciones) están vacíos, no hacemos
//     nada por ahora — las condiciones se modelarán como labels en string,
//     basta con que el frontend tenga la lista cerrada.
func BootstrapSeed() error {
	if err := seedAdmin(); err != nil {
		return err
	}
	if err := seedSkillsIfEmpty(); err != nil {
		return err
	}
	// Cargar clases SRD ANTES que ejemplos homebrew para que las skills
	// referenciadas estén ya en la base.
	if err := SeedClassesFromFile("data/seed/srd_classes.json", "SRD 5.1", false); err != nil {
		return err
	}
	if err := SeedClassesFromFile("data/seed/homebrew_examples.json", "homebrew-examples", true); err != nil {
		return err
	}
	if err := SeedSRDSpellsIfEmpty(); err != nil {
		return err
	}
	return nil
}

func seedAdmin() error {
	count, err := userRepo.Collection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	username := strings.ToLower(utils.GetEnvWithDefault("ADMIN_USERNAME", "admin"))
	email := strings.ToLower(utils.GetEnvWithDefault("ADMIN_EMAIL", "admin@dragonroll.local"))
	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		// Sólo en debug aceptamos un default; en release exigimos password explícita.
		if config.Cfg.Server.Mode == "release" {
			return errors.New("BootstrapSeed: ADMIN_PASSWORD obligatorio en release")
		}
		password = "change-me-please"
		utils.Info("BootstrapSeed: ADMIN_PASSWORD no definido, usando default de desarrollo")
	}

	admin := &models.User{
		Username:  username,
		Email:     email,
		Password:  password,
		Roles:     []models.Rol{models.ADMIN, models.DM, models.PLAYER},
		ActiveRol: models.ADMIN,
		CreatedAt: time.Now(),
	}
	if _, err := CreateUserService(admin); err != nil {
		return errors.New("seed admin: " + err.Error())
	}
	utils.Info("BootstrapSeed: admin creado con username=" + username)
	return nil
}

// seedSkillsFile es la representación JSON que esperamos en data/seed/skills.json.
type seedSkillsFile struct {
	Skills []models.Skill `json:"skills"`
}

func seedSkillsIfEmpty() error {
	count, err := skillRepo.Collection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	path := utils.GetEnvWithDefault("SEED_SKILLS_PATH", "data/seed/skills.json")
	data, err := os.ReadFile(path)
	if err != nil {
		utils.Info("BootstrapSeed: skills seed no encontrado en " + path + " (se omite)")
		return nil
	}

	var file seedSkillsFile
	if err := json.Unmarshal(data, &file); err != nil {
		return errors.New("seed skills parse: " + err.Error())
	}

	now := time.Now()
	for i := range file.Skills {
		s := file.Skills[i]
		s.Meta.Scope = models.ScopeGlobal
		s.Meta.IsHomebrew = false
		s.Meta.SourceLabel = "system-d20"
		s.Meta.CreatedAt = now
		s.Meta.UpdatedAt = now
		if _, err := skillRepo.InsertOne(&s); err != nil {
			return err
		}
	}
	utils.Info("BootstrapSeed: cargadas las skills universales")
	return nil
}
