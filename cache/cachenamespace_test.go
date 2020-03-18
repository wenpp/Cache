package cache

import (
	"testing"
)

func TestInitNameSpace(t *testing.T) {
	nsp := InitNameSapce("user", 0)
	err := nsp.Put("key1", String("hello"))
	if err != nil {
		t.Fatal(err.Error())
	}
	expect := String("hello")

	if value, _ := nsp.Get("key1"); value != expect {
		t.Fatal("expect value is hello actually got", value)
	}
}
