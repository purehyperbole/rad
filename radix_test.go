package rad

import (
	"fmt"
	"testing"

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

			fmt.Println(Graphviz(r))
		})
	}

	//r.Insert([]byte("test1234"), "test1234")
	// r.Insert([]byte("test"), "test")

	/*
		n1 := r.root.next(b("t"))
		require.NotNil(t, n1)
		assert.Equal(t, []byte("e"), n1.prefix)
	*/
}

/*
func TestRadixGraphviz(t *testing.T) {
	r := New()

	n1 := &Node{
		prefix: []byte("est"),
		value:  "test",
	}

	n2 := &Node{
		prefix: []byte("234"),
		value:  "test1234",
	}

	n1.setNext(b("1"), n2)
	r.root.setNext(b("t"), n1)

	// fmt.Println(r.Graphviz())
}

func TestRadixFindInsertionPoint(t *testing.T) {
	r := New()

	n1 := &Node{
		prefix: []byte("e"),
	}

	n2 := &Node{
		prefix: []byte("m"),
		value:  "team",
	}

	n3 := &Node{
		prefix: []byte("t"),
		value:  "test",
	}

	r.root.setNext(b("t"), n1)
	n1.setNext(b("a"), n2)
	n1.setNext(b("s"), n3)

	parent, i, dv := r.findInsertionPoint([]byte("toaster"))
	parent.print()
	fmt.Println(i)
	fmt.Println(dv)
}
*/

func b(char string) byte {
	return []byte(char)[0]
}
