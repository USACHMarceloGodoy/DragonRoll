<template>
  <div class="campaign-detail">
    <div v-if="pending" class="state">Cargando…</div>
    <div v-else-if="error || !c" class="state error">No se pudo cargar la campaña.</div>

    <template v-else>
      <header class="hero">
        <div class="hero-text">
          <NuxtLink to="/campaigns" class="back">← Mis campañas</NuxtLink>
          <h1>{{ c.name }}</h1>
          <p class="meta">
            <span class="badge" :class="isDM ? 'badge-dm' : 'badge-player'">
              {{ isDM ? 'Eres el DM' : 'Eres jugador' }}
            </span>
            <span class="dim">· {{ c.system }}</span>
            <span class="dim">· {{ c.player_ids?.length || 0 }} jugadores</span>
          </p>
          <p v-if="c.description" class="desc">{{ c.description }}</p>
        </div>
      </header>

      <Tabs v-model:value="activeTab">
        <TabList>
          <Tab value="party">Party</Tab>
          <Tab value="npcs">NPCs</Tab>
          <Tab value="sessions">Bitácora</Tab>
          <Tab v-if="isDM" value="players">Jugadores</Tab>
        </TabList>

        <TabPanels>
          <!-- Party (PJs en la campaña) -->
          <TabPanel value="party">
            <div v-if="charactersInCampaign.length === 0" class="state empty">
              Aún no hay personajes asociados a esta campaña.
              <p class="hint">Tus jugadores pueden asociar sus PJ desde la hoja del personaje.</p>
            </div>
            <div v-else class="grid">
              <article v-for="ch in charactersInCampaign" :key="ch._id" class="char-card" @click="$router.push(`/characters/${ch._id}`)">
                <div class="char-pic" :style="{ background: ch.color || '#6b4a16' }">
                  <span>{{ ch.pic || '⚔' }}</span>
                </div>
                <div class="char-body">
                  <h3>{{ ch.name }}</h3>
                  <p class="char-meta">Nv {{ ch.level }} · HP {{ ch.hp }}/{{ ch.hp_max }}</p>
                </div>
              </article>
            </div>
          </TabPanel>

          <!-- NPCs -->
          <TabPanel value="npcs">
            <div class="tab-header">
              <h2>NPCs de la campaña</h2>
              <Button v-if="isDM" label="Agregar NPC" icon="pi pi-plus" size="small" @click="openNpcDialog()" />
            </div>
            <div v-if="!c.npcs?.length" class="state empty">Sin NPCs todavía.</div>
            <div v-else class="grid">
              <article v-for="npc in c.npcs" :key="npc.id" class="npc-card">
                <div class="npc-pic" :style="{ background: npc.color || '#444' }">
                  <span>{{ npc.pic || '?' }}</span>
                  <span v-if="!npc.alive" class="dead-overlay">†</span>
                </div>
                <div class="npc-body">
                  <h3>{{ npc.name }}</h3>
                  <p class="npc-meta">{{ npc.role }} <span v-if="npc.faction">· {{ npc.faction }}</span></p>
                  <span class="attitude" :class="`att-${npc.attitude}`">{{ npc.attitude }}</span>
                  <p v-if="npc.summary" class="npc-summary">{{ npc.summary }}</p>
                </div>
                <Button v-if="isDM" icon="pi pi-pencil" text rounded severity="secondary" @click="openNpcDialog(npc)" />
              </article>
            </div>
          </TabPanel>

          <!-- Sessions -->
          <TabPanel value="sessions">
            <div class="tab-header">
              <h2>Bitácora de sesiones</h2>
              <Button v-if="isDM" label="Nueva sesión" icon="pi pi-plus" size="small" @click="openSessionDialog()" />
            </div>
            <div v-if="!c.sessions?.length" class="state empty">Sin sesiones registradas.</div>
            <div v-else class="sessions">
              <article v-for="(s, i) in c.sessions.slice().reverse()" :key="i" class="session-card">
                <header class="session-header">
                  <span class="session-num">#{{ s.number }}</span>
                  <h3>{{ s.title }}</h3>
                  <span class="dim">{{ s.date }}</span>
                </header>
                <p>{{ s.summary }}</p>
              </article>
            </div>
          </TabPanel>

          <!-- Players (admin del DM) -->
          <TabPanel v-if="isDM" value="players">
            <div class="tab-header">
              <h2>Jugadores</h2>
              <Button label="Invitar jugador" icon="pi pi-user-plus" size="small" @click="inviteDialog = true" />
            </div>
            <p v-if="!c.player_ids?.length" class="state empty">Sin jugadores aún.</p>
            <ul v-else class="players-list">
              <li v-for="pid in c.player_ids" :key="pid">
                <i class="pi pi-user" />
                <span class="player-id">{{ pid }}</span>
              </li>
            </ul>
            <p class="hint">Los jugadores se identifican por ID. Más adelante esta lista mostrará username + email.</p>
          </TabPanel>
        </TabPanels>
      </Tabs>
    </template>

    <!-- Dialogs -->
    <Dialog v-model:visible="npcDialog" modal :header="editingNpc?.id ? 'Editar NPC' : 'Nuevo NPC'" :style="{ width: '540px' }">
      <div v-if="editingNpc" class="dialog-form">
        <div class="row split">
          <div class="field flex-2">
            <label>Nombre</label>
            <InputText v-model="editingNpc.name" />
          </div>
          <div class="field">
            <label>Rol</label>
            <InputText v-model="editingNpc.role" placeholder="Mentor, Aliada, …" />
          </div>
        </div>
        <div class="row split">
          <div class="field">
            <label>Facción</label>
            <InputText v-model="editingNpc.faction" />
          </div>
          <div class="field">
            <label>Ubicación</label>
            <InputText v-model="editingNpc.location" />
          </div>
          <div class="field">
            <label>Actitud</label>
            <Select v-model="editingNpc.attitude" :options="ATTITUDES" />
          </div>
        </div>
        <div class="row split">
          <div class="field">
            <label>Emoji</label>
            <InputText v-model="editingNpc.pic" placeholder="🧙 🐉 ⚔" />
          </div>
          <div class="field">
            <label>Color</label>
            <InputText v-model="editingNpc.color" placeholder="#6b4a16" />
          </div>
          <div class="field check">
            <label><Checkbox v-model="editingNpc.alive" :binary="true" /> Vivo</label>
          </div>
        </div>
        <div class="field">
          <label>Resumen</label>
          <Textarea v-model="editingNpc.summary" rows="2" />
        </div>
        <div class="field">
          <label>Detalles (solo DM ve)</label>
          <Textarea v-model="editingNpc.details" rows="3" />
        </div>
      </div>
      <template #footer>
        <Button label="Cancelar" text @click="npcDialog = false" />
        <Button label="Eliminar" v-if="editingNpc?.id" severity="danger" text @click="onDeleteNpc" />
        <Button label="Guardar" :loading="saving" @click="onSaveNpc" />
      </template>
    </Dialog>

    <Dialog v-model:visible="sessionDialog" modal header="Sesión" :style="{ width: '540px' }">
      <div v-if="editingSession" class="dialog-form">
        <div class="row split">
          <div class="field">
            <label>Número</label>
            <InputNumber v-model="editingSession.number" :min="1" />
          </div>
          <div class="field flex-2">
            <label>Título</label>
            <InputText v-model="editingSession.title" />
          </div>
          <div class="field">
            <label>Fecha</label>
            <InputText v-model="editingSession.date" placeholder="2026-05-04" />
          </div>
        </div>
        <div class="field">
          <label>Resumen</label>
          <Textarea v-model="editingSession.summary" rows="6" />
        </div>
      </div>
      <template #footer>
        <Button label="Cancelar" text @click="sessionDialog = false" />
        <Button label="Guardar" :loading="saving" @click="onSaveSession" />
      </template>
    </Dialog>

    <Dialog v-model:visible="inviteDialog" modal header="Invitar jugador" :style="{ width: '420px' }">
      <p class="muted">Username o email del jugador. Debe estar registrado.</p>
      <InputText v-model="inviteValue" placeholder="username o email" class="full" />
      <template #footer>
        <Button label="Cancelar" text @click="inviteDialog = false" />
        <Button label="Invitar" :loading="saving" @click="onInvite" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import type { Campaign, NPC, SessionEntry } from '~/models/Campaign'
