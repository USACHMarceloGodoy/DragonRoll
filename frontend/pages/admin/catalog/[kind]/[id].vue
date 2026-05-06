<template>
  <div class="editor-page">
    <header class="page-header">
      <NuxtLink :to="`/admin/catalog/${kind}`" class="back">← Volver</NuxtLink>
      <h1>{{ isNew ? 'Nuevo' : 'Editar' }} · {{ kindLabel }}</h1>
    </header>

    <div v-if="loading" class="state">Cargando…</div>

    <form v-else class="form" @submit.prevent="onSave">
      <!-- Form específico al tipo. Cada uno emite update:modelValue. -->
      <div class="block">
        <h2>{{ kindLabel }}</h2>

        <SkillForm    v-if="kind === 'skills'"     v-model="doc" />
        <SpellForm    v-else-if="kind === 'spells'" v-model="doc" />
        <FeatureForm  v-else-if="kind === 'features'" v-model="doc" />
        <RaceForm     v-else-if="kind === 'races'"   v-model="doc" :features="(refs.features as any) || []" />
        <SubclassForm v-else-if="kind === 'subclasses'" v-model="doc" :classes="(refs.classes as any) || []" :spells="(refs.spells as any) || []" :features="(refs.features as any) || []" />
        <ClassForm    v-else-if="kind === 'classes'" v-model="doc" :skills="(refs.skills as any) || []" :spells="(refs.spells as any) || []" :features="(refs.features as any) || []" />
      </div>

      <!-- Meta común -->
      <div class="block">
        <CatalogMetaEditor v-model="meta" :campaigns="campaigns" :can-create-global="isAdmin" />
      </div>

      <div class="actions">
        <NuxtLink :to="`/admin/catalog/${kind}`" class="btn-secondary">Cancelar</NuxtLink>
        <Button type="submit" :label="isNew ? 'Crear' : 'Guardar cambios'" :loading="saving" severity="primary" />
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import type { CatalogKind, CatalogMeta } from '~/models/Catalog'
import SkillForm from '~/components/CatalogForm/SkillForm.vue'
import SpellForm from '~/components/CatalogForm/SpellForm.vue'
import FeatureForm from '~/components/CatalogForm/FeatureForm.vue'
import RaceForm from '~/components/CatalogForm/RaceForm.vue'
import SubclassForm from '~/components/CatalogForm/SubclassForm.vue'
import ClassForm from '~/components/CatalogForm/ClassForm.vue'
import CatalogMetaEditor from '~/components/CatalogForm/MetaEditor.vue'
import Button from 'primevue/button'

definePageMeta({ title: 'Editor catálogo', auth: true })

const route = useRoute()
const router = useRouter()
const { showSuccess, showError } = useAddToast()
const { data: authData } = useAuth()

const kind = computed(() => route.params.kind as CatalogKind)
const id = computed(() => route.params.id as string)
const isNew = computed(() => id.value === 'new' || !id.value)

const KIND_LABELS: Record<CatalogKind, string> = {
  classes: 'Clase',
  subclasses: 'Subclase',
  spells: 'Conjuro',
  features: 'Rasgo',
  races: 'Raza',
  skills: 'Habilidad',
}
const kindLabel = computed(() => KIND_LABELS[kind.value] || kind.value)

const isAdmin = computed(() => authData.value?.role === 'admin')

const doc = ref<any>({})
const meta = ref<CatalogMeta>({
  scope: 'user',
  is_homebrew: true,
  source_label: '',
  created_at: '',
  updated_at: '',
})

const refs = reactive<{ classes?: any[], spells?: any[], features?: any[], skills?: any[] }>({})
const campaigns = ref<Array<{ _id: string, name: string }>>([])
const loading = ref(true)
const saving = ref(false)

/**
 * Carga las referencias que el form puede necesitar (skills/spells/features/classes).
 * Lo hacemos en paralelo y le damos al usuario un loading inicial.
 */
async function loadRefs() {
  const svc = useCatalogService()
  const promises: Promise<any>[] = []
  // Sólo lo necesario por kind, pero por simplicidad cargamos lo común.
  if (['classes', 'subclasses', 'races'].includes(kind.value)) promises.push(svc.list('features').then(d => refs.features = d))
  if (['classes', 'subclasses'].includes(kind.value)) promises.push(svc.list('spells').then(d => refs.spells = d))
  if (kind.value === 'classes') promises.push(svc.list('skills').then(d => refs.skills = d))
  if (kind.value === 'subclasses') promises.push(svc.list('classes').then(d => refs.classes = d))
  await Promise.all(promises)
}

async function loadCampaigns() {
  // Si más adelante agregamos endpoint /campaigns, traemos aquí; por ahora
  // dejamos vacío salvo que el endpoint exista.
  try {
    const data = await ($fetch as any)('/campaigns/', { baseURL: useRuntimeConfig().public.backendUrl })
    campaigns.value = (data?.data || []).map((c: any) => ({ _id: c._id, name: c.name }))
  } catch { /* sin campañas, o endpoint no autenticado: ignorar silenciosamente */ }
}

async function loadDoc() {
  loading.value = true
  try {
    await loadRefs()
    if (!isNew.value) {
      const svc = useCatalogService()
      const fetched = await (svc.get as any)(kind.value, id.value)
      doc.value = { ...fetched }
      meta.value = fetched.meta || meta.value
    }
  } catch (e: any) {
    showError(e?.message || 'No se pudo cargar')
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await Promise.all([loadDoc(), loadCampaigns()])
})

async function onSave() {
  saving.value = true
  try {
    const payload = { ...doc.value, meta: meta.value }
    const svc = useCatalogService()
    if (isNew.value) {
      const created = await svc.create(kind.value, payload)
      showSuccess('Creado.')
      router.push(`/admin/catalog/${kind.value}/${created._id}`)
    } else {
      await svc.update(kind.value, id.value, payload)
      showSuccess('Cambios guardados.')
    }
  } catch (e: any) {
    showError(e?.message || 'No se pudo guardar')
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.editor-page { max-width: 1100px; margin: 0 auto; padding: 24px; }
.page-header { margin-bottom: 16px; }
.page-header h1 { margin: 4px 0 0; font-size: 24px; }
.back { color: #6b7280; font-size: 12px; text-decoration: none; }
.back:hover { color: #b48a3c; }

.state { padding: 48px; text-align: center; color: #6b7280; }
.form { display: flex; flex-direction: column; gap: 14px; }
.block {
  background: var(--p-content-background, #fff);
  border: 1px solid var(--p-content-border-color, #e5e7eb);
  border-radius: 12px;
  padding: 16px 18px;
}
.block h2 { margin: 0 0 12px; font-size: 14px; letter-spacing: 0.05em; }

.actions { display: flex; justify-content: flex-end; gap: 10px; }
.btn-secondary {
  padding: 10px 16px; border-radius: 8px;
  border: 1px solid var(--p-content-border-color, #e5e7eb);
  color: #4b5563; text-decoration: none;
}
</style>
