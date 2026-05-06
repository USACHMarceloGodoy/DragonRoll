<template>
  <div />
</template>

<script setup lang="ts">
definePageMeta({ title: 'Inicio', auth: true })

// El dashboard real es /characters. Esta página existe sólo como
// landing tras el login y redirige inmediatamente.
onMounted(() => {
  navigateTo('/characters')
})
</script>
