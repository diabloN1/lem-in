package Algorithms

import (
	"fmt"
	"lem-in/GlobVar"
	"lem-in/Helpers"
	"os"
)

// OrderAnts assigns ants to the shortest path and calculates the number of turns required.
// It returns the number of turns and the ordered list of ants.
func OrderAnts(indexValidPaths int) (int, []int) {

	foundPath := GlobVar.AllValidPaths[indexValidPaths]
	antsN := GlobVar.AntsNumber
	// Order Ants
	ants := make([]int, len(foundPath))
	indexShortestPath := 0
	shortestPathLen := 0

	for antsN > 0 {
		shortestPathLen = len(foundPath[indexShortestPath]) + ants[indexShortestPath]

		if len(foundPath) > 1 && indexShortestPath+1 < len(foundPath) && len(foundPath[indexShortestPath+1])+ants[indexShortestPath+1] < shortestPathLen {
			shortestPathLen = len(foundPath[indexShortestPath+1]) + ants[indexShortestPath+1]
			indexShortestPath++
		}

		ants[indexShortestPath]++
		antsN--

		if indexShortestPath == len(GlobVar.AllValidPaths[indexValidPaths])-1 {
            indexShortestPath = 0
        }
	}

	return shortestPathLen - 1, ants
}

// FindValidPaths uses BFS to find all valid paths from the start room to the end room.
// It handles backtracking and ensures that paths do not overlap.
func FindValidPaths() {

	linksToRemove := make(map[string][]string)

	for {
		hasFoundAny := BFS()

		if !hasFoundAny {
			break
		}
		Helpers.SaveBeforeInPath()

		isBackTracking, revNode, toRemove := CheckIfBackTrackingPath()

		if isBackTracking {

			linksToRemove[revNode] = append(linksToRemove[revNode], toRemove)
			GlobVar.AllValidPaths = append(GlobVar.AllValidPaths, GlobVar.ValidPaths[:len(GlobVar.ValidPaths)-1])

			GlobVar.ValidPaths = [][]string{}

			GlobVar.Rooms = Helpers.CopyRoomsMap(GlobVar.OriginalRooms)

			for rev, links := range linksToRemove {
				for _, toRm := range links {
					room := GlobVar.Rooms[rev]
					room.Links = Helpers.RemoveLink(room.Links, toRm)
					GlobVar.Rooms[rev] = room

					room2 := GlobVar.Rooms[toRm]
					room2.Links = Helpers.RemoveLink(room2.Links, rev)
					GlobVar.Rooms[toRm] = room2
				}
			}
		} else if GlobVar.AntsNumber == len(GlobVar.ValidPaths) {
			return
		} else {
			Helpers.RemovePathsLinks()
		}
	}

	if len(GlobVar.ValidPaths) == 0 {
		fmt.Println("ERROR: invalid data format; No valid paths found!")
		os.Exit(0)
	}
}

// CheckIfBackTrackingPath checks if the last found path is backtracking over an existing path.
// It returns the rooms involved in the backtracking.
func CheckIfBackTrackingPath() (bool, string, string) {
	lastPath := GlobVar.ValidPaths[len(GlobVar.ValidPaths)-1]
	pathRooms := GlobVar.ValidPaths[:len(GlobVar.ValidPaths)-1]
	links := make(map[string]string)
	// get all path links reversed
	for i := len(pathRooms) - 1; i >= 0; i-- {

		for j := len(pathRooms[i]) - 2; j >= 1; j-- {
			links[pathRooms[i][j]] = pathRooms[i][j-1]
		}
	}

	for i := 1; i < len(lastPath)-1; i++ {
		if links[lastPath[i]] == lastPath[i+1] {
			return true, lastPath[i], lastPath[i+1]
		}
	}

	return false, "", ""
}


// BFS performs a breadth-first search to find a valid path from the start room to the end room.
// It handles backtracking and ensures that rooms are not revisited.
func BFS() bool {

	startRoom := GlobVar.Rooms[GlobVar.Start]
	startRoom.IsChecked = true
	GlobVar.Rooms[GlobVar.Start] = startRoom

	alreadyInRevesedPath := false

	paths := [][]string{{GlobVar.Start}}
	for len(paths) != 0 {
		if len(GlobVar.Rooms[GlobVar.Start].Links) == 0 {
			return false
		}

		for i := 0; i < len(paths); i++ {
			validLinks := 0

			lastInPath := paths[i][len(paths[i])-1]

			if GlobVar.Rooms[lastInPath].BeforeInPath != "" && !alreadyInRevesedPath {
				beforeInPath := GlobVar.Rooms[lastInPath].BeforeInPath
				validLinks++

				room := GlobVar.Rooms[lastInPath]
				room.IsChecked = false
				GlobVar.Rooms[lastInPath] = room

				paths[i] = append(paths[i], beforeInPath)

				alreadyInRevesedPath = true

				continue
			}

			for j, link := range GlobVar.Rooms[lastInPath].Links {

				if link == GlobVar.End {

					if validLinks == 0 {
						paths[i] = append(paths[i], link)
					} else {
						paths = append(paths, append(paths[i][:len(paths[i])-1], link))
					}
					GlobVar.ValidPaths = append(GlobVar.ValidPaths, paths[i])
					Helpers.ResetIsChecked()
					return true
				}

				if !GlobVar.Rooms[link].IsChecked {
					room := GlobVar.Rooms[link]
					validLinks++
					room.IsChecked = true
					GlobVar.Rooms[link] = room
					if validLinks == 1 {
						paths[i] = append(paths[i], link)
					} else {
						var path = make([]string, len(paths[i]))
						copy(path, paths[i])
						paths = append(paths, append(path[:len(path)-1], link))
					}

				} else if j == len(GlobVar.Rooms[lastInPath].Links)-1 && validLinks == 0 {
					if i+1 < len(paths) {
						paths = append(paths[:i], paths[i+1:]...)
					} else {
						paths = paths[:i]
					}

				}
			}
		}
	}
	return false
}

