package services

import (
	"dragonroll/backend/models"
	"dragonroll/backend/utils"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateUserService(user *models.User) (*models.User, error) {
	user.Username = strings.TrimSpace(strings.ToLower(user.Username))
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))

	if err := utils.ValidateUsername(user.Username); err != nil {
		return nil, err
	}
	if err := utils.ValidateEmail(user.Email); err != nil {
		return nil, err
	}
	if user.Password == "" {
		return nil, errors.New("password obligatorio")
	}

	if existing, _ := userRepo.FindOne(bson.M{"username": user.Username}); existing != nil {
		return nil, errors.New("username ya existe")
	}
	if existing, _ := userRepo.FindOne(bson.M{"email": user.Email}); existing != nil {
		return nil, errors.New("email ya está registrado")
	}

	if len(user.Roles) == 0 {
		user.Roles = []models.Rol{models.PLAYER}
	}
	if user.ActiveRol == "" {
		user.ActiveRol = user.Roles[0]
	}

	user.Password = utils.GeneratePassword(user.Password)
	user.CreatedAt = time.Now()

	id, err := userRepo.InsertOne(user)
	if err != nil {
		return nil, err
	}
	user.ID = id
	return user, nil
}

func GetUserByLogin(login string) (*models.User, error) {
	login = strings.TrimSpace(strings.ToLower(login))
	if login == "" {
		return nil, errors.New("login vacío")
	}
	user, err := userRepo.FindOne(bson.M{
		"$or": []bson.M{
			{"username": login},
			{"email": login},
		},
	})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("usuario no encontrado")
	}
	return user, nil
}

func GetUserByID(id primitive.ObjectID) (*models.User, error) {
	user, err := userRepo.FindOne(bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("usuario no encontrado")
	}
	return user, nil
}

func GetAllUsersService() ([]*models.User, error) {
	return userRepo.GetAll()
}

// SetActiveRoleService cambia el rol activo (player ↔ dm) si lo tiene.
func SetActiveRoleService(userID primitive.ObjectID, newRole models.Rol) error {
	user, err := GetUserByID(userID)
	if err != nil {
		return err
	}
	allowed := false
	for _, r := range user.Roles {
		if r == newRole {
			allowed = true
			break
		}
	}
	if !allowed {
		return errors.New("rol no habilitado para este usuario")
	}
	return userRepo.PatchOne(userID, bson.M{"active_rol": newRole})
}

// UpdateProfileService permite que el usuario edite su username, email y avatar.
// Si cambia username/email, valida unicidad. Password se cambia con
// ChangePasswordService porque requiere la password actual.
func UpdateProfileService(userID primitive.ObjectID, username, email, avatar string) (*models.User, error) {
	current, err := GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	username = strings.TrimSpace(strings.ToLower(username))
	email = strings.TrimSpace(strings.ToLower(email))

	updates := bson.M{}
	if username != "" && username != current.Username {
		if err := utils.ValidateUsername(username); err != nil {
			return nil, err
		}
		// Unicidad — cualquier otro user con ese username falla.
		if existing, _ := userRepo.FindOne(bson.M{"username": username, "_id": bson.M{"$ne": userID}}); existing != nil {
			return nil, errors.New("ese username ya está tomado")
		}
		updates["username"] = username
	}
	if email != "" && email != current.Email {
		if err := utils.ValidateEmail(email); err != nil {
			return nil, err
		}
		if existing, _ := userRepo.FindOne(bson.M{"email": email, "_id": bson.M{"$ne": userID}}); existing != nil {
			return nil, errors.New("ese email ya está registrado")
		}
		updates["email"] = email
	}
	if avatar != current.Avatar {
		updates["avatar"] = avatar
	}

	if len(updates) == 0 {
		return current, nil
	}
	if err := userRepo.PatchOne(userID, updates); err != nil {
		return nil, err
	}
	return GetUserByID(userID)
}

// ChangePasswordService cambia la contraseña previa validación de la actual.
// Esto es importante: nunca permitas cambiar password sin reautenticación,
// porque si se filtra un JWT podrías secuestrar la cuenta para siempre.
func ChangePasswordService(userID primitive.ObjectID, currentPassword, newPassword string) error {
	if len(newPassword) < 8 {
		return errors.New("la contraseña nueva debe tener al menos 8 caracteres")
	}
	user, err := GetUserByID(userID)
	if err != nil {
		return err
	}
	if !utils.CompareHashToPassword(currentPassword, user.Password) {
		return errors.New("la contraseña actual no es correcta")
	}
	hashed := utils.GeneratePassword(newPassword)
	return userRepo.PatchOne(userID, bson.M{"password": hashed})
}
