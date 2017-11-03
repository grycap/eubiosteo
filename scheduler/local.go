package scheduler

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"pergamo/pergamo/structs"
)

type Local struct {
	Path string
}

func (l *Local) GetAddr() string {
	return "localhost"
}

func (l *Local) Apply(alloc *structs.Alloc) error {

	fmt.Println("-- apply --")
	fmt.Println(alloc.ScriptPath)

	go func() {
		_, err := exec.Command(filepath.Join(l.Path, alloc.ScriptPath)).CombinedOutput()
		if err != nil {
			fmt.Println("-- error dispatch local --")
			fmt.Println(err)
		}
	}()

	return nil
}

func NewLocalDispatcher(path string) (Dispatch, error) {
	l := &Local{path}

	return l, nil
}
