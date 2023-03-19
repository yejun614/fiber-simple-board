package main

import (
  "fmt"
  "errors"
  "strconv"
  "simple-board/model"
  "github.com/gofiber/fiber/v2"
  "github.com/gofiber/fiber/v2/middleware/session"
)

type BoardPostBody struct {
  Title   string    `json:"title"`
  Content string    `json:"content"`
}

type BoardPostSuccess struct {
  ID uint
}

func BoardGet(c *fiber.Ctx) error {
  postID := c.Params("id")

  // Query DB
  post := new(model.Post)
  result := DB.Where("id = ?", postID).First(post)

  // cannot found a post
  if result.Error != nil {
    return result.Error
  }

  // Response
  return c.JSON(post)
}

func BoardPost(c *fiber.Ctx) error {
  var sess *session.Session
  var err error
  var userID uint

  if sess, err = SessionStore.Get(c); err != nil || sess.Get("id") == nil {
    return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
  }

  switch t := sess.Get("id").(type) {
  case uint:
    userID = t
  default:
    return errors.New("Session UserID type is not uint")
  }

  body := new(BoardPostBody)

  if err := c.BodyParser(body); err != nil {
    return err
  }

  post := model.Post{
    Title: body.Title,
    Content: body.Content,
    UserID: userID,
  }

  DB.Create(&post);

  return c.JSON(BoardPostSuccess{
    ID: post.ID,
  })
}

func BoardPut(c *fiber.Ctx) error {
  var err error
  var sess *session.Session
  var userID uint
  postID := c.Params("id")
  post := new(model.Post)

  if result := DB.Where("id = ?", postID).First(&post); result.Error != nil {
    return result.Error
  }

  if sess, err = SessionStore.Get(c); err != nil {
    return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
  }

  switch t := sess.Get("id").(type) {
  case uint:
    userID = t
  default:
    return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
  }

  if userID != post.UserID {
    return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
  }

  body := new(BoardPostBody)

  if err := c.BodyParser(body); err != nil {
    return err
  }

  if body.Title != "" {
    post.Title = body.Title
  }

  if body.Content != "" {
    post.Content = body.Content
  }

  DB.Save(&post)

  return c.SendString(fmt.Sprintf("Updated a post (id: %s)", postID))
}

func BoardDelete(c *fiber.Ctx) error {
  var err error
  var sess *session.Session
  var userID uint
  var postID uint
  post := new(model.Post)

  if paramPostID, err := strconv.ParseUint(c.Params("id"), 10, 32); err != nil {
    return err
  } else {
    postID = uint(paramPostID)
  }

  if sess, err = SessionStore.Get(c); err != nil {
    return c.SendStatus(403)
  }

  switch t := sess.Get("id").(type) {
  case uint:
    userID = t
  default:
    return c.SendStatus(403)
  }

  if userID != postID {
    return c.SendStatus(403)
  }

  if result := DB.Where("id = ?", postID).First(&post); result.Error != nil {
    return result.Error
  }

  if result := DB.Delete(&post); result.Error != nil {
    return result.Error
  }

  return c.SendString(fmt.Sprintf("Deleted a post (id: %d)", postID))
}
