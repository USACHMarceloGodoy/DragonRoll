package services

import (
	"dragonroll/backend/models"
	"dragonroll/backend/repositories"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateCharacterService(ac AuthContext, c *models.Character) (*models.Character, error) {
	if c.Name == "" {
		return nil, errors.New("nombre obligatorio")
	}
	if c.ClassID.IsZero() {
		return nil, errors.New("class_id obligatorio")
	}
	if c.RaceID.IsZero() {
		return nil, errors.New("race_id obligatorio")
	}
	if c.Level <= 0 {
		c.Level = 1
	}
	if c.HPMax > 0 && c.HP == 0 {
		c.HP = c.HPMax
	}
	if c.OwnerID.IsZero() {
		c.OwnerID = ac.UserID
	}
	c.Active = true
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now

	id, err := characterRepo.InsertOne(c)
	if err != nil {
		return nil, err
	}
	c.ID = id
	return c, nil
}

// ListAccessibleCharacters: para la home del usuario. Le mostramos sus
// PJs + los de las campañas donde es DM o jugador.
func ListAccessibleCharacters(ac AuthContext) ([]*models.Character, error) {
	dm, player, err := MemberCampaigns(ac.UserID)
	if err != nil {
		return nil, err
	}
	return characterRepo.FindAccessible(ac.UserID, dm, player)
}

func GetCharacterByIDAuthorized(ac AuthContext, id primitive.ObjectID) (*models.Character, error) {
	c, err := characterRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, errors.New("personaje no encontrado")
	}
	ok, err := CanViewCharacter(ac, c)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("no autorizado")
	}
	return c, nil
}

// GetCharacterExpanded hace populate del catálogo y devuelve el documento listo
// para que el frontend renderice la hoja sin pedir nada más.
func GetCharacterExpanded(ac AuthContext, id primitive.ObjectID) (*models.CharacterExpanded, error) {
	c, err := GetCharacterByIDAuthorized(ac, id)
	if err != nil {
		return nil, err
	}
	exp := &models.CharacterExpanded{Character: c}

	if !c.ClassID.IsZero() {
		var class models.Class
		if err := classRepo.FindByID(c.ClassID, &class); err == nil && class.Name != "" {
			exp.Class = &class
		}
	}
	if !c.SubclassID.IsZero() {
		var sub models.Subclass
		if err := subclassRepo.FindByID(c.SubclassID, &sub); err == nil && sub.Name != "" {
			exp.Subclass = &sub
		}
	}
	if !c.RaceID.IsZero() {
		var race models.Race
		if err := raceRepo.FindByID(c.RaceID, &race); err == nil && race.Name != "" {
			exp.Race = &race
		}
	}

	// Skills: traemos TODAS las skills accesibles para que el frontend
	// pueda renderizar la lista completa marcando proficient/expert.
	dm, player, _ := MemberCampaigns(ac.UserID)
	visFilter := repositories.VisibilityFilter(ac.UserID, append(dm, player...))
	if rawSkills, err := skillRepo.FindManyRaw(visFilter); err == nil {
		var skills []*models.Skill
		for _, raw := range rawSkills {
			var s models.Skill
			if err := bson.Unmarshal(raw, &s); err == nil {
				skills = append(skills, &s)
			}
		}
		exp.Skills = skills
	}

	// Cantrips y conjuros conocidos
	if rawSpells, err := spellRepo.FindManyByIDs(c.CantripIDs); err == nil {
		exp.Cantrips = unmarshalSpells(rawSpells)
	}
	if rawSpells, err := spellRepo.FindManyByIDs(c.SpellsKnownIDs); err == nil {
		exp.Spells = unmarshalSpells(rawSpells)
	}

	// Features actuales: las que da la clase hasta el nivel actual + las
	// de la subclase hasta el nivel actual + las raciales.
	exp.Features = collectActiveFeatures(c, exp.Class, exp.Subclass, exp.Race)

	return exp, nil
}

func unmarshalSpells(raws []bson.Raw) []*models.Spell {
	var out []*models.Spell
	for _, r := range raws {
		var s models.Spell
		if err := bson.Unmarshal(r, &s); err == nil {
			out = append(out, &s)
		}
	}
	return out
}

func collectActiveFeatures(c *models.Character, class *models.Class, sub *models.Subclass, race *models.Race) []*models.Feature {
	var ids []primitive.ObjectID
	if class != nil {
		for _, step := range class.Levels {
			if step.Level <= c.Level {
				ids = append(ids, step.FeatureIDs...)
			}
		}
	}
	if sub != nil {
		for _, step := range sub.Levels {
			if step.Level <= c.Level {
				ids = append(ids, step.FeatureIDs...)
			}
		}
	}
	if race != nil {
		ids = append(ids, race.TraitIDs...)
	}
	if len(ids) == 0 {
		return nil
	}
	raws, err := featureRepo.FindManyByIDs(ids)
	if err != nil {
		return nil
	}
	var out []*models.Feature
	for _, r := range raws {
		var f models.Feature
		if err := bson.Unmarshal(r, &f); err == nil {
			out = append(out, &f)
		}
	}
	return out
}

