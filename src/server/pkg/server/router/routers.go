package router

import (
	"speak-sphere/pkg/server/controller"

	"github.com/gin-gonic/gin"
)

func RegisterVocabularySetRouter(route *gin.Engine) {
	r := route.Group("/vocabulary_set")

	r.POST("/add", controller.AddVocabularySet)

	r.DELETE("/delete/:id", controller.DeleteVocabularySetByID)

	r.PUT("/update", controller.UpdateVocabularySet)

	r.GET("/list/:page_size/:current_page", controller.ListVocabularySet)
	r.GET("/:id", controller.FindVocabularySetByID)
	r.GET("/category/:category/:page_size/:current_page", controller.FindVocabularySetByCategory)
	r.GET("/count", controller.CountVocabularySet)
}

func RegisterDictionaryRouter(route *gin.Engine) {
	r := route.Group("/dictionary")

	r.POST("/add", controller.AddDictionary)

	r.DELETE("/delete/:id", controller.DeleteDictionaryByID)

	r.PUT("/update", controller.UpdateDictionary)

	r.GET("/query", controller.FindDictionaryByID)
	r.GET("/list", controller.ListDictionary)
}

func RegisterVocabularyRouter(route *gin.Engine) {
	r := route.Group("/vocabulary")
	r.POST("/add", controller.AddVocabulary)

	r.DELETE("/delete/:id", controller.DeleteVocabularyByID)

	r.PUT("/update", controller.UpdateVocabulary)
	r.PUT("/update/unwanted/:id", controller.SetVocabularyUnwanted)

	r.PUT("/update/study/count/:id", controller.UpdateVocabularyStudyCount)
	r.PUT("/update/reset/:id", controller.ResetStudyCountToZero)

	r.GET("/query", controller.FindVocabularyByID)
	r.GET("/count/:vocabulary_set_id", controller.CountVocabulary)
	r.GET("/list", controller.ListVocabulary)
	r.GET("/learn", controller.GetVocabulariesToLearn)
	r.GET("/review", controller.GetVocabulariesToReview)
	r.GET("/check", controller.CheckVocabularyInVocabularySet)
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
