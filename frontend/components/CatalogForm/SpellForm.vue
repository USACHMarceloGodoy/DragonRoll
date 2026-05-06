<template>
  <div class="form">
    <div class="row split">
      <div class="field flex-2">
        <label>Nombre *</label>
        <InputText v-model="local.name" required />
      </div>
      <div class="field">
        <label>Nivel *</label>
        <InputNumber v-model="local.level" :min="0" :max="9" />
        <small class="muted">0 = truco</small>
      </div>
      <div class="field">
        <label>Escuela</label>
        <Select v-model="local.school" :options="SCHOOLS" />
      </div>
    </div>

    <div class="row split">
      <div class="field">
        <label>Tiempo de lanzamiento</label>
        <InputText v-model="local.casting_time" placeholder="action / bonus / 1 minute" />
      </div>
      <div class="field">
        <label>Alcance</label>
        <InputText v-model="local.range" placeholder="60 ft / Toque / Personal" />
      </div>
      <div class="field">
        <label>Componentes</label>
        <InputText v-model="local.components" placeholder="V/S/M (foco)" />
      </div>
      <div class="field">
        <label>Duración</label>
        <InputText v-model="local.duration" placeholder="Instantáneo / Concentración, 1 min" />
      </div>
    </div>

    <div class="field">
      <label>Descripción *</label>
      <Textarea v-model="local.description" rows="5" required />
    </div>
    <div class="field">
      <label>Efecto al subir el slot</label>
      <Textarea v-model="local.higher_level" rows="2" placeholder="+1d8 por nivel sobre 1." />
    </div>

    <div class="row split">
      <div class="field">
        <label>Daño</label>
        <InputText v-model="local.damage" placeholder="2d8" />
      </div>
      <div class="field">
        <label>Tipo de daño</label>
        <Select v-model="local.damage_type" :options="DAMAGE_TYPES" :show-clear="true" />
      </div>
      <div class="field">
        <label>TS contra</label>
        <Select v-model="local.save_ability" :options="STATS" :show-clear="true" />
      </div>
    </div>

    <div class="row inline">
      <label class="check">
        <Checkbox v-model="local.ritual" :binary="true" />
        Ritual
      </label>
      <label class="check">
        <Checkbox v-model="local.concentration" :binary="true" />
        Requiere concentración
      </label>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Spell } from '~/models/Catalog'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Select from 'primevue/select'
import Textarea from 'primevue/textarea'
import Checkbox from 'primevue/checkbox'

const props = defineProps<{ modelValue: Partial<Spell> }>()
const emit = defineEmits<{ 'update:modelValue': [Partial<Spell>] }>()

const SCHOOLS = ['abjuración', 'conjuración', 'adivinación', 'encantamiento', 'evocación', 'ilusión', 'nigromancia', 'transmutación']
const STATS = ['FUE', 'DES', 'CON', 'INT', 'SAB', 'CAR']
const DAMAGE_TYPES = ['cortante', 'perforante', 'contundente', 'ácido', 'frío', 'fuego', 'fuerza', 'rayo', 'necrótico', 'veneno', 'psíquico', 'radiante', 'tronante']

const local = reactive<Partial<Spell>>({
  name: '', level: 0, school: 'evocación',
  casting_time: 'action', range: '60 ft', components: 'V/S', duration: 'Instantáneo',
  description: '', higher_level: '', ritual: false, concentration: false,
  damage: '', damage_type: '', save_ability: '',
  ...props.modelValue,
})

watch(local, (v) => emit('update:modelValue', { ...v }), { deep: true })
</script>

<style scoped>
.form { display: flex; flex-direction: column; gap: 12px; }
.row.split { display: grid; grid-template-columns: repeat(auto-fit, minmax(160px, 1fr)); gap: 12px; }
.row.inline { display: flex; gap: 18px; align-items: center; }
.field { display: flex; flex-direction: column; gap: 4px; }
.field.flex-2 { grid-column: span 2; }
.field label { font-size: 12px; font-weight: 600; color: #4b5563; }
.check { display: flex; align-items: center; gap: 8px; font-size: 13px; }
.muted { color: #6b7280; font-size: 11px; }
</style>
