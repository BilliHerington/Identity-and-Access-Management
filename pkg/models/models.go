package models

import (
	"github.com/dgrijalva/jwt-go"
)

type RegisterData struct {
	Name             string `json:"name"`
	Email            string `json:"email" binding:"required,email"`
	VerificationCode string `json:"verificationCode"`
	Password         string `json:"password" binding:"required"`
	Role             string `json:"role"`
	JWT              string `json:"jwt"`
}
type EmailData struct {
	Email string `json:"email" binding:"required,email"`
}
type AuthData struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RolesData struct for data about Roles and Privileges
type RolesData struct {
	RoleName   string   `json:"roleName" binding:"required"`
	Privileges []string `json:"privileges" binding:"required"`
}

// UserRoleData struct for data about Users and Roles
type UserRoleData struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required"`
}
type Claims struct {
	Email       string `json:"email" binding:"required,email"`
	UserID      string `json:"userID" binding:"required"`
	UserVersion string `json:"userVersion" binding:"required"`
	jwt.StandardClaims
}
type ResetPass struct {
	Email         string `json:"email" binding:"required"`
	ResetPassCode string `json:"resetPassCode" binding:"required"`
	NewPassword   string `json:"newPassword" binding:"required"`
}
type AdminPrivilegeData struct {
	// -----Admin Privileges-----
	DeleteUser  string
	RedactUser  string
	GetUserData string
	//GetUserList     string
	CreateRole      string
	DeleteRole      string
	GetRoleData     string
	AccessSystem    string
	AccessLogs      string
	AccessAnalytics string

	//// -----Project Manager Privileges-----
	//CreateProject string
	//ManageProject string
	//AssignTasks   string
	//ViewReports   string
	//ViewUsers     string
	//
	//// -----Developer Privileges-----
	//AccessProjects string
	//CommitCode     string
	//ManageTasks    string
	//ViewTeamData   string
	//
	//// -----QA Tester Privileges-----
	//AccessTestProjects string
	//ReportBugs         string
	//EditBugs           string
	//ViewTestTasks      string
	//
	//// -----Support Privileges-----
	//HandleTickets string
	//ViewUserInfo  string
	//
	//// -----User Privileges-----
	//AccessResources    string
	//InteractWithSystem string
}

var AdminPrivileges = AdminPrivilegeData{
	// -----Admin Privileges-----
	DeleteUser:  "deleteUser",
	RedactUser:  "redactUser",
	GetUserData: "getUserData",
	//GetUserList:     "getUserList",
	CreateRole:      "createRole",
	DeleteRole:      "deleteRole",
	GetRoleData:     "getRoleData",
	AccessSystem:    "accessSystem",
	AccessLogs:      "accessLogs",
	AccessAnalytics: "accessAnalytics",
	//
	//// -----Project Manager Privileges-----
	//CreateProject: "createProject",
	//ManageProject: "manageProject",
	//AssignTasks:   "assignTasks",
	//ViewReports:   "viewReports",
	//ViewUsers:     "viewUsers",
	//
	//// -----Developer Privileges-----
	//AccessProjects: "accessProjects",
	//CommitCode:     "commitCode",
	//ManageTasks:    "manageTasks",
	//ViewTeamData:   "viewTeamData",
	//
	//// -----QA Tester Privileges-----
	//AccessTestProjects: "accessTestProjects",
	//ReportBugs:         "reportBugs",
	//EditBugs:           "editBugs",
	//ViewTestTasks:      "viewTestTasks",
	//
	//// -----Support Privileges-----
	//HandleTickets: "handleTickets",
	//ViewUserInfo:  "viewUserInfo",
	//
	//// -----User Privileges-----
	//AccessResources:    "accessResources",
	//InteractWithSystem: "interactWithSystem",
}

type UserPrivilegeData struct {
	AccessResources    string
	InteractWithSystem string
}

var UserPrivileges = UserPrivilegeData{
	AccessResources:    "accessResources",
	InteractWithSystem: "interactWithSystem",
}
