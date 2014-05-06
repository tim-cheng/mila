package routes

import (
  "github.com/tim-cheng/mila/server/models"
)

type Routes struct {
  Db *models.MyDb
}

func New(db *models.MyDb) *Routes {
  return &Routes{db}
}
