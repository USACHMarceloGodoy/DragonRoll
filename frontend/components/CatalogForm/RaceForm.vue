<template>
  <div class="form">
    <div class="row split">
      <div class="field flex-2">
        <label>Nombre *</label>
        <InputText v-model="local.name" required />
      </div>
      <div class="field">
        <label>Tamaño</label>
        <Select v-model="local.size" :options="SIZES" />
      </div>
      <div class="field">
        <label>Velocidad (ft)</label>
        <InputNumber v-model="local.speed" :min="0" />
      </div>
    </div>

    <div class="field">
      <label>Descripción corta</label>
      <Textarea v-model="local.short_desc" rows="2" />
    </div>

    <fieldset>
      <legend>Bonificadores de atributo</legend>
      <div v-for="(b, i) in local.ability_bonuses || []" :key="i" class="row inline">
        <Select v-model="b.stat" :options="STATS" />
        <InputNumber v-model="b.bonus" :min="-2" :max="3" :show-buttons="true" />
        <Button icon="pi pi-times" severity="secondary" text rounded @click="removeBonus(i)" />
      </div>
      <Button label="Agregar bono" icon="pi pi-plus" severity="secondary" size="small" @click="addBonus" />
    </fieldset>

    <div class="field">
      <label>Idiomas</label>
      <Chips v-model="local.languages" placeholder="Enter para agregar" />
    </div>

    <fieldset>
      <legend>Rasgos raciales (Features con origin=race)</legend>
      <p class="muted">Si tu lista está vacía, primero crea Features con origen "race".</p>
      <MultiSelect
        v-model="local.trait_ids"
        :options="raceFeatures || []"
        option-label="name"
        option-value="_id"
        placeholder="Elegí rasgos"
        display="chip"
        filter
      />
    </fieldset>
  </div>
</template>

<script setup lang="ts">
import type { Feature, Race } from '~/models/Catalog'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Select from 'primevue/select'
import Textarea from 'primevue/textarea'
import MultiSelect from 'primevue/multiselect'
import Button from 'primevue/button'
import Chips from 'primevue/chips'

const props = defineProps<{
  modelValue: Partial<Race>
  /** Features de origen=race disponibles para enlazar como traits. */
  features?: Feature[]
}>()
const emit = defineEmits<{ 'update:modelValue': [Partial<Race>] }>()

const SIZES = ['Pequeño', 'Mediano', 'Grande']
const STATS = ['FUE', 'DES', 'CON', 'INT', 'SAB', 'CAR']

const raceFeatures = computed(() => (props.features || []).filter(f => f.origin === 'race'))

const local = reactive<Partial<Race>>({
  name: '', short_desc: '', size: 'Mediano', speed: 30,
  ability_bonuses: [], languages: [], trait_ids: [],
  ...props.modelValue,
})

function addBonus() {
  local.ability_bonuses = [...(local.ability_bonuses || []), { stat: 'FUE', bonus: 1 }]
}
function removeBonus(i: number) {
  local.ability_bonuses = (local.ability_bonuses || []).filter((_, j) => j !== i)
}

watch(local, (v) => emit('update:modelValue', { ...v }), { deep: true })
</script>

<style scoped>
.form { display: flex; flex-direction: column; gap: 12px; }
.row.split { display: grid; grid-template-columns: repeat(auto-fit, minmax(160px, 1fr)); gap: 12px; }
.row.inline { display: flex; gap: 8px; align-items: center; margin-bottom: 6px; }
.field { display: flex; flex-direction: column; gap: 4px; }
.field.flex-2 { grid-column: span 2; }
.field label { font-size: 12px; font-weight: 600; color: #4b5563; }
fieldset { border: 1px solid var(--p-content-border-color, #e5e7eb); border-radius: 8px; padding: 10px 14px; }
fieldset legend { padding: 0 6px; font-size: 11px; font-weight: 600; color: #6b7280; }
.muted { color: #6b7280; font-size: 11px; margin: 0 0 8px; }
</style>
