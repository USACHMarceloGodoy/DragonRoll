import { CatalogService } from '~/services/Catalog/catalog.service'
import type { CatalogKind } from '~/models/Catalog'

/**
 * useCatalogList: helper para componentes que necesitan listar
 * clases, razas, conjuros, etc. Cachea por kind.
 */
export const useCatalogList = (kind: CatalogKind) => {
  const svc = new CatalogService()
  return useAsyncData(`catalog:${kind}`, () => {
    // El service usa overloads tipados: forzamos any aquí para no
    // duplicar el switch — el caller hace el casteo correcto.
    return (svc.list as (k: CatalogKind) => Promise<unknown[]>)(kind)
  })
}

export const useCatalogService = () => new CatalogService()
