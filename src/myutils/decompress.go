package myutils

import (
	"archive/tar"
	"archive/zip"
	"compress/bzip2"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func Rmdir(dirpath string) {
	var e error
	if Global_OS == "windows" {
		dirpath = strings.Replace(dirpath, "\\\\", "\\", -1)
		_, e = exec.Command("cmd", "/c", "del", "/f", "/s", "/q", dirpath).Output()
	} else {
		_, e = exec.Command("bash", "-c", "rm -rf "+dirpath).Output()
	}
	if e != nil {
		fmt.Println("Failed to remove directory :", dirpath)

		os.Exit(2)
	}
	os.RemoveAll(dirpath)
}

// Ungzip and untar from source file to destination directory
// you need check file exist before you call this function
func UnTarGz(srcFilePath string, destDirPath string) error {
	//fmt.Println("UnTarGzing " + srcFilePath + "...")
	// Create destination directory
	os.Mkdir(destDirPath, os.ModePerm)
	var tr *tar.Reader
	fr, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer fr.Close()

	if strings.HasSuffix(srcFilePath, ".bz2") {
		br := bzip2.NewReader(fr)
		tr = tar.NewReader(br)
	} else {
		// Gzip reader
		gr, err := gzip.NewReader(fr)
		if err != nil {
			return err
		}
		defer gr.Close()
		// Tar reader
		tr = tar.NewReader(gr)
	}

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// End of tar archive
			break
		}
		//handleError(err)
		//fmt.Println("UnTarGzing file..." + hdr.Name)
		// Check if it is diretory or file
		if hdr.Typeflag != tar.TypeDir {
			// Get files from archive
			// Create diretory before create file
			os.MkdirAll(destDirPath+"/"+path.Dir(hdr.Name), os.ModePerm)
			// Write data to file
			fw, _ := os.Create(destDirPath + "/" + hdr.Name)

			os.Chmod(destDirPath+"/"+hdr.Name, os.FileMode(hdr.Mode))

			if err != nil {
				return err
			}
			_, err = io.Copy(fw, tr)
			if err != nil {
				return err
			}
		}
	}
	//fmt.Println("Well done!")
	return nil
}

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		fmt.Println("error1")
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()

		if err != nil {
			fmt.Println("error2")
			return err
		}
		defer rc.Close()

		fpath := filepath.Join(dest, f.Name)

		var fmode os.FileMode
		fmode = f.Mode()
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
		} else {
			var fdir string
			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}

			err = os.MkdirAll(fdir, 0777)
			if err != nil {
				log.Fatal(err)
				return err
			}
			f, err := os.OpenFile(
				fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
			if err != nil {
				//fmt.Println("###", fpath)
				return err
			}

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}

			os.Chmod(f.Name(), fmode)
			f.Close()
		}
	}
	return nil
}

func Decompress(src, dest string) {
	Exist(src)
	e := errors.New("unknow file type")
	if strings.HasSuffix(src, ".bz2") || strings.HasSuffix(src, ".gz") || strings.HasSuffix(src, ".tgz") {
		//fmt.Println(".gz detected")
		e = UnTarGz(src, dest)
	} else if strings.HasSuffix(src, ".zip") {
		e = Unzip(src, dest)

	} else {
		fmt.Println("unknow file type")
		os.Exit(2)
	}
	CheckError(e)
}
