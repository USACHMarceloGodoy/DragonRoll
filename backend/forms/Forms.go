package forms

import "dragonroll/backend/models"

type UserForm struct {
	Username string       `json:"username" binding:"required,min=3,max=24"`
	Email    string       `json:"email"    binding:"required,email"`
	Password string       `json:"password" binding:"required,min=8"`
	Roles    []models.Rol `json:"roles"`
	Avatar   string       `json:"avatar"`
}

type SetActiveRoleForm struct {
	ActiveRol models.Rol `json:"active_rol" binding:"required"`
}

type UpdateProfileForm struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

type ChangePasswordForm struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password"     binding:"required,min=8"`
}

type AssignCampaignForm struct {
	CampaignID string `json:"campaign_id"` // hex; vacío = desasociar
}

type CharacterForm struct {
	Name       string `json:"name"        binding:"required"`
	Pronouns   string `json:"pronouns"`
	ClassID    string `json:"class_id"    binding:"required"` // hex
	SubclassID string `json:"subclass_id"`                    // hex, opcional
	RaceID     string `json:"race_id"     binding:"required"` // hex
	CampaignID string `json:"campaign_id"`                    // hex, opcional
	Background string `json:"background"`
	Alignment  string `json:"alignment"`
	Level      int    `json:"level"       binding:"min=1,max=20"`
	XP         int    `json:"xp"`
	NextXP     int    `json:"next_xp"`
	Pic        string `json:"pic"`
	Color      string `json:"color"`

	HP          int               `json:"hp"`
	HPMax       int               `json:"hp_max"`
	TempHP      int               `json:"temp_hp"`
	AC          int               `json:"ac"`
	Speed       int               `json:"speed"`
	Initiative  int               `json:"initiative"`
	Inspiration bool              `json:"inspiration"`
	HitDice     models.HitDice    `json:"hit_dice"`
	DeathSaves  models.DeathSaves `json:"death_saves"`

	Stats models.AbilityScores `json:"stats"`

	SkillProficiencyIDs []string `json:"skill_proficiency_ids"`
	SkillExpertiseIDs   []string `json:"skill_expertise_ids"`

	CantripIDs        []string `json:"cantrip_ids"`
	SpellsKnownIDs    []string `json:"spells_known_ids"`
	SpellsPreparedIDs []string `json:"spells_prepared_ids"`

	SlotsUsed [9]int                 `json:"slots_used"`
	Resources []models.ResourceState `json:"resources"`
	Attacks   []models.Attack        `json:"attacks"`

	Conditions []string         `json:"conditions"`
	Inventory  models.Inventory `json:"inventory"`

	Personality string `json:"personality"`
	Ideals      string `json:"ideals"`
	Bonds       string `json:"bonds"`
	Flaws       string `json:"flaws"`
}

type CombatPatchForm struct {
	HP         *int     `json:"hp,omitempty"`
	TempHP     *int     `json:"temp_hp,omitempty"`
	Conditions []string `json:"conditions,omitempty"`
	SlotsUsed  *[9]int  `json:"slots_used,omitempty"`
}

type CampaignForm struct {
	Name        string `json:"name"        binding:"required"`
	Description string `json:"description"`
	System      string `json:"system"`
}

type AddPlayerForm struct {
	UsernameOrEmail string `json:"user" binding:"required"`
}
