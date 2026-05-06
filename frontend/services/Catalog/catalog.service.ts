import type { CatalogKind, Class, Feature, Race, Skill, Spell, Subclass } from '~/models/Catalog'
import type { ApiResponse } from '~/models/apiResponse'

/**
 * Service genérico para los 6 catálogos. Tipado con overloads para que
 * `list('classes')` devuelva `Class[]` y `list('spells')` devuelva
 * `Spell[]` sin gimnasia de casteo en cada llamada.
 */
export class CatalogService {
  private get api() {
    return useNuxtApp().$api as typeof $fetch
  }

  async list(kind: 'classes'): Promise<Class[]>
  async list(kind: 'subclasses'): Promise<Subclass[]>
  async list(kind: 'spells'): Promise<Spell[]>
  async list(kind: 'features'): Promise<Feature[]>
  async list(kind: 'races'): Promise<Race[]>
  async list(kind: 'skills'): Promise<Skill[]>
  async list(kind: CatalogKind): Promise<unknown[]> {
    const res = await this.api<ApiResponse<unknown[]>>(`/${kind}`, { method: 'GET' })
    return res.data ?? []
  }

  async get(kind: 'classes', id: string): Promise<Class>
  async get(kind: 'subclasses', id: string): Promise<Subclass>
  async get(kind: 'spells', id: string): Promise<Spell>
  async get(kind: 'features', id: string): Promise<Feature>
  async get(kind: 'races', id: string): Promise<Race>
  async get(kind: 'skills', id: string): Promise<Skill>
  async get(kind: CatalogKind, id: string): Promise<unknown> {
    const res = await this.api<ApiResponse<unknown>>(`/${kind}/${id}`, { method: 'GET' })
    if (!res.data) throw new Error(res.message || 'No se encontró el documento')
    return res.data
  }

  /** Crea un documento de catálogo. El body lleva la Meta con scope/campaign_id/source_label. */
  async create<T = unknown>(kind: CatalogKind, payload: T): Promise<{ _id: string }> {
    const res = await this.api<ApiResponse<{ _id: string }>>(`/${kind}`, {
      method: 'POST',
      body: payload as object,
    })
    if (!res.data) throw new Error(res.message || 'No se pudo crear')
    return res.data
  }

  async update<T = unknown>(kind: CatalogKind, id: string, payload: T): Promise<void> {
    await this.api(`/${kind}/${id}`, { method: 'PUT', body: payload as object })
  }

  async remove(kind: CatalogKind, id: string): Promise<void> {
    await this.api(`/${kind}/${id}`, { method: 'DELETE' })
  }
}
