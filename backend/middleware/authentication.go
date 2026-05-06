package middleware

import (
	"dragonroll/backend/config"
	"dragonroll/backend/models"
	"dragonroll/backend/services"
	"dragonroll/backend/utils"
	"errors"
	"net/http"
	"slices"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserClaims es el subset de datos que viaja firmado dentro del JWT.
// No metemos la contraseña ni datos sensibles — sólo lo que necesitan los
// middlewares de autorización para tomar decisiones rápidas.
type UserClaims struct {
	ID        primitive.ObjectID `json:"_id"`
	Username  string             `json:"username"`
	Roles     []models.Rol       `json:"roles"`
	ActiveRol models.Rol         `json:"active_rol"`
}

// LoadJWTAuth construye el middleware gin-jwt.
// Para detalle de cada flag, ver la doc de gin-jwt/v2; lo importante es:
//   - SendCookie: sí, así el frontend Nuxt puede dejar la cookie HttpOnly
//   - TokenLookup: cookie en debug, header en release
//   - Authenticator: nuestra lógica de validar credenciales
func LoadJWTAuth() *jwt.GinJWTMiddleware {
	tokenLookup := "header: Authorization, cookie: token"

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            "DragonRoll",
		Key:              []byte(config.Cfg.JWT.Key),
		SigningAlgorithm: config.Cfg.JWT.SigningAlg,
		Timeout:          24 * time.Hour,
		MaxRefresh:       24 * time.Hour,
		Authenticator:    AuthenticatorFunc,
		Authorizator:     AuthorizatorFunc,
		PayloadFunc:      PayloadFunc,
		Unauthorized:     UnauthorizedHandlerFunc,
		LoginResponse:    LoginResponse,
		LogoutResponse:   LogoutResponse,
		IdentityHandler:  IdentityHandlerFunc,
		TokenLookup:      tokenLookup,
		TokenHeadName:    "Bearer",
		TimeFunc:         time.Now,
		SendCookie:       true,
		SecureCookie:     config.Cfg.Server.Mode == "release",
		CookieHTTPOnly:   true,
		CookieName:       "token",
		CookieSameSite:   http.SameSiteLaxMode,
	})
	if err != nil {
		utils.Fatal("Error inicializando JWT middleware", err)
	}
	return authMiddleware
}

// AuthenticatorFunc valida { user, password } y retorna las claims.
func AuthenticatorFunc(c *gin.Context) (interface{}, error) {
	var login models.Login
	if err := c.ShouldBind(&login); err != nil {
		return nil, jwt.ErrMissingLoginValues
	}

	user, err := services.GetUserByLogin(login.User)
	if err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	if !utils.CompareHashToPassword(login.Password, user.Password) {
		return nil, errors.New("credenciales inválidas")
	}

	claims := UserClaims{
		ID:        user.ID,
		Username:  user.Username,
		Roles:     user.Roles,
		ActiveRol: user.ActiveRol,
	}
	utils.Info("Login exitoso: " + user.Username + " desde " + c.ClientIP())
	c.Set("user", user)
	return claims, nil
}

// PayloadFunc inyecta las claims en el JWT.
func PayloadFunc(data interface{}) jwt.MapClaims {
	if claims, ok := data.(UserClaims); ok {
		return jwt.MapClaims{
			"_id":        claims.ID.Hex(),
			"username":   claims.Username,
			"roles":      claims.Roles,
			"active_rol": claims.ActiveRol,
		}
	}
	return jwt.MapClaims{}
}

// IdentityHandlerFunc reconstituye claims desde el JWT en cada request
// y deja el user_id en el contexto de Gin para que los controladores
// puedan llamar c.Get("user_id") sin volver a parsear nada.
func IdentityHandlerFunc(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	if claims == nil {
		return nil
	}
	idHex, _ := claims["_id"].(string)
	id, err := primitive.ObjectIDFromHex(idHex)
	if err == nil {
		c.Set("user_id", id)
	}
	c.Set("username", claims["username"])
	c.Set("active_rol", claims["active_rol"])
	return claims
}

// SetRoles restringe una ruta a una lista de roles.
// Uso típico:  router.GET("/admin", LoadJWTAuth().MiddlewareFunc(), SetRoles(models.ADMIN), Handler)
func SetRoles(roles ...models.Rol) gin.HandlerFunc {
	if slices.Contains(roles, models.ALL) {
		roles = models.ALL_ROLES
	}
	return func(c *gin.Context) {
		c.Set("roles", roles)
		c.Next()
	}
}

// AuthorizatorFunc se invoca después de la validación del JWT.
// Permite el acceso si no se restringió por roles, o si el rol activo
// del usuario está en la lista permitida.
func AuthorizatorFunc(data interface{}, c *gin.Context) bool {
	claims, ok := data.(map[string]interface{})
	if !ok {
		return false
	}
	required, exists := c.Get("roles")
	if !exists {
		return true
	}
	activeRolStr, _ := claims["active_rol"].(string)
	activeRol := models.Rol(activeRolStr)
	for _, r := range required.([]models.Rol) {
		if activeRol == r {
			return true
		}
	}
	return false
}

func UnauthorizedHandlerFunc(c *gin.Context, code int, message string) {
	if message == "cookie token is empty" || message == "auth header is empty" {
		message = "No autorizado: token ausente o inválido."
	}
	utils.Info("Petición no autorizada: " + c.Request.URL.Path + " desde " + c.ClientIP())
	utils.JsonResponse(c, code, message, nil)
}

func LoginResponse(c *gin.Context, code int, token string, expire time.Time) {
	user, _ := c.Get("user")
	utils.JWTResponse(c, code, "Login exitoso.", token, expire, user)
}

func LogoutResponse(c *gin.Context, code int) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": "Logout realizado.",
		"status":  "success",
	})
}
