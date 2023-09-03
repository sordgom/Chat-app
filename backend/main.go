package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/redis/go-redis/v9"

	"github.com/sordgom/jwt-go/chat"
	"github.com/sordgom/jwt-go/controllers"
	"github.com/sordgom/jwt-go/initializers"
	"github.com/sordgom/jwt-go/middleware"
	"github.com/sordgom/jwt-go/redisrepo"
	"github.com/sordgom/jwt-go/videochat"
	"github.com/sordgom/jwt-go/voicechat"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("Failed to load environment variables! \n", err.Error())
	}
	initializers.ConnectDB(&config)
	initializers.ConnectRedis(&config)
}

func main() {
	server := flag.String("server", "", "http,websocket")
	flag.Parse()
	fmt.Println("Distributed Chat App v0.01")
	if *server == "text" {
		fmt.Println("Chat server is starting on :8082")
		chat.StartWebsocketServer()
	} else if *server == "http" {
		fmt.Println("http server is starting on :8080")
		startLoginServer()
	} else if *server == "video" {
		fmt.Println("Video Chat server is starting on :8081")
		videochat.SetupVideoChat()
	} else if *server == "voice" {
		fmt.Println("Voice Chat server is starting on :8083")
		voicechat.SetupVoiceChat()
	} else {
		fmt.Println("invalid server. Available server: text or video")
	}
}

func startLoginServer() {
	//Create the index for fetching chat between two users
	redisrepo.CreateFetchChatBetweenIndex()

	app := fiber.New()
	api := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST",
		AllowCredentials: true,
	}))
	app.Mount("/api", api)
	app.Use(logger.New())

	api.Route("/auth", func(auth fiber.Router) {
		auth.Post("/register", controllers.SignUpUser)
		auth.Post("/login", controllers.SignInUser)
		auth.Get("/logout", middleware.DeserializeUser, controllers.LogoutUser)
		auth.Get("/verify", controllers.CheckUserExists)
		auth.Get("/refresh", controllers.RefreshAccessToken)

		auth.Post("/add-contact", controllers.UpdateContactList)
		auth.Get("/contact-list", middleware.DeserializeUser, controllers.ContactList)
		auth.Get("/chat-history", middleware.DeserializeUser, controllers.ChatHistory)
	})

	api.Get("/users/me", middleware.DeserializeUser, controllers.GetUser)
	api.Post("/token", controllers.CreateAccessToken)

	ctx := context.TODO()
	value, err := initializers.RedisClient.Get(ctx, "test").Result()

	if err == redis.Nil {
		fmt.Println("key: test does not exist")
	} else if err != nil {
		panic(err)
	}

	api.Get("/healthchecker", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "all good",
			"message": value,
		})
	})

	api.All("*", func(c *fiber.Ctx) error {
		path := c.Path()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": fmt.Sprintf("Path: %v does not exists on this server", path),
		})
	})

	log.Fatal(app.Listen(":8080"))
}
