package proto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"tidbyt.dev/pixlet/runtime"
)

func TestLoadModule(t *testing.T) {
	src := `
load("encoding/proto.star", "proto")
load("render.star", "render")

def main():
    return render.Root(child=render.Box())
`
	app := &runtime.Applet{}
	err := app.Load("test.star", []byte(src), nil)
	assert.NoError(t, err)
	roots, err := app.Run(map[string]string{})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(roots))
}

func TestRegisterFile(t *testing.T) {
	src := `
load("encoding/proto.star", "proto")

def main():
    proto.register_files("testdata/test.proto")  # relative to test file
    return []
`
	app := &runtime.Applet{}
	err := app.Load("test.star", []byte(src), nil)
	assert.NoError(t, err)
	_, err = app.Run(map[string]string{})
	assert.NoError(t, err)
}

func TestRegisterFileNotFound(t *testing.T) {
	src := `
load("encoding/proto.star", "proto")

def main():
    proto.register_files("does_not_exist.proto")
    return []
`
	app := &runtime.Applet{}
	err := app.Load("test.star", []byte(src), nil)
	assert.NoError(t, err)
	_, err = app.Run(map[string]string{})
	wantErr := "no such file or directory"
	assert.Containsf(t, err.Error(), wantErr, "expected error containing %q, got %s", wantErr, err)
}

func TestCannotRegisterNonStringInputElements(t *testing.T) {
	src := `
load("encoding/proto.star", "proto")

def main():
    proto.register_files(123)
    return []
`
	app := &runtime.Applet{}
	err := app.Load("test.star", []byte(src), nil)
	assert.NoError(t, err)
	_, err = app.Run(map[string]string{})
	wantErr := "non-string type"
	assert.Containsf(t, err.Error(), wantErr, "expected error containing %q, got %s", wantErr, err)
}

func TestRegisterProtoTwiceAtOnceIsOK(t *testing.T) {
	src := `
load("encoding/proto.star", "proto")

def main():
    proto.register_files("testdata/test.proto", "testdata/test.proto")
    return []
`
	app := &runtime.Applet{}
	err := app.Load("test.star", []byte(src), nil)
	assert.NoError(t, err)
	_, err = app.Run(map[string]string{})
	assert.NoError(t, err)
}

func TestRegisterProtoTwiceSeparatelyIsOK(t *testing.T) {
	src := `
load("encoding/proto.star", "proto")

def main():
    proto.register_files("testdata/test.proto")
    proto.register_files("testdata/test.proto")
    return []
`
	app := &runtime.Applet{}
	err := app.Load("test.star", []byte(src), nil)
	assert.NoError(t, err)
	_, err = app.Run(map[string]string{})
	assert.NoError(t, err)
}

func TestCannotRegisterTwoConflictingFiles(t *testing.T) {
	src := `
load("encoding/proto.star", "proto")

def main():
    proto.register_files("testdata/test.proto", "testdata/copy.proto")
    return []
`
	app := &runtime.Applet{}
	err := app.Load("test.star", []byte(src), nil)
	assert.NoError(t, err)
	_, err = app.Run(map[string]string{})
	wantErr := "symbol \"foo.TestMessage\" already defined"
	assert.Containsf(t, err.Error(), wantErr, "expected error containing %q, got %s", wantErr, err)
}

func TestCannotRegisterNonProtoFile(t *testing.T) {
	src := `
load("encoding/proto.star", "proto")

def main():
    proto.register_files("proto_test.go")
    return []
`
	app := &runtime.Applet{}
	err := app.Load("test.star", []byte(src), nil)
	assert.NoError(t, err)
	_, err = app.Run(map[string]string{})
	wantErr := "syntax error"
	assert.Containsf(t, err.Error(), wantErr, "expected error containing %q, got %s", wantErr, err)
}

