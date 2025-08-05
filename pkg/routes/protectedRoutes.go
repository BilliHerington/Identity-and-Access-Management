package routes

import (
	"IAM/pkg/ginHandlers/requestLimiterHandler"
	"IAM/pkg/ginHandlers/rolesHandlers"
	"IAM/pkg/ginHandlers/usersHandlers"
	"IAM/pkg/middlewares"
	"IAM/pkg/models"
	"github.com/gin-gonic/gin"
)

func RegisterProtectedRoutes(router *gin.Engine) {
	protected := router.Group("/")
	{
		protected.Use(middlewares.AuthMiddleware(), requestLimiterHandler.RequestLimiter(15, 30))
		//---Users----
		protected.GET("/get-users", usersHandlers.GetUsersList())
		protected.GET("/get-all-users-data", middlewares.CheckPrivileges(models.AdminPrivileges.GetUserData), usersHandlers.GetAllUsersData())
		protected.DELETE("/delete-user", middlewares.CheckPrivileges(models.AdminPrivileges.DeleteUser), usersHandlers.DeleteUser())
		//---Roles----
		protected.GET("/get-roles", rolesHandlers.GetRolesList())
		protected.GET("/get-all-roles-data", rolesHandlers.GetAllRolesData())
		protected.POST("/assign-role", middlewares.CheckPrivileges(models.AdminPrivileges.CreateRole), rolesHandlers.AssignRole())
		protected.POST("/create-role", middlewares.CheckPrivileges(models.AdminPrivileges.CreateRole), rolesHandlers.CreateRole())
		protected.DELETE("/delete-role", middlewares.CheckPrivileges(models.AdminPrivileges.DeleteRole), rolesHandlers.DeleteRole())
		protected.POST("/redact-role", middlewares.CheckPrivileges(models.AdminPrivileges.CreateRole), rolesHandlers.RedactRole())
	}
}
