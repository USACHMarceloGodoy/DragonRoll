<template>
  <div class="sheet-page">
    <div v-if="pending" class="state">Cargando…</div>
    <div v-else-if="error || !ch" class="state error">No se pudo cargar el personaje.</div>

    <template v-else>
      <div class="action-bar">
        <NuxtLink to="/characters" class="back">← Mis personajes</NuxtLink>
        <div class="actions">
          <Button
            :label="campaignLabel"
            icon="pi pi-flag"
            severity="secondary"
            outlined
            size="small"
            @click="campaignDialog = true"
          />
          <Button
            label="Exportar JSON"
            icon="pi pi-download"
            severity="secondary"
            outlined
            size="small"
            @click="onExport"
          />
          <Button
            icon="pi pi-trash"
            severity="danger"
            outlined
            size="small"
            @click="onDelete"
          />
        </div>
      </div>

      <header class="hero">
        <div class="hero-pic" :style="{ background: ch.color || '#6b4a16' }">
          <span>{{ ch.pic || '⚔' }}</span>
        </div>
        <div class="hero-info">
          <h1>{{ ch.name }}</h1>
          <p class="line">
            <strong>{{ ch.race?.name || '—' }}</strong>
            · {{ ch.class?.name || '—' }}<span v-if="ch.subclass">: {{ ch.subclass.name }}</span>
            · Nivel {{ ch.level }}
          </p>
          <p class="line muted">{{ ch.background || 'Aventurero' }} · {{ ch.alignment || '—' }}</p>
        </div>
        <div class="hero-vitals">
          <div class="vital">
            <label>HP</label>
            <strong>{{ ch.hp }}<span class="of">/ {{ ch.hp_max }}</span></strong>
          </div>
          <div class="vital">
            <label>CA</label>
            <strong>{{ ch.ac }}</strong>
          </div>
          <div class="vital">
            <label>Velocidad</label>
            <strong>{{ ch.speed }}<span class="of">ft</span></strong>
          </div>
          <div class="vital">
            <label>Iniciativa</label>
            <strong>{{ formatMod(ch.initiative) }}</strong>
          </div>
        </div>
      </header>

      <section class="grid">
        <!-- Atributos -->
        <article class="block stats">
          <h2>Atributos</h2>
          <div class="stat-grid">
            <div v-for="s in STATS" :key="s" class="stat">
              <label>{{ s }}</label>
              <strong>{{ ch.stats[s].score }}</strong>
              <span class="mod">{{ formatMod(abilityMod(ch.stats[s].score)) }}</span>
            </div>
          </div>
        </article>

        <!-- Habilidades -->
        <article class="block skills">
          <h2>Habilidades</h2>
          <ul>
            <li v-for="sk in resolvedSkills" :key="sk._id" :class="{ prof: sk.proficient, exp: sk.expert }">
              <span class="prof-dot" />
              <span class="sk-name">{{ sk.name }}</span>
              <span class="sk-stat">{{ sk.stat }}</span>
              <span class="sk-mod">{{ formatMod(sk.bonus) }}</span>
            </li>
          </ul>
        </article>

        <!-- Ataques -->
        <article class="block attacks">
          <h2>Ataques</h2>
          <table v-if="ch.attacks?.length">
            <thead><tr><th>Nombre</th><th>Bono</th><th>Daño</th><th>Costo</th></tr></thead>
            <tbody>
              <tr v-for="(a, i) in ch.attacks" :key="i">
                <td>{{ a.name }}</td>
                <td class="num">{{ formatMod(a.bonus) }}</td>
                <td>{{ a.damage }} <span class="dim">{{ a.damage_type }}</span></td>
                <td><span class="badge">{{ a.cost }}</span></td>
              </tr>
            </tbody>
          </table>
          <p v-else class="empty">Sin ataques registrados.</p>
        </article>

        <!-- Rasgos / Features -->
        <article class="block features">
          <h2>Rasgos</h2>
          <div v-for="f in ch.features_resolved || []" :key="f._id" class="feature">
            <header>
              <strong>{{ f.name }}</strong>
              <span class="badge">{{ f.cost }}</span>
            </header>
            <p>{{ f.description }}</p>
          </div>
          <p v-if="!ch.features_resolved?.length" class="empty">Sin rasgos cargados.</p>
        </article>

        <!-- Slots de conjuro (si aplica) -->
        <article v-if="hasSpellcasting" class="block spellcasting">
          <h2>Conjuros</h2>
          <div class="slots">
            <div v-for="(slots, lvl) in slotsByLevel" :key="lvl" class="slot-row">
              <label>Nivel {{ lvl + 1 }}</label>
              <div class="slot-pips">
                <span
                  v-for="i in slots.total"
                  :key="i"
                  class="pip"
                  :class="{ used: i <= slots.used }"
                />
              </div>
              <span class="slot-count">{{ slots.total - slots.used }} / {{ slots.total }}</span>
            </div>
          </div>

          <h3 v-if="ch.cantrips_resolved?.length">Trucos</h3>
          <ul class="spell-list">
            <li v-for="sp in ch.cantrips_resolved || []" :key="sp._id">
              <strong>{{ sp.name }}</strong>
              <span class="dim">· {{ sp.casting_time }} · {{ sp.range }}</span>
            </li>
          </ul>
          <h3 v-if="ch.spells_resolved?.length">Conjuros conocidos</h3>
          <ul class="spell-list">
            <li v-for="sp in ch.spells_resolved || []" :key="sp._id">
              <strong>Nv{{ sp.level }} {{ sp.name }}</strong>
              <span class="dim">· {{ sp.casting_time }} · {{ sp.range }}</span>
            </li>
          </ul>
        </article>

        <!-- Inventario -->
        <article class="block inventory">
          <h2>Inventario</h2>
          <p><strong>Oro:</strong> {{ ch.inventory.gold }}</p>
          <div v-if="ch.inventory.weapons?.length"><strong>Armas:</strong> {{ ch.inventory.weapons.join(', ') }}</div>
          <div v-if="ch.inventory.armor?.length"><strong>Armadura:</strong> {{ ch.inventory.armor.join(', ') }}</div>
          <div v-if="ch.inventory.items?.length"><strong>Otros:</strong> {{ ch.inventory.items.join(', ') }}</div>
        </article>

        <!-- Roleplay -->
        <article class="block roleplay">
          <h2>Roleplay</h2>
          <p v-if="ch.personality"><strong>Personalidad:</strong> {{ ch.personality }}</p>
          <p v-if="ch.ideals"><strong>Ideales:</strong> {{ ch.ideals }}</p>
          <p v-if="ch.bonds"><strong>Vínculos:</strong> {{ ch.bonds }}</p>
          <p v-if="ch.flaws"><strong>Defectos:</strong> {{ ch.flaws }}</p>
        </article>
      </section>
    </template>

    <Dialog v-model:visible="campaignDialog" modal header="Asociar a campaña" :style="{ width: '460px' }">
      <p class="muted-small">
        Elige una campaña para este personaje. Al asociarlo, el DM podrá ver y editar la hoja durante las sesiones.
      </p>
      <Select
        v-model="selectedCampaignId"
        :options="campaignOptions"
        option-label="name"
        option-value="_id"
        placeholder="(Sin campaña)"
        :show-clear="true"
        class="full"
      />
      <template #footer>
        <Button label="Cancelar" text @click="campaignDialog = false" />
        <Button label="Guardar" :loading="assigning" @click="onAssignCampaign" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { abilityMod, profBonus, STATS, type Stat } from '~/models/Character'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import Select from 'primevue/select'

