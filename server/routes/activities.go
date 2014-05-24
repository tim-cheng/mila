package routes

import (
  "github.com/codegangsta/martini"
  "github.com/codegangsta/martini-contrib/render"
  "github.com/tim-cheng/mila/server/models"
)

func (rt *Routes) GetActivities(params martini.Params, r render.Render) {
  activities, err := rt.Db.GetActivities(params["id"])
  if err == nil && len(activities) > 0 {
    retActs := make([]map[string]interface{}, len(activities))
    for i := range activities {
      a := activities[i].(*models.Activity)
      u, err := rt.Db.GetUserName(a.FriendId)
      if err == nil {
        typeStr := "";
        switch a.Type {
          case models.ActivityTypePost:
            typeStr = "posted"
          case models.ActivityTypeComment:
            typeStr = "commented on your post"
          case models.ActivityTypeLike:
            typeStr = "liked your post"
          case models.ActivityTypeInvite:
            typeStr = "invited you"
        }
        retActs[i] = map[string]interface{}{
          "friend_id" : a.FriendId,
          "post_id" : a.PostId,
          "activity" : u.FirstName + " " + u.LastName + " " + typeStr + " " + a.Message,
        }
      }
    }
    r.JSON(200, retActs)
  } else {
    r.JSON(404, map[string]interface{}{
      "message": "Failed to get activities",
    })
  }
}

// helpers
func (rt *Routes) postActivityInvite(u1Id int64, u2Id int64, msg string) {
  a, err := rt.Db.NewActivity(u1Id, u2Id, 0, models.ActivityTypeInvite, msg)
  if err == nil {
    rt.Db.PostActivity(a)
  }
}

func (rt *Routes) postActivityComment(u1Id int64, u2Id int64, postId int64, msg string) {
  a, err := rt.Db.NewActivity(u1Id, u2Id, postId, models.ActivityTypeComment, msg)
  if err == nil {
    rt.Db.PostActivity(a)
  }
}

func (rt *Routes) postActivityPost(u1Id int64, u2Id int64, postId int64, msg string) {
  a, err := rt.Db.NewActivity(u1Id, u2Id, postId, models.ActivityTypePost, msg)
  if err == nil {
    rt.Db.PostActivity(a)
  }
}

func (rt *Routes) postActivityLike(u1Id int64, u2Id int64, postId int64) {
  a, err := rt.Db.NewActivity(u1Id, u2Id, postId, models.ActivityTypeLike, "")
  if err == nil {
    rt.Db.PostActivity(a)
  }
}