func TestCreateProtoByName(t *testing.T) {
	src := `
load("encoding/proto.star", "proto")
load("math.star", "math")

def main():
    proto.register_files("testdata/test.proto")
    m = proto.new("foo.TestMessage")(bar="abc",baz=4,qux=[3.14,2.71])
    
    if m.bar != "abc":
        fail("wrong bar field: ", m.bar)
    if m.baz != 4:
        fail("wrong baz field: ", m.baz)
    epsilon = 0.000001
    if math.fabs(m.qux[0] - 3.14) > epsilon:
        fail("wrong qux[0] value:", m.qux[0])
    if math.fabs(m.qux[1] - 2.71) > epsilon:
        fail("wrong qux[1] value: ", m.qux[1])
    return []
`
	app := &runtime.Applet{}
	err := app.Load("test.star", []byte(src), nil)
	assert.NoError(t, err)
	_, err = app.Run(map[string]string{})
	assert.NoError(t, err)
}

func TestRegisterTwoProtoFiles(t *testing.T) {
    src := `
load("encoding/proto.star", "proto")

def main():
    proto.register_files("testdata/test.proto", "testdata/other.proto")
    m = proto.new("foo.TestMessage")(bar="abc")
    n = proto.new("zzz.Different")(text="houpla")
    
    if m.bar != "abc":
        fail("wrong bar field: ", m.bar)
    if n.text != "houpla":
        fail("wrong text field: ", n.bar)
    
    return []
`
    app := &runtime.Applet{}
    err := app.Load("test.star", []byte(src), nil)
    assert.NoError(t, err)
    _, err = app.Run(map[string]string{})
    assert.NoError(t, err)
}

func TestCreateProtoByFile(t *testing.T) {
	src := `
load("encoding/proto.star", "proto")
load("math.star", "math")

def main():
    proto.register_files("testdata/test.proto")
    f = proto.file("testdata/test.proto")
    m = f.TestMessage(bar="abc",baz=4,qux=[3.14,2.71])

    if m.bar != "abc":
        fail("wrong bar field: ", m.bar)
    if m.baz != 4:
        fail("wrong baz field: ", m.baz)
    epsilon = 0.000001
    if math.fabs(m.qux[0] - 3.14) > epsilon:
        fail("wrong qux[0] value:", m.qux[0])
    if math.fabs(m.qux[1] - 2.71) > epsilon:
        fail("wrong qux[1] value: ", m.qux[1])
    return []
`
	app := &runtime.Applet{}
	err := app.Load("test.star", []byte(src), nil)
	assert.NoError(t, err)
	_, err = app.Run(map[string]string{})
	assert.NoError(t, err)
}

func TestReadWireProto(t *testing.T) {
	src := `
load("encoding/base64.star", "base64")
load("encoding/proto.star", "proto")
load("math.star", "math")

def main():
    proto.register_files("testdata/test.proto")
    f = proto.file("testdata/test.proto")
    m = proto.new("foo.TestMessage")()

    proto.unmarshal(base64.decode("CgNhYmMQBB3D9UhAHaRwLUA="), m)

    if m.bar != "abc":
        fail("wrong bar field: ", m.bar)
    if m.baz != 4:
        fail("wrong baz field: ", m.baz)
    epsilon = 0.000001
    if math.fabs(m.qux[0] - 3.14) > epsilon:
        fail("wrong qux[0] value:", m.qux[0])
    if math.fabs(m.qux[1] - 2.71) > epsilon:
        fail("wrong qux[1] value: ", m.qux[1])
    return []
`
	app := &runtime.Applet{}
	err := app.Load("test.star", []byte(src), nil)
	assert.NoError(t, err)
	_, err = app.Run(map[string]string{})
	assert.NoError(t, err)
}

func TestRoundTripWireProto(t *testing.T) {
	src := `
load("encoding/base64.star", "base64")
load("encoding/proto.star", "proto")
load("math.star", "math")

def main():
    proto.register_files("testdata/test.proto")
    f = proto.file("testdata/test.proto")
    orig = proto.new("foo.TestMessage")(bar="abc",baz=4,qux=[3.14,2.71])
    res = proto.new("foo.TestMessage")()
    proto.unmarshal(proto.marshal(orig), res)

    if res != orig:
        fail("Cannot round trip successfully")
    return []
`
	app := &runtime.Applet{}
	err := app.Load("test.star", []byte(src), nil)
	assert.NoError(t, err)
	_, err = app.Run(map[string]string{})
	assert.NoError(t, err)
}
