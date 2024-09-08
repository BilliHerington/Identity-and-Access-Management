package middlewares

//
//func CheckPrivileges(requiredPrivilege string) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		userID := c.GetString("user_id") // Получаем ID пользователя из контекста
//
//		//ищем роль пользователя по ID
//		var userRole string
//		for _, ur := range models.AllUserRolesList {
//			if ur.UserID == userID {
//				userRole = ur.Role
//			}
//		}
//		//ишем привилегии для этой роли
//		var privileges []string
//		for _, role := range models.AllRolesList {
//			if role.Name == userRole {
//				privileges = role.Privileges
//				break
//			}
//		}
//		//проверяем имеет ли пользователь необходимую привилегию
//		hasPrivilege := false
//		for _, privilege := range privileges {
//			if privilege == requiredPrivilege {
//				hasPrivilege = true
//				break
//			}
//		}
//		//если у пользователя нет привилегий
//		if !hasPrivilege {
//			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have the required privileges"})
//			c.Abort()
//			return
//		}
//		c.Next()
//	}
//}
