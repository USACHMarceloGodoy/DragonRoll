// Espejo TS de models/Campaign.go.

export interface NPC {
  id: string
  name: string
  role: string
  faction?: string
  location?: string
  alive: boolean
  attitude: 'aliado' | 'neutral' | 'sospechoso' | 'enemigo'
  pic?: string
  color?: string
  summary: string
  details?: string
  tags: string[]
  comments: NPCNote[]
}

export interface NPCNote {
  author: string
  avatar?: string
  text: string
}

export interface Location {
  id: string
  name: string
  type: string
  description: string
  tags: string[]
}

export interface SessionEntry {
  number: number
  date: string
  title: string
  summary: string
}

export interface Campaign {
  _id?: string
  dm_id: string
  player_ids: string[]
  name: string
  description?: string
  system: string
  npcs: NPC[]
  locations: Location[]
  sessions: SessionEntry[]
  active: boolean
  created_at?: string
  updated_at?: string
}
