package rad

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
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
				r.Insert([]byte(kv.Key), String(kv.Value))
			}

			for _, kv := range tc.Lookups {
				value := r.Lookup([]byte(kv.Key))
				require.NotNil(t, value)
				assert.Equal(t, String(kv.Value), value)
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
		r.Insert([]byte(k), Bytes(k))
	}

	err := r.Iterate([]byte("hypot"), func(key []byte, value Comparable) error {
		results = append(results, key)
		return nil
	})

	assert.Nil(t, err)
	assert.Len(t, results, 13)

	for i := range results {
		assert.True(t, bytes.HasPrefix(results[i], []byte("hypot")))
	}

	err = r.Iterate([]byte("hypot"), func(key []byte, value Comparable) error {
		return errors.New("hello")
	})

	assert.NotNil(t, err)

	err = r.Iterate([]byte("does-not-exist"), func(key []byte, value Comparable) error {
		return nil
	})

	assert.Nil(t, err)

	r = New()

	var found []byte

	r.Insert([]byte("hello"), String("test"))

	r.Iterate([]byte("hel"), func(key []byte, value Comparable) error {
		found = key
		return nil
	})

	assert.Equal(t, "hello", string(found))
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
				r.MustInsert(batch[b][x], Bytes(batch[b][x]))
			}
			wg.Done()
		}(i)
	}

	wg.Wait()

	for i := 0; i < 8; i++ {
		for x := 0; x < 1000; x++ {
			value := r.Lookup(batch[i][x])
			require.NotNil(t, value)
			assert.True(t, bytes.Equal(value.(Bytes), batch[i][x]))
		}
	}
}

func TestConcurrentInsertInt(t *testing.T) {
	var wg sync.WaitGroup

	w := 32

	r := New()

	batch := make([][][]byte, w)

	for i := 0; i < w; i++ {
		batch[i] = make([][]byte, 10000)

		for x := 0; x < 10000; x++ {
			batch[i][x] = []byte(strconv.Itoa(x))
		}
	}

	wg.Add(w)

	for i := 0; i < w; i++ {
		go func(b int) {
			for x := range batch[b] {
				r.Insert(batch[b][x], Bytes(batch[b][x]))
			}
			wg.Done()
		}(i)
	}

	wg.Wait()

	for x := 0; x < 1000; x++ {
		value := r.Lookup(batch[0][x])
		require.NotNil(t, value)
		assert.True(t, bytes.Equal(value.(Bytes), batch[0][x]))
	}
}

func TestSwap(t *testing.T) {
	uuids := make([][]byte, 10000)

	for i := 0; i < 10000; i++ {
		uuids[i] = []byte(uuid.New().String())
	}

	r := New()

	// swap empty
	for x := 0; x < 10000; x++ {
		success := r.Swap(uuids[x], nil, Bytes(uuids[x]))
		require.True(t, success)
	}

	v := []byte("new-value")

	// swap existing
	for x := 0; x < 10000; x++ {
		success := r.Swap(uuids[x], Bytes(uuids[x]), Bytes(v))
		require.True(t, success)
	}

	// test swapping a value that has been inserted, but its value is nil
	r.Insert([]byte("hello"), nil)
	success := r.Swap([]byte("hello"), nil, String("HELLO"))
	require.True(t, success)
}

func TestConcurrentSwap(t *testing.T) {
	for x := 0; x < 100; x++ {
		var wg sync.WaitGroup
		var failures int64

		w := 32

		r := New()

		wg.Add(w)

		for i := 0; i < w; i++ {
			go func(b int) {
				val := fmt.Sprintf("test-value-%d", b)

				if !r.Swap([]byte("test-key"), nil, Bytes(val)) {
					atomic.AddInt64(&failures, 1)
				}

				wg.Done()
			}(i)
		}

		wg.Wait()

		assert.Equal(t, int64(w-1), failures)
	}

	for x := 0; x < 100; x++ {
		var wg sync.WaitGroup
		var failures int64

		w := 32

		r := New()

		wg.Add(w)

		r.Insert([]byte("test-key"), Bytes("test-value"))

		for i := 0; i < w; i++ {
			go func(b int) {
				val := fmt.Sprintf("test-value-%d", b)

				if !r.Swap([]byte("test-key"), Bytes("test-value"), Bytes(val)) {
					atomic.AddInt64(&failures, 1)
				}

				wg.Done()
			}(i)
		}

		wg.Wait()

		assert.Equal(t, int64(w-1), failures)
	}
}

func BenchmarkConcurrentInsert(b *testing.B) {
	ids := make([][]byte, 1000000)

	for i := 0; i < 1000000; i++ {
		ids[i] = []byte(uuid.New().String())
	}

	r := New()

	b.ResetTimer()

	var counter int64

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			i := atomic.AddInt64(&counter, 1) - 1
			r.Insert(ids[i], Bytes{})
		}
	})
}

func BenchmarkConcurrentLookup(b *testing.B) {
	ids := make([][]byte, 1000000)

	r := New()

	for i := 0; i < 1000000; i++ {
		ids[i] = []byte(uuid.New().String())
		r.Insert(ids[i], Bytes{})
	}

	b.ResetTimer()

	var counter int64

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			i := atomic.AddInt64(&counter, 1) - 1
			r.Lookup(ids[i])
		}
	})
}

func b(char string) byte {
	return []byte(char)[0]
}
