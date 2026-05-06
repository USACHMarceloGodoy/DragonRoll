import tailwindcss from '@tailwindcss/vite'
import Aura from '@primeuix/themes/aura'
import myPreset from './theme/primePreset'

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  modules: [
    '@nuxt/eslint',
    '@nuxt/fonts',
    '@primevue/nuxt-module',
    'nuxt-security',
  ],
  devtools: { enabled: true },
  app: {
    head: {
      style: [
        {
          textContent: '@layer reset, primevue, app;',
          key: 'css-layer-order',
        }],
    },
  },
  css: [
    '~/assets/css/main.css',
    '~/assets/scss/styles.scss',
  ],
  runtimeConfig: {
    public: {
      baseUrl: 'http://localhost:3000',
      backendUrl: 'http://localhost:8080/api/v1',
      authUrl: 'http://localhost:8080/api/v1/auth',
    },
    sessionPassword: 'secret',
    cookieName: '__session',
    cookieExpires: 86400,
  },
  compatibilityDate: '2025-05-15',
  vite: {
    plugins: [
      tailwindcss(),
    ],
    css: {
      preprocessorOptions: {
        scss: {
          additionalData:
            '@use "@/assets/scss/_variables.scss" as *;',
        },
      },
    },
  },
  postcss: {
    plugins: {
      autoprefixer: {},
    },
  },
  eslint: {
    config: {
      stylistic: true,
    },
  },
  primevue: {
    options: {
      ripple: true,
      unstyled: false,
      theme: {
        preset: myPreset,
        options: {
          // cssLayer: { name: 'primevue', order: 'reset, primevue, app' },
          darkModeSelector: 'none',
        },
      },
    },
  },
})
