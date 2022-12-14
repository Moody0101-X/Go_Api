package main

import (   	
    "fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"log"
	"github.com/Moody0101-X/Go_Api/socketOperations"
	"github.com/Moody0101-X/Go_Api/routes"
	"github.com/Moody0101-X/Go_Api/database"
	"github.com/Moody0101-X/Go_Api/models"
)

func RequestCancelRecover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			
			if err := recover(); err != nil {
				fmt.Println("The Request was cancelled because an unexpected error interupted.")
				fmt.Println("err:\n")
				log.Fatal(err);
				
				c.Request.Context().Done()
			}

		}()
		
		c.Next()
	}	
}

func run() {
	// GET THE PORT :)
	var PORT string = models.GetEnv("PORT")
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(cors.Default())	
	router.Use(gin.Logger(), RequestCancelRecover())
	// HTML/JS/CSS/IMG loaders
	
	router.Static("/static", "./public/static")
	router.Static("/img", "./public/img")
	router.Static("/Global.css", "./public/Global.css")
	
	router.LoadHTMLGlob("public/*.html")

	// POST routes.
	router.POST("/v2/login", routes.Login) // login and get a token for the updating/creation/deletion of personal data.
	router.POST("/v2/update", routes.Update) // Updating user's information by token
	router.POST("/v2/NewPost", routes.NewPost) // adding a post by token.
	router.POST("/v2/DeletePost", routes.DeletePost) // Deleting a post by token
	router.POST("/v2/signup", routes.SignUp) // Making new account
	router.POST("/v2/comment", routes.AddCommentRoute) // For likes
	router.POST("/v2/like", routes.AddLikeRoute) // For comments
	router.POST("/v2/like/remove", routes.RemoveLikeRoute)
	router.POST("/v2/follow", routes.FollowRoute)
	router.POST("/v2/unfollow", routes.UnfollowRoute)
	
	// Get routes.
	router.GET("/v2/getUserPosts", routes.GetUserPostsRoute) // gettting user post by id
	router.GET("/v2/GetAllPosts", routes.GetAllPostsRoute) // getting all posts
	router.GET("/v2/query", routes.GetUsersRoute) // user look up by name
	router.GET("/v2/getUser", routes.GetUserByIdRoute) // get user by id
	router.GET("/v2/getFollowers/:uuid", routes.GetUserFollowersById)
	router.GET("/v2/getFollowings/:uuid", routes.GetUserFollowingsById)
	router.GET("/v2/getComments/:pid", routes.GetPostComments)
	router.GET("/v2/getLikes/:pid", routes.GetPostLikes)
	router.GET("/v2/getUserNotifications/:uuid", routes.GetAllNotificationsRoute)
	// router.Static("/", "./public")
	router.GET("/", routes.Index)
    router.NoRoute(routes.Index)

    // Socket routes.
    router.GET("/v2/NotificationSock", socketOperations.NotificationServer)
	
	// running the server.
	fmt.Println("Serving in port ", PORT)
	router.Run(PORT)
}

func main() {
	err, path := database.InitializeDb();

	if err != nil {
        fmt.Println("Error opening the database! ", err.Error())
        return
    }

    fmt.Println("connected to db: ", path)

	run()
}

