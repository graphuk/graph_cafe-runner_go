package utils

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func fixFilenameSlashes(filename string) string {
	filename = strings.Replace(filename, `/`, string(os.PathSeparator), -1)
	filename = strings.Replace(filename, `\`, string(os.PathSeparator), -1)
	return filename
}

func Unzip(src, dest string) error {
	src = fixFilenameSlashes(src)
	dest = fixFilenameSlashes(dest)

	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(dest, f.Name)
		path = fixFilenameSlashes(path)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			f, err := os.OpenFile(
				path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func UnzipWithoutRootFolder(src, dest string) error {
	src = fixFilenameSlashes(src)
	dest = fixFilenameSlashes(dest)

	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		filename := fixFilenameSlashes(f.Name)
		log.Println(filename)
		firstSlashOffset := strings.Index(filename, string(os.PathSeparator))
		if firstSlashOffset != -1 {
			filename = filename[firstSlashOffset+1:]
		}

		path := filepath.Join(dest, filename)
		path = fixFilenameSlashes(path)

		if f.FileInfo().IsDir() {
			log.Println(`folder: ` + path)
			os.MkdirAll(path, f.Mode())
		} else {
			log.Println(`file: ` + path)
			f, err := os.OpenFile(
				path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func Zipit(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}
