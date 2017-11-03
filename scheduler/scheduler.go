package scheduler

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"pergamo/pergamo/structs"
	"pergamo/scheduler/script"
	"strconv"
)

type Dispatch interface {
	Apply(alloc *structs.Alloc) error
	GetAddr() string
}

func createIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, os.ModePerm)
	}

	return nil
}

type Scheduler struct {
	Dispatcher    Dispatch
	Scripter      script.Script
	schedulerPath string
	onedataPath   string
	apiRestPort   int
}

func NewScheduler(schedulerPath, onedataPath string, apiRestPort int, dispatcher Dispatch) (*Scheduler, error) {
	err := createIfNotExists(schedulerPath)
	if err != nil {
		return nil, err
	}

	/*
		// dispatcher
		dispatch, err := NewLocalDispatcher(schedulerPath)
		if err != nil {
			return nil, err
		}
	*/

	// scripter
	scripter, err := script.NewSync()
	if err != nil {
		return nil, err
	}

	scheduler := &Scheduler{
		schedulerPath: schedulerPath,
		onedataPath:   onedataPath,
		apiRestPort:   apiRestPort,
		Dispatcher:    dispatcher,
		Scripter:      scripter,
	}

	return scheduler, nil
}

func (s *Scheduler) Apply(job *structs.Job, alloc *structs.Alloc) error {
	addrs := s.Dispatcher.GetAddr()

	// crear el archivo segun el scripter
	scriptContent, err := s.Scripter.Plan(s.onedataPath, job, alloc, addrs+":"+strconv.Itoa(s.apiRestPort))
	if err != nil {
		return fmt.Errorf("Scheduler.Apply. crear el script %s", err.Error())
	}

	// Save file
	err = ioutil.WriteFile(filepath.Join(s.schedulerPath, alloc.ScriptPath), []byte(scriptContent), 0777)
	if err != nil {
		return fmt.Errorf("Scheduler.Apply. guardar el script %s", err.Error())
	}

	// Run dispatcher
	err = s.Dispatcher.Apply(alloc)
	if err != nil {
		return fmt.Errorf("Scheduler.Apply. aplicar el script %s", err.Error())
	}

	return nil
}
