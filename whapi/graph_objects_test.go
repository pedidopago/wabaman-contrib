package whapi

import (
	"encoding/json"
	"github.com/pedidopago/wabaman-contrib/v2/fbgraph"
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

func TestEntryTimeUnmarshal(t *testing.T) {
	input := `{
		"object": "whatsapp_business_account",
		"entry": [{
			"id": "110729168243459",
			"time": 1790883985,
			"changes": [{"field": "account_update", "value": {"event": "ACCOUNT_OFFBOARDED"}}]
		}]
	}`

	var w WebhookObject
	if err := json.Unmarshal([]byte(input), &w); err != nil {
		t.Fatal(err)
	}
	if len(w.Entry) != 1 {
		t.Fatalf("len(Entry) = %d, want 1", len(w.Entry))
	}
	if w.Entry[0].Time != 1790883985 {
		t.Errorf("Time = %d, want 1790883985", w.Entry[0].Time)
	}
}

func TestSystemMessageNumberChange(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name: "user_changed_number (v12+): new number in wa_id",
			input: `{
				"from": "5511900000001",
				"from_user_id": "BR.1000000000000001",
				"id": "wamid.X",
				"system": {
					"type": "user_changed_number",
					"body": "User A changed from 5511900000001 to 5511900000002",
					"identity": "",
					"new_wa_id": "",
					"wa_id": "5511900000002",
					"customer": "",
					"user_id": "BR.1000000000000001"
				},
				"timestamp": "1791221040",
				"type": "system"
			}`,
			want: "5511900000002",
		},
		{
			name: "customer_changed_number (v11-): new number in new_wa_id",
			input: `{
				"from": "5511900000001",
				"id": "wamid.Y",
				"system": {
					"type": "customer_changed_number",
					"body": "User A changed from 5511900000001 to 5511900000002",
					"new_wa_id": "5511900000002",
					"customer": "5511900000001"
				},
				"timestamp": "1791221040",
				"type": "system"
			}`,
			want: "5511900000002",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m MessageObject
			if err := json.Unmarshal([]byte(tt.input), &m); err != nil {
				t.Fatal(err)
			}
			if m.System == nil || !m.System.IsNumberChange() {
				t.Fatalf("IsNumberChange() = false for %+v", m.System)
			}
			if got := m.System.NewNumber(); got != tt.want {
				t.Errorf("NewNumber() = %q, want %q", got, tt.want)
			}
		})
	}

	if (MessageObjectSystem{Type: SysMsgTypeUserChangedUserID}).IsNumberChange() {
		t.Error("user_changed_user_id must not be a number change")
	}
}

func TestValueObjectMessageTemplateIDNumberOrString(t *testing.T) {
	for _, input := range []string{
		`{"event":"APPROVED","message_template_id":1234567890123,"metadata":{"phone_number_id":"1"}}`,
		`{"event":"APPROVED","message_template_id":"1234567890123","metadata":{"phone_number_id":"1"}}`,
	} {
		var v ValueObject
		if err := json.Unmarshal([]byte(input), &v); err != nil {
			t.Fatalf("unmarshal %s: %v", input, err)
		}
		if v.MessageTemplateID != 1234567890123 || v.Event != "APPROVED" || v.Metadata.PhoneNumberID != "1" {
			t.Errorf("got id=%d event=%q phone=%q from %s", v.MessageTemplateID, v.Event, v.Metadata.PhoneNumberID, input)
		}
	}

	var v ValueObject
	if err := json.Unmarshal([]byte(`{"message_template_id":"abc"}`), &v); err == nil {
		t.Error("want error for non-numeric message_template_id")
	}
}

func TestPartnerClientCertificationUnmarshal(t *testing.T) {
	input := `{"event":"PARTNER_CLIENT_CERTIFICATION_STATUS_UPDATE","partner_client_certification_info":{"client_business_id":"2729063490586005","status":"FAILED","rejection_reasons":["LEGAL_NAME_NOT_FOUND_IN_DOCUMENTS"]}}`
	var v ValueObject
	if err := json.Unmarshal([]byte(input), &v); err != nil {
		t.Fatal(err)
	}
	info := v.PartnerClientCertificationInfo
	if AccountUpdateEvent(v.Event) != AccountUpdatePartnerClientCertificationStatusUpdate || info == nil ||
		info.ClientBusinessID != "2729063490586005" || info.Status != PartnerClientCertificationFailed ||
		len(info.RejectionReasons) != 1 || info.RejectionReasons[0] != "LEGAL_NAME_NOT_FOUND_IN_DOCUMENTS" {
		t.Errorf("got event=%q info=%+v", v.Event, info)
	}
}

