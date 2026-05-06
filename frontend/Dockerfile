# syntax=docker/dockerfile:1.7
# Etapa 1: build de Nuxt.
FROM node:20-alpine AS builder

WORKDIR /src

# Cache de dependencias.
COPY package.json package-lock.json* ./
RUN npm ci --no-audit --no-fund

# Código completo (incluye nuxt.config.ts).
COPY . .

# Build SSR. Nuxt deja todo en .output/.
RUN npm run build

# Etapa 2: runtime mínimo.
FROM node:20-alpine AS runtime

WORKDIR /app
ENV NODE_ENV=production
ENV NITRO_HOST=0.0.0.0
ENV NITRO_PORT=3000

# Solo el output de Nuxt; sin node_modules de dev.
COPY --from=builder /src/.output ./.output

EXPOSE 3000
USER node
CMD ["node", ".output/server/index.mjs"]
