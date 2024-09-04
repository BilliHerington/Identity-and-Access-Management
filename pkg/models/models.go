package models

type RegisterData struct {
	ID       string `json:"id" binding:"required"`          // Уникальный идентификатор пользователя, обязательно
	Name     string `json:"name" binding:"required"`        // Имя пользователя, обязательно
	Email    string `json:"email" binding:"required,email"` // Email пользователя, обязательно и должен быть валидным email
	Password string `json:"password" binding:"required"`    // Пароль, обязательно
}
type AuthData struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
