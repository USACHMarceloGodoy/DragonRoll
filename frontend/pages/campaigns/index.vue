<template>
  <div class="campaigns-page">
    <header class="page-header">
      <div>
        <h1 class="title">Mis Campañas</h1>
        <p class="subtitle">Campañas que dirijo o en las que juego.</p>
      </div>
      <Button label="Nueva campaña" icon="pi pi-plus" @click="openNewDialog" />
    </header>

    <div v-if="pending" class="state">Cargando…</div>
    <div v-else-if="error" class="state error">No se pudieron cargar las campañas.</div>
    <div v-else-if="!data?.length" class="state empty">
      Aún no formas parte de ninguna campaña. Crea una para empezar a jugar.
    </div>

    <div v-else class="grid">
      <article
        v-for="c in data"
        :key="c._id"
        class="card"
        @click="goTo(c._id!)"
      >
        <div class="card-header">
          <h2>{{ c.name }}</h2>
          <span class="badge" :class="isDM(c) ? 'badge-dm' : 'badge-player'">
            {{ isDM(c) ? 'DM' : 'Jugador' }}
          </span>
        </div>
        <p class="meta">{{ c.system }}</p>
        <p v-if="c.description" class="desc">{{ c.description }}</p>
        <p class="counts">
          <i class="pi pi-users" /> {{ (c.player_ids?.length || 0) + 1 }} miembros
          <span class="dot">·</span>
          <i class="pi pi-calendar" /> {{ c.sessions?.length || 0 }} sesiones
        </p>
      </article>
    </div>

    <Dialog v-model:visible="newDialog" modal header="Nueva campaña" :style="{ width: '480px' }">
      <div class="dialog-form">
        <div class="field">
          <label>Nombre *</label>
          <InputText v-model="newCampaign.name" />
        </div>
        <div class="field">
          <label>Sistema</label>
          <InputText v-model="newCampaign.system" placeholder="D&D 5e" />
        </div>
        <div class="field">
          <label>Descripción</label>
          <Textarea v-model="newCampaign.description" rows="3" placeholder="Sinopsis breve, ambientación…" />
        </div>
      </div>
      <template #footer>
        <Button label="Cancelar" text @click="newDialog = false" />
        <Button label="Crear" :loading="creating" @click="onCreate" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import type { Campaign } from '~/models/Campaign'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'

definePageMeta({ title: 'Campañas', auth: true })

const router = useRouter()
const { showSuccess, showError } = useAddToast()
const { data: authData } = useAuth()

const { data, pending, error, refresh } = await useCampaigns()

const newDialog = ref(false)
const creating = ref(false)
const newCampaign = ref<{ name: string; system: string; description: string }>({
  name: '', system: 'D&D 5e', description: '',
})

function openNewDialog() {
  newCampaign.value = { name: '', system: 'D&D 5e', description: '' }
  newDialog.value = true
}

function goTo(id: string) { router.push(`/campaigns/${id}`) }

function isDM(c: Campaign): boolean {
  return authData.value?._id === c.dm_id
}

async function onCreate() {
  if (!newCampaign.value.name.trim()) {
    showError('El nombre es obligatorio')
    return
  }
  creating.value = true
  try {
    const svc = useCampaignService()
    const created = await svc.create(newCampaign.value)
    showSuccess(`Campaña "${created.name}" creada.`)
    newDialog.value = false
    await refresh()
    router.push(`/campaigns/${created._id}`)
  } catch (e: any) {
    showError(e?.data?.message || e?.message || 'No se pudo crear la campaña')
  } finally {
    creating.value = false
  }
}
</script>

<style scoped>
.campaigns-page { max-width: 1100px; margin: 0 auto; padding: 24px; }
.page-header { display: flex; justify-content: space-between; align-items: flex-end; gap: 16px; margin-bottom: 24px; flex-wrap: wrap; }
.title { font-size: 28px; font-weight: 700; margin: 0; }
.subtitle { color: #6b7280; margin: 4px 0 0; }

.state { padding: 48px; text-align: center; color: #6b7280; }
.state.error { color: #b91c1c; }
.state.empty { font-style: italic; }

.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
}
.card {
  padding: 16px;
  background: var(--p-content-background, #fff);
  border: 1px solid var(--p-content-border-color, #e5e7eb);
  border-radius: 12px;
  cursor: pointer;
  transition: transform .15s, box-shadow .15s, border-color .15s;
}
.card:hover {
  transform: translateY(-1px);
  box-shadow: 0 8px 18px rgba(0,0,0,.08);
  border-color: var(--p-primary-color, #b48a3c);
}
.card-header { display: flex; align-items: center; gap: 8px; margin-bottom: 6px; }
.card h2 { margin: 0; font-size: 16px; font-weight: 700; flex: 1; }
.badge {
  font-size: 10px; padding: 2px 8px; border-radius: 999px;
  font-weight: 600; text-transform: uppercase; letter-spacing: 0.06em;
}
.badge-dm { background: #fef3c7; color: #92400e; }
.badge-player { background: #dbeafe; color: #1e40af; }

.meta { font-size: 12px; color: #6b7280; margin: 0 0 8px; }
.desc { font-size: 13px; color: #374151; margin: 0 0 8px; line-height: 1.4; }
.counts { font-size: 12px; color: #6b7280; margin: 0; display: flex; align-items: center; gap: 6px; }
.dot { opacity: .5; }

.dialog-form { display: flex; flex-direction: column; gap: 12px; }
.field { display: flex; flex-direction: column; gap: 4px; }
.field label { font-size: 12px; font-weight: 600; color: #4b5563; }
</style>
