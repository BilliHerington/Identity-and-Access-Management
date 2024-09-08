package models

type RegisterData struct {
	//	ID       string `json:"id" binding:"required"`          // Уникальный идентификатор пользователя, обязательно
	Name     string `json:"name"`                           // Имя пользователя
	Email    string `json:"email" binding:"required,email"` // Email пользователя, обязательно и должен быть валидным email
	Password string `json:"password" binding:"required"`    // Пароль, обязательно
	Role     string `json:"role"`
	JWT      string `json:"jwt"`
}

type AuthData struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
type UserData struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
	JWT   string `json:"jwt"`
}

//var AllRolesList []RoleStruct

// RoleStruct структура для хранения ролей и их привилегий
type RolesData struct {
	Name       string   `json:"name" binding:"required"`
	Privileges []string `json:"privileges" binding:"required"`
}

//var AllUserRolesList []UserRoleStruct

// UserRoleStruct структура для хранения информации о ролях пользователей
//type UserRoleStruct struct {
//	UserID string `json:"user_id" binding:"required"`
//	Role   string `json:"role" binding:"required"`
//}
