package database

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func initDataDirIfNotExists(dataDir string) error {
	genesisFilePath := getGenesisJsonFilePath(dataDir)
	if fileExist(genesisFilePath) {
		return nil
	}

	databaseDirPath := getDatabaseDirPath(dataDir)
	if err := os.MkdirAll(databaseDirPath, os.ModePerm); err != nil {
		return err
	}

	genesisJsonFilePath := getGenesisJsonFilePath(dataDir)
	if err := writeGenesisToDisk(genesisJsonFilePath); err != nil {
		return err
	}

	blocksDbFilePath := getBlocksDbFilePath(dataDir)
	if err := writeEmptyBlocksDbToDisk(blocksDbFilePath); err != nil {
		return err
	}

	return nil
}

func fileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

func dirExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return true, err
}

func getDatabaseDirPath(dataDir string) string {
	return filepath.Join(dataDir, "database")
}

func getGenesisJsonFilePath(dataDir string) string {
	return filepath.Join(getDatabaseDirPath(dataDir), "genesis.json")
}

func getBlocksDbFilePath(dataDir string) string {
	return filepath.Join(getDatabaseDirPath(dataDir), "block.db")
}

func writeEmptyBlocksDbToDisk(path string) error {
	return ioutil.WriteFile(path, []byte(""), os.ModePerm)
}
