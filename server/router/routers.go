package router

import (
    "encoding/base64"
    "github.com/0xJacky/Nginx-UI/server/api"
    "github.com/0xJacky/Nginx-UI/server/model"
    "github.com/gin-gonic/gin"
    "github.com/gin-gonic/gin/binding"
    "github.com/go-playground/locales/en"
    "github.com/go-playground/locales/zh"
    ut "github.com/go-playground/universal-translator"
    "github.com/go-playground/validator/v10"
    en_translations "github.com/go-playground/validator/v10/translations/en"
    zh_translations "github.com/go-playground/validator/v10/translations/zh"
    "net/http"
)

func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			tmp, _ := base64.StdEncoding.DecodeString(c.Query("token"))
			token = string(tmp)
			if token == "" {
				c.JSON(http.StatusForbidden, gin.H{
					"message": "auth fail",
				})
				c.Abort()
				return
			}
		}

		n := model.CheckToken(token)

		if n < 1 {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "auth fail",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func Translations() gin.HandlerFunc {
	return func(c *gin.Context) {
		uni := ut.New(en.New(), zh.New())
		locale := c.GetHeader("locale")
		trans, _ := uni.GetTranslator(locale)
		v, ok := binding.Validator.Engine().(*validator.Validate)
		if ok {
			switch locale {
			case "zh":
				_ = zh_translations.RegisterDefaultTranslations(v, trans)
				break
			case "en":
				_ = en_translations.RegisterDefaultTranslations(v, trans)
				break
			default:
				_ = zh_translations.RegisterDefaultTranslations(v, trans)
				break
			}
			c.Set("trans", trans)
		}

		c.Next()
	}
}

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())

	r.Use(gin.Recovery())
	r.Use(Translations())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})

	r.GET("install", api.InstallLockCheck)
	r.POST("install", api.InstallNginxUI)

	r.POST("/login", api.Login)
	r.DELETE("/logout", api.Logout)

	g := r.Group("/", authRequired())
	{
		r.GET("/analytic", api.Analytic)

		g.GET("/users", api.GetUsers)
		g.GET("/user/:id", api.GetUser)
		g.POST("/user", api.AddUser)
		g.POST("/user/:id", api.EditUser)
		g.DELETE("/user/:id", api.DeleteUser)

		g.GET("domains", api.GetDomains)
		g.GET("domain/:name", api.GetDomain)
		g.POST("domain/:name", api.EditDomain)
		g.POST("domain/:name/enable", api.EnableDomain)
		g.POST("domain/:name/disable", api.DisableDomain)
		g.DELETE("domain/:name", api.DeleteDomain)

		g.GET("configs", api.GetConfigs)
		g.GET("config/:name", api.GetConfig)
		g.POST("config", api.AddConfig)
		g.POST("config/:name", api.EditConfig)

		g.GET("backups", api.GetFileBackupList)
		g.GET("backup/:id", api.GetFileBackup)

		g.GET("template/:name", api.GetTemplate)

		g.GET("cert/issue/:domain", api.IssueCert)
		g.GET("cert/:domain/info", api.CertInfo)
	}

	return r
}
