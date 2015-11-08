package config_test

import "bufio"
import "fmt"

import "github.com/twoleds/config"
import "io/ioutil"
import "os"
import "testing"

func TestWriter(t *testing.T) {

	f, err := ioutil.TempFile(os.TempDir(), "golang-config-")
	if err != nil {
		t.Errorf("Cannot create temp file: %s", err.Error())
		t.FailNow()
	}
	defer f.Close()
	defer os.Remove(f.Name())

	b := config.NewWriter(f)

	b.Bool("BoolValue", true)
	b.Float("FloatValue", 3.14)
	b.Int("IntValue", 123)
	b.String("StringValue", "Test String 123")
	b.Line()

	b.Section("Section", "One")
	b.Bool("BoolValue", true)
	b.Float("FloatValue", 2.745)
	b.Int("IntValue", 124)
	b.String("StringValue", "Test String - First Section")
	b.CloseSection()
	b.Line()

	b.Section("Nested", "One")
	b.Section("Nested", "Two")
	b.Section("Nested", "Three")
	b.Bool("BoolValue", false)
	b.Float("FloatValue", -3.14)
	b.Int("IntValue", -123456)
	b.String("StringValue", "Test String - Nested Section")
	b.CloseSection()
	b.CloseSection()
	b.CloseSection()
	b.Line()

	b.Section("Section", "Two")
	b.Bool("BoolValue", false)
	b.Float("FloatValue", -2.745)
	b.Int("IntValue", 123456)
	b.String("StringValue", "Test String 123")
	b.CloseSection()
	b.Line()

	b.Section("Section", "Three")
	b.Bool("BoolValue", true)
	b.Float("FloatValue", 3.14)
	b.Int("IntValue", -150)
	b.String("StringValue", "Test String - Last Section")
	b.CloseSection()
	b.Line()

	b.Flush()
	f.Seek(0, 0)

	x := bufio.NewReader(f)
	z, _ := x.ReadString('\f')
	fmt.Println(z)
	f.Seek(0, 0)

	c, err := config.ParseFromFile(f.Name())
	if err != nil {
		t.Errorf("Cannot parse temp file: %s", err.Error())
		t.FailNow()
	}

	// Test flat values

	if val, ok := c.Bool("BoolValue"); ok == false || val != true {
		t.Error("Invalid value for query 'BoolValue'")
		t.Fail()
	}

	if val, ok := c.Float("FloatValue"); ok == false || val != 3.14 {
		t.Error("Invalid value for query 'FloatValue'")
		t.Fail()
	}

	if val, ok := c.Int("IntValue"); ok == false || val != 123 {
		t.Error("Invalid value for query 'IntValue'")
		t.Fail()
	}

	if val, ok := c.String("StringValue"); ok == false || val != "Test String 123" {
		t.Error("Invalid value for query 'StringValue'")
		t.Fail()
	}

	// Test first section

	if val, ok := c.Bool("Section/BoolValue"); ok == false || val != true {
		t.Error("Invalid value for query 'Section/BoolValue'")
		t.Fail()
	}

	if val, ok := c.Float("Section/FloatValue"); ok == false || val != 2.745 {
		t.Error("Invalid value for query 'Section/FloatValue'")
		t.Fail()
	}

	if val, ok := c.Int("Section/IntValue"); ok == false || val != 124 {
		t.Error("Invalid value for query 'Section/IntValue'")
		t.Fail()
	}

	if val, ok := c.String("Section/StringValue"); ok == false || val != "Test String - First Section" {
		t.Error("Invalid value for query 'Section/StringValue'")
		t.Fail()
	}

	// Test nested section

	if val, ok := c.Bool("Nested/Nested/Nested/BoolValue"); ok == false || val != false {
		t.Error("Invalid value for query 'Nested/Nested/Nested/BoolValue'")
		t.Fail()
	}

	if val, ok := c.Float("Nested/Nested/Nested/FloatValue"); ok == false || val != -3.14 {
		t.Error("Invalid value for query 'Nested/Nested/Nested/FloatValue'")
		t.Fail()
	}

	if val, ok := c.Int("Nested/Nested/Nested/IntValue"); ok == false || val != -123456 {
		t.Error("Invalid value for query 'Nested/Nested/Nested/IntValue'")
		t.Fail()
	}

	if val, ok := c.String("Nested/Nested/Nested/StringValue"); ok == false || val != "Test String - Nested Section" {
		t.Error("Invalid value for query 'Nested/Nested/Nested/StringValue'")
		t.Fail()
	}

	// Test second section

	if val, ok := c.Bool("Section:Two/BoolValue"); ok == false || val != false {
		t.Error("Invalid value for query 'Section:Two/BoolValue'")
		t.Fail()
	}

	if val, ok := c.Float("Section:Two/FloatValue"); ok == false || val != -2.745 {
		t.Error("Invalid value for query 'Section:Two/FloatValue'")
		t.Fail()
	}

	if val, ok := c.Int("Section:Two/IntValue"); ok == false || val != 123456 {
		t.Error("Invalid value for query 'Section:Two/IntValue'")
		t.Fail()
	}

	if val, ok := c.String("Section:Two/StringValue"); ok == false || val != "Test String 123" {
		t.Error("Invalid value for query 'Section:Two/StringValue'")
		t.Fail()
	}

	// Test last section

	if val, ok := c.Bool("Section:Three/BoolValue"); ok == false || val != true {
		t.Error("Invalid value for query 'Section:Three/BoolValue'")
		t.Fail()
	}

	if val, ok := c.Float("Section:Three/FloatValue"); ok == false || val != 3.14 {
		t.Error("Invalid value for query 'Section:Three/FloatValue'")
		t.Fail()
	}

	if val, ok := c.Int("Section:Three/IntValue"); ok == false || val != -150 {
		t.Error("Invalid value for query 'Section:Three/IntValue'")
		t.Fail()
	}

	if val, ok := c.String("Section:Three/StringValue"); ok == false || val != "Test String - Last Section" {
		t.Error("Invalid value for query 'Section:Three/StringValue'")
		t.Fail()
	}

}
