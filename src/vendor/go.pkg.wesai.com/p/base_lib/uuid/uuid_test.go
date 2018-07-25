package uuid

import (
	"regexp"
	"testing"
)

func TestUUIDForFormat(t *testing.T) {
	re := regexp.MustCompile("([0-9a-f]){8}-([0-9a-f]){4}-([0-9a-f]){4}-([0-9a-f]){4}-([0-9a-f]){12}")
	for i := 0; i < 100000; i++ {
		uuid1 := UUID(true)
		if re.FindString(uuid1) != uuid1 {
			t.Logf("Malformed UUID: %q!!", uuid1)
		}
		uuid2 := UUID(false)
		if re.FindString(uuid2) != uuid2 {
			t.Logf("Malformed UUID: %q!!", uuid2)
		}
	}
}

func TestUUIDForUniformity(t *testing.T) {
	rn := 100000
	m := make(map[string]bool, rn)
	for i := 0; i < rn; i++ {
		uuid1 := UUID(true)
		uuid2 := UUID(false)
		if uuid1 == uuid2 || m[uuid1] || m[uuid2] {
			t.Logf("Duplicate UUID %q!!", uuid1)
		}
		m[uuid1] = true
		m[uuid2] = true
	}
}

func BenchmarkUUIDWithSysTrue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UUID(true)
	}
}
func BenchmarkUUIDWithSysFalse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UUID(false)
	}
}

func BenchmarkUUIDWithSysTrueParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			UUID(true)
		}
	})
}

func BenchmarkUUIDWithSysFalseParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			UUID(false)
		}
	})
}
