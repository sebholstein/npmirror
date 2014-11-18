package storage

import (
	"io/ioutil"
	"log"
	"os"
)

type FileStorage struct {
	fileDir string
}

func (f *FileStorage) PackageInfo(pkgName string) (string, error) {
	log.Println("get package info for", pkgName)

	if _, err := os.Stat(f.fileDir + "/pkginfo/" + pkgName + ".json"); os.IsNotExist(err) {
		log.Println("package info not found for", pkgName)
		return "", ErrNotFound
	}

	content, err := ioutil.ReadFile(f.fileDir + "/pkginfo/" + pkgName + ".json")

	if err != nil {
		return "", err
	}

	log.Println("found package info file for", pkgName)
	return string(content), nil
}

func (f *FileStorage) SetPackageInfo(pkgName string, pkgInfo string) error {
	log.Println("set package info file for", pkgName)
	file, err := os.Create(f.fileDir + "/pkginfo/" + pkgName + ".json")
	if err != nil {
		log.Println("error opening file", err.Error())
	}
	defer file.Close()
	if err != nil {
		return err
	}

	_, err = file.WriteString(pkgInfo)
	file.Sync()
	if err != nil {
		log.Println("error writing package info content: " + err.Error())
	}

	return err
}

func (f *FileStorage) PackageFile(pkgName string, fileName string) (fileContent []byte, err error) {
	log.Println("get package file for", pkgName, f.fileDir+"/pkgfiles/"+pkgName+"/"+fileName)
	if _, err := os.Stat(f.fileDir + "/pkgfiles/" + pkgName + "/" + fileName); os.IsNotExist(err) {
		log.Println("no package file found for", pkgName)
		return nil, ErrNotFound
	}
	content, err := ioutil.ReadFile(f.fileDir + "/pkgfiles/" + pkgName + "/" + fileName)

	if err != nil {
		return nil, err
	}

	return content, nil
}

func (f *FileStorage) SetPackageFile(pkgName string, fileName string, fileContent []byte) error {
	log.Println("set package file for", pkgName)
	os.MkdirAll(f.fileDir+"/pkgfiles/"+pkgName, 0666)
	file, err := os.Create(f.fileDir + "/pkgfiles/" + pkgName + "/" + fileName)
	if err != nil {
		return err
	}
	file.Write(fileContent)
	defer file.Close()
	file.Sync()
	return nil
}

func (f *FileStorage) createDirs() error {
	err := os.MkdirAll(f.fileDir+"/pkginfo", 0666)
	if err != nil {
		return err
	}
	err = os.MkdirAll(f.fileDir+"/pkgfiles", 0666)
	if err != nil {
		return err
	}
	return nil
}

func NewFileStorage(fileDir string) (*FileStorage, error) {
	f := &FileStorage{
		fileDir: fileDir,
	}
	return f, f.createDirs()
}
