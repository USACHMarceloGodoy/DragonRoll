<template>
  <div class="form">
    <div class="row split">
      <div class="field flex-2">
        <label>Nombre *</label>
        <InputText v-model="local.name" required />
      </div>
      <div class="field">
        <label>Clase padre *</label>
        <Select
          v-model="local.parent_class"
          :options="classes || []"
          option-label="name"
          option-value="_id"
          placeholder="Elegí clase"
          required
        />
      </div>
    </div>

    <div class="field">
      <label>Descripción corta</label>
      <Textarea v-model="local.short_desc" rows="2" />
    </div>
    <div class="field">
      <label>Sabor narrativo</label>
      <Textarea v-model="local.flavor" rows="3" />
    </div>

    <fieldset>
      <legend>Niveles que aporta features</legend>
      <p class="muted">Sólo los niveles donde la subclase desbloquea algo. Típico: 3, 7, 15, 18.</p>
      <div v-for="(step, i) in local.levels" :key="i" class="lvl-row">
        <div class="field">
          <label>Nv</label>
          <InputNumber v-model="step.level" :min="1" :max="20" class="w-cell" />
        </div>
        <div class="field flex-2">
          <label>Features</label>
          <MultiSelect
            v-model="step.feature_ids"
            :options="subFeatures"
            option-label="name"
            option-value="_id"
            placeholder="Elegí features"
            display="chip"
            filter
          />
        </div>
        <div class="field">
          <label>Notas</label>
          <InputText v-model="step.notes" />
        </div>
        <Button icon="pi pi-times" text rounded severity="secondary" @click="removeStep(i)" />
      </div>
      <Button label="Agregar nivel" icon="pi pi-plus" severity="secondary" size="small" @click="addStep" />
    </fieldset>

    <fieldset>
      <legend>Conjuros otorgados automáticamente</legend>
      <p class="muted">Lista de pares (nivel, conjuro). Estos conjuros se agregan al lanzar al nivel correspondiente y no cuentan contra el límite de preparados.</p>
      <div v-for="(g, i) in local.granted_spells" :key="i" class="row inline">
        <InputNumber v-model="g.level" :min="1" :max="20" class="w-cell" />
        <Select v-model="g.spell_id" :options="spells || []" option-label="name" option-value="_id" placeholder="Elegí conjuro" filter class="flex-1" />
        <Button icon="pi pi-times" text rounded severity="secondary" @click="removeGranted(i)" />
      </div>
      <Button label="Agregar conjuro otorgado" icon="pi pi-plus" severity="secondary" size="small" @click="addGranted" />
    </fieldset>
  </div>
</template>

<script setup lang="ts">
import type { Class, Feature, GrantedSpell, LevelStep, Spell, Subclass } from '~/models/Catalog'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Select from 'primevue/select'
import MultiSelect from 'primevue/multiselect'
import Textarea from 'primevue/textarea'
import Button from 'primevue/button'

const props = defineProps<{
  modelValue: Partial<Subclass>
  classes?: Class[]
  spells?: Spell[]
  features?: Feature[]
}>()
const emit = defineEmits<{ 'update:modelValue': [Partial<Subclass>] }>()

const subFeatures = computed(() => (props.features || []).filter(f => f.origin === 'class' || f.origin === 'subclass'))

const local = reactive<Partial<Subclass>>({
  name: '', short_desc: '', flavor: '',
  parent_class: '',
  levels: [],
  granted_spells: [],
  ...props.modelValue,
})

function newStep(): LevelStep {
  return {
    level: 3,
    proficiency_bonus: 0,    // se ignora en subclase, vive en la clase padre
    feature_ids: [],
    spell_slots: new Array(9).fill(0),
    cantrips_known: 0,
    spells_known: 0,
    notes: '',
  }
}

function addStep() { local.levels = [...(local.levels || []), newStep()] }
function removeStep(i: number) { local.levels = (local.levels || []).filter((_, j) => j !== i) }

function addGranted() {
  const g: GrantedSpell = { level: 3, spell_id: '' }
  local.granted_spells = [...(local.granted_spells || []), g]
}
function removeGranted(i: number) {
  local.granted_spells = (local.granted_spells || []).filter((_, j) => j !== i)
}

watch(local, (v) => emit('update:modelValue', { ...v }), { deep: true })
</script>

<style scoped>
.form { display: flex; flex-direction: column; gap: 14px; }
.row.split { display: grid; grid-template-columns: repeat(auto-fit, minmax(160px, 1fr)); gap: 12px; }
.row.inline { display: flex; align-items: end; gap: 10px; margin-bottom: 8px; }
.field { display: flex; flex-direction: column; gap: 4px; }
.field.flex-1 { flex: 1; }
.field.flex-2 { grid-column: span 2; }
.field label { font-size: 12px; font-weight: 600; color: #4b5563; }
fieldset { border: 1px solid var(--p-content-border-color, #e5e7eb); border-radius: 8px; padding: 10px 14px; }
fieldset legend { padding: 0 6px; font-size: 11px; font-weight: 600; color: #6b7280; }
.muted { color: #6b7280; font-size: 11px; margin: 0 0 8px; }
.lvl-row { display: grid; grid-template-columns: 70px 1fr 200px auto; gap: 10px; align-items: end; margin-bottom: 8px; }
.w-cell { width: 70px; }
.flex-1 { flex: 1; }
</style>
