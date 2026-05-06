<template>
  <fieldset class="meta-fieldset">
    <legend>Visibilidad y origen</legend>
    <div class="grid">
      <div class="field">
        <label>Alcance (scope)</label>
        <Select
          v-model="local.scope"
          :options="scopeOptions"
          option-label="label"
          option-value="value"
          @update:model-value="emitChange"
        />
        <small class="muted">{{ scopeHelp }}</small>
      </div>

      <div v-if="local.scope === 'campaign'" class="field">
        <label>Campaña</label>
        <Select
          v-model="local.campaign_id"
          :options="campaigns || []"
          option-label="name"
          option-value="_id"
          placeholder="Elegí campaña"
          @update:model-value="emitChange"
        />
      </div>

      <div class="field">
        <label>Etiqueta de origen</label>
        <InputText
          v-model="local.source_label"
          placeholder="ej. SRD 5.1, Castoñegro v2, Homebrew Marcelo"
          @input="emitChange"
        />
        <small class="muted">Sólo descriptiva. Útil para filtrar después.</small>
      </div>

      <div class="field checkbox-field">
        <label>
          <Checkbox v-model="local.is_homebrew" :binary="true" @change="emitChange" />
          Es contenido homebrew
        </label>
        <small class="muted">Marca para distinguir de contenido oficial importado.</small>
      </div>
    </div>
  </fieldset>
</template>

<script setup lang="ts">
import type { CatalogMeta, Scope } from '~/models/Catalog'
import InputText from 'primevue/inputtext'
import Select from 'primevue/select'
import Checkbox from 'primevue/checkbox'

const props = defineProps<{
  modelValue: CatalogMeta
  /** Lista de campañas donde el usuario puede crear contenido (DM/admin). */
  campaigns?: Array<{ _id: string, name: string }>
  /** Si el usuario es admin, se permite crear con scope=global. */
  canCreateGlobal: boolean
}>()

const emit = defineEmits<{ 'update:modelValue': [CatalogMeta] }>()

const local = reactive<CatalogMeta>({
  scope: props.modelValue.scope || 'user',
  campaign_id: props.modelValue.campaign_id,
  owner_id: props.modelValue.owner_id,
  is_homebrew: props.modelValue.is_homebrew ?? true,
  source_label: props.modelValue.source_label || '',
  created_at: props.modelValue.created_at || '',
  updated_at: props.modelValue.updated_at || '',
})

watch(() => props.modelValue, (v) => {
  Object.assign(local, v)
}, { deep: true })

const scopeOptions = computed(() => {
  const opts: Array<{ value: Scope, label: string }> = [
    { value: 'user', label: 'Personal (sólo tú)' },
    { value: 'campaign', label: 'Campaña (DM y jugadores)' },
  ]
  if (props.canCreateGlobal) {
    opts.unshift({ value: 'global', label: 'Global (todos)' })
  }
  return opts
})

const scopeHelp = computed(() => {
  switch (local.scope) {
    case 'global':   return 'Visible para todos los usuarios. Sólo el admin lo edita.'
    case 'campaign': return 'Visible para los miembros de la campaña. Lo edita el DM.'
    default:         return 'Sólo tú lo ves. Sandbox personal.'
  }
})

function emitChange() {
  emit('update:modelValue', { ...local })
}
</script>

<style scoped>
.meta-fieldset {
  border: 1px solid var(--p-content-border-color, #e5e7eb);
  border-radius: 10px;
  padding: 12px 16px;
  margin: 12px 0;
}
.meta-fieldset legend {
  padding: 0 6px;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: #6b7280;
}
.grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 12px;
}
.field { display: flex; flex-direction: column; gap: 4px; }
.field label { font-size: 12px; font-weight: 600; color: #4b5563; }
.checkbox-field label { display: flex; align-items: center; gap: 8px; font-weight: 500; }
.muted { color: #6b7280; font-size: 11px; }
</style>
