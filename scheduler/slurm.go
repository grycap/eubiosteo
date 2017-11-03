package scheduler

import (
	"log"
	"net"
	"os/exec"
	"pergamo/pergamo/structs"
)

type Slurm struct {
	Path string
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func (s *Slurm) GetAddr() string {
	return GetOutboundIP().String()
}

func (s *Slurm) Apply(alloc *structs.Alloc) error {

	script := "sbatch " + alloc.ScriptPath

	_, err := exec.Command("/bin/sh", script).CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

func NewSlurmDispatcher(path string) (Dispatch, error) {
	s := &Slurm{
		Path: path,
	}

	return s, nil
}
