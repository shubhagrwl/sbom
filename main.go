package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type CycloneDX struct {
	Components []Component `json:"components"`
}

type Component struct {
	Name    string   `json:"name"`
	Version string   `json:"version"`
	Hashes  []Hashes `json:"hashes"`
}

type Hashes struct {
	Alg     string `json:"alg"`
	Content string `json:"content"`
}

func main() {
	rootPath := "./sbom-files"
	fileList := findSBOMFiles(rootPath)

	for _, file := range fileList {
		fmt.Printf("Processing SBOM file: %s\n", file)

		sbomChecksum := computeChecksum(file, "SHA-1")
		fmt.Printf("SBOM file checksum: %s\n", sbomChecksum)

		if !validateChecksumAPI(sbomChecksum) {
			fmt.Printf("Checksum validation failed for %s\n", file)
			continue
		}

		sbomData := readSBOMFile(file)
		verifyArtifacts(rootPath, sbomData)
	}
}

func findSBOMFiles(rootPath string) []string {
	var files []string
	filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".json") {
			files = append(files, path)
		}
		return nil
	})
	return files
}

func computeChecksum(filename, alg string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var checksum []byte
	if alg == "SHA-1" {
		hash := sha1.Sum(data)
		checksum = hash[:]
	} else {
		panic("Unsupported hash algorithm")
	}
	return hex.EncodeToString(checksum)
}

func validateChecksumAPI(checksum string) bool {
	// Fake API validation
	return true
}

func readSBOMFile(filename string) CycloneDX {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var sbomData CycloneDX
	err = json.Unmarshal(data, &sbomData)
	if err != nil {
		panic(err)
	}
	return sbomData
}

func verifyArtifacts(rootPath string, sbomData CycloneDX) {
	for _, component := range sbomData.Components {
		artifactPath := filepath.Join(rootPath, component.Name)
		artifactChecksum := computeChecksum(artifactPath, "SHA-1")

		var expectedChecksum string
		for _, hash := range component.Hashes {
			if hash.Alg == "SHA-1" {
				expectedChecksum = hash.Content
				break
			}
		}

		if artifactChecksum != expectedChecksum {
			fmt.Printf("Checksum mismatch for artifact %s\n", component.Name)
		} else {
			fmt.Printf("Checksum verified for artifact %s\n", component.Name)
		}
	}
}
