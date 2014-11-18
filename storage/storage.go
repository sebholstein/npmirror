package storage

import (
	"errors"
)

var ErrNotFound = errors.New("package not found")

type Storage interface {
	PackageInfo(pkgName string) (packageInfo string, err error)
	SetPackageInfo(pkgName string, packageInfo string) error
	PackageFile(pkgName string, fileName string) (fileContent []byte, err error)
	SetPackageFile(pkgName string, fileName string, fileContent []byte) error
}
