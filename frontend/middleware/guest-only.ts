export default defineNuxtRouteMiddleware((to, from) => {
    if (import.meta.server) {
      return
    }
    const { status } = useAuth()
    const isAuthenticated = status.value === 'authenticated'
    console.log({ to, from, isAuthenticated })
    if (isAuthenticated) {
      if (import.meta.server) return navigateTo({ name: "index" });
      return abortNavigation();
  
    }
})
