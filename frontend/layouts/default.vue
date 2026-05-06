<template>
  <div>
    <Toast />
    <Header />
    <div class="layout">
      <slot />
    </div>
    <Footer />
  </div>
</template>

<style scoped>
.layout {
  min-height: 100vh;
  background-color: var(--color-grey-1);
  padding: 8px !important;
}
</style>
