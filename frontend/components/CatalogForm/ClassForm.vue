<template>
  <div class="form">
    <!-- Identidad básica -->
    <div class="row split">
      <div class="field flex-2">
        <label>Nombre *</label>
        <InputText v-model="local.name" required />
      </div>
      <div class="field">
        <label>Dado de golpe</label>
        <Select v-model="local.hit_die" :options="HIT_DICE" />
      </div>
    </div>

    <div class="field">
      <label>Descripción corta</label>
      <Textarea v-model="local.short_desc" rows="2" />
    </div>

    <!-- Stats y proficiencias -->
    <fieldset>
      <legend>Atributos y competencias</legend>
      <div class="row split">
        <div class="field">
          <label>Stats primarias</label>
          <MultiSelect v-model="local.primary_stats" :options="STATS" placeholder="2 stats típicas" />
        </div>
        <div class="field">
          <label>Tiradas de salvación competentes</label>
          <MultiSelect v-model="local.save_proficient" :options="STATS" placeholder="Stat × 2" />
        </div>
      </div>
      <div class="row split">
        <div class="field">
          <label>Armaduras</label>
          <Chips v-model="local.armor_profs" placeholder="Enter para añadir" />
        </div>
        <div class="field">
          <label>Armas</label>
          <Chips v-model="local.weapon_profs" placeholder="Enter para añadir" />
        </div>
        <div class="field">
          <label>Herramientas</label>
          <Chips v-model="local.tool_profs" placeholder="Enter para añadir" />
        </div>
      </div>
    </fieldset>

    <!-- Skill picks -->
    <fieldset>
      <legend>Elección de habilidades</legend>
      <div class="row split">
        <div class="field">
          <label>Cantidad a elegir</label>
          <InputNumber v-model="local.skill_choice!.count" :min="0" :max="6" />
        </div>
        <div class="field flex-2">
          <label>Pool de habilidades disponibles</label>
          <MultiSelect
            v-model="local.skill_choice!.from_skill_ids"
            :options="skills || []"
            option-label="name"
            option-value="_id"
            placeholder="Elegí las skills permitidas"
            display="chip"
            filter
          />
        </div>
      </div>
    </fieldset>

    <!-- Subclase -->
    <div class="row split">
      <div class="field">
        <label>Nivel de elección de subclase</label>
        <InputNumber v-model="local.subclass_level" :min="1" :max="20" />
        <small class="muted">A qué nivel se desbloquea la subclase (clásico: 1, 2 o 3).</small>
      </div>
    </div>

    <!-- Casting (opcional) -->
    <fieldset>
      <legend>Lanzamiento de conjuros (opcional)</legend>
      <div v-if="local.casting" class="casting-block">
        <div class="row split">
          <div class="field">
            <label>Atributo</label>
            <Select v-model="local.casting.ability" :options="['INT','SAB','CAR']" />
          </div>
          <div class="field">
            <label>Tipo de preparación</label>
            <Select v-model="local.casting.preparation" :options="['prepared','known','spellbook','pact']" />
          </div>
          <div class="field">
            <label class="check">
              <Checkbox v-model="local.casting.ritual" :binary="true" />
              Permite ritual
            </label>
          </div>
        </div>
        <div class="field">
          <label>Lista de conjuros disponibles</label>
          <MultiSelect
            v-model="local.casting.spell_list_ids"
            :options="spells || []"
            option-label="name"
            option-value="_id"
            placeholder="Conjuros que esta clase puede usar"
            display="chip"
            filter
          />
        </div>
        <Button label="Quitar lanzamiento de conjuros" severity="secondary" size="small" @click="local.casting = undefined" />
      </div>
      <Button v-else label="Agregar lanzamiento de conjuros" icon="pi pi-plus" severity="secondary" size="small" @click="addCasting" />
    </fieldset>

    <!-- Tabla de niveles -->
    <fieldset>
      <legend>Tabla de progresión (niveles 1–20)</legend>
      <p class="muted">
        Cada fila es un nivel. PB = bono de proficiencia. Slots: cuántos slots de cada nivel se tienen <em>al alcanzar este nivel de clase</em> (no acumulativos extra). Conoc = trucos / conjuros conocidos. Click en "Features" para enlazar rasgos a este nivel.
      </p>
      <div class="table-scroll">
        <table class="lvl-table">
          <thead>
            <tr>
              <th>Nv</th>
              <th>PB</th>
              <th colspan="9" class="slots-head">Slots por nivel del conjuro (1-9)</th>
              <th>Trucos</th>
              <th>Conoc</th>
              <th>Notas</th>
              <th>Features</th>
            </tr>
            <tr class="sub-head">
              <th></th><th></th>
              <th v-for="i in 9" :key="i">{{ i }}</th>
              <th></th><th></th><th></th><th></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="step in local.levels" :key="step.level">
              <td class="num">{{ step.level }}</td>
              <td><InputNumber v-model="step.proficiency_bonus" :min="2" :max="6" :show-buttons="false" class="w-cell" /></td>
              <td v-for="i in 9" :key="i">
                <InputNumber v-model="step.spell_slots[i-1]" :min="0" :max="9" :show-buttons="false" class="w-cell" />
              </td>
              <td><InputNumber v-model="step.cantrips_known" :min="0" :max="10" :show-buttons="false" class="w-cell" /></td>
              <td><InputNumber v-model="step.spells_known" :min="-1" :max="50" :show-buttons="false" class="w-cell" /></td>
              <td><InputText v-model="step.notes" placeholder="ej. ASI" class="w-notes" /></td>
              <td>
                <Button
                  :badge="String((step.feature_ids || []).length)"
                  icon="pi pi-list"
                  size="small"
                  text
                  severity="secondary"
                  @click="openFeatures(step)"
                />
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <small class="muted">"Conoc" = -1 significa "todos los preparables" (ej: Clérigo, Druida, Paladín).</small>
    </fieldset>

    <!-- Diálogo de features por nivel -->
    <Dialog v-model:visible="featuresDialog" modal header="Features de este nivel" :style="{ width: '520px' }">
      <p class="muted">Elegí features (rasgos) que se desbloquean al alcanzar este nivel.</p>
      <MultiSelect
        v-if="editingStep"
        v-model="editingStep.feature_ids"
        :options="classFeatures"
        option-label="name"
        option-value="_id"
        placeholder="Elegí features"
        display="chip"
        filter
        class="full"
      />
      <template #footer>
        <Button label="Listo" @click="featuresDialog = false" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import type { Class, Feature, LevelStep, Skill, Spell } from '~/models/Catalog'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Select from 'primevue/select'