import type { Character } from '~/models/Character'
import Tabs from 'primevue/tabs'
import TabList from 'primevue/tablist'
import Tab from 'primevue/tab'
import TabPanels from 'primevue/tabpanels'
import TabPanel from 'primevue/tabpanel'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Select from 'primevue/select'
import Textarea from 'primevue/textarea'
import Checkbox from 'primevue/checkbox'

definePageMeta({ title: 'Campaña', auth: true })

const route = useRoute()
const id = computed(() => route.params.id as string)
const { showSuccess, showError } = useAddToast()
const { data: authData } = useAuth()

const { data: c, pending, error, refresh } = await useCampaign(id)
const { data: allChars } = await useCharacters()

const activeTab = ref<string>('party')

const isDM = computed(() => authData.value?._id === c.value?.dm_id)

const charactersInCampaign = computed<Character[]>(() => {
  if (!allChars.value || !c.value?._id) return []
  return (allChars.value as Character[]).filter(ch => ch.campaign_id === c.value!._id)
})

// --- NPC dialog ---
const ATTITUDES = ['aliado', 'neutral', 'sospechoso', 'enemigo']
const npcDialog = ref(false)
const editingNpc = ref<NPC | null>(null)
const saving = ref(false)

function openNpcDialog(npc?: NPC) {
  if (npc) {
    editingNpc.value = JSON.parse(JSON.stringify(npc))
  } else {
    editingNpc.value = {
      id: '', name: '', role: '', faction: '', location: '',
      alive: true, attitude: 'neutral', pic: '🧙', color: '#6b4a16',
      summary: '', details: '', tags: [], comments: [],
    }
  }
  npcDialog.value = true
}

