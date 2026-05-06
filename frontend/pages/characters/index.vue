<template>
  <div class="characters-page">
    <header class="page-header">
      <div>
        <h1 class="title">Mis Personajes</h1>
        <p class="subtitle">Personajes propios y de las campañas en las que participas.</p>
      </div>
      <div class="header-actions">
        <Button label="Importar JSON" icon="pi pi-upload" severity="secondary" outlined @click="openImport" />
        <NuxtLink to="/characters/new" class="cta">
          <i class="pi pi-plus" /> Nuevo personaje
        </NuxtLink>
      </div>
    </header>

    <div v-if="pending" class="state">Cargando…</div>
    <div v-else-if="error" class="state error">No se pudieron cargar los personajes.</div>
    <div v-else-if="!data?.length" class="state empty">
      Aún no tienes personajes. Crea uno para empezar.
    </div>
    <div v-else class="grid">
      <article
        v-for="ch in data"
        :key="ch._id"
        class="card"
        @click="goTo(ch._id!)"
      >
        <div class="card-pic" :style="{ background: ch.color || 'var(--gold-deep, #6b4a16)' }">
          <span>{{ ch.pic || '⚔' }}</span>
        </div>
        <div class="card-body">
          <h2>{{ ch.name }}</h2>
          <p class="meta">
            <span>Nv {{ ch.level }}</span>
            <span class="dot">·</span>
            <span>{{ ch.background || 'Aventurero' }}</span>
          </p>
          <div class="hp-bar">
            <div class="hp-fill" :style="{ width: hpPct(ch) + '%' }" />
          </div>
          <p class="hp-label">{{ ch.hp }} / {{ ch.hp_max }} HP</p>
        </div>
      </article>
    </div>

    <!-- Input file oculto, disparado por el botón Importar -->
    <input
      ref="fileInput"
      type="file"
      accept="application/json,.json"
      style="display: none"
      @change="onFileSelected"
    />
  </div>
</template>

<script setup lang="ts">
import type { Character } from '~/models/Character'
import Button from 'primevue/button'

definePageMeta({ title: 'Mis personajes', auth: true })

const { data, pending, error, refresh } = await useCharacters()
const router = useRouter()
const { showSuccess, showError } = useAddToast()

const fileInput = ref<HTMLInputElement | null>(null)

function openImport() {
  fileInput.value?.click()
}

async function onFileSelected(ev: Event) {
  const target = ev.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return
  try {
    const text = await file.text()
    const portable = JSON.parse(text)
    const svc = useCharacterService()
    const created = await svc.importPortable(portable)
    showSuccess(`"${created.name}" importado.`)
    await refresh()
    router.push(`/characters/${created._id}`)
  } catch (e: any) {
    showError(e?.data?.message || e?.message || 'JSON inválido o referencias no encontradas')
  } finally {
    // Reseteamos el input para que se pueda subir el mismo archivo dos veces.
    if (fileInput.value) fileInput.value.value = ''
  }
}

function goTo(id: string) { router.push(`/characters/${id}`) }
function hpPct(ch: Character): number {
  if (!ch.hp_max) return 0
  return Math.max(0, Math.min(100, Math.round((ch.hp / ch.hp_max) * 100)))
}
</script>

<style scoped>
.characters-page {
  max-width: 1100px;
  margin: 0 auto;
  padding: 24px;
}
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  gap: 16px;
  margin-bottom: 24px;
  flex-wrap: wrap;
}
.header-actions { display: flex; gap: 8px; align-items: center; }
.title {
  font-size: 28px;
  font-weight: 700;
  letter-spacing: 0.02em;
  margin: 0;
}
.subtitle {
  color: var(--p-text-muted-color, #6b7280);
  margin: 4px 0 0;
}
.cta {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 18px;
  border-radius: 8px;
  background: var(--p-primary-color, #b48a3c);
  color: var(--p-primary-contrast-color, #fff);
  text-decoration: none;
  font-weight: 600;
}
.cta:hover { filter: brightness(1.1); }

.state {
  padding: 48px;
  text-align: center;
  color: var(--p-text-muted-color, #6b7280);
}
.state.error { color: #b91c1c; }
.state.empty { font-style: italic; }

.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}
.card {
  display: flex;
  gap: 14px;
  padding: 16px;
  background: var(--p-content-background, #fff);
  border: 1px solid var(--p-content-border-color, #e5e7eb);
  border-radius: 12px;
  cursor: pointer;
  transition: transform .15s, box-shadow .15s, border-color .15s;
}
.card:hover {
  transform: translateY(-1px);
  box-shadow: 0 8px 18px rgba(0,0,0,.08);
  border-color: var(--p-primary-color, #b48a3c);
}
.card-pic {
  flex: 0 0 64px;
  height: 64px;
  border-radius: 12px;
  display: grid;
  place-items: center;
  font-size: 30px;
  color: #fff;
  text-shadow: 0 1px 2px rgba(0,0,0,.4);
}
.card-body { flex: 1; min-width: 0; }
.card-body h2 {
  margin: 0 0 4px;
  font-size: 16px;
  font-weight: 700;
}
.meta {
  margin: 0 0 8px;
  font-size: 12px;
  color: var(--p-text-muted-color, #6b7280);
  display: flex;
  gap: 6px;
}
.dot { opacity: .5; }
.hp-bar {
  height: 6px;
  background: var(--p-content-border-color, #e5e7eb);
  border-radius: 999px;
  overflow: hidden;
}
.hp-fill {
  height: 100%;
  background: linear-gradient(90deg, #ef4444, #f59e0b, #10b981);
  transition: width .25s;
}
.hp-label {
  margin: 4px 0 0;
  font-size: 11px;
  color: var(--p-text-muted-color, #6b7280);
  font-variant-numeric: tabular-nums;
}
</style>
