package whapi

import (
	"encoding/json"
	"testing"
)

func TestUserIDUpdateUnmarshal(t *testing.T) {
	input := `{
		"messaging_product": "whatsapp",
		"metadata": {
			"display_phone_number": "16505551111",
			"phone_number_id": "123456123"
		},
		"user_id_update": [
			{
				"wa_id": "16315551181",
				"detail": "User changed phone number",
				"user_id": {
					"previous": "xxxxxxxxxxxxxxxxxxxxxxxxxxx",
					"current": "yyyyyyyyyyyyyyyyyyyyyyyyyyy"
				},
				"parent_user_id": {
					"previous": "aaaaaaaaaaaaaaaaaaaaaaaaaaa",
					"current": "bbbbbbbbbbbbbbbbbbbbbbbbbbb"
				},
				"timestamp": "1750100472"
			}
		]
	}`

	var v ValueObject
	if err := json.Unmarshal([]byte(input), &v); err != nil {
		t.Fatal(err)
	}

	if len(v.UserIDUpdate) != 1 {
		t.Fatalf("len(UserIDUpdate) = %d, want 1", len(v.UserIDUpdate))
	}
	u := v.UserIDUpdate[0]
	if u.WAID != "16315551181" {
		t.Errorf("WAID = %q, want 16315551181", u.WAID)
	}
	if u.Detail != "User changed phone number" {
		t.Errorf("Detail = %q", u.Detail)
	}
	if u.UserID.Previous != "xxxxxxxxxxxxxxxxxxxxxxxxxxx" || u.UserID.Current != "yyyyyyyyyyyyyyyyyyyyyyyyyyy" {
		t.Errorf("UserID = %+v", u.UserID)
	}
	if u.ParentUserID == nil || u.ParentUserID.Previous != "aaaaaaaaaaaaaaaaaaaaaaaaaaa" || u.ParentUserID.Current != "bbbbbbbbbbbbbbbbbbbbbbbbbbb" {
		t.Errorf("ParentUserID = %+v", u.ParentUserID)
	}
	if u.Timestamp != "1750100472" {
		t.Errorf("Timestamp = %q", u.Timestamp)
	}
}

func TestUserChangedUserIDSystemMessageUnmarshal(t *testing.T) {
	input := `{
		"from": "16315551181",
		"id": "wamid.HBgLMTY1MDM4Nzk0MzkVAgASGBQzQTdCNTk3RjgzMzM5RTGRDMzRDcA",
		"timestamp": "1750100472",
		"system": {
			"body": "User changed their user ID",
			"type": "user_changed_user_id",
			"user_id": "yyyyyyyyyyyyyyyyyyyyyyyyyyy",
			"parent_user_id": "bbbbbbbbbbbbbbbbbbbbbbbbbbb"
		},
		"type": "system"
	}`

	var m MessageObject
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		t.Fatal(err)
	}

	if m.System == nil {
		t.Fatal("System is nil")
	}
	if m.System.Type != SysMsgTypeUserChangedUserID {
		t.Errorf("System.Type = %q, want %q", m.System.Type, SysMsgTypeUserChangedUserID)
	}
	if m.System.UserID != "yyyyyyyyyyyyyyyyyyyyyyyyyyy" {
		t.Errorf("System.UserID = %q", m.System.UserID)
	}
	if m.System.ParentUserID != "bbbbbbbbbbbbbbbbbbbbbbbbbbb" {
		t.Errorf("System.ParentUserID = %q", m.System.ParentUserID)
	}
}