definePageMeta({ title: 'Hoja', auth: true })

const route = useRoute()
const router = useRouter()
const { showSuccess, showError } = useAddToast()
const id = computed(() => route.params.id as string)

const { data: ch, pending, error, refresh } = await useCharacter(id)

// Campañas para el selector. Falla silencioso si el endpoint no responde
// (ej. usuario sin campañas todavía).
const { data: campaigns } = await useCampaigns()

const campaignOptions = computed(() => (campaigns.value || []).map(c => ({ _id: c._id, name: c.name })))

const campaignDialog = ref(false)
const assigning = ref(false)
const selectedCampaignId = ref<string>('')

watch(ch, (v) => {
  selectedCampaignId.value = v?.campaign_id || ''
}, { immediate: true })

const campaignLabel = computed(() => {
  const c = ch.value
  if (!c?.campaign_id) return 'Sin campaña'
  const found = campaignOptions.value.find(opt => opt._id === c.campaign_id)
  return found ? found.name : 'Campaña desconocida'
})

async function onAssignCampaign() {
  if (!ch.value?._id) return
  assigning.value = true
  try {
    const svc = useCharacterService()
    await svc.assignCampaign(ch.value._id, selectedCampaignId.value || '')
    showSuccess(selectedCampaignId.value ? 'Asociado a la campaña.' : 'Desasociado de la campaña.')
    campaignDialog.value = false
    await refresh()
  } catch (e: any) {
    showError(e?.data?.message || e?.message || 'No se pudo asociar')
  } finally {
    assigning.value = false
  }
}

