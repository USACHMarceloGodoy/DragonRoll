<template>
  <div class="new-page">
    <header class="page-header">
      <h1>Nuevo Personaje</h1>
      <p class="muted">Construye tu hoja a partir del catálogo disponible.</p>
    </header>

    <div v-if="loading" class="state">Cargando catálogo…</div>

    <form v-else class="form" @submit.prevent="onSubmit">
      <!-- Paso 1: Identidad -->
      <section class="block">
        <h2>1. Identidad</h2>
        <div class="row">
          <label>Nombre *</label>
          <InputText v-model="form.name" required />
        </div>
        <div class="row split">
          <div>
            <label>Pronombres</label>
            <InputText v-model="form.pronouns" placeholder="él/ella/elle" />
          </div>
          <div>
            <label>Trasfondo</label>
            <InputText v-model="form.background" placeholder="Soldado, Erudito, …" />
          </div>
          <div>
            <label>Alineamiento</label>
            <InputText v-model="form.alignment" placeholder="Legal Bueno" />
          </div>
        </div>
        <div class="row split">
          <div>
            <label>Emoji / Icono</label>
            <InputText v-model="form.pic" placeholder="⚔" />
          </div>
          <div>
            <label>Color (oklch)</label>
            <InputText v-model="form.color" placeholder="oklch(0.78 0.13 78)" />
          </div>
        </div>
      </section>

      <!-- Paso 2: Clase / Raza / Nivel -->
      <section class="block">
        <h2>2. Clase y Raza</h2>
        <div class="row split">
          <div>
            <label>Clase *</label>
            <Select v-model="form.class_id" :options="classOptions" optionLabel="name" optionValue="_id" placeholder="Elegir clase" required />
            <small v-if="selectedClass" class="muted">{{ selectedClass.short_desc }}</small>
          </div>
          <div>
            <label>Raza *</label>
            <Select v-model="form.race_id" :options="races || []" optionLabel="name" optionValue="_id" placeholder="Elegir raza" required />
            <small v-if="selectedRace" class="muted">{{ selectedRace.short_desc }}</small>
          </div>
        </div>
        <div class="row split">
          <div>
            <label>Subclase</label>
            <Select v-model="form.subclass_id" :options="subclassOptions" optionLabel="name" optionValue="_id" :placeholder="canPickSubclass ? 'Elegir subclase' : `Disponible al nivel ${selectedClass?.subclass_level ?? '-'}`" :disabled="!canPickSubclass" />
          </div>
          <div>
            <label>Nivel *</label>
            <InputNumber v-model="form.level" :min="1" :max="20" />
          </div>
        </div>
      </section>

      <!-- Paso 3: Atributos y combate -->
      <section class="block">
        <h2>3. Atributos</h2>
        <p class="muted small">Asigna las puntuaciones (recomendado 8-15 base; los bonificadores raciales se aplican aparte).</p>
        <div class="stat-grid">
          <div v-for="s in STATS" :key="s" class="stat-input">
            <label>{{ s }} <span class="muted">{{ STAT_NAMES[s] }}</span></label>
            <InputNumber v-model="form.stats[s].score" :min="3" :max="22" />
            <small class="mod">mod {{ formatMod(abilityMod(form.stats[s].score)) }}</small>
          </div>
        </div>
        <div class="row split">
          <div>
            <label>HP máx</label>
            <InputNumber v-model="form.hp_max" :min="1" />
          </div>
          <div>
            <label>HP actual</label>
            <InputNumber v-model="form.hp" :min="0" />
          </div>
          <div>
            <label>CA</label>
            <InputNumber v-model="form.ac" :min="1" :max="30" />
          </div>
          <div>
            <label>Velocidad</label>
            <InputNumber v-model="form.speed" :min="0" suffix=" ft" />
          </div>
        </div>
      </section>

      <!-- Paso 4: Habilidades -->
      <section class="block">
        <h2>4. Habilidades</h2>
        <p class="muted small">
          Marca habilidades con competencia.
          <span v-if="selectedClass">Tu clase permite elegir <strong>{{ selectedClass.skill_choice?.count || 0 }}</strong>.</span>
        </p>
        <ul class="skills-pick">
          <li v-for="sk in skills || []" :key="sk._id">
            <label>
              <Checkbox v-model="profSet" :value="sk._id" />
              <span class="sk-name">{{ sk.name }}</span>
              <span class="sk-stat">{{ sk.stat }}</span>
            </label>
          </li>
        </ul>
      </section>

      <div class="actions">
        <NuxtLink to="/characters" class="btn-secondary">Cancelar</NuxtLink>
        <Button type="submit" label="Crear personaje" :loading="saving" severity="primary" />
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { abilityMod, STATS, STAT_NAMES, type Character } from '~/models/Character'
import type { Class, Race, Skill, Subclass } from '~/models/Catalog'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Select from 'primevue/select'
import Checkbox from 'primevue/checkbox'
import Button from 'primevue/button'