import MultiSelect from 'primevue/multiselect'
import Textarea from 'primevue/textarea'
import Checkbox from 'primevue/checkbox'
import Chips from 'primevue/chips'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'

const props = defineProps<{
  modelValue: Partial<Class>
  skills?: Skill[]
  spells?: Spell[]
  features?: Feature[]
}>()
const emit = defineEmits<{ 'update:modelValue': [Partial<Class>] }>()

const STATS = ['FUE', 'DES', 'CON', 'INT', 'SAB', 'CAR']
const HIT_DICE = ['d6', 'd8', 'd10', 'd12']

/**
 * Inicializa una tabla 1..20 si viene vacía. Si el modelo viene con
 * algunos niveles, los respeta y rellena los faltantes con valores por
 * defecto razonables (PB de 5e).
 */
function buildLevels(existing: LevelStep[] = []): LevelStep[] {
  const map = new Map(existing.map(l => [l.level, l]))
  const out: LevelStep[] = []
  for (let lv = 1; lv <= 20; lv++) {
    const existing = map.get(lv)
    if (existing) {
      out.push({
        ...existing,
        spell_slots: existing.spell_slots && existing.spell_slots.length === 9
          ? [...existing.spell_slots]
          : new Array(9).fill(0),
        feature_ids: existing.feature_ids || [],
      })
    } else {
      out.push({
        level: lv,
        proficiency_bonus: defaultProfBonus(lv),
        feature_ids: [],
        spell_slots: new Array(9).fill(0),
        cantrips_known: 0,
        spells_known: 0,
        notes: '',
      })
    }
  }
  return out
}
function defaultProfBonus(lv: number): number {
  if (lv >= 17) return 6
  if (lv >= 13) return 5
  if (lv >= 9) return 4
  if (lv >= 5) return 3
  return 2
}