async function onExport() {
  if (!ch.value?._id) return
  try {
    const svc = useCharacterService()
    const portable = await svc.exportPortable(ch.value._id)
    // Disparar descarga en el navegador.
    const blob = new Blob([JSON.stringify(portable, null, 2)], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${ch.value.name.replace(/\s+/g, '_')}.dragonroll.json`
    document.body.appendChild(a)
    a.click()
    a.remove()
    URL.revokeObjectURL(url)
    showSuccess('Personaje exportado.')
  } catch (e: any) {
    showError(e?.message || 'No se pudo exportar')
  }
}

async function onDelete() {
  if (!ch.value?._id) return
  if (!confirm(`¿Eliminar a ${ch.value.name}? Esta acción no se puede deshacer.`)) return
  try {
    const svc = useCharacterService()
    await svc.remove(ch.value._id)
    showSuccess('Personaje eliminado.')
    router.push('/characters')
  } catch (e: any) {
    showError(e?.message || 'No se pudo eliminar')
  }
}

function formatMod(n: number): string {
  return n >= 0 ? `+${n}` : `${n}`
}

/** Habilidades resueltas: cruza el catálogo de skills con las del PJ. */
const resolvedSkills = computed(() => {
  if (!ch.value?.skills_resolved) return []
  const profSet = new Set(ch.value.skill_proficiency_ids || [])
  const expSet = new Set(ch.value.skill_expertise_ids || [])
  const pb = profBonus(ch.value.class, ch.value.level)
  return ch.value.skills_resolved.map(s => {
    const proficient = profSet.has(s._id)
    const expert = expSet.has(s._id)
    const stat = s.stat as Stat
    const mod = abilityMod(ch.value!.stats[stat].score)
    const bonus = mod + (expert ? pb * 2 : proficient ? pb : 0)
    return { ...s, proficient, expert, bonus }
  })
})

const hasSpellcasting = computed(() => !!ch.value?.class?.casting)

/** Slots por nivel a partir de la tabla de la clase + slots usados del PJ. */
const slotsByLevel = computed(() => {
  if (!ch.value?.class) return []
  const step = ch.value.class.levels.find(l => l.level === ch.value!.level)
  const totals = step?.spell_slots ?? new Array(9).fill(0)
  const used = ch.value.slots_used || new Array(9).fill(0)
  return totals
    .map((total, i) => ({ total, used: used[i] || 0 }))
    .filter(s => s.total > 0)
})
</script>

<style scoped>
.sheet-page { max-width: 1100px; margin: 0 auto; padding: 24px; }

.action-bar {
  display: flex; justify-content: space-between; align-items: center;
  margin-bottom: 14px;
  flex-wrap: wrap;
  gap: 10px;
}
.back { color: #6b7280; font-size: 12px; text-decoration: none; }
.back:hover { color: #b48a3c; }
.actions { display: flex; gap: 8px; flex-wrap: wrap; }

.muted-small { font-size: 13px; color: #6b7280; margin: 0 0 12px; }
.full { width: 100%; }

.state { padding: 48px; text-align: center; color: #6b7280; }
.state.error { color: #b91c1c; }

.hero {
  display: flex;
  gap: 16px;
  align-items: center;
  background: var(--p-content-background, #fff);
  border: 1px solid var(--p-content-border-color, #e5e7eb);
  border-radius: 14px;
  padding: 18px;
  margin-bottom: 18px;
}
.hero-pic {
  width: 72px; height: 72px; border-radius: 14px;
  display: grid; place-items: center;
  font-size: 36px; color: #fff;
  text-shadow: 0 1px 2px rgba(0,0,0,.4);
}
.hero-info { flex: 1; min-width: 0; }
.hero-info h1 { margin: 0 0 4px; font-size: 22px; }
.line { margin: 0; }
.line.muted { color: #6b7280; font-size: 13px; }
.hero-vitals {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 10px;
}
.vital {
  text-align: center;
  padding: 8px 14px;
  border: 1px solid var(--p-content-border-color, #e5e7eb);
  border-radius: 10px;
  min-width: 70px;
}
.vital label { display: block; font-size: 10px; color: #6b7280; letter-spacing: .15em; text-transform: uppercase; }
.vital strong { font-size: 22px; font-weight: 700; }
.vital .of { font-size: 12px; color: #6b7280; margin-left: 2px; }

.grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 14px;
}
.block {
  background: var(--p-content-background, #fff);
  border: 1px solid var(--p-content-border-color, #e5e7eb);
  border-radius: 12px;
  padding: 14px 16px;
}
.block h2 {
  font-size: 13px; font-weight: 700; letter-spacing: .14em; text-transform: uppercase;
  color: #6b7280;
  margin: 0 0 10px;
}
.block h3 { font-size: 12px; margin: 12px 0 6px; color: #6b7280; text-transform: uppercase; letter-spacing: .12em; }

.stat-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 8px; }
.stat {
  display: flex; flex-direction: column; align-items: center;
  padding: 10px 8px;
  border: 1px solid var(--p-content-border-color, #e5e7eb);
  border-radius: 10px;
}
.stat label { font-size: 10px; letter-spacing: .14em; color: #6b7280; }
.stat strong { font-size: 22px; font-weight: 700; }
.stat .mod { font-size: 12px; color: #6b7280; font-variant-numeric: tabular-nums; }

.skills ul { list-style: none; padding: 0; margin: 0; }
.skills li {
  display: grid;
  grid-template-columns: 12px 1fr auto auto;
  gap: 8px;
  align-items: center;
  padding: 4px 0;
  font-size: 13px;
}
.skills li .prof-dot {
  width: 8px; height: 8px; border-radius: 50%;
  background: transparent; border: 1px solid #d1d5db;
}
.skills li.prof .prof-dot { background: var(--p-primary-color, #b48a3c); border-color: transparent; }
.skills li.exp .prof-dot { background: linear-gradient(135deg, var(--p-primary-color, #b48a3c) 50%, #6b7280 50%); }
.skills li .sk-stat { font-size: 10px; color: #6b7280; }
.skills li .sk-mod { font-variant-numeric: tabular-nums; font-weight: 600; }

table { width: 100%; border-collapse: collapse; font-size: 13px; }
table th { text-align: left; padding: 6px 4px; color: #6b7280; font-weight: 600; font-size: 11px; text-transform: uppercase; letter-spacing: .1em; }
table td { padding: 6px 4px; border-top: 1px solid var(--p-content-border-color, #e5e7eb); }
.num { font-variant-numeric: tabular-nums; text-align: right; }
.dim { color: #9ca3af; font-size: 12px; }
.badge {
  font-size: 10px; padding: 2px 8px; border-radius: 999px;
  background: var(--p-content-border-color, #e5e7eb); color: #4b5563;
}

.feature { padding: 8px 0; border-top: 1px solid var(--p-content-border-color, #e5e7eb); }
.feature:first-of-type { border-top: 0; }
.feature header { display: flex; align-items: center; gap: 8px; margin-bottom: 4px; }
.feature p { margin: 0; font-size: 13px; color: #374151; }
.empty { color: #9ca3af; font-style: italic; font-size: 13px; }

.slots .slot-row {
  display: grid;
  grid-template-columns: 80px 1fr auto;
  gap: 10px;
  align-items: center;
  padding: 4px 0;
}
.slots label { font-size: 11px; color: #6b7280; }
.slot-pips { display: flex; gap: 4px; }
.pip { width: 14px; height: 14px; border-radius: 50%; border: 1px solid var(--p-primary-color, #b48a3c); }
.pip.used { background: var(--p-primary-color, #b48a3c); }
.slot-count { font-size: 11px; color: #6b7280; font-variant-numeric: tabular-nums; }

.spell-list { list-style: none; padding: 0; margin: 0; }
.spell-list li { padding: 3px 0; font-size: 13px; }
</style>
