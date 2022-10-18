package rest

import (
	"github.com/gin-contrib/cors"
	"github.com/nabilsea/hacktiv8-final-project/database"
	"github.com/nabilsea/hacktiv8-final-project/docs"
	"github.com/nabilsea/hacktiv8-final-project/repository/comment_repository/comment_pg"
	"github.com/nabilsea/hacktiv8-final-project/repository/photo_repository/photo_pg"
	"github.com/nabilsea/hacktiv8-final-project/repository/social_media_repository/social_media_pg"
	"github.com/nabilsea/hacktiv8-final-project/repository/user_repository/user_pg"
	"github.com/nabilsea/hacktiv8-final-project/service"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func StartApp() {
	database.StartDB()
	db := database.GetDB()

	userRepo := user_pg.NewUserPG(db)
	userService := service.NewUserService(userRepo)
	userRestHandler := NewUserRestHandler(userService)

	photoRepo := photo_pg.NewPhotoPG(db)
	photoService := service.NewPhotoService(photoRepo)
	photoRestHandler := NewPhotoRestHandler(photoService)

	commentRepo := comment_pg.NewCommentPG(db)
	commentService := service.NewCommentService(commentRepo)
	commentRestHandler := NewCommentRestHandler(commentService)

	socialMediaRepo := social_media_pg.NewSocialMediaPG(db)
	socialMediaService := service.NewSocialMediaService(socialMediaRepo)
	socialMediaRestHandler := NewSocialMediaRestHandler(socialMediaService)

	authService := service.NewAuthService(userRepo, photoRepo, commentRepo, socialMediaRepo)

	// ! Routing
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	v1 := router.Group("/api/v1")

	v1.POST("/login", userRestHandler.Login)
	v1.POST("/register", userRestHandler.Register)
	userRoute := v1.Group("/users")
	userRoute.PUT("/", authService.Authentication(), userRestHandler.UpdateUserData)
	userRoute.DELETE("/", authService.Authentication(), userRestHandler.DeleteUser)

	photoRoute := v1.Group("/photos")
	photoRoute.Use(authService.Authentication())
	photoRoute.POST("/", photoRestHandler.PostPhoto)
	photoRoute.GET("/", photoRestHandler.GetAllPhotos)
	photoRoute.PUT("/:photoID", authService.PhotoAuthorization(), photoRestHandler.UpdatePhoto)
	photoRoute.DELETE("/:photoID", authService.PhotoAuthorization(), photoRestHandler.DeletePhoto)

	commentRoute := v1.Group("/comments")
	commentRoute.Use(authService.Authentication())
	commentRoute.POST("/", commentRestHandler.PostComment)
	commentRoute.GET("/", commentRestHandler.GetAllComments)
	commentRoute.PUT("/:commentID", authService.CommentAuthorization(), commentRestHandler.UpdateComment)
	commentRoute.DELETE("/:commentID", authService.CommentAuthorization(), commentRestHandler.DeleteComment)

	socialMediaRoute := v1.Group("/social-medias")
	socialMediaRoute.Use(authService.Authentication())
	socialMediaRoute.POST("/", socialMediaRestHandler.AddSocialMedia)
	socialMediaRoute.GET("/", socialMediaRestHandler.GetAllSocialMedias)
	socialMediaRoute.PUT("/:socialMediaID", authService.SocialMediaAuthorization(), socialMediaRestHandler.EditSocialMediaData)
	socialMediaRoute.DELETE("/:socialMediaID", authService.SocialMediaAuthorization(), socialMediaRestHandler.DeleteSocialMedia)

	docs.SwaggerInfo.Title = "MyGrams API"
	docs.SwaggerInfo.Description = "MyGrams project is a simple REST API for social media application"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "http://localhost:8000"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	router.Run(":8000")
}
