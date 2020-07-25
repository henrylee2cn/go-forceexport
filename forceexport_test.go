package forceexport

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSymtabNamesOfActiveFunc(t *testing.T) {
	t.Log(SymtabNamesOfActiveFunc())
}

func TestTimeNow(t *testing.T) {
	var timeNowFunc func() (int64, int32)
	GetFunc(&timeNowFunc, "time.now")
	sec, nsec := timeNowFunc()
	assert.NotEqual(t, 0, sec)
	assert.NotEqual(t, 0, nsec)
	t.Logf("sec:%d, nsec:%d", sec, nsec)
}

// Note that we need to disable inlining here, or else the function won't be
// compiled into the binary. We also need to call it from the test so that the
// compiler doesn't remove it because it's unused.
//go:noinline
func addOne(x int) int {
	return x + 1
}

func TestAddOne(t *testing.T) {
	assert.Equal(t, 4, addOne(3))
	var addOneFunc func(x int) int
	err := GetFunc(&addOneFunc, "github.com/henrylee2cn/go-forceexport.addOne")
	assert.NoError(t, err)
	assert.Equal(t, 4, addOneFunc(3))
}

func TestGetSelf(t *testing.T) {
	var getFunc func(interface{}, string) error
	err := GetFunc(&getFunc, "github.com/henrylee2cn/go-forceexport.GetFunc")
	assert.NoError(t, err)
	// The two functions should share the same code pointer, so they should
	// have the same string representation.
	assert.Equal(t, fmt.Sprintf("%p", GetFunc), fmt.Sprintf("%p", getFunc))
	// if fmt.Sprintf("%p", getFunc) != fmt.Sprintf("%p", GetFunc) {
	// 	t.Fatalf("Expected ")
	// }
	// Call it again on itself!
	err = getFunc(&getFunc, "github.com/henrylee2cn/go-forceexport.GetFunc")
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("%p", GetFunc), fmt.Sprintf("%p", getFunc))
}

func TestInvalidFunc(t *testing.T) {
	var invalidFunc func()
	err := GetFunc(&invalidFunc, "invalidpackage.invalidfunction")
	assert.EqualError(t, err, "invalid function name: invalidpackage.invalidfunction")
	assert.Nil(t, invalidFunc)
}
