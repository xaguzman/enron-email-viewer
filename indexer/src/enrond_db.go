package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const enron_db_url = "http://www.cs.cmu.edu/~enron/enron_mail_20110402.tgz"

// download
func downloadEnronDb(filepath string) error {

	out, err := os.Create(filepath + ".tmp") // tmp file to download, rename later
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(enron_db_url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	// rename tmp file
	err = os.Rename(filepath+".tmp", filepath)
	return err
}

// untarGz decompresses a .tar.gz file and extracts it to the specified destination directory.
func untarGz(src, dest string) error {
	gzFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer gzFile.Close()

	// gzip reader
	gzReader, err := gzip.NewReader(gzFile)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	// tar reader
	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // End of tar archive
		}

		if err != nil {
			return err
		}

		// destination file path
		path := filepath.Join(dest, header.Name)

		switch header.Typeflag {

		case tar.TypeDir: // if it's a directory, create it
			if err := os.MkdirAll(path, 0755); err != nil {
				return err
			}

		case tar.TypeReg: // if it's a file, write it
			outFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		}
	}
	return nil
}