func TestHandoverUnmarshal(t *testing.T) {
	passed := `{"messaging_product":"whatsapp","sender":{"phone_number":"5511999990000"},"recipient":{"phone_number_id":"123","display_phone_number":"+55 11 3000-0000"},"type":"control_passed","timestamp":"1794917311","control_passed":{"previous_owner_app_id":"1","previous_owner_app_role":"meta_business_agent","previous_owner_role":"ai_agent","new_owner_role":"escalation","metadata":"user asked for a human","conversation_context":{"type":"summary","summary":{"text":"wants a refund"}}}}`
	var v ValueObject
	if err := json.Unmarshal([]byte(passed), &v); err != nil {
		t.Fatal(err)
	}
	cp := v.ControlPassed
	if v.Type != HandoverControlPassed || v.Sender == nil || v.Sender.PhoneNumber != "5511999990000" ||
		v.Recipient == nil || v.Recipient.PhoneNumberID != "123" || v.Timestamp != "1794917311" || cp == nil ||
		cp.PreviousOwnerRole != fbgraph.RoutingRoleAIAgent || cp.NewOwnerRole != fbgraph.RoutingRoleEscalation ||
		cp.Metadata != "user asked for a human" || cp.ConversationContext.SummaryText() != "wants a refund" {
		t.Errorf("control_passed: got %+v / %+v", v, cp)
	}

	taken := `{"messaging_product":"whatsapp","sender":{"phone_number":"5511999990000"},"recipient":{"phone_number_id":"123","display_phone_number":"x"},"type":"control_taken","timestamp":"1794917312","control_taken":{"previous_owner_role":"customer_service","new_owner_role":"escalation"}}`
	v = ValueObject{}
	if err := json.Unmarshal([]byte(taken), &v); err != nil {
		t.Fatal(err)
	}
	if v.Type != HandoverControlTaken || v.ControlTaken == nil || v.ControlTaken.PreviousOwnerRole != fbgraph.RoutingRoleCustomerService ||
		v.ControlTaken.ConversationContext.SummaryText() != "" {
		t.Errorf("control_taken: got %+v", v.ControlTaken)
	}
}

func TestStandbyUnmarshal(t *testing.T) {
	in := `{"messaging_product":"whatsapp","metadata":{"display_phone_number":"x","phone_number_id":"123"},"standby":{"contacts":[{"profile":{"name":"Ana"},"wa_id":"5511999990000","user_id":"BR.1"}],"messages":[{"from":"5511999990000","from_user_id":"BR.1","id":"wamid.in","timestamp":"1","type":"text","text":{"body":"oi"}}]}}`
	var v ValueObject
	if err := json.Unmarshal([]byte(in), &v); err != nil {
		t.Fatal(err)
	}
	if v.Standby == nil || len(v.Standby.Messages) != 1 || v.Standby.Messages[0].ID != "wamid.in" || len(v.Standby.Contacts) != 1 {
		t.Fatalf("standby messages: got %+v", v.Standby)
	}

	echo := `{"messaging_product":"whatsapp","metadata":{"display_phone_number":"x","phone_number_id":"123"},"standby":{"message_echoes":[{"id":"wamid.out","timestamp":"2","message":{"messaging_product":"whatsapp","to":"5511999990000","recipient_type":"individual","type":"text","text":{"body":"seu pedido saiu","preview_url":true},"context":{"message_id":"wamid.in"}},"template":{"name":"t"}}]}}`
	v = ValueObject{}
	if err := json.Unmarshal([]byte(echo), &v); err != nil {
		t.Fatal(err)
	}
	if v.Standby == nil || len(v.Standby.MessageEchoes) != 1 {
		t.Fatalf("standby echoes: got %+v", v.Standby)
	}
	e := v.Standby.MessageEchoes[0]
	if e.ID != "wamid.out" || e.Message == nil || e.Message.To != "5511999990000" || e.Message.Text == nil ||
		e.Message.Text.Body != "seu pedido saiu" || e.Message.Context == nil || e.Message.Context.MessageID != "wamid.in" || len(e.Template) == 0 {
		t.Errorf("echo: got %+v", e)
	}
}

func TestConversationContextOnMessages(t *testing.T) {
	in := `{"messaging_product":"whatsapp","messages":[{"from":"1","id":"wamid.x","timestamp":"1","type":"text","text":{"body":"oi"}}],"conversation_context":{"type":"summary","summary":{"text":"resumo"}}}`
	var v ValueObject
	if err := json.Unmarshal([]byte(in), &v); err != nil {
		t.Fatal(err)
	}
	if v.ConversationContext.SummaryText() != "resumo" || len(v.Messages) != 1 {
		t.Errorf("got %+v", v.ConversationContext)
	}
}
