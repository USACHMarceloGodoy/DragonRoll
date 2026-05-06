export default defineEventHandler(async (event) => {
  const activeSession = event.context.session
  if (!activeSession) {
    return null
  }
  return activeSession
})
