package rad

import (
	"fmt"
	"testing"
)

func TestRadixInsert(t *testing.T) {
	r := New()

	r.Insert([]byte("test1234"), "test1234")
	r.Insert([]byte("test"), "test")

	fmt.Println(r.Graphviz())

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
