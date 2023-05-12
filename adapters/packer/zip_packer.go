package packer

import (
	"archive/zip"
	"github.com/ClickerAI/ClickerAI/core/ports"
	"io"
	"os"
	"path/filepath"
)

type ZipPacker struct{}

func (p *ZipPacker) Pack(session ports.LoggingSession) (string, error) {
	// Get the directory path for the logging session
	dirPath := session.GetDirPath()

	// Open the directory
	dir, err := os.Open(dirPath)
	if err != nil {
		return "", err
	}
	defer dir.Close()

	// Get the file info for the directory
	info, err := dir.Stat()
	if err != nil {
		return "", err
	}

	// Create a zip file with the same name as the directory
	zipPath := filepath.Join(filepath.Dir(dirPath), info.Name()+".zip")
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return "", err
	}
	defer zipFile.Close()

	// Create a zip writer and add the files in the directory to the archive
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create a new file header for the file
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Set the name for the file to be relative to the directory
		header.Name, err = filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}

		// If the file is a directory, don't add it to the archive
		if info.IsDir() {
			return nil
		}

		// Create a new file in the archive and copy the contents of the original file
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})

	if err != nil {
		return "", err
	}

	// Remove the original directory
	err = os.RemoveAll(dirPath)
	if err != nil {
		return "", err
	}

	// Return the path to the zip file
	return zipPath, nil
}
