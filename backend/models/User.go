package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User representa una cuenta de la app DragonRoll.
// Un mismo usuario puede ser DM en una campaña y Player en otra; por eso
// se separan los Roles habilitados del ActiveRol seleccionado en sesión.
type User struct {
	ID        primitive.ObjectID `json:"_id"        bson:"_id,omitempty"`
	Username  string             `json:"username"   bson:"username"`         // único, login
	Email     string             `json:"email"      bson:"email"`            // único
	Password  string             `json:"-"          bson:"password"`         // hash bcrypt; nunca se serializa hacia el cliente
	Roles     []Rol              `json:"roles"      bson:"roles"`            // roles habilitados
	ActiveRol Rol                `json:"active_rol" bson:"active_rol"`       // rol activo de la sesión
	Avatar    string             `json:"avatar"     bson:"avatar,omitempty"` // emoji o URL pequeña
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

// Rol es el rol de un usuario en la app.
type Rol string

const (
	ADMIN  Rol = "admin"  // mantenimiento de la plataforma
	DM     Rol = "dm"     // Dungeon Master
	PLAYER Rol = "player" // jugador
	ALL    Rol = "all"
)

var ADMIN_ROLES = []Rol{ADMIN}
var DM_ROLES = []Rol{ADMIN, DM}
var ALL_ROLES = []Rol{ADMIN, DM, PLAYER}
