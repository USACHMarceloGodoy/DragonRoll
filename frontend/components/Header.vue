<template>
  <header class="dr-header">
    <NuxtLink to="/" class="brand">
      <span class="brand-mark">20</span>
      <span class="brand-text">DRAGONROLL</span>
    </NuxtLink>

    <nav v-if="status === 'authenticated'" class="nav">
      <NuxtLink to="/characters">Personajes</NuxtLink>
      <NuxtLink to="/campaigns">Campañas</NuxtLink>
      <NuxtLink v-if="isDM || isAdmin" to="/admin/catalog">Catálogo</NuxtLink>
    </nav>

    <div class="spacer" />

    <div v-if="status === 'authenticated'" class="profile">
      <NuxtLink to="/profile" class="username-link">
        <span class="avatar">{{ (data as any)?.avatar || '👤' }}</span>
        <span class="username">{{ data?.username }}</span>
      </NuxtLink>
      <Button icon="pi pi-sign-out" text rounded severity="secondary" aria-label="Cerrar sesión" @click="logout" />
    </div>
    <NuxtLink v-else to="/login" class="login-link">Ingresar</NuxtLink>
  </header>
</template>

<script setup lang="ts">
import Button from 'primevue/button'

const { status, data, logout } = useAuth()

const isAdmin = computed(() => data.value?.role === 'admin')
const isDM = computed(() => data.value?.role === 'dm' || isAdmin.value)
</script>

<style scoped>
.dr-header {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 10px 20px;
  background: linear-gradient(180deg, #1f1a14, #15110c);
  border-bottom: 1px solid #6b4a16;
  color: #e6dfd2;
}
.brand { display: flex; align-items: center; gap: 10px; text-decoration: none; }
.brand-mark {
  display: grid; place-items: center;
  width: 28px; height: 28px;
  background: radial-gradient(circle, #d4a04a, #6b4a16);
  clip-path: polygon(50% 0%, 100% 25%, 100% 75%, 50% 100%, 0 75%, 0 25%);
  font-weight: 800; color: #1a1410; font-size: 13px;
}
.brand-text {
  color: #d4a04a;
  font-weight: 700;
  letter-spacing: 0.14em;
  font-size: 14px;
}
.nav { display: flex; gap: 4px; }
.nav a {
  padding: 8px 14px;
  border-radius: 6px;
  color: #c9c0ac;
  text-decoration: none;
  font-size: 12px;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  font-weight: 600;
}
.nav a:hover { color: #d4a04a; background: rgba(212,160,74,0.08); }
.nav a.router-link-active { color: #d4a04a; background: rgba(212,160,74,0.12); }
.spacer { flex: 1; }
.profile { display: flex; align-items: center; gap: 8px; }
.username-link {
  display: flex; align-items: center; gap: 8px;
  text-decoration: none;
  padding: 4px 10px; border-radius: 6px;
}
.username-link:hover { background: rgba(212,160,74,0.08); }
.avatar { font-size: 18px; }
.username { font-size: 13px; color: #c9c0ac; }
.login-link {
  padding: 6px 14px;
  border: 1px solid #d4a04a;
  border-radius: 6px;
  color: #d4a04a;
  text-decoration: none;
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.1em;
}
</style>
