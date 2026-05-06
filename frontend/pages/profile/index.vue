<template>
  <div class="profile-page">
    <header class="page-header">
      <h1>Mi Perfil</h1>
      <p class="muted">Ajusta tu información, contraseña y rol activo.</p>
    </header>

    <section class="block">
      <h2>Información</h2>
      <form @submit.prevent="onSaveProfile">
        <div class="field">
          <label>Username</label>
          <InputText v-model="profile.username" />
          <small class="muted">Letras minúsculas, números, punto/guión/barra-baja. 3-24 caracteres.</small>
        </div>
        <div class="field">
          <label>Email</label>
          <InputText v-model="profile.email" type="email" />
        </div>
        <div class="field">
          <label>Avatar (emoji)</label>
          <InputText v-model="profile.avatar" placeholder="🧙 ⚔ 🐉" />
          <small class="muted">Por ahora aceptamos solo emoji o un texto corto.</small>
        </div>
        <div class="actions">
          <Button type="submit" label="Guardar cambios" :loading="savingProfile" />
        </div>
      </form>
    </section>

    <section class="block">
      <h2>Cambiar contraseña</h2>
      <form @submit.prevent="onChangePassword">
        <div class="field">
          <label>Contraseña actual *</label>
          <Password v-model="pwd.current" :feedback="false" toggleMask required />
        </div>
        <div class="field">
          <label>Nueva contraseña *</label>
          <Password v-model="pwd.next" toggleMask required />
          <small class="muted">Mínimo 8 caracteres.</small>
        </div>
        <div class="field">
          <label>Confirmar nueva contraseña *</label>
          <Password v-model="pwd.confirm" :feedback="false" toggleMask required />
        </div>
        <div class="actions">
          <Button type="submit" label="Actualizar contraseña" severity="warning" :loading="savingPwd" />
        </div>
      </form>
    </section>

    <section class="block">
      <h2>Rol activo</h2>
      <p class="muted">
        Tu cuenta tiene los roles: <strong>{{ profile.roles?.join(', ') || '—' }}</strong>.
        Elige cuál usas en esta sesión para ajustar la UI.
      </p>
      <div class="roles">
        <Button
          v-for="r in profile.roles || []"
          :key="r"
          :label="r.toUpperCase()"
          :severity="r === profile.active_rol ? 'primary' : 'secondary'"
          :outlined="r !== profile.active_rol"
          @click="onSetRole(r)"
        />
      </div>
    </section>

    <section class="block danger">
      <h2>Cerrar sesión</h2>
      <p class="muted">Cierra la sesión en este dispositivo.</p>
      <Button label="Cerrar sesión" icon="pi pi-sign-out" severity="danger" outlined @click="onLogout" />
    </section>
  </div>
</template>

<script setup lang="ts">
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import Button from 'primevue/button'

definePageMeta({ title: 'Perfil', auth: true })

const { showSuccess, showError } = useAddToast()
const { data, logout, refresh } = useAuth()

interface ProfileForm {
  username: string
  email: string
  avatar: string
  roles?: string[]
  active_rol?: string
}

const profile = reactive<ProfileForm>({
  username: data.value?.username || '',
  email: data.value?.email || '',
  avatar: data.value?.avatar || '',
  roles: (data.value as any)?.roles,
  active_rol: (data.value as any)?.active_rol,
})

watch(data, (v) => {
  if (!v) return
  profile.username = v.username || ''
  profile.email = v.email || ''
  profile.avatar = (v as any).avatar || ''
  profile.roles = (v as any).roles
  profile.active_rol = (v as any).active_rol
})

const pwd = reactive({ current: '', next: '', confirm: '' })

const savingProfile = ref(false)
const savingPwd = ref(false)

async function onSaveProfile() {
  savingProfile.value = true
  try {
    const api = useNuxtApp().$api as typeof $fetch
    await api('/user/me', {
      method: 'PUT',
      body: { username: profile.username, email: profile.email, avatar: profile.avatar },
    })
    showSuccess('Perfil actualizado.')
    await refresh()
  } catch (e: any) {
    showError(e?.data?.message || e?.message || 'Error al actualizar')
  } finally {
    savingProfile.value = false
  }
}

async function onChangePassword() {
  if (pwd.next !== pwd.confirm) {
    showError('Las contraseñas no coinciden')
    return
  }
  if (pwd.next.length < 8) {
    showError('La nueva contraseña debe tener al menos 8 caracteres')
    return
  }
  savingPwd.value = true
  try {
    const api = useNuxtApp().$api as typeof $fetch
    await api('/user/me/password', {
      method: 'PUT',
      body: { current_password: pwd.current, new_password: pwd.next },
    })
    showSuccess('Contraseña actualizada.')
    pwd.current = ''
    pwd.next = ''
    pwd.confirm = ''
  } catch (e: any) {
    showError(e?.data?.message || e?.message || 'Error al actualizar contraseña')
  } finally {
    savingPwd.value = false
  }
}

async function onSetRole(role: string) {
  try {
    const api = useNuxtApp().$api as typeof $fetch
    await api('/user/me/active-role', { method: 'PUT', body: { active_rol: role } })
    showSuccess(`Rol activo: ${role}`)
    await refresh()
  } catch (e: any) {
    showError(e?.message || 'Error al cambiar rol')
  }
}

async function onLogout() {
  await logout()
  navigateTo('/login')
}
</script>

<style scoped>
.profile-page { max-width: 720px; margin: 0 auto; padding: 24px; }
.page-header { margin-bottom: 18px; }
.page-header h1 { margin: 0 0 4px; font-size: 24px; }
.muted { color: #6b7280; font-size: 13px; }

.block {
  background: var(--p-content-background, #fff);
  border: 1px solid var(--p-content-border-color, #e5e7eb);
  border-radius: 12px;
  padding: 18px;
  margin-bottom: 14px;
}
.block.danger { border-color: #fecaca; }
.block h2 { font-size: 14px; letter-spacing: 0.05em; margin: 0 0 12px; }

.field { display: flex; flex-direction: column; gap: 4px; margin-bottom: 12px; }
.field label { font-size: 12px; font-weight: 600; color: #4b5563; }
.field :deep(.p-password) { width: 100%; }
.field :deep(.p-password-input) { width: 100%; }

.actions { display: flex; justify-content: flex-end; }
.roles { display: flex; gap: 8px; flex-wrap: wrap; }
</style>