const local = reactive<Partial<Class>>({
  name: '', short_desc: '', hit_die: 'd8',
  primary_stats: [], save_proficient: [],
  armor_profs: [], weapon_profs: [], tool_profs: [],
  skill_choice: { count: 2, from_skill_ids: [] },
  subclass_level: 3,
  subclasses: [],
  levels: buildLevels(props.modelValue.levels),
  casting: undefined,
  ...props.modelValue,
  // forzar consistencia tras spread:
})

if (!local.skill_choice) local.skill_choice = { count: 2, from_skill_ids: [] }
if (!local.levels || local.levels.length !== 20) local.levels = buildLevels(local.levels)

const classFeatures = computed(() =>
  (props.features || []).filter(f => f.origin === 'class' || f.origin === 'subclass'))

function addCasting() {
  local.casting = { ability: 'CAR', preparation: 'prepared', ritual: false, spell_list_ids: [] }
}

const featuresDialog = ref(false)
const editingStep = ref<LevelStep | null>(null)
function openFeatures(step: LevelStep) {
  editingStep.value = step
  featuresDialog.value = true
}

watch(local, (v) => emit('update:modelValue', { ...v }), { deep: true })
</script>

<style scoped>
.form { display: flex; flex-direction: column; gap: 14px; }
.row.split { display: grid; grid-template-columns: repeat(auto-fit, minmax(160px, 1fr)); gap: 12px; }
.field { display: flex; flex-direction: column; gap: 4px; }
.field.flex-2 { grid-column: span 2; }
.field label { font-size: 12px; font-weight: 600; color: #4b5563; }
.check { display: flex; align-items: center; gap: 8px; font-size: 13px; }
fieldset { border: 1px solid var(--p-content-border-color, #e5e7eb); border-radius: 8px; padding: 10px 14px; }
fieldset legend { padding: 0 6px; font-size: 11px; font-weight: 600; color: #6b7280; }
.muted { color: #6b7280; font-size: 11px; margin: 0 0 8px; }

.table-scroll { overflow-x: auto; max-width: 100%; }
.lvl-table {
  width: 100%; border-collapse: collapse;
  font-size: 12px; font-variant-numeric: tabular-nums;
}
.lvl-table th, .lvl-table td {
  border: 1px solid var(--p-content-border-color, #e5e7eb);
  padding: 2px 4px;
  text-align: center;
}
.lvl-table th {
  background: var(--p-content-hover-background, #f9fafb);
  font-size: 10px; font-weight: 700; color: #4b5563;
  letter-spacing: 0.04em; text-transform: uppercase;
}
.lvl-table td.num { font-weight: 700; background: var(--p-content-hover-background, #f9fafb); }
.lvl-table .slots-head { color: #6b7280; }
.lvl-table .sub-head th { font-size: 10px; padding: 2px; }
.w-cell { width: 56px; }
.w-cell :deep(input) { padding: 2px 4px; text-align: center; }
.w-notes { width: 110px; }
.w-notes :deep(input) { padding: 2px 4px; font-size: 11px; }
.casting-block { display: flex; flex-direction: column; gap: 10px; }
.full { width: 100%; }
</style>
