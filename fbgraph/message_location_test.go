package fbgraph

import (
	"encoding/json"
	"testing"
)

func TestLocationObjectUnmarshalBothTypes(t *testing.T) {
	cases := map[string]string{
		"meta inbound (numbers)": `{"latitude":37.44216251868683,"longitude":-122.16153582049394,"name":"Philz","address":"101 Forest Ave"}`,
		"send api (strings)":     `{"latitude":"37.44216251868683","longitude":"-122.16153582049394","name":"Philz","address":"101 Forest Ave"}`,
	}
	for name, in := range cases {
		var l LocationObject
		if err := json.Unmarshal([]byte(in), &l); err != nil {
			t.Fatalf("%s: unmarshal: %v", name, err)
		}
		if l.Latitude != "37.44216251868683" || l.Longitude != "-122.16153582049394" {
			t.Fatalf("%s: got lat=%q lng=%q", name, l.Latitude, l.Longitude)
		}
		if l.Name != "Philz" || l.Address != "101 Forest Ave" {
			t.Fatalf("%s: got name=%q address=%q", name, l.Name, l.Address)
		}
	}

	// absent / null lat-lng must not panic and must be empty
	var l LocationObject
	if err := json.Unmarshal([]byte(`{"name":"x","latitude":null}`), &l); err != nil {
		t.Fatalf("null: unmarshal: %v", err)
	}
	if l.Latitude != "" || l.Longitude != "" {
		t.Fatalf("null: expected empty, got lat=%q lng=%q", l.Latitude, l.Longitude)
	}
}
