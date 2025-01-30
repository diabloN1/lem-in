package main

import (
	"fmt"
	"lem-in/Algorithms"
	"lem-in/GlobVar"
	"lem-in/Helpers"
	"lem-in/Utils"
	"os"
	"sort"
)



func main() {
	args := os.Args[1:]
    if len(args) != 1 {
		fmt.Println("USAGE: go run . \"data.txt\"")
        return // Exit if no input file is provided
    }

	// Read the input file
    dataBytes, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Println("ERROR: invalid data format;", err)
		return
	}

	// Parse the input data and initialize global variables
	err = Utils.ParsingData(string(dataBytes))


	GlobVar.OriginalRooms = Helpers.CopyRoomsMap(GlobVar.Rooms)
	if err != nil {
		fmt.Println("ERROR: invalid data format;", err)
		return
	}

	Algorithms.FindValidPaths()
	GlobVar.AllValidPaths = append(GlobVar.AllValidPaths, GlobVar.ValidPaths)

	// Sort valid paths by length (shortest first)
	sort.Slice(GlobVar.ValidPaths, func(i, j int) bool {
		return len(GlobVar.ValidPaths[i]) < len(GlobVar.ValidPaths[j])
	})

	// Assign ants to the shortest path and calculate the number of turns
	shortestPathIndex := 0
	lessTurns, antsOrdred := Algorithms.OrderAnts(0)

	// Check other paths to find the one with the least number of turns
	for i := 1; i < len(GlobVar.AllValidPaths); i++ {
		turns, ants := Algorithms.OrderAnts(i)
		if turns < lessTurns {
			antsOrdred = ants
			lessTurns = turns
			shortestPathIndex = i
		}
	}

	// Print the results, including the input data and ant movements
	Utils.HandleExport(antsOrdred,lessTurns,shortestPathIndex, string(dataBytes))
	fmt.Println()
}