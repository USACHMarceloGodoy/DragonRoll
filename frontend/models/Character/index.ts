import type { Class, Feature, Race, Skill, Spell, Subclass } from '../Catalog'

export interface HitDice {
  current: number
  max: number
  type: 'd6' | 'd8' | 'd10' | 'd12'
}

export interface DeathSaves {
  successes: number
  failures: number
}

export interface Ability {
  score: number
}

export interface AbilityScores {
  FUE: Ability
  DES: Ability
  CON: Ability
  INT: Ability
  SAB: Ability
  CAR: Ability
}

export interface ResourceState {
  feature_id?: string
  name: string
  current: number
  max: number
  recharge: string
}

export interface Attack {
  name: string
  type: 'attack' | 'save'
  bonus: number
  damage: string
  damage_type: string
  range: string
  cost: 'action' | 'bonus' | 'reaction' | 'spell-slot'
  description: string
  source_kind?: 'weapon' | 'feature' | 'spell'
  source_ref?: string
}

export interface Inventory {
  armor: string[]
  weapons: string[]
  tools: string[]
  languages: string[]
  proficiencies: string[]
  items: string[]
  gold: number
}

/**
 * Character: lo que persistimos en Mongo. Las referencias son strings
 * (ObjectID en hex). Para la vista preferimos CharacterExpanded.
 */
export interface Character {
  _id?: string
  owner_id?: string
  campaign_id?: string
  active: boolean

  name: string
  pronouns?: string
  class_id: string
  subclass_id?: string
  race_id: string
  background?: string
  alignment?: string
  level: number
  xp: number
  next_xp: number

  pic?: string
  color?: string

  hp: number
  hp_max: number
  temp_hp: number
  ac: number
  speed: number
  initiative: number
  inspiration: boolean
  hit_dice: HitDice
  death_saves: DeathSaves

  stats: AbilityScores

  skill_proficiency_ids: string[]
  skill_expertise_ids: string[]

  cantrip_ids: string[]
  spells_known_ids: string[]
  spells_prepared_ids: string[]

  /** Slots usados por nivel del conjuro (índice 0 = nivel 1). */
  slots_used: number[]
  resources: ResourceState[]
  attacks: Attack[]

  conditions: string[]
  inventory: Inventory

  personality?: string
  ideals?: string
  bonds?: string
  flaws?: string

  created_at?: string
  updated_at?: string
}

/** Lo que devuelve /characters/:id?expand=true (default). */
export interface CharacterExpanded extends Character {
  class?: Class
  subclass?: Subclass
  race?: Race
  skills_resolved?: Skill[]
  cantrips_resolved?: Spell[]
  spells_resolved?: Spell[]
  features_resolved?: Feature[]
}

export type Stat = 'FUE' | 'DES' | 'CON' | 'INT' | 'SAB' | 'CAR'

export const STAT_NAMES: Record<Stat, string> = {
  FUE: 'Fuerza',
  DES: 'Destreza',
  CON: 'Constitución',
  INT: 'Inteligencia',
  SAB: 'Sabiduría',
  CAR: 'Carisma',
}

export const STATS: Stat[] = ['FUE', 'DES', 'CON', 'INT', 'SAB', 'CAR']

/** mod = floor((score - 10) / 2). */
export function abilityMod(score: number): number {
  return Math.floor((score - 10) / 2)
}

/** Bono de proficiencia desde la tabla de la clase al nivel actual. */
export function profBonus(klass: Class | undefined, level: number): number {
  if (!klass) return 2
  const step = klass.levels.find(l => l.level === level)
  return step?.proficiency_bonus ?? 2
}
