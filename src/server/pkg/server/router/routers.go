package router

import (
	"speak-sphere/pkg/server/controller"

	"github.com/gin-gonic/gin"
)

func RegisterBookRouter(route *gin.Engine) {
	r := route.Group("/book")

	r.POST("/add", controller.AddBook)

	r.DELETE("/delete/:id", controller.DeleteBookByID)

	r.PUT("/update", controller.UpdateBook)

	r.GET("/list/:page_size/:current_page", controller.ListBook)
	r.GET("/:id", controller.FindBookByID)
	r.GET("/category/:category/:page_size/:current_page", controller.FindBookByCategory)
	r.GET("/count", controller.CountBook)
}

func RegisterDictionaryRouter(route *gin.Engine) {
	r := route.Group("/dictionary")

	r.POST("/add", controller.AddDictionary)

	r.DELETE("/delete/:id", controller.DeleteDictionaryByID)

	r.PUT("/update", controller.UpdateDictionary)

	r.GET("/query", controller.FindDictionaryByID)
	r.GET("/list", controller.ListDictionary)
}

func RegisterEntryRouter(route *gin.Engine) {
	r := route.Group("/entry")
	r.POST("/add", controller.AddEntry)

	r.DELETE("/delete/:id", controller.DeleteEntryByID)

	r.PUT("/update", controller.UpdateEntry)
	r.PUT("/update/unwanted/:id", controller.SetEntryUnwanted)

	r.PUT("/update/study/count/:id", controller.UpdateEntryStudyCount)
	r.PUT("/update/reset/:id", controller.ResetStudyCountToZero)

	r.GET("/query", controller.FindEntryByID)
	r.GET("/count/:book_id", controller.CountEntry)
	r.GET("/list", controller.ListEntry)
	r.GET("/learn", controller.GetEntriesToLearn)
	r.GET("/review", controller.GetEntriesToReview)
	r.GET("/check", controller.CheckWordInBook)
}

func RegisterUserRouter(route *gin.Engine) {
	r := route.Group("/user")

	r.POST("/register", controller.Register)
	r.POST("/login", controller.Login)

	r.DELETE("/delete/:id", controller.DeleteUserByID)

	r.PUT("/update", controller.UpdateUser)

	r.GET("/profile", controller.GetUserProfile)
}

func RegisterLLMServiceRouter(route *gin.Engine) {
	r := route.Group("/llm")

	r.POST("/add", controller.AddLLMService)

	r.DELETE("/delete/:id", controller.DeleteLLMService)

	r.PUT("/update/:id", controller.UpdateLLMService)

	r.GET("/list", controller.GetLLMServices)
	r.GET("/default", controller.GetDefaultLLMService)
	r.GET("/:id", controller.GetLLMServiceByID)
}

func RegisterAIPromptRouter(route *gin.Engine) {
	r := route.Group("/aiprompt")

	r.POST("/add", controller.AddAIPrompt)

	r.DELETE("/delete/:id", controller.DeleteAIPrompt)

	r.PUT("/update/:id", controller.UpdateAIPrompt)

	r.GET("/list", controller.GetAIPrompts)
	r.GET("/default", controller.GetDefaultAIPrompt)
	r.GET("/:id", controller.GetAIPromptByID)
}
