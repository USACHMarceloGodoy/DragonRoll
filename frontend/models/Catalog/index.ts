// Tipos compartidos del catálogo. Espejan models/*.go del backend.
// Mantenerlos sincronizados es manual; cuando crezca conviene generarlos
// con un go-to-typescript codegen, pero por ahora la duplicación es OK.

export type Scope = 'global' | 'campaign' | 'user'

export interface CatalogMeta {
  scope: Scope
  campaign_id?: string
  owner_id?: string
  is_homebrew: boolean
  source_label?: string
  created_at: string
  updated_at: string
}

// --- Skill --------------------------------------------------------
export interface Skill {
  _id: string
  name: string
  code: string
  stat: 'FUE' | 'DES' | 'CON' | 'INT' | 'SAB' | 'CAR'
  description?: string
  meta: CatalogMeta
}

// --- Spell --------------------------------------------------------
export interface Spell {
  _id: string
  name: string
  level: number             // 0..9 (0 = truco)
  school: string
  casting_time: string
  range: string
  components: string
  duration: string
  description: string
  higher_level?: string
  ritual: boolean
  concentration: boolean
  damage?: string
  damage_type?: string
  save_ability?: string
  meta: CatalogMeta
}

// --- Feature ------------------------------------------------------
export interface FeatureUses {
  max: number
  formula?: string
  recharge: 'short_rest' | 'long_rest' | 'per_day' | 'none'
}

export interface Feature {
  _id: string
  name: string
  description: string
  origin: 'class' | 'subclass' | 'race' | 'background' | 'feat' | 'item'
  origin_ref?: string
  cost: 'action' | 'bonus' | 'reaction' | 'passive' | 'spell-slot' | 'free'
  uses?: FeatureUses
  roll_hint?: string
  meta: CatalogMeta
}

// --- Class --------------------------------------------------------
export interface SkillChoice {
  count: number
  from_skill_ids: string[]
}

export interface CastingInfo {
  ability: 'INT' | 'SAB' | 'CAR'
  preparation: 'prepared' | 'known' | 'spellbook' | 'pact'
  ritual: boolean
  spell_list_ids: string[]
}

export interface LevelStep {
  level: number
  proficiency_bonus: number
  feature_ids: string[]
  /**
   * Tabla de slots por nivel del conjuro. Índice 0 = nivel 1 del conjuro.
   * Largo siempre 9.
   */
  spell_slots: number[]
  cantrips_known: number
  spells_known: number
  notes?: string
}

export interface Class {
  _id: string
  name: string
  short_desc: string
  hit_die: 'd6' | 'd8' | 'd10' | 'd12'
  primary_stats: string[]
  save_proficient: string[]
  armor_profs: string[]
  weapon_profs: string[]
  tool_profs: string[]
  skill_choice: SkillChoice
  subclass_level: number
  subclasses: string[]
  levels: LevelStep[]
  casting?: CastingInfo
  meta: CatalogMeta
}

// --- Subclass -----------------------------------------------------
export interface GrantedSpell {
  level: number
  spell_id: string
}

export interface Subclass {
  _id: string
  parent_class: string
  name: string
  short_desc: string
  flavor?: string
  levels: LevelStep[]
  granted_spells: GrantedSpell[]
  meta: CatalogMeta
}

// --- Race ---------------------------------------------------------
export interface AbilityBonus {
  stat: string
  bonus: number
}

export interface Race {
  _id: string
  name: string
  short_desc: string
  speed: number
  size: 'Pequeño' | 'Mediano' | 'Grande'
  ability_bonuses: AbilityBonus[]
  languages: string[]
  trait_ids: string[]
  parent_race?: string
  meta: CatalogMeta
}

export type CatalogKind = 'classes' | 'subclasses' | 'spells' | 'features' | 'races' | 'skills'
