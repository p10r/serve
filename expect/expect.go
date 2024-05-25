package expect

import (
	"context"
	"encoding/json"
	"github.com/p10r/serve/db"
	"github.com/p10r/serve/domain"
	"reflect"
	"sort"
	"testing"
)

func Equal[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func SliceEqual[T comparable](t *testing.T, got, want []T) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("got %v, want %v", got, want)
		return
	}

	for i := range got {
		if got[i] != want[i] {
			t.Errorf("got %v, want %v", got, want)
			return
		}
	}
}

func DeepEqual[T any](t *testing.T, got, want T) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func NotEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got == want {
		t.Errorf("didn't want %v", got)
	}
}

func Len[T any](t *testing.T, got []T, want int) {
	if len(got) != want {
		t.Errorf("got length %d, want %d", len(got), want)
	}
}

func True(t *testing.T, got bool) {
	t.Helper()
	if !got {
		t.Error("got false, want true")
	}
}

func False(t *testing.T, got bool) {
	t.Helper()
	if got {
		t.Error("got true, want false")
	}
}

func NoErr(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func Err(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("wanted error")
	}
}

// JsonEqual compares the JSON in two byte slices.
func JsonEqual(t *testing.T, a, b []byte) bool {
	t.Helper()

	var j, j2 interface{}
	if err := json.Unmarshal(a, &j); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(b, &j2); err != nil {
		t.Fatal(err)
	}

	return reflect.DeepEqual(j2, j)
}

type matchWithoutID struct {
	HomeName  string
	AwayName  string
	StartTime int64
	Country   string
	League    string
}

func MatchesEqual(t *testing.T, got, want domain.Matches) {
	t.Helper()

	var gotten []matchWithoutID
	for _, match := range got {
		m := matchWithoutID{
			match.HomeName,
			match.AwayName,
			match.StartTime,
			match.Country,
			match.League,
		}
		gotten = append(gotten, m)
	}

	var wanted []matchWithoutID
	for _, match := range want {
		m := matchWithoutID{
			match.HomeName,
			match.AwayName,
			match.StartTime,
			match.Country,
			match.League,
		}
		wanted = append(wanted, m)
	}

	sort.Slice(gotten, func(i, j int) bool {
		return len(gotten[i].HomeName) > len(gotten[j].HomeName)
	})

	sort.Slice(wanted, func(i, j int) bool {
		return len(wanted[i].HomeName) > len(wanted[j].HomeName)
	})

	DeepEqual(t, gotten, wanted)
}

func MatchStoreContains(t *testing.T, store *db.MatchStore, want domain.Matches) {
	matches, err := store.All(context.Background())
	NoErr(t, err)

	MatchesEqual(t, matches, want)
}
