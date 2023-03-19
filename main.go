package main

import (
  "gorm.io/gorm"
  "simple-board/model"
  "gorm.io/driver/mysql"
  "github.com/gofiber/fiber/v2"
  "github.com/gofiber/fiber/v2/middleware/logger"
  "github.com/gofiber/fiber/v2/middleware/session"
)

var DB *gorm.DB
var SessionStore *session.Store

func main() {
  // Create Fiber App
  app := fiber.New()

  // App Middleware
  app.Use(logger.New())

  // Create Session Store
  SessionStore = session.New()

  // GORM Open
  dialector := mysql.Open("root:toor@tcp(localhost:3306)/simple_board")

  if db, err := gorm.Open(dialector, &gorm.Config{}); err != nil {
    panic(err)
  } else {
    DB = db
  }

  // Auto Migration
  DB.AutoMigrate(
    &model.User{},
    &model.Post{},
  )

  // Routing (/)
  app.Get("/", func(c *fiber.Ctx) error {  // Hello, World!
    return c.SendString("Hello, World!")
  })

  // Routing (/auth)
  auth := app.Group("/auth")
  auth.Get("/", AuthGet)          // Auth Check
  auth.Post("/", AuthPost)        // Login
  auth.Delete("/", AuthDelete)    // Logout

  // Routing (/user)
  user := app.Group("/user")
  user.Get("/:user", UserGet)        // Get posts by a user
  user.Post("/", UserPost)           // Create a user
  user.Put("/:user", UserPut)        // Update user password
  user.Delete("/:user", UserDelete)  // Delete a user

  // Routing (/board)
  board := app.Group("/board")
  board.Get("/post/:id", BoardGet)        // Read a Post
  board.Post("/post", BoardPost)          // Create a Post
  board.Put("/post/:id", BoardPut)        // Update a Post
  board.Delete("/post/:id", BoardDelete)  // Delete a Post

  // Server Start!
  app.Listen(":3000")
}
