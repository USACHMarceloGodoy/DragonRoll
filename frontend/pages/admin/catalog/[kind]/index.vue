<template>
  <div class="list-page">
    <header class="page-header">
      <div>
        <NuxtLink to="/admin/catalog" class="back">← Catálogo</NuxtLink>
        <h1>{{ pageLabel }}</h1>
      </div>
      <NuxtLink :to="`/admin/catalog/${kind}/new`" class="cta">
        <i class="pi pi-plus" /> Nuevo
      </NuxtLink>
    </header>

    <div class="filters">
      <InputText v-model="filter" placeholder="Filtrar por nombre…" />
      <Select v-model="scopeFilter" :options="SCOPES" option-label="label" option-value="value" placeholder="Cualquier scope" :show-clear="true" />
    </div>

    <div v-if="pending" class="state">Cargando…</div>
    <div v-else-if="error" class="state error">No se pudieron cargar.</div>
    <div v-else-if="!filtered.length" class="state empty">Sin resultados.</div>

    <div v-else class="cards">
      <article v-for="doc in filtered" :key="doc._id" class="card">
        <div class="card-main">
          <h2>
            <NuxtLink :to="`/admin/catalog/${kind}/${doc._id}`">
              {{ doc.name || doc.code || '(sin nombre)' }}
            </NuxtLink>
          </h2>
          <p class="meta">
            <span class="badge" :class="`badge-${doc.meta?.scope}`">{{ scopeLabel(doc.meta?.scope) }}</span>
            <span v-if="doc.meta?.source_label" class="src">{{ doc.meta.source_label }}</span>
            <span v-if="doc.meta?.is_homebrew" class="badge badge-homebrew">homebrew</span>
            <span v-if="extraDetail(doc)" class="extra">{{ extraDetail(doc) }}</span>
          </p>
        </div>
        <Button icon="pi pi-trash" text rounded severity="danger" @click="onDelete(doc)" />
      </article>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { CatalogKind } from '~/models/Catalog'
import InputText from 'primevue/inputtext'
import Select from 'primevue/select'
import Button from 'primevue/button'

definePageMeta({ title: 'Catálogo', auth: true })

const route = useRoute()
const router = useRouter()
const { showSuccess, showError } = useAddToast()

const kind = computed(() => route.params.kind as CatalogKind)

const LABELS: Record<CatalogKind, string> = {
  classes: 'Clases',
  subclasses: 'Subclases',
  spells: 'Conjuros',
  features: 'Rasgos',
  races: 'Razas',
  skills: 'Habilidades',
}
const pageLabel = computed(() => LABELS[kind.value] || kind.value)

const SCOPES = [
  { value: 'global',   label: 'Global' },
  { value: 'campaign', label: 'Campaña' },
  { value: 'user',     label: 'Personal' },
]

const filter = ref('')
const scopeFilter = ref<string | null>(null)

const { data, pending, error, refresh } = await useCatalogList(kind.value)

// Cuando cambia kind, recarga.
watch(kind, () => refresh())

const filtered = computed(() => {
  const arr = (data.value as any[] | null) || []
  return arr.filter((d) => {
    const matchesText = !filter.value || (d.name || d.code || '').toLowerCase().includes(filter.value.toLowerCase())
    const matchesScope = !scopeFilter.value || d.meta?.scope === scopeFilter.value
    return matchesText && matchesScope
  })
})

function scopeLabel(s?: string): string {
  return SCOPES.find(x => x.value === s)?.label || s || ''
}

/** Detalle secundario por tipo (ej. nivel del conjuro, dado de golpe). */
function extraDetail(doc: any): string | null {
  switch (kind.value) {
    case 'spells':   return doc.level === 0 ? 'Truco' : `Nivel ${doc.level}`
    case 'classes':  return `Hit die ${doc.hit_die}`
    case 'subclasses': return doc.parent_class ? 'subclase' : ''
    case 'features': return doc.cost ? `cost ${doc.cost}` : ''
    case 'races':    return doc.size || ''
    case 'skills':   return doc.stat || ''
  }
  return null
}

async function onDelete(doc: any) {
  if (!confirm(`¿Borrar "${doc.name || doc.code}"? Esta acción no se puede deshacer.`)) return
  try {
    const svc = useCatalogService()
    await svc.remove(kind.value, doc._id)
    showSuccess('Eliminado.')
    await refresh()
  } catch (e: any) {
    showError(e?.message || 'No se pudo eliminar')
  }
}
</script>

<style scoped>
.list-page { max-width: 1100px; margin: 0 auto; padding: 24px; }
.page-header { display: flex; justify-content: space-between; align-items: flex-end; gap: 16px; margin-bottom: 16px; flex-wrap: wrap; }
.page-header h1 { margin: 4px 0 0; font-size: 24px; }
.back { color: #6b7280; font-size: 12px; text-decoration: none; }
.back:hover { color: #b48a3c; }

.cta {
  display: inline-flex; align-items: center; gap: 8px;
  padding: 10px 18px; border-radius: 8px;
  background: var(--p-primary-color, #b48a3c); color: #fff;
  text-decoration: none; font-weight: 600;
}

.filters { display: flex; gap: 10px; margin-bottom: 14px; flex-wrap: wrap; }

.state { padding: 48px; text-align: center; color: #6b7280; }
.state.error { color: #b91c1c; }
.state.empty { font-style: italic; }

.cards { display: flex; flex-direction: column; gap: 8px; }
.card {
  display: flex; align-items: center; gap: 12px;
  padding: 12px 14px;
  background: var(--p-content-background, #fff);
  border: 1px solid var(--p-content-border-color, #e5e7eb);
  border-radius: 10px;
}
.card-main { flex: 1; min-width: 0; }
.card h2 { margin: 0; font-size: 15px; font-weight: 600; }
.card h2 a { color: inherit; text-decoration: none; }
.card h2 a:hover { color: var(--p-primary-color, #b48a3c); }
.card .meta { margin: 4px 0 0; display: flex; gap: 8px; align-items: center; flex-wrap: wrap; font-size: 12px; }
.badge {
  padding: 2px 8px; border-radius: 999px; font-size: 10px; font-weight: 600;
  text-transform: uppercase; letter-spacing: 0.06em;
}
.badge-global { background: #dbeafe; color: #1e40af; }
.badge-campaign { background: #fef3c7; color: #92400e; }
.badge-user { background: #f3f4f6; color: #4b5563; }
.badge-homebrew { background: #ede9fe; color: #5b21b6; }
.src { color: #6b7280; }
.extra { color: #6b7280; margin-left: auto; }
</style>
