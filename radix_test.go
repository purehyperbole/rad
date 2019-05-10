package rad

import (
	"bytes"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testvalue struct {
	Key    string
	Value  string
	Prefix string
}

func TestRadixInsertLookup(t *testing.T) {
	cases := []struct {
		Name          string
		ExpectedNodes int
		Existing      []testvalue
		Lookups       []testvalue
	}{
		{
			"simple",
			2,
			[]testvalue{{"test", "1234", "est"}},
			[]testvalue{{"test", "1234", "est"}},
		},
		{
			"normal",
			2,
			[]testvalue{{"too", "1234", "oo"}, {"bad", "5678", "ad"}, {"you'll", "9101112", "ou'll"}, {"never", "13141516", "ever"}, {"be", "17181920", "e"}, {"rad", "21222324", "ad"}},
			[]testvalue{{"too", "1234", "oo"}, {"bad", "5678", "ad"}, {"you'll", "9101112", "ou'll"}, {"never", "13141516", "ever"}, {"be", "17181920", "e"}, {"rad", "21222324", "ad"}},
		},
		{
			"derivative",
			3,
			[]testvalue{{"test", "1234", "est"}, {"test1234", "bacon", "est"}},
			[]testvalue{{"test1234", "bacon", "234"}},
		},
		{
			"split",
			3,
			[]testvalue{{"test1234", "bacon", "234"}, {"test", "1234", "est"}},
			[]testvalue{{"test1234", "bacon", "234"}, {"test1234", "bacon", "234"}},
		},
		{
			"split-single-shared-character",
			5,
			[]testvalue{{"test", "1234", "est"}, {"test1234", "bacon", "est"}, {"test1000", "egg", "est"}},
			[]testvalue{{"test", "1234", "est"}, {"test1234", "bacon", "34"}, {"test1000", "egg", "00"}},
		},
		{
			"complex",
			13,
			[]testvalue{{"test", "1234", "st"}, {"test1234", "bacon", "234"}, {"tomato", "egg", "ato"}, {"tamale", "hash browns", "male"}, {"todo", "beans", ""}, {"todos", "mushrooms", "s"}, {"abalienate", "toast", ""}, {"abalienated", "onions", ""}, {"abalienating", "sausage", "ng"}},
			[]testvalue{{"test", "1234", "st"}, {"test1234", "bacon", "234"}, {"tomato", "egg", "ato"}, {"tamale", "hash browns", "male"}, {"todo", "beans", "o"}, {"todos", "mushrooms", ""}, {"abalienate", "toast", ""}, {"abalienated", "onions", ""}, {"abalienating", "sausage", "ng"}},
		},
		{
			"single-character",
			3,
			[]testvalue{{"todo", "toast", "odo"}, {"todos", "bacon", ""}},
			[]testvalue{{"todo", "toast", "odo"}, {"todos", "bacon", ""}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			r := New()

			for _, kv := range tc.Existing {
				r.Insert([]byte(kv.Key), kv.Value)
			}

			for _, kv := range tc.Lookups {
				value := r.Lookup([]byte(kv.Key))
				require.NotNil(t, value)
				assert.Equal(t, kv.Value, value)
			}
		})
	}
}

func TestIterate(t *testing.T) {
	keys := []string{
		"hypotensive",
		"hyposulfurous",
		"hypotensor",
		"hypotension",
		"hypotaxia",
		"hypotaxic",
		"hyposulfite",
		"hyposuprarenalism",
		"hyposulphurous",
		"hyposulphuric",
		"hypostomatic",
		"hyposulphate",
		"hypostomial",
		"hypotenuses",
		"hypostome",
		"hypostoma",
		"hypotensions",
		"hypotarsal",
		"hypostomatous",
		"hypotaxis",
		"hypostrophe",
		"hyposulphite",
		"hypostomous",
		"hypotarsus",
		"hypostyptic",
		"hypotactic",
		"hypostypsis",
		"hypotenusal",
		"hypotenuse",
	}

	r := New()

	var results [][]byte

	for _, k := range keys {
		r.Insert([]byte(k), []byte(k))
	}

	r.Iterate([]byte("hypot"), func(key []byte, value interface{}) {
		results = append(results, key)
	})

	assert.Len(t, results, 13)

	for i := range results {
		assert.True(t, bytes.HasPrefix(results[i], []byte("hypot")))
	}
}

func TestConcurrentInsert(t *testing.T) {
	var wg sync.WaitGroup

	r := New()

	batch := make([][][]byte, 8)

	for i := 0; i < 8; i++ {
		batch[i] = make([][]byte, 1000)

		for x := 0; x < 1000; x++ {
			batch[i][x] = []byte(uuid.New().String())
		}
	}

	wg.Add(8)

	for i := 0; i < 8; i++ {
		for x := range batch[i] {
			r.MustInsert(batch[i][x], batch[i][x])
		}
		wg.Done()
	}

	wg.Wait()

	for i := 0; i < 8; i++ {
		for x := 0; x < 1000; x++ {
			value := r.Lookup(batch[i][x])
			require.NotNil(t, value)
			assert.True(t, bytes.Equal(value.([]byte), batch[i][x]))
		}
	}
}

func b(char string) byte {
	return []byte(char)[0]
}
