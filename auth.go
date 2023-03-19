package main;

import (
  "fmt"
  "simple-board/model"
  "github.com/gofiber/fiber/v2"
  "github.com/go-playground/validator/v10"
)

type AuthPostBody struct {
  Username string `json:"username" validate:"required,min=3"`
  Password string `json:"password" validate:"required,min=3"`
}

func AuthGet(c *fiber.Ctx) error {
  if sess, err := SessionStore.Get(c); err != nil {
    return err
  } else {
    username := sess.Get("username")

    if username == nil {
      return c.SendString("Unauthorized")
    } else {
      return c.SendString(fmt.Sprintf("Username: %s", username))
    }
  }
}

func AuthPost(c *fiber.Ctx) error {
  body := new(AuthPostBody)
  user := new(model.User)

  if err := c.BodyParser(body); err != nil {
    return c.Status(fiber.StatusBadRequest).SendString(err.Error())
  }

  validate := validator.New()
  if err := validate.Struct(body); err != nil {
    return c.Status(fiber.StatusBadRequest).SendString(err.Error())
  }

  if result := DB.Where("username = ?", body.Username).First(&user); result.Error != nil {
    message := fmt.Sprintf("Bad username or password (username: %s)", body.Username)
    return c.Status(fiber.StatusUnauthorized).SendString(message)
  }

  if body.Password != user.Password {
    message := fmt.Sprintf("Bad username or password (username: %s)", body.Username)
    return c.Status(fiber.StatusUnauthorized).SendString(message)
  }

  if sess, err := SessionStore.Get(c); err != nil {
    return err
  } else {
    sess.Set("id", user.ID)
    sess.Set("username", user.Username)

    if err := sess.Save(); err != nil {
      return err
    }
  }

  return c.SendString(fmt.Sprintf("Success authentication (username: %s)", body.Username))
}

func AuthDelete(c *fiber.Ctx) error {
  if sess, err := SessionStore.Get(c); err != nil {
    return err
  } else {
    username := sess.Get("username")

    if username == nil {
      return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
    } else {
      if err := sess.Destroy(); err != nil {
        return err
      }

      if err := sess.Save(); err != nil {
        return err
      }
    }
  }

  return c.SendString("sess Destroyed")
}
