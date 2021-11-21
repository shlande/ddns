package ddns

import (
	"reflect"
	"testing"
)

func TestFindBestSolution(t *testing.T) {
	has := []*Record{
		{
			Type:       "A",
			DomainName: "colaha.tech",
			RecordID:   "1",
			Value:      "10.10.10.1",
			Prefix:     "test",
		},
	}
	want := []*Record{
		{
			Type:       "A",
			DomainName: "colaha.tech",
			RecordID:   "1",
			Value:      "10.0.0.1",
			Prefix:     "test",
		},
	}
	crt, del, upd, err := findBestSolution(has, "A", "colaha.tech", "test", []string{"10.0.0.1"})
	if err != nil {
		panic(err)
	}
	if !check(crt, del, upd, nil, nil, want) {
		panic("not equal")
	}

	has = []*Record{
		{
			Type:       "AAAA",
			DomainName: "colaha.tech",
			RecordID:   "1",
			Value:      "10.10.10.1",
			Prefix:     "test",
		},
		{
			Type:       "A",
			DomainName: "colaha.tech",
			RecordID:   "2",
			Value:      "10.10.10.2",
			Prefix:     "test",
		},
	}

	crt, del, upd, err = findBestSolution(has, "A", "colaha.tech", "test", []string{"10.0.0.1"})
	if err != nil {
		panic(err)
	}
	if !check(crt, del, upd, nil, has[1:], want) {
		panic("not equal")
	}

	has = []*Record{}

	crt, del, upd, err = findBestSolution(has, "A", "colaha.tech", "test", []string{"10.0.0.1"})
	if err != nil {
		panic(err)
	}
	if !check(crt, del, upd, []*Record{
		{
			Type:       "A",
			DomainName: "colaha.tech",
			Value:      "10.0.0.1",
			Prefix:     "test",
		},
	}, nil, nil) {
		panic("not equal")
	}
}

func check(crt, del, udp, crtw, delw, udpw []*Record) bool {
	deepcheck := func(a, b []*Record) bool {
		if len(a) != len(b) {
			return false
		}
		for i, v := range a {
			if !reflect.DeepEqual(b[i], v) {
				return false
			}
		}
		return true
	}
	return deepcheck(crt, crtw) && deepcheck(del, delw) && deepcheck(udp, udpw)
}
