import { useToast } from 'primevue/usetoast'

export const useAddToast = () => {
  const toast = useToast()

  const showSuccess = (detail: string, summary: string = 'Éxito') => {
    toast.add({ severity: 'success', summary, detail, life: 3000 })
  }

  const showError = (detail: string, summary: string = 'Error') => {
    toast.add({ severity: 'error', summary, detail, life: 4000 })
  }

  const showInfo = (detail: string, summary: string = 'Información') => {
    toast.add({ severity: 'info', summary, detail, life: 3000 })
  }

  const showWarn = (detail: string, summary: string = 'Advertencia') => {
    toast.add({ severity: 'warn', summary, detail, life: 3000 })
  }

  return {
    showSuccess,
    showError,
    showInfo,
    showWarn,
  }
}