definePageMeta({ title: 'Nuevo personaje', auth: true })

const router = useRouter()
const { showSuccess, showError } = useAddToast()

// Cargamos los 4 catálogos en paralelo. Si alguno falla, el wizard
// muestra estado de carga; en producción habría que reintentar/mostrar error.
const [{ data: classes }, { data: subclasses }, { data: races }, { data: skills }] = await Promise.all([
  useCatalogList('classes'),
  useCatalogList('subclasses'),
  useCatalogList('races'),
  useCatalogList('skills'),
])

const loading = computed(() => !classes.value || !races.value || !skills.value)

// PrimeVue Select trabaja mejor con un array tipado.
const classOptions = computed<Class[]>(() => (classes.value as Class[]) || [])
const subclassOptionsAll = computed<Subclass[]>(() => (subclasses.value as Subclass[]) || [])
const racesList = computed<Race[]>(() => (races.value as Race[]) || [])
const skillsList = computed<Skill[]>(() => (skills.value as Skill[]) || [])

// Reactividad sobre la clase elegida.
const form = reactive<Partial<Character> & { stats: Character['stats'] }>({
  name: '',
  pronouns: '',
  class_id: '',
  subclass_id: '',
  race_id: '',
  background: '',
  alignment: '',
  level: 1,
  pic: '⚔',
  color: 'oklch(0.78 0.13 78)',
  hp: 10,
  hp_max: 10,
  ac: 10,
  speed: 30,
  initiative: 0,
  inspiration: false,
  hit_dice: { current: 1, max: 1, type: 'd8' },
  death_saves: { successes: 0, failures: 0 },
  stats: {
    FUE: { score: 10 }, DES: { score: 10 }, CON: { score: 10 },
    INT: { score: 10 }, SAB: { score: 10 }, CAR: { score: 10 },
  },
  skill_proficiency_ids: [],
  skill_expertise_ids: [],
  cantrip_ids: [],
  spells_known_ids: [],
  spells_prepared_ids: [],
  slots_used: new Array(9).fill(0),
  resources: [],
  attacks: [],
  conditions: [],
  inventory: { armor: [], weapons: [], tools: [], languages: [], proficiencies: [], items: [], gold: 0 },
  active: true,
})

const profSet = ref<string[]>([])

const selectedClass = computed<Class | undefined>(() =>
  classOptions.value.find(c => c._id === form.class_id),
)
const selectedRace = computed<Race | undefined>(() =>
  racesList.value.find(r => r._id === form.race_id),
)

// Subclase: filtramos las que apuntan a la clase elegida.
const subclassOptions = computed<Subclass[]>(() => {
  if (!form.class_id) return []
  return subclassOptionsAll.value.filter(s => s.parent_class === form.class_id)
})

const selectedSubclass = computed<Subclass | undefined>(() =>
  subclassOptionsAll.value.find(s => s._id === form.subclass_id),
)

const canPickSubclass = computed(() => {
  return !!selectedClass.value && (form.level || 1) >= (selectedClass.value.subclass_level || 99)
})

// Cuando cambia la clase: ajustar hit_dice.type y limpiar subclase si no aplica.
watch(selectedClass, (k) => {
  if (k && form.hit_dice) {
    form.hit_dice.type = k.hit_die
  }
  if (!canPickSubclass.value) {
    form.subclass_id = ''
  }
})

