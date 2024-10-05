package rpmbuild

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/rpmpack"
	"github.com/jesseduffield/go-git/v5"
)

type Builder struct {
	rpmpack.RPMMetaData
	BinDir  string
	DistDir string
	Files   []PackageFile
}

func (b *Builder) genBinName() string {
	return b.BinDir + "/" + b.Name
}

func (b *Builder) genRPMName() (string, error) {
	if b.DistDir == "" {
		b.DistDir = "dist"
	}

	if b.Name == "" {
		b.SetNameFromRepo()
		if b.Name == "" {
			return "", errors.New("unable to find a suitable name, please add to config")
		}
	}

	if b.Version == "" {
		b.Version = "0.0.1"
	}

	if b.Release == "" {
		b.Release = "1"
	}

	if b.Arch == "" {
		b.Arch = "noarch"
	}

	return strings.Join([]string{b.DistDir, "/", b.Name, "-", b.Version, "-", b.Release, ".", b.Arch, ".rpm"}, ""), nil
}

func (b *Builder) SetNameFromRepo() error {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return err
	}
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}
	b.Name = filepath.Base(worktree.Filesystem.Root())
	return nil
}

func (b *Builder) Build() error {
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

func (b *Builder) Package() error {
	if err := os.RemoveAll(b.DistDir); err != nil {
		return err
	}

	if err := os.MkdirAll(b.DistDir, os.ModePerm); err != nil {
		return err
	}

	rpmName, err := b.genRPMName()
	if err != nil {
		return err
	}

	out, err := os.Create(rpmName)
	if err != nil {
		return err
	}
	defer out.Close()

	r, err := rpmpack.NewRPM(b.RPMMetaData)
	if err != nil {
		return err
	}

	for _, file := range b.Files {
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
