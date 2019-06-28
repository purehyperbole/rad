package rad

import (
	"bytes"
	"fmt"
	"strconv"
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
		Name     string
		Existing []testvalue
		Lookups  []testvalue
	}{
		{
			"simple",
			[]testvalue{{"test", "1234", "est"}},
			[]testvalue{{"test", "1234", "est"}},
		},
		{
			"normal",
			[]testvalue{{"too", "1234", "oo"}, {"bad", "5678", "ad"}, {"you'll", "9101112", "ou'll"}, {"never", "13141516", "ever"}, {"be", "17181920", "e"}, {"rad", "21222324", "ad"}},
			[]testvalue{{"too", "1234", "oo"}, {"bad", "5678", "ad"}, {"you'll", "9101112", "ou'll"}, {"never", "13141516", "ever"}, {"be", "17181920", "e"}, {"rad", "21222324", "ad"}},
		},
		{
			"derivative",
			[]testvalue{{"test", "1234", "est"}, {"test1234", "bacon", "est"}},
			[]testvalue{{"test1234", "bacon", "234"}},
		},
		{
			"split",
			[]testvalue{{"test1234", "bacon", "234"}, {"test", "1234", "est"}},
			[]testvalue{{"test1234", "bacon", "234"}, {"test1234", "bacon", "234"}},
		},
		{
			"split-single-shared-character",
			[]testvalue{{"test", "1234", "est"}, {"test1234", "bacon", "est"}, {"test1000", "egg", "est"}},
			[]testvalue{{"test", "1234", "est"}, {"test1234", "bacon", "34"}, {"test1000", "egg", "00"}},
		},
		{
			"complex",
			[]testvalue{{"test", "1234", "st"}, {"test1234", "bacon", "234"}, {"tomato", "egg", "ato"}, {"tamale", "hash browns", "male"}, {"todo", "beans", ""}, {"todos", "mushrooms", "s"}, {"abalienate", "toast", ""}, {"abalienated", "onions", ""}, {"abalienating", "sausage", "ng"}},
			[]testvalue{{"test", "1234", "st"}, {"test1234", "bacon", "234"}, {"tomato", "egg", "ato"}, {"tamale", "hash browns", "male"}, {"todo", "beans", "o"}, {"todos", "mushrooms", ""}, {"abalienate", "toast", ""}, {"abalienated", "onions", ""}, {"abalienating", "sausage", "ng"}},
		},
		{
			"single-character",
			[]testvalue{{"todo", "toast", "odo"}, {"todos", "bacon", ""}},
			[]testvalue{{"todo", "toast", "odo"}, {"todos", "bacon", ""}},
		},
		{
			"mixed",
			[]testvalue{{"unsophisticatedness", "0", "-"}, {"unsophisticate", "1", "-"}, {"unsophisticatedly", "2", "-"}, {"unsophisticated", "3", "-"}, {"unsophistication", "4", "-"}, {"unsophistic", "5", "-"}, {"unsophistically", "6", "-"}, {"unsophistical", "7", "-"}},
			[]testvalue{{"unsophisticatedness", "0", "-"}, {"unsophisticate", "1", "-"}, {"unsophisticatedly", "2", "-"}, {"unsophisticated", "3", "-"}, {"unsophistication", "4", "-"}, {"unsophistic", "5", "-"}, {"unsophistically", "6", "-"}, {"unsophistical", "7", "-"}},
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

	fmt.Println(Graphviz(r))

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
		go func(b int) {
			for x := range batch[b] {
				r.MustInsert(batch[b][x], batch[b][x])
			}
			wg.Done()
		}(i)
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

func TestConcurrentInsertInt(t *testing.T) {
	var wg sync.WaitGroup

	r := New()

	batch := make([][][]byte, 32)

	for i := 0; i < 32; i++ {
		batch[i] = make([][]byte, 10000)

		for x := 0; x < 10000; x++ {
			batch[i][x] = []byte(strconv.Itoa(x))
		}
	}

	wg.Add(32)

	for i := 0; i < 32; i++ {
		go func(b int) {
			for x := range batch[b] {
				r.Insert(batch[b][x], batch[b][x])
			}
			wg.Done()
		}(i)
	}

	wg.Wait()

	for x := 0; x < 1000; x++ {
		value := r.Lookup(batch[0][x])
		require.NotNil(t, value)
		assert.True(t, bytes.Equal(value.([]byte), batch[0][x]))
	}
}

func b(char string) byte {
	return []byte(char)[0]
}
