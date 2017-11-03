package script

import "pergamo/pergamo/structs"

// Script copy only shares with the volume the necessary file

type Copy struct {
}

func (c *Copy) Plan(onedata string, job *structs.Job, alloc *structs.Alloc, addrs string) (string, error) {
	return "", nil
}

func NewCopy() (Script, error) {
	c := Copy{}
	return &c, nil
}
