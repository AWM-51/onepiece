package Fileutils

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// GetFileStream 获取文件的字节流。
func GetFileStream(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}

// DeleteFile 删除本地临时文件。
func DeleteFile(filePath string) error {
	return os.Remove(filePath)
}

// UnzipFile 解压文件到指定目录。
func UnzipFile(zipFilePath string, destDir string) error {
	r, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		err := extractFile(f, destDir)
		if err != nil {
			return err
		}
	}

	return nil
}

func extractFile(f *zip.File, destDir string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	path := filepath.Join(destDir, f.Name)

	if f.FileInfo().IsDir() {
		os.MkdirAll(path, f.Mode())
	} else {
		os.MkdirAll(filepath.Dir(path), f.Mode())
		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, rc)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetRelativePath 获取文件相对路径。
func GetRelativePath(filePath string, basePath string) (string, error) {
	relPath, err := filepath.Rel(basePath, filePath)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(relPath, "\\", "/"), nil
}

// GetAbsolutePath 获取文件绝对路径。
func GetAbsolutePath(filePath string) (string, error) {
	return filepath.Abs(filePath)
}
