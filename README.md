# DragonRoll

Compañero de mesa **D&D 5e self-hosted**, **data-driven** y libre de copyright. La aplicación viene **vacía** de contenido oficial: tú o tu DM crean las clases, subclases, conjuros, razas y rasgos como entradas en la base de datos. Esto te permite tener tu propio universo homebrew sin tocar reglas oficiales.

> **Nota legal:** Las mecánicas universales del sistema d20 (las 6 stats, las 18 categorías de habilidades, las condiciones, los tipos de daño) son ideas — no se copyrightean — y vienen sembradas en el primer arranque. Las clases, subclases, hechizos y rasgos los creas tú.

## Stack

| Capa | Tecnología |
|---|---|
| Frontend | Nuxt 3 + Vue 3 + PrimeVue + Tailwind |
| Backend | Go 1.24 + Gin + gin-jwt + bcrypt |
| Base de datos | MongoDB 7 |
| Reverse proxy / HTTPS | Caddy 2 (Let's Encrypt automático) |
| Orquestación | Docker Compose |
| Scripts auxiliares | Python 3 (parser SRD opcional) |

## Estructura

```
dragonroll/
├── backend/              # API Go (Gin + Mongo + JWT)
│   ├── models/           # Class, Subclass, Spell, Feature, Race, Skill, Character, Campaign, User
│   ├── repositories/     # CRUD Mongo (uno genérico para los 6 catálogos)
│   ├── services/         # Lógica + autorización (admin / DM / player)
│   ├── controllers/      # HTTP handlers
│   ├── routes/           # /api/v1/...
│   ├── middleware/       # JWT + CORS
│   ├── data/seed/        # JSON de bootstrap (skills universales + ejemplos homebrew)
│   ├── main.go
│   └── Dockerfile
├── frontend/             # Nuxt 3
│   ├── pages/characters/ # /characters, /characters/[id], /characters/new
│   ├── models/           # Character, Catalog (TS espejo del backend)
│   ├── services/         # CharacterService, CatalogService
│   ├── composables/      # useCharacters, useCatalog, useAuth
│   └── Dockerfile
├── scripts/              # parse_srd.py (placeholder)
├── docker-compose.yml    # mongo + backend + frontend + caddy
├── Caddyfile             # reverse proxy con HTTPS automático
├── .env.example
└── README.md
```

## Modelo de datos (resumen)

```
Class ─────────────────────┐
  • hit_die                ├── Character.class_id
  • level_progression[20]  │   .subclass_id
  • feature_ids por nivel  │   .race_id
  • spell_slots[9] por nv  │   .skill_proficiency_ids[]
                           │   .cantrip_ids[], .spells_known_ids[]
Subclass ──────────────────┘   .features (resueltos por nivel)

Spell, Feature, Race, Skill — entradas de catálogo independientes.
```

Cada entrada de catálogo tiene un **scope**:

- `global`: visible para todos. Solo el **admin** crea/edita.
- `campaign`: visible dentro de una campaña. El **DM** de esa campaña crea/edita.
- `user`: sandbox personal. Cualquier usuario.

El admin tiene control total sobre el contenido global. El DM tiene **control total sobre los personajes y catálogo de su campaña**, incluyendo editar HP en tiempo real durante el combate.

## Despliegue rápido (PC casero + DuckDNS)

### 1. DNS dinámico gratis

Andá a [duckdns.org](https://www.duckdns.org), iniciá sesión (con Google, GitHub, Twitter…) y reclamá un subdominio, p.ej. `marcelo-dragonroll`. Te queda `marcelo-dragonroll.duckdns.org`. Anotá tu **token DuckDNS** del dashboard.

Configurá la actualización automática de IP. La forma más simple: un cron en tu PC que haga `curl` al endpoint de DuckDNS cada 5 minutos:

```bash
# /etc/cron.d/duckdns
*/5 * * * * root curl -s "https://www.duckdns.org/update?domains=marcelo-dragonroll&token=TU_TOKEN&ip=" > /dev/null
```

### 2. Abrir puertos en tu router

Redirigí los puertos **80** y **443** del router a la IP local de tu PC casero. (Buscá "port forwarding" en el manual de tu router; la mayoría tiene una sección Virtual Server / Forwarding.)

### 3. Instalar Docker en el PC casero

Si es Ubuntu Server (recomendado):

```bash
sudo apt update && sudo apt install -y curl
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER
newgrp docker
```

### 4. Clonar y configurar

```bash
git clone TU_REPO dragonroll
cd dragonroll
cp .env.example .env
```

Editá `.env` y completá los secretos. Para generarlos:

```bash
echo "JWT_KEY=$(openssl rand -hex 32)"            >> .env
echo "SESSION_PASSWORD=$(openssl rand -hex 32)"   >> .env
echo "DB_PASS_MONGO=$(openssl rand -base64 24)"   >> .env
```

Y configurá:

```env
DRAGONROLL_DOMAIN=marcelo-dragonroll.duckdns.org
ACME_EMAIL=tu-email@ejemplo.com
PUBLIC_BASE_URL=https://marcelo-dragonroll.duckdns.org
CORS_URLS=https://marcelo-dragonroll.duckdns.org
ADMIN_PASSWORD=una-contraseña-fuerte
```

### 5. Levantar

```bash
docker compose --profile prod up -d --build
docker compose logs -f
```

La primera vez tarda 5-10 min (build de Go + Nuxt + descarga de Caddy).

Caddy va a obtener un certificado Let's Encrypt automáticamente. Si todo salió bien, en 1-2 minutos podés entrar a `https://marcelo-dragonroll.duckdns.org` y ver el login.

Hacé login con `admin` / la `ADMIN_PASSWORD` que pusiste, y empezá a crear clases.

## Despliegue de desarrollo (sin HTTPS)

Para iterar en tu propio Mac/PC sin Caddy:

```bash
# Solo Mongo + backend + frontend (sin perfil prod = sin Caddy):
docker compose up mongo backend frontend
```

Abrí `http://localhost:3000`. CORS_URLS por defecto ya permite localhost.

## Endpoints principales

### Auth
- `POST /api/v1/auth/login` `{ user, password }` → cookie `token` HttpOnly + body con JWT
- `POST /api/v1/auth/logout`
- `POST /api/v1/auth/refresh`

### Usuarios
- `POST /api/v1/user/` registro público
- `GET  /api/v1/user/me`
- `PUT  /api/v1/user/me/active-role` cambiar entre player/dm
- `GET  /api/v1/user/` listar (solo admin)

### Personajes
- `GET    /api/v1/characters/` lista accesible (propios + de campañas como DM/jugador)
- `POST   /api/v1/characters/`
- `GET    /api/v1/characters/:id?expand=true` hoja con catálogo poblado
- `PUT    /api/v1/characters/:id` reemplazo completo
- `PATCH  /api/v1/characters/:id/combat` parche rápido HP/conditions/slots
- `DELETE /api/v1/characters/:id`

### Campañas
- `GET    /api/v1/campaigns/`
- `POST   /api/v1/campaigns/`
- `GET    /api/v1/campaigns/:id`
- `PUT    /api/v1/campaigns/:id`
- `POST   /api/v1/campaigns/:id/players` agregar jugador por username/email
- `DELETE /api/v1/campaigns/:id`

### Catálogo (idéntico en cada uno)
- `GET    /api/v1/{classes,subclasses,spells,features,races,skills}` lista visible
- `POST   /api/v1/{...}` crear
- `GET    /api/v1/{...}/:id`
- `PUT    /api/v1/{...}/:id`
- `DELETE /api/v1/{...}/:id`

## Roadmap

- [x] Backend Go con CRUD de personajes y catálogo data-driven
- [x] Autorización por scope (admin global, DM de campaña, sandbox personal)
- [x] Frontend con dashboard y hoja de personaje renderizada con populate
- [x] Wizard de creación de PJ que lee del catálogo
- [x] Docker Compose con Caddy + HTTPS automático
- [x] **Editor visual del catálogo** (`/admin/catalog`) — listado, creación y edición de Class/Subclass/Spell/Feature/Race/Skill con formularios específicos por tipo
- [x] **Parser del SRD 5.1 oficial** (Python) → 319 hechizos importables al primer arranque
- [ ] Modo combate con grilla, iniciativa y tirada de dados
- [ ] Wiki de NPCs y bitácora de sesión por campaña
- [ ] Manager de mapas
- [ ] Hooks de fórmula en `Feature.Uses` (evaluar `5 × nivel` en runtime)
- [ ] Tirada de dados embebida + log compartido en tiempo real (WebSocket)

## Importar el SRD 5.1 al primer arranque

El SRD 5.1 (System Reference Document oficial de WotC) está bajo **CC-BY-4.0** y es libremente usable con atribución. El proyecto incluye un parser Python que convierte el SRD JSON oficial a un seed que el backend carga automáticamente.

### Pasos

1. Descargá el SRD JSON del repo público de Tabyltop:

   ```bash
   git clone https://github.com/Tabyltop/CC-SRD vendor/CC-SRD
   ```

2. Generá el seed de hechizos:

   ```bash
   python3 scripts/parse_srd.py \
     --input vendor/CC-SRD/SRD5.1-CCBY4.0License-TT.json \
     --output backend/data/seed/srd_spells.json
   ```

   Esto extrae **~319 hechizos** con casting time, range, components, duration, descripción y "At Higher Levels".

3. Reconstruí el backend y reiniciá:

   ```bash
   docker compose build backend
   docker compose up -d
   ```

   Al primer arranque (con la base de spells vacía), el backend inserta los hechizos con `meta.source_label = "SRD 5.1"` y `meta.scope = "global"`. Idempotente: no los duplica si ya están.

### Atribución obligatoria por la licencia

CC-BY-4.0 exige incluir la atribución del SRD en cualquier uso. El footer del frontend ya la incluye, y el campo `meta.source_label` deja registro en cada documento. **No quites el footer** mientras uses contenido del SRD.

## Editor visual del catálogo

Una vez logueado como admin/DM, andá a **`/admin/catalog`**. Tenés un CRUD completo para los 6 tipos:

- **Clases** con tabla editable de progresión 1-20 (PB, slots por nivel, features por nivel, casting)
- **Subclases** con clase padre, niveles donde aporta features, y conjuros otorgados
- **Conjuros** con todos los campos (componentes, escuela, duración, daño, TS, ritual, concentración)
- **Rasgos (Features)** con cost (action/bonus/passive), origen, recurso de uso (max + recharge)
- **Razas** con bonos de stat, idiomas y traits enlazados
- **Habilidades** (las 18 universales + las que crees)

Cada documento se crea con un **scope**:

- `global` → solo admin (visible para todos)
- `campaign` → DM de la campaña (visible para sus jugadores)
- `user` → cualquiera (sandbox personal)

El editor enlaza referencias automáticamente (al crear una clase, su pool de skills es un MultiSelect con las skills cargadas; al editar una subclase, el dropdown de "Clase padre" lista todas las clases visibles).

## Troubleshooting

**El certificado HTTPS no se emite.** Caddy necesita poder llegar al puerto 80 desde internet para el desafío HTTP-01. Verificá que el port forwarding funcione: desde otra red (4G del celular, p.ej.) probá `curl -I http://marcelo-dragonroll.duckdns.org`. Si da timeout, el router está bloqueando.

**El backend no levanta y se queja de `JWT_KEY`.** Necesita al menos 32 caracteres. Generalo con `openssl rand -hex 32`.

**"username ya existe" al hacer login con admin la primera vez.** Significa que el bootstrap no corrió porque la base no estaba realmente vacía (quizás un build anterior dejó datos). Borrá el volumen: `docker compose down -v` y volvé a levantar.

**Quiero ver qué hay en Mongo.** Conectá `mongosh`:
```bash
docker compose exec mongo mongosh -u $DB_USER_MONGO -p $DB_PASS_MONGO --authenticationDatabase admin
> use dragonroll
> show collections
> db.classes.find().pretty()
```

## Licencia

MIT para el código. El contenido que tú/tus DMs creen en el catálogo es de su autoría.
