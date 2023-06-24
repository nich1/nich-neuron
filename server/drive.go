package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

const driveLocation = "database/drive/"

type folder struct {
	Name    string   `json:"name"`
	Files   []string `json:"files"`
	Folders []folder `json:"folders"`
}

// Recursive function for getDrive
func buildFileData(root folder, path string) folder {
	// Get list of all files and directories in root
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < len(entries); i++ {
		// Is a directory
		if entries[i].IsDir() {
			var newFolder folder
			newFolder.Name = entries[i].Name()
			newFolder = buildFileData(newFolder, path+newFolder.Name+"/")
			root.Folders = append(root.Folders, newFolder)
		} else { // Is a file
			root.Files = append(root.Files, entries[i].Name())
		}

	}
	return root
}

func getDrive(context *gin.Context) {
	// Populate struct
	var root folder
	root.Name = "drive"
	root = buildFileData(root, driveLocation)

	// Return json
	context.JSON(http.StatusOK, root)

}

func postDrive(context *gin.Context) {
	// single file
	file, err := context.FormFile("file")
	if err != nil {
		fmt.Println(err)
	}
	log.Println(file.Filename)

	location := context.GetHeader("location")

	if strings.Contains(location, "..") {
		fmt.Println("Unsafe path")
		return
	}

	// Make Any New Necessary Directories
	dirs := strings.Split(location, "/")
	fullpath := driveLocation
	for i := 0; i < len(dirs); i++ {
		fullpath = fullpath + dirs[i] + "/"
		// Check if valid path
		if _, err := os.Stat(fullpath); err != nil {
			if os.IsNotExist(err) {
				// Dir does not exist, create dir
				err := os.Mkdir((fullpath), 0755)
				if err != nil {
					fmt.Println("Error in making directory")
				}
			}
		}
	}

	// Upload the file to specific driveLocation.
	err = context.SaveUploadedFile(file, driveLocation+location+"/"+file.Filename)
	if err != nil {
		fmt.Println(err)
	}

	context.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}

func deleteDrive(context *gin.Context) {
	// Pass in string which is the path to what is to delete
	// Convert rawJson to calendarEvent struct
	var deletePath string
	if err := context.BindJSON(&deletePath); err != nil {
		context.IndentedJSON(http.StatusConflict, "Error occured during file conversion")
		fmt.Println(err)
		return
	}
	if strings.Contains(deletePath, "..") || deletePath == "" {
		fmt.Println("Unsafe path")
		return
	}
	fmt.Println(1)
	fmt.Println(driveLocation + deletePath)

	err := os.RemoveAll(driveLocation + deletePath)
	if err != nil {
		fmt.Println(err)
	}

}
