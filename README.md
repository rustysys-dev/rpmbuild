# rpmbuild

Simple RPM Build library on top of rpmpack to make CI/CD processing easy

:> [!WARNING]

> I will be continuing to work on this so long as I need it. It will likely be at least
> a short while until its API is fully stable.

## Example

Currently go still uses the go binary to build the application, and unless there is a good reason why not
I will likely stick with that for convenience... (we need go anyway to generate)

As you may be able to guess by now, this can give you packaging ability using only our existing go binary
and a very little bit of code.

`main.go`

```go
//go:generate go run cmd/package_app.go

...
```

`cmd/package_app.go`

```go
package main

import (
	"github.com/google/rpmpack"
	"github.com/rustysys-dev/rpmbuild"
)

var app = rpmbuild.Builder{
	BinDir:  "bin",
	DistDir: "dist",

	RPMMetaData: rpmpack.RPMMetaData{
		Name:        "your_package",
		Summary:     "your package's summary",
		Description: "your package's description",
		Version:     "1.0.0", // version used in the output filename
		Release:     "1", // revision used in the output filename
		Arch:        "x86_64", // CPU architecture used in the output filename
		Packager:    "your name <your@email.dev>",
		Licence:     "MIT",
		Compressor:  "zstd",
		Provides: []*rpmpack.Relation{{
			Name:    "some_binary_you_provide",
			Version: "1.0.0",
		}},
	},

	// Scripts should be `/bin/sh` executable contents and not a filename
	Scripts: rpmbuild.Scripts{
		PreTransact:   "",
		PostTransact:  "",
		PreInstall:    "",
		PostInstall:   "",
		PreUninstall:  "",
		PostUninstall: "",
	},

	Files: []rpmbuild.PackageFile{
		{
			Source:      "bin/your_package",
			Destination: "/usr/bin/your_package",
		},
		{
			Source:      "scripts/systemd/your_package.service",
			Destination: "/usr/lib/systemd/user/your_package.service",
		},
	},
}

func main() {
	if err := app.Build(); err != nil {
		panic(err)
	}

	if err := app.Package(); err != nil {
		panic(err)
	}
}
```
