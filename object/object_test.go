package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{ Value: "Hello World" }
	hello2 := &String{ Value: "Hello World" }
	diff1 := &String{ Value: "My name is jhonny" }
	diff2 := &String{ Value: "My name is jhonny" }

	if hello1.MapKey() != hello2.MapKey() {
		t.Errorf("strings with the same content have different map keys")
	}

	if diff1.MapKey() != diff2.MapKey() {
		t.Errorf("strings with the same content have different map keys")
	}

	if hello1.MapKey() == diff1.MapKey() {
		t.Errorf("strings with different content have same map keys")
	}

}
