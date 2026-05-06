import type { Character, CharacterExpanded } from '~/models/Character'
import type { ApiResponse } from '~/models/apiResponse'

/**
 * Service para personajes. Usa $api (creado en plugins/fetch.client.ts)
 * que ya inyecta el JWT y apunta a backendUrl.
 *
 * Nota: el backend usa rutas con barra final (POST /characters/),
 * mientras que GET /characters/:id no la usa. Mantenemos la convención
 * del backend para no toparnos con redirects que rompan POST.
 */
export class CharacterService {
  private get api() {
    return useNuxtApp().$api as typeof $fetch
  }

  async list(): Promise<Character[]> {
    const res = await this.api<ApiResponse<Character[]>>('/characters/', { method: 'GET' })
    return res.data ?? []
  }

  /** GET por defecto trae el catálogo expandido (clase, raza, hechizos). */
  async get(id: string, expand = true): Promise<CharacterExpanded> {
    const res = await this.api<ApiResponse<CharacterExpanded>>(`/characters/${id}`, {
      method: 'GET',
      query: { expand: expand ? 'true' : 'false' },
    })
    if (!res.data) throw new Error(res.message || 'No se pudo obtener el personaje')
    return res.data
  }

  async create(payload: Partial<Character>): Promise<Character> {
    const res = await this.api<ApiResponse<Character>>('/characters/', {
      method: 'POST',
      body: payload,
    })
    if (!res.data) throw new Error(res.message || 'No se pudo crear el personaje')
    return res.data
  }

  async update(id: string, payload: Partial<Character>): Promise<Character> {
    const res = await this.api<ApiResponse<Character>>(`/characters/${id}`, {
      method: 'PUT',
      body: payload,
    })
    if (!res.data) throw new Error(res.message || 'No se pudo actualizar el personaje')
    return res.data
  }

  /** Patch rápido para cambios de combate (HP, conditions, slots). */
  async patchCombat(id: string, fields: {
    hp?: number
    temp_hp?: number
    conditions?: string[]
    slots_used?: number[]
  }): Promise<void> {
    await this.api(`/characters/${id}/combat`, { method: 'PATCH', body: fields })
  }

  async remove(id: string): Promise<void> {
    await this.api(`/characters/${id}`, { method: 'DELETE' })
  }

  /** Asocia el PJ a una campaña (o lo desasocia si pasas vacío). */
  async assignCampaign(id: string, campaignId: string): Promise<void> {
    await this.api(`/characters/${id}/campaign`, {
      method: 'PUT',
      body: { campaign_id: campaignId },
    })
  }

  /** Devuelve el JSON portable del PJ (incluye nombres para reimport). */
  async exportPortable(id: string): Promise<unknown> {
    const res = await this.api<ApiResponse<unknown>>(`/characters/${id}/export`, { method: 'GET' })
    if (!res.data) throw new Error(res.message || 'No se pudo exportar')
    return res.data
  }

  /** Crea un PJ nuevo a partir de un JSON portable previamente exportado. */
  async importPortable(payload: unknown): Promise<Character> {
    const res = await this.api<ApiResponse<Character>>('/characters/import', {
      method: 'POST',
      body: payload as object,
    })
    if (!res.data) throw new Error(res.message || 'No se pudo importar el personaje')
    return res.data
  }
}
