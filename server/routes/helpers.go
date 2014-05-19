package routes

import (
  "github.com/mostafah/mandrill"
  "fmt"
)

func sendNewUserEmail(email, firstName string) {
  mandrill.Key = "izQqlSTrNP4ZKZQ_rtM3-Q"
  msg := mandrill.NewMessageTo(email, firstName)
  msg.HTML = "<p>Welcome to Parent2D</p>"
  msg.Text = "Welcome to Parent2D"
  msg.Subject = "Welcome to Parent2D"
  msg.FromEmail = "sherry@parent2d.com"
  msg.FromName = "Parent2D"
  res, err := msg.Send(false)
  fmt.Printf("res = %v, err = %v\n", res, err)
}

func sendUserInviteEmail(email, firstName, inviterName string) {
  mandrill.Key = "izQqlSTrNP4ZKZQ_rtM3-Q"
  msg := mandrill.NewMessageTo(email, firstName)
  msg.HTML = "<p>" + inviterName + " invites you to join Parent2D</p>"
  msg.Text = inviterName + " invites you to join Parent2D"
  msg.Subject = inviterName + " invites you to join Parent2D"
  msg.FromEmail = "sherry@parent2d.com"
  msg.FromName = "Parent2D"
  res, err := msg.Send(false)
  fmt.Printf("res = %v, err = %v\n", res, err)
}