async function onSaveNpc() {
  if (!editingNpc.value || !c.value) return
  saving.value = true
  try {
    const updated: Campaign = JSON.parse(JSON.stringify(c.value))
    if (editingNpc.value.id) {
      // Update existente
      const idx = updated.npcs.findIndex(n => n.id === editingNpc.value!.id)
      if (idx >= 0) updated.npcs[idx] = editingNpc.value
    } else {
      // Nuevo: id basado en timestamp
      editingNpc.value.id = `npc-${Date.now()}`
      updated.npcs.push(editingNpc.value)
    }
    const svc = useCampaignService()
    await svc.update(c.value._id!, updated)
    showSuccess('NPC guardado.')
    npcDialog.value = false
    await refresh()
  } catch (e: any) {
    showError(e?.data?.message || e?.message || 'Error al guardar')
  } finally {
    saving.value = false
  }
}

async function onDeleteNpc() {
  if (!editingNpc.value?.id || !c.value) return
  if (!confirm(`¿Eliminar a ${editingNpc.value.name}?`)) return
  saving.value = true
  try {
    const updated: Campaign = JSON.parse(JSON.stringify(c.value))
    updated.npcs = updated.npcs.filter(n => n.id !== editingNpc.value!.id)
    const svc = useCampaignService()
    await svc.update(c.value._id!, updated)
    showSuccess('NPC eliminado.')
    npcDialog.value = false
    await refresh()
  } catch (e: any) {
    showError(e?.message || 'Error al eliminar')
  } finally {
    saving.value = false
  }
}

// --- Session dialog ---
const sessionDialog = ref(false)
const editingSession = ref<SessionEntry | null>(null)

function openSessionDialog() {
  const nextNum = (c.value?.sessions?.length || 0) + 1
  editingSession.value = {
    number: nextNum,
    date: new Date().toISOString().split('T')[0],
    title: '',
    summary: '',
  }
  sessionDialog.value = true
}

async function onSaveSession() {
  if (!editingSession.value || !c.value) return
  saving.value = true
  try {
    const updated: Campaign = JSON.parse(JSON.stringify(c.value))
    updated.sessions = [...(updated.sessions || []), editingSession.value]
    const svc = useCampaignService()
    await svc.update(c.value._id!, updated)
    showSuccess('Sesión guardada.')
    sessionDialog.value = false
    await refresh()
  } catch (e: any) {
    showError(e?.message || 'Error al guardar')
  } finally {
    saving.value = false
  }
}

