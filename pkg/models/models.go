package models

import (
	"github.com/dgrijalva/jwt-go"
)

type RegisterData struct {
	//	ID       string `json:"id" binding:"required"`          // Уникальный идентификатор пользователя, обязательно
	Name     string `json:"name"`                           // Имя пользователя
	Email    string `json:"email" binding:"required,email"` // Email пользователя, обязательно и должен быть валидным email
	Password string `json:"password" binding:"required"`    // Пароль, обязательно
	Role     string `json:"role"`
	JWT      string `json:"jwtHandlers"`
}

type AuthData struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RolesData структура для хранения ролей и их привилегий
type RolesData struct {
	Name       string   `json:"name" binding:"required"`
	Privileges []string `json:"privileges" binding:"required"`
}

// UserRoleData структура для хранения информации о ролях пользователей
type UserRoleData struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required"`
}
type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}
type DeleteUserData struct {
	UserID string `json:"user_id"`
}
type DeleteRoleData struct {
	Name string `json:"name"`
}

type PrivilegeData struct {
	// -----Admin Privileges-----
	DeleteUser      string
	RedactUser      string
	GetUserData     string
	CreateRole      string
	DeleteRole      string
	GetRoleData     string
	AccessSystem    string
	AccessLogs      string
	AccessAnalytics string

	// -----Project Manager Privileges-----
	CreateProject string
	ManageProject string
	AssignTasks   string
	ViewReports   string
	ViewUsers     string

	// -----Developer Privileges-----
	AccessProjects string
	CommitCode     string
	ManageTasks    string
	ViewTeamData   string

	// -----QA Tester Privileges-----
	AccessTestProjects string
	ReportBugs         string
	EditBugs           string
	ViewTestTasks      string

	// -----Support Privileges-----
	HandleTickets string
	ViewUserInfo  string

	// -----User Privileges-----
	AccessResources    string
	InteractWithSystem string
}

var Privileges = PrivilegeData{
	// -----Admin Privileges-----
	DeleteUser:      "deleteUser",
	RedactUser:      "redactUser",
	GetUserData:     "getUserData",
	CreateRole:      "createRole",
	DeleteRole:      "deleteRole",
	GetRoleData:     "getRoleData",
	AccessSystem:    "accessSystem",
	AccessLogs:      "accessLogs",
	AccessAnalytics: "accessAnalytics",

	// -----Project Manager Privileges-----
	CreateProject: "createProject",
	ManageProject: "manageProject",
	AssignTasks:   "assignTasks",
	ViewReports:   "viewReports",
	ViewUsers:     "viewUsers",

	// -----Developer Privileges-----
	AccessProjects: "accessProjects",
	CommitCode:     "commitCode",
	ManageTasks:    "manageTasks",
	ViewTeamData:   "viewTeamData",

	// -----QA Tester Privileges-----
	AccessTestProjects: "accessTestProjects",
	ReportBugs:         "reportBugs",
	EditBugs:           "editBugs",
	ViewTestTasks:      "viewTestTasks",

	// -----Support Privileges-----
	HandleTickets: "handleTickets",
	ViewUserInfo:  "viewUserInfo",

	// -----User Privileges-----
	AccessResources:    "accessResources",
	InteractWithSystem: "interactWithSystem",
}
