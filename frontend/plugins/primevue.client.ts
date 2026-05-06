import ToastService from 'primevue/toastservice'
import PrimeVue from 'primevue/config'

export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.vueApp.use(ToastService)
  nuxtApp.vueApp.use(PrimeVue)
})
