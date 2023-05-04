package proto

import (
	"fmt"
	"sync"

	"github.com/emcfarlane/starlarkproto"
	"github.com/jhump/protoreflect/desc/protoparse"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"google.golang.org/protobuf/reflect/protoregistry"
)

const (
	ModuleName = "proto"
)

var (
	once   sync.Once
	module starlark.StringDict
)

// 
func LoadModule() (starlark.StringDict, error) {
	once.Do(func() {
		members := starlarkproto.NewModule(protoregistry.GlobalFiles).Members
		members["register_files"] = starlark.NewBuiltin("register_files", register_files)
		module = starlark.StringDict{
			ModuleName: &starlarkstruct.Module{
				Name:    ModuleName,
				Members: members,
			},
		}
	})

	return module, nil
}

func register_files(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var filepaths []string
	for i := 0; i < args.Len(); i++ {
		if val, ok := args.Index(i).(starlark.String); ok {
			filepaths = append(filepaths, string(val))
		} else {
			return nil, fmt.Errorf("non-string type for arg %d", i)
		}
	}

	descriptors, err := protoparse.Parser{}.ParseFiles(filepaths...)
	if err != nil {
		return nil, fmt.Errorf("parsing proto files: %v", err)
	}

	for _, desc := range descriptors {
		filedesc := desc.UnwrapFile()
		_, err = protoregistry.GlobalFiles.FindFileByPath(filedesc.Path())
		if err == nil {
			continue  // proto file already registered, so skip it
		}
		if err != protoregistry.NotFound {
			return nil, err  // New proto files won't be found. Other errors are problems.
		}

		err = protoregistry.GlobalFiles.RegisterFile(filedesc)
		if err == nil {
			continue
		}
		return nil, fmt.Errorf("registering proto file descriptors: %v", err)
	}

	return starlark.None, nil
}