func UpdateCharacterService(ac AuthContext, id primitive.ObjectID, updated *models.Character) (*models.Character, error) {
	current, err := characterRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, errors.New("personaje no encontrado")
	}
	ok, err := CanEditCharacter(ac, current)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("no autorizado para editar este personaje")
	}
	updated.ID = current.ID
	updated.OwnerID = current.OwnerID // ownership es inmutable acá
	updated.CreatedAt = current.CreatedAt
	updated.UpdatedAt = time.Now()
	if err := characterRepo.UpdateOne(id, updated); err != nil {
		return nil, err
	}
	return updated, nil
}

// PatchCombat: HP/conditions/slots — útil durante una sesión sin reescribir todo.
func PatchCombat(ac AuthContext, id primitive.ObjectID, fields bson.M) error {
	current, err := characterRepo.FindByID(id)
	if err != nil {
		return err
	}
	if current == nil {
		return errors.New("personaje no encontrado")
	}
	ok, err := CanEditCharacter(ac, current)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("no autorizado")
	}
	fields["updated_at"] = time.Now()
	return characterRepo.PatchFields(id, fields)
}

func DeleteCharacterService(ac AuthContext, id primitive.ObjectID) error {
	current, err := characterRepo.FindByID(id)
	if err != nil {
		return err
	}
	if current == nil {
		return errors.New("personaje no encontrado")
	}
	ok, err := CanEditCharacter(ac, current)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("no autorizado")
	}
	return characterRepo.DeleteOne(id)
}

// AssignCharacterToCampaign asocia (o desasocia con campaignID = primitive.NilObjectID)
// un PJ a una campaña. Solo el dueño del PJ puede asociarlo.
// Para asociar, el usuario debe ser miembro (jugador o DM) de la campaña destino.
func AssignCharacterToCampaign(ac AuthContext, charID primitive.ObjectID, campaignID primitive.ObjectID) error {
	current, err := characterRepo.FindByID(charID)
	if err != nil {
		return err
	}
	if current == nil {
		return errors.New("personaje no encontrado")
	}
	if !ac.IsAdmin() && current.OwnerID != ac.UserID {
		return errors.New("solo el dueño del personaje puede asociarlo a una campaña")
	}
	if !campaignID.IsZero() {
		camp, err := campaignRepo.FindByID(campaignID)
		if err != nil || camp == nil {
			return errors.New("campaña destino no encontrada")
		}
		// El user debe ser miembro de la campaña.
		if camp.DMID != ac.UserID && !ac.IsAdmin() {
			isPlayer := false
			for _, p := range camp.PlayerIDs {
				if p == ac.UserID {
					isPlayer = true
					break
				}
			}
			if !isPlayer {
				return errors.New("no eres miembro de esa campaña")
			}
		}
	}
	return characterRepo.PatchFields(charID, bson.M{
		"campaign_id": campaignID,
		"updated_at":  time.Now(),
	})
}

// ExportCharacterJSON devuelve una copia "portable" del PJ. La idea:
// el JSON exportado puede ser reimportado en otra instancia de DragonRoll
// (por ejemplo si Marcelo cambia de servidor o un jugador quiere migrar).
//
// Nota importante: el PJ apunta al catálogo por ObjectID. Esos IDs no
// existen en la otra instancia. Por eso el export INCLUYE los nombres
// (class_name, race_name, etc.) — el importador hace match por nombre
// contra el catálogo destino. Si no encuentra match, queda como string
// pendiente que el usuario debe resolver manualmente en la UI.
type CharacterPortable struct {
	*models.Character
	ClassName    string `json:"class_name"`
	SubclassName string `json:"subclass_name,omitempty"`
	RaceName     string `json:"race_name"`
	// Nombres de skills/spells/features por ID, para resolver al importar.
	SkillProfNames      []string `json:"skill_proficiency_names"`
	SkillExpertiseNames []string `json:"skill_expertise_names"`
	CantripNames        []string `json:"cantrip_names"`
	SpellsKnownNames    []string `json:"spells_known_names"`
	SpellsPreparedNames []string `json:"spells_prepared_names"`
}

func ExportCharacterPortable(ac AuthContext, charID primitive.ObjectID) (*CharacterPortable, error) {
	c, err := GetCharacterByIDAuthorized(ac, charID)
	if err != nil {
		return nil, err
	}
	out := &CharacterPortable{Character: c}

	if !c.ClassID.IsZero() {
		var cls models.Class
		_ = classRepo.FindByID(c.ClassID, &cls)
		out.ClassName = cls.Name
	}
	if !c.SubclassID.IsZero() {
		var sub models.Subclass
		_ = subclassRepo.FindByID(c.SubclassID, &sub)
		out.SubclassName = sub.Name
	}
	if !c.RaceID.IsZero() {
		var race models.Race
		_ = raceRepo.FindByID(c.RaceID, &race)
		out.RaceName = race.Name
	}

	out.SkillProfNames = resolveNamesByIDs(skillRepo, c.SkillProficiencyIDs)
	out.SkillExpertiseNames = resolveNamesByIDs(skillRepo, c.SkillExpertiseIDs)
	out.CantripNames = resolveNamesByIDs(spellRepo, c.CantripIDs)
	out.SpellsKnownNames = resolveNamesByIDs(spellRepo, c.SpellsKnownIDs)
	out.SpellsPreparedNames = resolveNamesByIDs(spellRepo, c.SpellsPreparedIDs)

	return out, nil
}

