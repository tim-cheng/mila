package routes

import (
	"bytes"
	"fmt"
	"github.com/mostafah/mandrill"
	"net/http"
	"strconv"
)

func sendNewUserEmail(email, firstName string) {
	mandrill.Key = "izQqlSTrNP4ZKZQ_rtM3-Q"
	msg := mandrill.NewMessageTo(email, firstName)
	msg.Text = "Dear " + firstName + ",\n\n" +
		"Welcome to Parent2D, a great way to connect and interact with your trusted fellow parents within your 2 degree network. Getting answers and sharing parenting insights have never been so easy, quick, and fun!\n\n" +
		"Go ahead a give a try. Invite more friends to join. Connect and be better informed parents!\n\n" +
		"Sherry\n"
	msg.HTML = "<p>" + msg.Text + "</p>"
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

func sendUserPushMsg(userId int64, pushMsg string) {
	msg := fmt.Sprintf("{\"channels\":[\"user_%d\"],\"data\":{\"alert\":\"%s\"}}", userId, pushMsg)
	c := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.parse.com/1/push", bytes.NewBufferString(msg))
	if err != nil {
		return
	}
	req.Header.Add("X-Parse-Application-Id", "hR6yp7Uqz7B0JL8mflpbGKiQa9jsZS4IFFfToHxC")
	req.Header.Add("X-Parse-REST-API-Key", "58nBXGHOuQkIJeC0nfTxq3tCVWDMxOIZgg2g910J")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(msg)))
	resp, err := c.Do(req)
	fmt.Println("push response is: ", resp, err)
}
