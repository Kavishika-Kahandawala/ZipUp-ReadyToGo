package main

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Struct to hold the file name mappings from JSON
type Config struct {
	FileNameMap map[string]string `json:"fileNameMap"`
}

func main() {

	teal := "\033[1;36m"
	red := "\033[31m"
	logColor := "\033[1;34m"
	resetColor := "\033[0m"

	// Load config from JSON
	config, err := loadConfig("renameList.json")
	if err != nil {
		fmt.Printf(red+"Error loading config: %v\n"+resetColor, err)
		os.Exit(1)
	}

	// Displaying initial information
	fmt.Println(teal + "\nThis is like a public-beta version")
	fmt.Println("So, if you have any feedback, let me know. :)")
	fmt.Print("If you need to exclude specific files or folders, please place them inside a folder named ‘DNU’, stands for ‘Do Not Use’ xd.\n\n" + resetColor)
	fmt.Println(logColor + "!! Please note that once files are renamed, you can no longer revert them.")
	fmt.Println("!! So if you need them, please have a backup first :)" + resetColor)
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter the path to the main folder: ")
	scanner.Scan()
	path := scanner.Text()

	fmt.Print("Go auto? Otherwise will be in manual mode (more configurable) (y/n)")
	scanner.Scan()
	autoBoolean := strings.ToLower(scanner.Text()) == "y"

	suffix := ""
	keepOriginal := false
	needZip := true
	needRename := true

	if autoBoolean {
		currentDate := time.Now()

		// Format the date as _yyyymmdd
		formattedDate := currentDate.Format("_20060102")
		suffix = formattedDate

	} else {
		fmt.Print("Enter the suffix to append to each zip file: ")
		scanner.Scan()
		suffix = scanner.Text()

		fmt.Print("Do you want to keep the original folder names? (y/n): ")
		scanner.Scan()
		keepOriginal = strings.ToLower(scanner.Text()) == "y"

		fmt.Print("Do you want to have the zip file in each folder? (n/y): ")
		scanner.Scan()
		needZip = strings.ToLower(scanner.Text()) == "y"

		fmt.Print("Is renaming required during formatting? (n/y): ")
		scanner.Scan()
		needRename = strings.ToLower(scanner.Text()) == "y"

	}

	// Renaming confirmation
	// if needRename {
	// fmt.Println("Pausing for 3 seconds for the note...")
	// time.Sleep(3 * time.Second)

	// fmt.Print("Proceed? (y/n): ")
	// scanner.Scan()
	// proceed := strings.ToLower(scanner.Text()) == "y"

	// 	if !needRename {
	// 		fmt.Println(logColor + "Exiting..." + resetColor)
	// 		os.Exit(0)
	// 	}
	// }

	// Create "Formatted" folder
	newFolderPath := filepath.Join(path, "Formatted")
	os.MkdirAll(newFolderPath, os.ModePerm)

	// Process subfolders
	subfolders, _ := os.ReadDir(path)
	for _, subfolder := range subfolders {
		if subfolder.IsDir() && subfolder.Name() != "Formatted" && subfolder.Name() != "DNU" {
			subfolderPath := filepath.Join(path, subfolder.Name())

			subSubfolders, _ := os.ReadDir(subfolderPath)
			for _, subSubfolder := range subSubfolders {
				if subSubfolder.IsDir() {
					oldPath := filepath.Join(subfolderPath, subSubfolder.Name())
					tempPath := filepath.Join(subfolderPath, subSubfolder.Name()+suffix)
					newPath := filepath.Join(newFolderPath, subSubfolder.Name()+suffix+".zip")

					// Print progress
					fmt.Printf(logColor+"\nProcessing %s...\n", tempPath+resetColor)

					os.Rename(oldPath, tempPath)
					if needRename {
						renameAll(subfolderPath, config.FileNameMap) // Use config's fileNameMap
					}
					fmt.Printf("Zipping ...\n")
					zipFolder(tempPath, newPath)
					if keepOriginal {
						os.Rename(tempPath, oldPath)
					}

					if needZip {
						zipFileName := subSubfolder.Name() + suffix + ".zip"
						copyPath := filepath.Join(subfolderPath, zipFileName)
						err := copyFile(newPath, copyPath)
						if err != nil {
							fmt.Printf(red+"Error copying %s to %s: %v\n"+resetColor, newPath, subfolderPath, err)
						}
					}
					fmt.Printf(logColor+"Finished processing %s\n"+resetColor, tempPath)
				}
			}
		}
	}

	// Final message
	fmt.Println(logColor + "All done!")
	fmt.Print("Press 'Enter' to exit..." + resetColor)
	scanner.Scan() // Wait for Enter Key
}

// Function to load config from a JSON file
func loadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// Function to zip a folder
func zipFolder(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	filepath.WalkDir(source, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = filepath.Join(filepath.Base(source), path[len(source):])

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}
		return err
	})

	return err
}

// Function to copy a file from source to destination
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	err = destFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

// Function to rename files based on the loaded fileNameMap
func renameAll(dirPath string, fileNameMap map[string]string) {
	red := "\033[31m"
	resetColor := "\033[0m"

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf(red+"Error visiting path %v: %v\n"+resetColor, path, err)
			return err
		}

		if !info.IsDir() {
			oldName := info.Name()
			for oldPrefix, newName := range fileNameMap {
				if strings.HasPrefix(oldName, oldPrefix) {
					oldPath := path
					newPath := filepath.Join(filepath.Dir(path), newName)
					err := os.Rename(oldPath, newPath)
					if err != nil {
						fmt.Printf(red+"Error renaming %s to %s: %v\n"+resetColor, oldPath, newPath, err)
					} else {
						fmt.Printf("Successfully renamed %s to %s\n", oldPath, newPath)
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf(red+"Error walking the path %v: %v\n"+resetColor, dirPath, err)
	}
}
