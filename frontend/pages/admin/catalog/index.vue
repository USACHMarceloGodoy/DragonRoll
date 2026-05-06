<template>
  <div class="admin-page">
    <header class="page-header">
      <h1>Catálogo</h1>
      <p class="muted">Crea y edita el contenido del juego: clases, subclases, conjuros, rasgos, razas y habilidades.</p>
    </header>

    <div class="grid">
      <NuxtLink v-for="t in TYPES" :key="t.kind" :to="`/admin/catalog/${t.kind}`" class="card">
        <span class="icon">{{ t.icon }}</span>
        <div>
          <h2>{{ t.label }}</h2>
          <p>{{ t.desc }}</p>
        </div>
      </NuxtLink>
    </div>

    <div v-if="!isAdmin && !isDM" class="warning">
      <i class="pi pi-info-circle" />
      Tu cuenta es de jugador. Puedes crear contenido <strong>personal</strong> (sólo tú lo verás). Para crear contenido de campaña, pídele a tu DM que te invite a una.
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ title: 'Catálogo', auth: true })

const { data } = useAuth()
const isAdmin = computed(() => data.value?.role === 'admin')
const isDM = computed(() => data.value?.role === 'dm' || isAdmin.value)

const TYPES = [
  { kind: 'classes',    label: 'Clases',     icon: '⚔', desc: 'Define hit die, progresión 1-20, slots, casting.' },
  { kind: 'subclasses', label: 'Subclases',  icon: '✦', desc: 'Variantes de cada clase con sus features y conjuros.' },
  { kind: 'spells',     label: 'Conjuros',   icon: '✶', desc: 'Trucos y conjuros de nivel 1 a 9.' },
  { kind: 'features',   label: 'Rasgos',     icon: '◈', desc: 'Habilidades de clase, subclase, raza, fondo o feat.' },
  { kind: 'races',      label: 'Razas',      icon: '☥', desc: 'Velocidad, tamaño, bonos de stat, idiomas, rasgos.' },
  { kind: 'skills',     label: 'Habilidades', icon: '◇', desc: 'Las 18 universales más las que crees tú.' },
]
</script>

<style scoped>
.admin-page { max-width: 1100px; margin: 0 auto; padding: 24px; }
.page-header { margin-bottom: 18px; }
.page-header h1 { margin: 0 0 4px; font-size: 24px; }
.muted { color: #6b7280; }

.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 14px;
}
.card {
  display: flex;
  gap: 14px;
  padding: 16px;
  background: var(--p-content-background, #fff);
  border: 1px solid var(--p-content-border-color, #e5e7eb);
  border-radius: 12px;
  text-decoration: none;
  color: inherit;
  transition: transform .15s, box-shadow .15s, border-color .15s;
}
.card:hover {
  transform: translateY(-1px);
  box-shadow: 0 8px 18px rgba(0,0,0,.08);
  border-color: var(--p-primary-color, #b48a3c);
}
.card .icon {
  flex: 0 0 56px;
  width: 56px; height: 56px;
  border-radius: 12px;
  background: linear-gradient(135deg, #6b4a16, #d4a04a);
  color: #fff;
  display: grid; place-items: center;
  font-size: 28px;
  text-shadow: 0 1px 2px rgba(0,0,0,.4);
}
.card h2 { margin: 0; font-size: 16px; }
.card p { margin: 4px 0 0; font-size: 12px; color: #6b7280; }

.warning {
  margin-top: 18px;
  padding: 12px 16px;
  background: var(--p-yellow-50, #fefce8);
  border: 1px solid var(--p-yellow-200, #fde68a);
  border-radius: 8px;
  color: #713f12;
  font-size: 13px;
}
</style>
