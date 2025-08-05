package routes

import (
	"IAM/pkg/ginHandlers/authenticationHandlers"
	"IAM/pkg/ginHandlers/authenticationHandlers/googleAuth"
	"IAM/pkg/ginHandlers/requestLimiterHandler"
	"github.com/gin-gonic/gin"
	"os"
)

func RegisterPublicRoutes(router *gin.Engine) {
	public := router.Group("/")
	TestModeWithoutGoogle := os.Getenv("USE_TEST_MODE_WITHOUT_GOOGLE")

	{
		public.Use(requestLimiterHandler.RequestLimiter(5, 30))

		if TestModeWithoutGoogle != "true" {
			public.GET("/oauth", googleAuth.OauthRedirect())
			public.GET("auth/callback", googleAuth.GoogleLogin())
		}

		public.POST("/registration", authenticationHandlers.StartRegistration())
		public.POST("/approve-registration", authenticationHandlers.ApproveRegistration())

		public.POST("/forgetPassword", authenticationHandlers.StartResetPassword())
		public.POST("/updatePassword", authenticationHandlers.ApproveResetPassword())

		public.POST("/auth", authenticationHandlers.Authenticate())
	}
}
