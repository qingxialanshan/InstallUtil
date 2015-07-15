package myutils

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Exist(fname string) {
	_, e := os.Stat(fname)
	CheckError(e)
}

func CopyFile(src, dst string) (int64, error) {
	sf, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	sfi, _ := sf.Stat()
	mode := sfi.Mode()
	if mode.IsDir() {
		//fmt.Println("Directory")
		e := os.Rename(src, dst)
		return 2, e
	}
	defer sf.Close()
	df, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer df.Close()
	return io.Copy(df, sf)
}

func Write_To_File(value, cfile string) {

	f, err := os.OpenFile(cfile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)

	CheckError(err)
	defer f.Close()

	if _, err = f.WriteString(value); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

}

func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

func Delline(name, value, cfile string) {

	fin, e := os.Open(cfile)
	if e != nil {
		fmt.Println(e)
		os.Exit(2)
	}
	r := bufio.NewReader(fin)
	s, err := Readln(r)
	var output []string
	for err == nil {

		if strings.Contains(s, name) && strings.Contains(s, value) || s == "\n" {
			output = append(output, "")
		} else {
			output = append(output, s)
		}
		s, err = Readln(r)
	}
	defer fin.Close()
	fout, oute := os.OpenFile(".bash_tmp", os.O_RDWR|os.O_CREATE, 0660)
	if oute != nil {
		fmt.Println(oute)
		os.Exit(2)
	}
	if _, err = fout.WriteString(strings.Join(output, "\n")); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer fout.Close()

	var config_file string
	if Global_OS == "linux" {
		config_file = "/.bashrc"
	} else {
		config_file = "/.bash_profile"
	}

	_, e1 := CopyFile(".bash_tmp", Global_Home+config_file)
	CheckError(e1)
	os.Remove(".bash_tmp")
}

func Substitude(fname, ostr, nstr string) {
	//substitude the ostr(old string) to nstr(new string) in file that fname pointed
	if _, e := os.Stat(fname); e != nil {
		os.Create(fname)
	}
	inf, err := os.Open(fname)
	CheckError(err)

	fd, err := ioutil.ReadAll(inf)
	CheckError(err)
	instr := string(fd)

	new := strings.Replace(instr, ostr, nstr, 1)

	//fmt.Println(new)
	inf.Close()
	os.Remove(fname)
	tmp_file := filepath.Join(filepath.Dir(fname), ".tmp_file")
	fout, oute := os.OpenFile(tmp_file, os.O_RDWR|os.O_CREATE, 0660)
	CheckError(oute)
	if _, err = fout.WriteString(new); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer fout.Close()
	os.Chmod(fname, 0777)
}

func Readfile(fname string) string {

	inf, err := os.Open(fname)
	CheckError(err)
	defer inf.Close()

	fd, err := ioutil.ReadAll(inf)
	CheckError(err)
	instr := string(fd)
	return instr
}

func Chmod(path string) error {
	//fmt.Println(path)
	files, e1 := ioutil.ReadDir(path)
	for _, f := range files {
		if f.IsDir() {
			new_path := filepath.Join(path, f.Name())
			Chmod(new_path)
		} else {
			e := os.Chmod(filepath.Join(path, f.Name()), 0777)
			if e != nil {
				return e
			}
		}

	}
	return e1
}
