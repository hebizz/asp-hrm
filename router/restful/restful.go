package restful

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.jiangxingai.com/asp-hrm/router/restful/apiV1"
)

func InitRestFul(r *gin.Engine) *gin.Engine {
	r.Use(Cors())
	router := r.Group("/hrm/api/v1")
	{
		router.GET("/department/query", apiV1.GetDepartmentHandler)
		router.GET("/departmentList/query", apiV1.GetDepartmentListHandler)
		router.POST("/department/create", apiV1.CreateDepartmentHandler)
		router.POST("department/update", apiV1.UpdateDepartmentHandler)
		router.POST("/department/remove", apiV1.DeleteDepartmentHandler)
		router.GET("/human/query", apiV1.GetHumanHandler)
		router.POST("/human/register", apiV1.RegisterHumanHandler)
		router.POST("/human/remove", apiV1.RemoveHumanHandler)
		router.POST("/human/update", apiV1.UpdateHumanHandler)
		router.GET("/identify/log", apiV1.GetIdentifyLogHandler)
	}
	return r
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
