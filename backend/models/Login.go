package models

// Login: el frontend envía { user, password } — mantenemos esa convención.
type Login struct {
	User     string `json:"user"     bson:"user"     binding:"required"`
	Password string `json:"password" bson:"password" binding:"required"`
}
