package Helpers

import "lem-in/GlobVar"

// RemovePathsLinks removes links between rooms that are part of a valid path.
// This ensures that paths do not overlap.
func RemovePathsLinks() {
	for index := 0; index < len(GlobVar.ValidPaths); index++ {
		path := GlobVar.ValidPaths[index]
		for i := 0; i < len(path)-1; i++ {
			node := GlobVar.Rooms[path[i]]
			node.Links = RemoveLink(node.Links, path[i+1])
			GlobVar.Rooms[path[i]] = node
		}
	}
}

// SaveBeforeInPath stores the previous room in the path for each room.
// This is used during backtracking to avoid revisiting rooms.
func SaveBeforeInPath() {
	lastPath := GlobVar.ValidPaths[len(GlobVar.ValidPaths)-1]
	for i := 1; i < len(lastPath)-1; i++ { // see if the link to the end should be removed
		room := GlobVar.Rooms[lastPath[i]]
		room.BeforeInPath = lastPath[i-1]
		GlobVar.Rooms[lastPath[i]] = room
	}
}


// RemoveLink removes a specific link from a room's list of links.
func RemoveLink(links []string, conflictRoom string) []string {
	for i := 0; i < len(links); i++ {
		if links[i] == conflictRoom {
			if i+1 < len(links) {
				links = append(links[:i], links[i+1:]...)
			} else {
				links = links[:i]
			}
		}
	}
	return links
}


// CopyRoomsMap creates a deep copy of the Rooms map.
// This is used to reset the state of rooms during pathfinding.
func CopyRoomsMap(original map[string]GlobVar.Room) map[string]GlobVar.Room {
	copied := make(map[string]GlobVar.Room)

	for key, room := range original {
		// Deep copy the links slice
		newLinks := make([]string, len(room.Links))
		copy(newLinks, room.Links)

		// Create a new Room struct with the copied slice
		copied[key] = GlobVar.Room{
			Links:        newLinks,
			IsChecked:    room.IsChecked,
			BeforeInPath: room.BeforeInPath,
		}
	}

	return copied
}

// ResetIsChecked resets the IsChecked flag for all rooms.
// This prepares the rooms for a new pathfinding iteration.
func ResetIsChecked() {
	for index := range GlobVar.Rooms {
		room := GlobVar.Rooms[index]
		room.IsChecked = false
		GlobVar.Rooms[index] = room
	}
}