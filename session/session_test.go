package session

import (
	"fmt"
	"testing"
)

func Test_NewSession(*testing.T) {
	s := NewSession("123", "lidawei", "20")

	fmt.Printf("NewSession: %v\n", s)
}

func Test_NewSessionStore(*testing.T) {
	ss := NewSessionStore()

	fmt.Printf("NewSessionStore: %p\n", ss)
}

func Test_SessionStore(t *testing.T) {

	s := NewSession("123", "lidawei", "20")

	ss := NewSessionStore()

	// Set
	err := ss.Set(s)

	if err != nil {
		t.Error(err)
	}

	// Get
	session, err := ss.Get("ffadsfjalksj")

	if session == nil && err == nil {
		fmt.Println("Can't find session")
	}

	session, err = ss.Get(s.SessionId)

	if session != nil {
		fmt.Printf("sessionId=%s, name=%s, expiration=%s", session.SessionId, session.UserName, session.Expiration)
	} else {
		t.Errorf("Error: session=%p, err=%s", session, err)
	}

	/*
		// Delete
		err = ss.Delete(s.SessionId)

		if err != nil {
			t.Error(err)
		}
	*/

	//ValidateOrgId
	if ok, err := ss.ValidateOrgId(s.SessionId, s.OrgId); ok {
		fmt.Printf("sessionId=%s, orgId=%s\n", s.SessionId, s.OrgId)
	} else {
		fmt.Println(err)
		t.Errorf("Error: sessionId=%s, orgId=%s\n", s.SessionId, s.OrgId)
	}

	if ok, err := ss.ValidateOrgId(s.SessionId, "456"); ok {
		fmt.Printf("sessionId=%s, orgId=%s\n", s.SessionId, s.OrgId)
	} else {
		fmt.Println(err)
	}

	if ok, err := ss.ValidateOrgId("adfadfadfa", s.OrgId); ok {
		fmt.Printf("sessionId=%s, orgId=%s\n", s.SessionId, s.OrgId)
	} else {
		fmt.Println(err)
	}

	if ok, err := ss.ValidateOrgId("adfadfad", "456"); ok {
		fmt.Printf("sessionId=%s, orgId=%s\n", s.SessionId, s.OrgId)
	} else {
		fmt.Println(err)
	}

	// ValidateUserId

	if ok, err := ss.ValidateUserId(s.SessionId, s.UserId); ok {
		fmt.Printf("sessionId=%s, userId=%s\n", s.SessionId, s.UserId)
	} else {
		fmt.Println(err)
		t.Errorf("Error: sessionId=%s, userId=%s\n", s.SessionId, s.UserId)
	}

	if ok, err := ss.ValidateUserId(s.SessionId, "456"); ok {
		fmt.Printf("sessionId=%s, userId=%s\n", s.SessionId, s.UserId)
	} else {
		fmt.Println(err)
	}

	if ok, err := ss.ValidateUserId("adfadfadfa", s.UserId); ok {
		fmt.Printf("sessionId=%s, userId=%s\n", s.SessionId, s.UserId)
	} else {
		fmt.Println(err)
	}

	if ok, err := ss.ValidateUserId("adfadfad", "456"); ok {
		fmt.Printf("sessionId=%s, userId=%s\n", s.SessionId, s.UserId)
	} else {
		fmt.Println(err)
	}
}