/**
 * Auto-grant de hechizos de subclase: cuando el jugador elige una subclase
 * (Subclass.granted_spells), los hechizos cuyo level requerido sea <= el
 * nivel del PJ se agregan a `spells_known_ids` automáticamente. Si el
 * jugador cambia de subclase, removemos los anteriores y aplicamos los
 * nuevos. Esto evita que el jugador tenga que ir uno por uno.
 *
 * Trade-off: si el jugador agregó manualmente un conjuro que coincidía
 * con un grant, al cambiar la subclase ese conjuro se va. Razonable para
 * un wizard de creación; revisar si crece el flujo de level-up.
 */
const previousGrants = ref<string[]>([])
watch([selectedSubclass, () => form.level], ([sub, lv]) => {
  // Quitar grants previos de la lista de known.
  if (previousGrants.value.length && form.spells_known_ids) {
    form.spells_known_ids = form.spells_known_ids.filter(id => !previousGrants.value.includes(id))
  }
  if (!sub || !sub.granted_spells) {
    previousGrants.value = []
    return
  }
  const newGrants = sub.granted_spells
    .filter(g => g.level <= (lv || 1))
    .map(g => g.spell_id)
  // Agregar los nuevos sin duplicar.
  const known = new Set(form.spells_known_ids || [])
  newGrants.forEach(id => known.add(id))
  form.spells_known_ids = [...known]
  previousGrants.value = newGrants
})

function formatMod(n: number): string {
  return n >= 0 ? `+${n}` : `${n}`
}

const saving = ref(false)
async function onSubmit() {
  if (!form.name || !form.class_id || !form.race_id) {
    showError('Completa nombre, clase y raza.')
    return
  }
  saving.value = true
  try {
    form.skill_proficiency_ids = [...profSet.value]
    const svc = useCharacterService()
    const created = await svc.create(form as Partial<Character>)
    showSuccess('Personaje creado.')
    router.push(`/characters/${created._id}`)
  } catch (e: any) {
    showError(e?.message || 'No se pudo crear el personaje')
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.new-page { max-width: 900px; margin: 0 auto; padding: 24px; }
.page-header { margin-bottom: 18px; }
.page-header h1 { margin: 0 0 4px; font-size: 24px; }
.muted { color: #6b7280; }
.muted.small { font-size: 12px; }
.state { padding: 48px; text-align: center; color: #6b7280; }

.form { display: flex; flex-direction: column; gap: 16px; }
.block {
  background: var(--p-content-background, #fff);
  border: 1px solid var(--p-content-border-color, #e5e7eb);
  border-radius: 12px;
  padding: 18px;
}
.block h2 { font-size: 14px; margin: 0 0 12px; letter-spacing: .05em; }
.row { margin-bottom: 12px; display: flex; flex-direction: column; gap: 4px; }
.row > label { font-size: 12px; font-weight: 600; color: #4b5563; }
.row.split {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 12px;
}
.row.split > div { display: flex; flex-direction: column; gap: 4px; }
.row.split label { font-size: 12px; font-weight: 600; color: #4b5563; }

.stat-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 10px; margin-bottom: 14px; }
.stat-input { display: flex; flex-direction: column; gap: 4px; padding: 10px; border: 1px solid var(--p-content-border-color, #e5e7eb); border-radius: 8px; }
.stat-input label { font-size: 11px; font-weight: 600; }
.stat-input .mod { font-size: 11px; color: #6b7280; }

.skills-pick {
  list-style: none; padding: 0; margin: 0;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 8px;
}
.skills-pick li label {
  display: flex; align-items: center; gap: 8px;
  padding: 8px 10px; border-radius: 8px;
  border: 1px solid var(--p-content-border-color, #e5e7eb);
  cursor: pointer;
  font-size: 13px;
}
.skills-pick .sk-stat { margin-left: auto; font-size: 10px; color: #6b7280; }

.actions { display: flex; justify-content: flex-end; gap: 10px; }
.btn-secondary {
  padding: 10px 16px; border-radius: 8px;
  border: 1px solid var(--p-content-border-color, #e5e7eb);
  color: #4b5563; text-decoration: none;
}
</style>
