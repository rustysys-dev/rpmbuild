package rpmbuild

// Scripts are shell (/bin/sh) executable commands and not a filename
//
// TODO: add post processing if necessary...
type Scripts struct {
	PreTransact   string
	PostTransact  string
	PreInstall    string
	PostInstall   string
	PreUninstall  string
	PostUninstall string
}
