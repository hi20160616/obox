package data

import (
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/alexmullins/zip"
)

func ZipFiles(filesPath, zipPath, passwd string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipw := zip.NewWriter(zipFile)
	defer zipw.Close()
	var w io.Writer
	return filepath.WalkDir(filesPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Type().IsRegular() {
			r, err := os.Open(path)
			if err != nil {
				return err
			}
			defer r.Close()
			if passwd == "" {
				w, err = zipw.Create(path)
			} else {
				w, err = zipw.Encrypt(path, passwd)
			}
			if err != nil {
				return err
			}

			_, err = io.Copy(w, r)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// ZipWriter zip files in inPath to outPath, notice: outPath contains zip file name
func ZipFiles2(inPath, outPath, passwd string) error {
	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	w := zip.NewWriter(outFile)
	defer w.Close()

	// walk folder and add stuffs
	return addFiles(w, inPath, "", passwd)
}

func addFiles(w *zip.Writer, basePath, baseInZip, passwd string) error {
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		baseName := filepath.Join(basePath, file.Name())
		baseInZipName := filepath.Join(baseInZip, file.Name())
		if !file.IsDir() {
			dat, err := ioutil.ReadFile(baseName)
			if err != nil {
				return err
			}

			// Add some files to the archive
			var fw io.Writer
			// add basePath as sub directory in zip
			baseInZipName = filepath.Join(basePath, baseInZipName)
			if passwd == "" {
				fw, err = w.Create(baseInZipName)
			} else {
				fw, err = w.Encrypt(baseInZipName, passwd)
			}
			if err != nil {
				return err
			}
			_, err = fw.Write(dat)
			if err != nil {
				return err
			}
		} else if file.IsDir() {
			// Recurse
			newBase := baseName + "/"
			addFiles(w, newBase, filepath.Join(baseInZip, file.Name())+"/", passwd)
		}
	}
	return nil
}

func UnzipFiles(inPath, outPath, passwd string) error {
	r, err := zip.OpenReader(inPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.IsEncrypted() {
			f.SetPassword(passwd)
		}
		dir, _ := filepath.Split(f.Name)
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
		w, err := os.Create(f.Name)
		if err != nil {
			return err
		}
		defer w.Close()

		fr, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(w, fr)
		if err != nil {
			return err
		}
	}
	return nil
}
