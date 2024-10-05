package rpmbuild

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/google/rpmpack"
)

type Builder struct {
	rpmpack.RPMMetaData
	binDir  string
	distDir string
	files   []PackageFile
}

func (b Builder) genBinName() string {
	return b.binDir + "/" + b.Name
}

func (b Builder) genRPMName() string {
	return b.distDir + "/" + b.Name + ".rpm"
}

func (b Builder) Build() error {
	if err := os.RemoveAll("build"); err != nil {
		return err
	}
	if err := os.MkdirAll("build", os.ModePerm); err != nil {
		return err
	}

	stdout, err := exec.Command("go", "build", "-o", b.genBinName(), ".").Output()
	if err != nil {
		fmt.Println(stdout)
		return err
	}

	fmt.Println(stdout)

	return nil
}

func (b Builder) Package() error {
	if err := os.RemoveAll(b.distDir); err != nil {
		return err
	}

	if err := os.MkdirAll(b.distDir, os.ModePerm); err != nil {
		return err
	}

	out, err := os.Create(b.genRPMName())
	if err != nil {
		return err
	}
	defer out.Close()

	r, err := rpmpack.NewRPM(b.RPMMetaData)
	if err != nil {
		return err
	}

	for _, file := range b.files {
		f, err := file.ToRPMFile()
		if err != nil {
			return err
		}

		if f != nil {
			r.AddFile(*f)
		}
	}

	// TODO: need to verify before write?
	if err := r.Write(out); err != nil {
		return err
	}

	return nil
}
