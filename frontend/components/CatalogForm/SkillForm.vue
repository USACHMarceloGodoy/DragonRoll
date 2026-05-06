<template>
  <div class="form">
    <div class="row split">
      <div class="field">
        <label>Nombre *</label>
        <InputText v-model="local.name" required />
      </div>
      <div class="field">
        <label>Código (slug) *</label>
        <InputText v-model="local.code" placeholder="acrobatics, athletics, …" required />
        <small class="muted">Identificador estable, no cambies si ya tienes PJ que la usen.</small>
      </div>
      <div class="field">
        <label>Atributo asociado</label>
        <Select v-model="local.stat" :options="STATS" />
      </div>
    </div>
    <div class="field">
      <label>Descripción</label>
      <Textarea v-model="local.description" rows="3" />
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Skill } from '~/models/Catalog'
import InputText from 'primevue/inputtext'
import Select from 'primevue/select'
import Textarea from 'primevue/textarea'

const props = defineProps<{ modelValue: Partial<Skill> }>()
const emit = defineEmits<{ 'update:modelValue': [Partial<Skill>] }>()

const STATS = ['FUE', 'DES', 'CON', 'INT', 'SAB', 'CAR']

const local = reactive<Partial<Skill>>({
  name: '', code: '', stat: 'FUE', description: '',
  ...props.modelValue,
})

watch(local, (v) => emit('update:modelValue', { ...v }), { deep: true })
</script>

<style scoped>
.form { display: flex; flex-direction: column; gap: 12px; }
.row.split { display: grid; grid-template-columns: repeat(auto-fit, minmax(180px, 1fr)); gap: 12px; }
.field { display: flex; flex-direction: column; gap: 4px; }
.field label { font-size: 12px; font-weight: 600; color: #4b5563; }
.muted { color: #6b7280; font-size: 11px; }
</style>
