package rpmbuild

import (
	"io"
	"os"

	"github.com/google/rpmpack"
)

type PackageFile struct {
	Source      string
	Destination string
}

func (f PackageFile) ToRPMFile() (*rpmpack.RPMFile, error) {
	file, err := os.Open(f.Source)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return &rpmpack.RPMFile{
		Name:  f.Destination,
		Body:  b,
		Mode:  uint(info.Mode()),
		MTime: uint32(info.ModTime().Unix()),
		// TODO: support more files
		Type: rpmpack.GenericFile,
	}, nil
}
