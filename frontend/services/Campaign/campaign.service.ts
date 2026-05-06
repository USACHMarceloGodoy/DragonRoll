import type { Campaign } from '~/models/Campaign'
import type { ApiResponse } from '~/models/apiResponse'

export class CampaignService {
  private get api() {
    return useNuxtApp().$api as typeof $fetch
  }

  async list(): Promise<Campaign[]> {
    const res = await this.api<ApiResponse<Campaign[]>>('/campaigns/', { method: 'GET' })
    return res.data ?? []
  }

  async get(id: string): Promise<Campaign> {
    const res = await this.api<ApiResponse<Campaign>>(`/campaigns/${id}`, { method: 'GET' })
    if (!res.data) throw new Error(res.message || 'No se pudo cargar la campaña')
    return res.data
  }

  async create(payload: Pick<Campaign, 'name' | 'description' | 'system'>): Promise<Campaign> {
    const res = await this.api<ApiResponse<Campaign>>('/campaigns/', {
      method: 'POST',
      body: payload,
    })
    if (!res.data) throw new Error(res.message || 'No se pudo crear la campaña')
    return res.data
  }

  /** Reemplazo total de la campaña (NPCs, ubicaciones, sesiones, etc.). */
  async update(id: string, payload: Partial<Campaign>): Promise<Campaign> {
    const res = await this.api<ApiResponse<Campaign>>(`/campaigns/${id}`, {
      method: 'PUT',
      body: payload,
    })
    if (!res.data) throw new Error(res.message || 'No se pudo actualizar la campaña')
    return res.data
  }

  /** El DM agrega un jugador por username o email. */
  async addPlayer(id: string, userIdentifier: string): Promise<string> {
    const res = await this.api<ApiResponse<string>>(`/campaigns/${id}/players`, {
      method: 'POST',
      body: { user: userIdentifier },
    })
    return res.data || ''
  }

  async remove(id: string): Promise<void> {
    await this.api(`/campaigns/${id}`, { method: 'DELETE' })
  }
}

export const useCampaigns = () => {
  const svc = new CampaignService()
  return useAsyncData('campaigns:list', () => svc.list())
}

export const useCampaign = (id: string | (() => string)) => {
  const svc = new CampaignService()
  const idRef = computed(() => (typeof id === 'function' ? id() : id))
  return useAsyncData(
    () => `campaign:${idRef.value}`,
    () => svc.get(idRef.value),
    { watch: [idRef] },
  )
}

export const useCampaignService = () => new CampaignService()
