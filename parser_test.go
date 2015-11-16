package config_test

import "github.com/twoleds-golang/config"
import "testing"

func TestParser(t *testing.T) {

	var str = `
		BoolValue T
		FloatValue 3.14
		IntValue 123
		StringValue "Test String 123"
	`

	cfg, err := config.ParseFromString(str)
	if err != nil {
		t.Errorf("Cannot parse config: %s", err.Error())
		t.FailNow()
	}

	if val, ok := cfg.Bool("BoolValue"); ok == false || val != true {
		t.Error("Invalid value for query 'BoolValue'")
		t.Fail()
	}

	if val, ok := cfg.Float("FloatValue"); ok == false || val != 3.14 {
		t.Error("Invalid value for query 'FloatValue'")
		t.Fail()
	}

	if val, ok := cfg.Int("IntValue"); ok == false || val != 123 {
		t.Error("Invalid value for query 'IntValue'")
		t.Fail()
	}

	if val, ok := cfg.String("StringValue"); ok == false || val != "Test String 123" {
		t.Error("Invalid value for query 'StringValue'")
		t.Errorf("Msg: %s", val)
		t.Fail()
	}

}