// --- Invite dialog ---
const inviteDialog = ref(false)
const inviteValue = ref('')
async function onInvite() {
  if (!inviteValue.value.trim() || !c.value) return
  saving.value = true
  try {
    const svc = useCampaignService()
    const username = await svc.addPlayer(c.value._id!, inviteValue.value.trim())
    showSuccess(`Jugador "${username}" invitado.`)
    inviteDialog.value = false
    inviteValue.value = ''
    await refresh()
  } catch (e: any) {
    showError(e?.data?.message || e?.message || 'No se pudo invitar')
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.campaign-detail { max-width: 1100px; margin: 0 auto; padding: 24px; }
.state { padding: 48px; text-align: center; color: #6b7280; }
.state.error { color: #b91c1c; }
.state.empty { font-style: italic; }
.hint { font-size: 12px; color: #9ca3af; font-style: italic; margin-top: 8px; }

.hero { margin-bottom: 18px; }
.back { color: #6b7280; font-size: 12px; text-decoration: none; }
.back:hover { color: #b48a3c; }
.hero h1 { margin: 4px 0 8px; font-size: 28px; }
.meta { display: flex; gap: 8px; align-items: center; margin: 0 0 8px; flex-wrap: wrap; font-size: 13px; }
.badge { padding: 2px 8px; border-radius: 999px; font-size: 10px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.06em; }
.badge-dm { background: #fef3c7; color: #92400e; }
.badge-player { background: #dbeafe; color: #1e40af; }
.dim { color: #6b7280; }
.desc { font-size: 14px; color: #374151; line-height: 1.5; margin: 0; }

.tab-header { display: flex; justify-content: space-between; align-items: center; margin: 12px 0 14px; }
.tab-header h2 { margin: 0; font-size: 14px; letter-spacing: 0.05em; text-transform: uppercase; color: #6b7280; }

.grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 12px; }

/* Char cards */
.char-card {
  display: flex; gap: 12px; padding: 12px;
  background: var(--p-content-background, #fff);
  border: 1px solid var(--p-content-border-color, #e5e7eb);
  border-radius: 10px; cursor: pointer;
}
.char-card:hover { border-color: #b48a3c; }
.char-pic {
  flex: 0 0 48px; width: 48px; height: 48px; border-radius: 10px;
  display: grid; place-items: center; font-size: 22px; color: #fff;
  text-shadow: 0 1px 2px rgba(0,0,0,.4);
}
.char-card h3 { margin: 0 0 2px; font-size: 14px; }
.char-meta { font-size: 12px; color: #6b7280; margin: 0; }

/* NPC cards */
.npc-card {
  display: flex; gap: 12px; padding: 12px;
  background: var(--p-content-background, #fff);
  border: 1px solid var(--p-content-border-color, #e5e7eb);
  border-radius: 10px;
}
.npc-pic {
  position: relative;
  flex: 0 0 56px; width: 56px; height: 56px; border-radius: 10px;
  display: grid; place-items: center; font-size: 24px; color: #fff;
}
.dead-overlay {
  position: absolute; inset: 0; background: rgba(0,0,0,.5);
  border-radius: 10px; color: #ef4444; font-size: 32px;
  display: grid; place-items: center;
}
.npc-body { flex: 1; min-width: 0; }
.npc-card h3 { margin: 0 0 2px; font-size: 14px; }
.npc-meta { font-size: 12px; color: #6b7280; margin: 0; }
.npc-summary { font-size: 12px; color: #4b5563; margin: 6px 0 0; line-height: 1.4; }
.attitude {
  display: inline-block; margin-top: 4px;
  font-size: 10px; padding: 2px 6px; border-radius: 999px;
  text-transform: uppercase; letter-spacing: 0.06em; font-weight: 600;
}
.att-aliado { background: #dcfce7; color: #14532d; }
.att-neutral { background: #f3f4f6; color: #4b5563; }
.att-sospechoso { background: #fef3c7; color: #92400e; }
.att-enemigo { background: #fee2e2; color: #991b1b; }

/* Sessions */
.sessions { display: flex; flex-direction: column; gap: 12px; }
.session-card {
  padding: 14px 16px;
  background: var(--p-content-background, #fff);
  border: 1px solid var(--p-content-border-color, #e5e7eb);
  border-radius: 10px;
}
.session-header { display: flex; align-items: center; gap: 10px; margin-bottom: 6px; }
.session-num {
  background: #6b4a16; color: #fff;
  padding: 2px 8px; border-radius: 999px; font-size: 11px; font-weight: 700;
}
.session-card h3 { margin: 0; font-size: 14px; flex: 1; }
.session-card p { margin: 0; font-size: 13px; color: #374151; line-height: 1.5; }

.players-list { list-style: none; padding: 0; margin: 0; }
.players-list li { padding: 8px 12px; border-bottom: 1px solid var(--p-content-border-color, #e5e7eb); display: flex; gap: 10px; align-items: center; }
.player-id { font-family: monospace; font-size: 12px; color: #6b7280; }

.dialog-form { display: flex; flex-direction: column; gap: 12px; }
.row.split { display: grid; grid-template-columns: repeat(auto-fit, minmax(140px, 1fr)); gap: 10px; }
.field { display: flex; flex-direction: column; gap: 4px; }
.field.flex-2 { grid-column: span 2; }
.field label { font-size: 12px; font-weight: 600; color: #4b5563; }
.field.check { justify-content: center; }
.field.check label { display: flex; align-items: center; gap: 8px; font-weight: 500; }
.muted { color: #6b7280; font-size: 13px; margin: 0 0 10px; }
.full { width: 100%; }
</style>
