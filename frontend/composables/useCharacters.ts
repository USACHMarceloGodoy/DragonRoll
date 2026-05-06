import { CharacterService } from '~/services/Character/character.service'

/**
 * useCharacters: lista los personajes accesibles del usuario actual.
 * Usa useAsyncData para caché por SSR/CSR.
 */
export const useCharacters = () => {
  const svc = new CharacterService()
  return useAsyncData('characters:list', () => svc.list())
}

export const useCharacter = (id: string | (() => string)) => {
  const svc = new CharacterService()
  const idRef = computed(() => (typeof id === 'function' ? id() : id))
  return useAsyncData(
    () => `characters:${idRef.value}`,
    () => svc.get(idRef.value),
    { watch: [idRef] },
  )
}

export const useCharacterService = () => new CharacterService()
