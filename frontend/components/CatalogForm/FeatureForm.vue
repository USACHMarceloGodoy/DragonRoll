<template>
  <div class="form">
    <div class="row split">
      <div class="field flex-2">
        <label>Nombre *</label>
        <InputText v-model="local.name" required />
      </div>
      <div class="field">
        <label>Origen</label>
        <Select v-model="local.origin" :options="ORIGINS" />
      </div>
      <div class="field">
        <label>Costo de acción</label>
        <Select v-model="local.cost" :options="COSTS" />
      </div>
    </div>

    <div class="field">
      <label>Descripción *</label>
      <Textarea v-model="local.description" rows="4" required />
    </div>

    <div class="field">
      <label>Pista de tirada</label>
      <InputText v-model="local.roll_hint" placeholder="2d8 radiante / cura X / TS DES" />
    </div>

    <fieldset>
      <legend>Recurso de uso (opcional)</legend>
      <p class="muted">Sólo para features con usos limitados. Déjalo vacío para rasgos pasivos o ilimitados.</p>
      <div v-if="local.uses" class="row split">
        <div class="field">
          <label>Usos máximos</label>
          <InputNumber v-model="local.uses.max" :min="0" />
        </div>
        <div class="field">
          <label>Fórmula (descriptiva)</label>
          <InputText v-model="local.uses.formula" placeholder="5 × nivel" />
        </div>
        <div class="field">
          <label>Recarga</label>
          <Select v-model="local.uses.recharge" :options="RECHARGES" option-label="label" option-value="value" />
        </div>
        <div class="field action">
          <Button label="Quitar usos" severity="secondary" size="small" @click="local.uses = undefined" />
        </div>
      </div>
      <Button v-else label="Añadir recurso de uso" severity="secondary" size="small" @click="addUses" />
    </fieldset>
  </div>
</template>

<script setup lang="ts">
import type { Feature } from '~/models/Catalog'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Select from 'primevue/select'
import Textarea from 'primevue/textarea'
import Button from 'primevue/button'

const props = defineProps<{ modelValue: Partial<Feature> }>()
const emit = defineEmits<{ 'update:modelValue': [Partial<Feature>] }>()

const ORIGINS = ['class', 'subclass', 'race', 'background', 'feat', 'item']
const COSTS = ['action', 'bonus', 'reaction', 'passive', 'spell-slot', 'free']
const RECHARGES = [
  { value: 'short_rest', label: 'Descanso corto' },
  { value: 'long_rest', label: 'Descanso largo' },
  { value: 'per_day', label: 'Por día' },
  { value: 'none', label: 'Sin recarga' },
]

const local = reactive<Partial<Feature>>({
  name: '', description: '', origin: 'class', cost: 'passive',
  roll_hint: '', uses: undefined,
  ...props.modelValue,
})

function addUses() {
  local.uses = { max: 1, formula: '', recharge: 'long_rest' }
}

watch(local, (v) => emit('update:modelValue', { ...v }), { deep: true })
</script>

<style scoped>
.form { display: flex; flex-direction: column; gap: 12px; }
.row.split { display: grid; grid-template-columns: repeat(auto-fit, minmax(160px, 1fr)); gap: 12px; }
.field { display: flex; flex-direction: column; gap: 4px; }
.field.flex-2 { grid-column: span 2; }
.field.action { justify-content: flex-end; }
.field label { font-size: 12px; font-weight: 600; color: #4b5563; }
fieldset { border: 1px solid var(--p-content-border-color, #e5e7eb); border-radius: 8px; padding: 10px 14px; }
fieldset legend { padding: 0 6px; font-size: 11px; font-weight: 600; color: #6b7280; }
.muted { color: #6b7280; font-size: 11px; margin: 0 0 8px; }
</style>