// resolveNamesByIDs hace lookup batch de los nombres dado un set de IDs.
// Si un ID ya no existe en el catálogo, se omite del array.
func resolveNamesByIDs(repo *repositories.CatalogRepository, ids []primitive.ObjectID) []string {
	if len(ids) == 0 {
		return nil
	}
	raws, err := repo.FindManyByIDs(ids)
	if err != nil {
		return nil
	}
	out := make([]string, 0, len(raws))
	for _, raw := range raws {
		var probe struct {
			Name string `bson:"name"`
		}
		if err := bson.Unmarshal(raw, &probe); err == nil && probe.Name != "" {
			out = append(out, probe.Name)
		}
	}
	return out
}

// ImportCharacterPortable crea un PJ nuevo a partir del JSON exportado.
// Hace match por nombre contra el catálogo del servidor destino.
// Lo que no resuelve queda como referencia vacía — el usuario completará
// manualmente desde la UI.
func ImportCharacterPortable(ac AuthContext, p *CharacterPortable) (*models.Character, error) {
	if p == nil || p.Character == nil {
		return nil, errors.New("payload inválido")
	}
	c := p.Character
	// Reseteamos identidad — es un PJ nuevo en este servidor.
	c.ID = primitive.NilObjectID
	c.OwnerID = ac.UserID
	c.CampaignID = primitive.NilObjectID

	// Match por nombre.
	c.ClassID = lookupIDByName(classRepo, p.ClassName, ac)
	c.SubclassID = lookupIDByName(subclassRepo, p.SubclassName, ac)
	c.RaceID = lookupIDByName(raceRepo, p.RaceName, ac)
	c.SkillProficiencyIDs = lookupIDsByNames(skillRepo, p.SkillProfNames, ac)
	c.SkillExpertiseIDs = lookupIDsByNames(skillRepo, p.SkillExpertiseNames, ac)
	c.CantripIDs = lookupIDsByNames(spellRepo, p.CantripNames, ac)
	c.SpellsKnownIDs = lookupIDsByNames(spellRepo, p.SpellsKnownNames, ac)
	c.SpellsPreparedIDs = lookupIDsByNames(spellRepo, p.SpellsPreparedNames, ac)

	// La clase es obligatoria — si no existe, abortamos.
	if c.ClassID.IsZero() && p.ClassName != "" {
		return nil, errors.New("la clase '" + p.ClassName + "' no existe en este servidor; créala antes de importar")
	}
	if c.RaceID.IsZero() && p.RaceName != "" {
		return nil, errors.New("la raza '" + p.RaceName + "' no existe en este servidor; créala antes de importar")
	}

	return CreateCharacterService(ac, c)
}

// lookupIDByName hace match case-insensitive contra el catálogo. Solo
// considera entradas visibles para el usuario.
func lookupIDByName(repo *repositories.CatalogRepository, name string, ac AuthContext) primitive.ObjectID {
	if name == "" {
		return primitive.NilObjectID
	}
	dm, player, _ := MemberCampaigns(ac.UserID)
	visible := repositories.VisibilityFilter(ac.UserID, append(dm, player...))
	visible["name"] = bson.M{"$regex": "^" + regexpQuoteMeta(name) + "$", "$options": "i"}
	raws, err := repo.FindManyRaw(visible)
	if err != nil || len(raws) == 0 {
		return primitive.NilObjectID
	}
	var probe struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	if err := bson.Unmarshal(raws[0], &probe); err != nil {
		return primitive.NilObjectID
	}
	return probe.ID
}

func lookupIDsByNames(repo *repositories.CatalogRepository, names []string, ac AuthContext) []primitive.ObjectID {
	out := make([]primitive.ObjectID, 0, len(names))
	for _, n := range names {
		id := lookupIDByName(repo, n, ac)
		if !id.IsZero() {
			out = append(out, id)
		}
	}
	return out
}

// regexpQuoteMeta es un escape mínimo para regex de Mongo. No usamos
// regexp.QuoteMeta del stdlib porque tiene reglas distintas en algunos chars.
func regexpQuoteMeta(s string) string {
	special := `\.+*?()|[]{}^$`
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		for j := 0; j < len(special); j++ {
			if s[i] == special[j] {
				out = append(out, '\\')
				break
			}
		}
		out = append(out, s[i])
	}
	return string(out)
}
