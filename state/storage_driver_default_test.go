//go:build !js && !wasm
// +build !js,!wasm

package state

import (
	"fmt"
	"os"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || !os.IsNotExist(err)
}

var _ = Describe("StorageDriverDefault", func() {
	var storageDriver StorageDriver
	var testByte []byte

	BeforeEach(func() {
		// 一意のファイル名を生成
		testDescription := CurrentSpecReport()
		testFileName := fmt.Sprintf("%s.json", strings.ReplaceAll(testDescription.FullText(), " ", "_"))
		storageDriver = NewStorageDriver(testFileName)
		testByte = []byte("test")
	})

	AfterEach(func() {
		// テスト終了後にファイルを削除
		testDescription := CurrentSpecReport()
		testFileName := fmt.Sprintf("%s.json", strings.ReplaceAll(testDescription.FullText(), " ", "_"))
		if FileExists(testFileName) {
			err := os.Remove(testFileName)
			Expect(err).ToNot(HaveOccurred())
		}
	})

	Describe("SaveData and LoadData", func() {
		It("should save and load the game state correctly", func() {
			// Save the game state
			err := storageDriver.SaveData(testByte)
			Expect(err).ToNot(HaveOccurred())

			// Load the game state
			loadedByte, err := storageDriver.LoadData()
			Expect(err).ToNot(HaveOccurred())
			Expect(loadedByte).To(Equal(testByte))
		})

		It("should return an empty state if the save file does not exist", func() {
			// Load the game state
			loadedByte, err := storageDriver.LoadData()
			Expect(err).To(HaveOccurred())
			Expect(loadedByte).To(Equal([]byte{}))
		})

		It("should overwrite the existing save file when saving new data", func() {
			// Save initial data
			err := storageDriver.SaveData(testByte)
			Expect(err).ToNot(HaveOccurred())

			// Save new data
			newTestByte := []byte("new test")
			err = storageDriver.SaveData(newTestByte)
			Expect(err).ToNot(HaveOccurred())

			// Load the game state
			loadedByte, err := storageDriver.LoadData()
			Expect(err).ToNot(HaveOccurred())
			Expect(loadedByte).To(Equal(newTestByte))
		})

		It("should handle saving an empty byte slice", func() {
			// Save empty data
			err := storageDriver.SaveData([]byte{})
			Expect(err).ToNot(HaveOccurred())

			// Load the game state
			loadedByte, err := storageDriver.LoadData()
			Expect(err).ToNot(HaveOccurred())
			Expect(loadedByte).To(Equal([]byte{}))
		})

		It("should return an error when saving fails", func() {
			storageDriver = NewStorageDriver("invalid_path/test_save_file.json")
			err := storageDriver.SaveData([]byte("test"))
			Expect(err).To(HaveOccurred())
		})
	})
})
