package main;

import (
  "fmt"
  "simple-board/model"
  "github.com/gofiber/fiber/v2"
  "github.com/go-playground/validator/v10"
)

type UserGetSuccess struct {
  Username string
  Posts    []model.Post
}

type UserPostBody struct {
  Username string    `json:"username"`
  Password string    `json:"password"`
}

type UserPutBody struct {
  Password string    `json:"password" validate:"required,min=3,max=20"`
}

func UserGet(c *fiber.Ctx) error {
  username := c.Params("user")
  user := new(model.User)

  result := DB.Where("username = ?", username).First(&user)

  if result.Error != nil {
    return result.Error
  }

  return c.JSON(UserGetSuccess{
    Username: user.Username,
    Posts: user.Posts,
  })
}

func UserPost(c *fiber.Ctx) error {
  body := new(UserPostBody)

  if err := c.BodyParser(body); err != nil {
    return err
  }

  user := model.User{
    Username: body.Username,
    Password: body.Password,
  }

  result := DB.Create(&user)

  if result.Error != nil {
    return result.Error
  }

  return c.SendString(fmt.Sprintf("Create a user (Username: %s)", user.Username))
}

func UserPut(c *fiber.Ctx) error {
  username := c.Params("user")

  if sess, err := SessionStore.Get(c); err != nil {
    return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
  } else if sessUsername := sess.Get("username"); sessUsername != username {
    return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
  }

  body := new(UserPutBody)
  validate := validator.New()

  if err := c.BodyParser(&body); err != nil {
    return c.Status(fiber.StatusBadRequest).SendString(err.Error())
  }

  if err := validate.Struct(body); err != nil {
    return c.Status(fiber.StatusBadRequest).SendString(err.Error())
  }

  user := new(model.User)

  if result := DB.Where("username = ?", username).First(&user); result.Error != nil {
    return result.Error
  }

  user.Password = body.Password

  if result := DB.Save(&user); result.Error != nil {
    return result.Error
  }

  return c.SendString(fmt.Sprintf("Updated user data (Username: %s)", user.Username))
}

func UserDelete(c *fiber.Ctx) error {
  username := c.Params("user")

  if sess, err := SessionStore.Get(c); err != nil {
    return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
  } else if sessUsername := sess.Get("username"); sessUsername != username {
    return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
  } else {
    user := new(model.User)

    if result := DB.Where("username = ?", username).First(&user); result.Error != nil {
      return result.Error
    }

    if result := DB.Delete(&user); result.Error != nil {
      return result.Error
    }

    sess.Destroy()
    sess.Save()

    return c.SendString(fmt.Sprintf("Deleted a user (username: %s)", user.Username))
  }
}
