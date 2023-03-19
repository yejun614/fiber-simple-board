package model;

type Post struct {
  ID        uint
  Title     string
  Content   string
  UserID    uint
  CreatedAt int64    `gorm:"autoCreateTime"`
  UpdatedAt int64    `gorm:"autoUpdateTime:milli"`
}
