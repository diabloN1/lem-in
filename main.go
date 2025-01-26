package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)	

type Room struct {
	links []string
	isChecked bool
	beforeInPath string
}

var (
	antsNumber int
	originalRooms = make(map[string]Room)
	Rooms = make(map[string]Room)
	start string
	end string
	validPaths [][]string
	allValidPaths [][][]string
)

func main() {
	args := os.Args[1:]
    if len(args) != 1 {
        return
    }

    dataBytes, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Println("ERROR: invalid data format;", err)
		return
	}

	err = ParsingData(string(dataBytes))


	originalRooms = copyRoomsMap(Rooms)
	if err != nil {
		fmt.Println("ERROR: invalid data format;", err)
		return
	}

	FindValidPaths()
	allValidPaths = append(allValidPaths, validPaths)

	// Sorting strings by length
	sort.Slice(validPaths, func(i, j int) bool {
		return len(validPaths[i]) < len(validPaths[j])
	})

	shortestPathIndex := 0
	lessTurns, antsOrdred := orderAnts(0)

	for i := 1; i < len(allValidPaths); i++ {
		turns, ants := orderAnts(i)
		if turns < lessTurns {
			antsOrdred = ants
			lessTurns = turns
			shortestPathIndex = i
		}
	}

	HandleExport(antsOrdred,lessTurns,shortestPathIndex, string(dataBytes))
	fmt.Println()
}

func orderAnts(indexValidPaths int) (int, []int) {

	foundPath := allValidPaths[indexValidPaths]
	antsN := antsNumber
	// Order Ants
	ants := make([]int, len(foundPath))
	indexShortestPath := 0
	shortestPathLen := 0

	for antsN > 0 {
		shortestPathLen = len(foundPath[indexShortestPath])+ants[indexShortestPath] // Should set to make int64 value because it is max value that can be reeturned by len()


		if len(foundPath) > 1 && indexShortestPath+1 < len(foundPath) && len(foundPath[indexShortestPath+1]) + ants[indexShortestPath+1] < shortestPathLen {
			shortestPathLen = len(foundPath[indexShortestPath+1]) + ants[indexShortestPath+1]
			indexShortestPath++
		}

		ants[indexShortestPath]++
		antsN--

		if indexShortestPath == len(validPaths)-1 {
			indexShortestPath = 0
		}
	}


	return shortestPathLen-1, ants
}

func FindValidPaths() {

	linksToRemove := make(map[string][]string)

	for {
		// time.Sleep(time.Second)

		hasFoundAny := bfs()
		
		if !hasFoundAny {
			break
		}
		saveBeforeInPath()
		

		isBackTracking, revNode, toRemove := checkIfBackTrackingPath()

		if isBackTracking {
		
			linksToRemove[revNode] = append(linksToRemove[revNode], toRemove)
			allValidPaths = append(allValidPaths, validPaths[:len(validPaths)-1])
		
			validPaths = [][]string{}

			Rooms = copyRoomsMap(originalRooms)
			// real    3m14.442s
			// user    5m42.522s
			// sys     0m13.767s
			// Rooms = make(map[string]Room)
			// ParsingData(data, false)

			for rev, links := range linksToRemove {
				for _, toRm := range links {
					room := Rooms[rev]
					room.links = removeLink(room.links, toRm)
					Rooms[rev] = room
			
					
					room2 := Rooms[toRm]
					room2.links = removeLink(room2.links, rev)
					Rooms[toRm] = room2
				}
			}
		} else if antsNumber == len(validPaths) {
			return
		} else {
			removePathsLinks()
		}
	}

	if len(validPaths) == 0 {
		fmt.Println("ERROR: invalid data format; No valid paths found!")
		os.Exit(0)
	}
}

func checkIfBackTrackingPath() (bool, string, string) {
	lastPath := validPaths[len(validPaths)-1]
	pathRooms := validPaths[:len(validPaths)-1]
	links := make(map[string]string)
	// get all path links reversed
	for i := len(pathRooms)-1; i >= 0; i-- {

		for j := len(pathRooms[i])-2; j >= 1; j-- {
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

func removePathsLinks() {
	for index := 0; index < len(validPaths); index++ {
		path := validPaths[index]
		for i := 0; i < len(path)-1; i++ {
			node := Rooms[path[i]]
			node.links = removeLink(node.links, path[i+1])
			Rooms[path[i]] = node
		}
	}
}

func saveBeforeInPath() {
	lastPath := validPaths[len(validPaths)-1]
	for i := 1; i < len(lastPath)-1; i++ { // see if the link to the end should be removed
		room := Rooms[lastPath[i]]
		room.beforeInPath = lastPath[i-1]
		Rooms[lastPath[i]] = room
	}
}

// Removes a link from a vertex
func removeLink(links []string, conflictRoom string) []string {
	for i := 0; i < len(links); i++ {
		if links[i] == conflictRoom {
			if i + 1 < len(links) {
				links = append(links[:i], links[i+1:]...)
			} else {
				links = links[:i]	
			}
		}
	}
	return links
}


// Helper function to deep copy a map of Room structs
func copyRoomsMap(original map[string]Room) map[string]Room {
    copied := make(map[string]Room)

    for key, room := range original {
        // Deep copy the links slice
        newLinks := make([]string, len(room.links))
        copy(newLinks, room.links)

        // Create a new Room struct with the copied slice
        copied[key] = Room{
            links:        newLinks,
            isChecked:    room.isChecked,
            beforeInPath: room.beforeInPath,
        }
    }

    return copied
}
func bfs() (bool) {


	startRoom := Rooms[start]
	startRoom.isChecked = true
	Rooms[start] = startRoom

	alreadyInRevesedPath := false


	paths := [][]string{{start}}
	for len(paths) != 0 {
		if len(Rooms[start].links) == 0 {
			return false
		}
		
		for i:= 0; i < len(paths); i++ {
			validLinks := 0


			lastInPath := paths[i][len(paths[i])-1]

			if Rooms[lastInPath].beforeInPath != "" && !alreadyInRevesedPath {
				beforeInPath := Rooms[lastInPath].beforeInPath
				validLinks++

				room := Rooms[lastInPath]
				room.isChecked = false
				Rooms[lastInPath] = room
				
				paths[i] = append(paths[i], beforeInPath)

				alreadyInRevesedPath = true

				continue
			}

			for j, link := range Rooms[lastInPath].links {

				if link == end {

					if validLinks == 0 {
						paths[i] = append(paths[i], link)
					} else {
						paths = append(paths, append(paths[i][:len(paths[i])-1], link))
					}
					validPaths = append(validPaths, paths[i])
					resetIsChecked()
					return true
				}

				if !Rooms[link].isChecked {
					room := Rooms[link]
					validLinks++
					room.isChecked = true
					Rooms[link] = room
					if validLinks == 1 {
						paths[i] = append(paths[i], link)
					} else {
						var path = make([]string, len(paths[i]))
						copy(path, paths[i])
						paths = append(paths, append(path[:len(path)-1], link))
					}

				} else if j == len(Rooms[lastInPath].links)-1 && validLinks == 0 {
					if i + 1 < len(paths) {
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

func resetIsChecked() {

	for index := range Rooms {
		room := Rooms[index]
		room.isChecked = false
		Rooms[index] = room
	}
}

func ParsingData(str string) error {
    var err error
    roomCordinations := make(map[string]bool)

    split := strings.Split(str, "\n")
	roomNames := make(map[string]bool)

	antsNumber, err = strconv.Atoi(split[0])
	if  err != nil || antsNumber <= 0 {
		return fmt.Errorf("ERROR: invalid number of ants")
	}

    for i:= 1; i < len(split); i++ {
             if  spacesplit := strings.Split(split[i], " "); len(spacesplit) == 3 {
				
				err := checkIsValideRoomInitialisation(roomNames, roomCordinations,spacesplit,i)
				if err != nil {
					return err
				}
				
            } else if split[i] == "##start" {
				if i+1 >= len(split) {
                    return fmt.Errorf("ERROR: start-flag trailling in the end at line %d", i+1)
                }

                if start != "" {
                    return fmt.Errorf("ERROR: multiple ##start found at line %d", i+1)
                }
				spacesplit := strings.Split(split[i+1], " ")
				
				err := checkIsValideRoomInitialisation(roomNames, roomCordinations,spacesplit,i+1)
				if err != nil {
					return err
				}
                start = spacesplit[0]
                i++
                
            } else if  split[i] == "##end" {
                if i+1 >= len(split) {
                    return fmt.Errorf("ERROR: end-flag trailling in the end at line %d", i+1)
                }
                
				if end != "" {
                    return fmt.Errorf("ERROR: multiple ##end found at line %d", i+1)
                }
				spacesplit := strings.Split(split[i+1], " ")

                end = spacesplit[0]    
				if len(spacesplit) != 3 {
                    return fmt.Errorf("ERROR: No valid room initialization after start flag at line %d", i+1)
				}
				
				err := checkIsValideRoomInitialisation(roomNames, roomCordinations,spacesplit,i+1)
				if err != nil {
					return err
				}
                roomCordinations[strings.Join(spacesplit[1:], " ")] = true
				i++
				
            } else if strings.HasPrefix(split[i], "#") {
                continue
            } else if dashSplit := strings.Split(split[i], "-"); len(dashSplit) == 2 {
				if !roomNames[dashSplit[0]] || !roomNames[dashSplit[1]] {
                    return fmt.Errorf("ERROR: Trying linking an non-initialized room at line %d", i+1)
				}
                if err := fillRoomData(split[i]); err != nil {
                    return fmt.Errorf("ERROR: invaliCd link format at line %d", i+1)
                }
            } else if split[i] == "" {
                continue
            } else {
                fmt.Println("ERROR: invalid data format, at line "+strconv.Itoa(i+1))
                os.Exit(0)
            }   
    }

	if (start == "" || end == "") {
		return fmt.Errorf("ERROR: missing start or end room;")
	} else if start == end {
		return fmt.Errorf("ERROR: start is end;")
	}
    return nil
}

func checkIsValideRoomInitialisation(roomNames map[string]bool, roomCordinations map[string]bool, spaceSplit []string, i int) error {
	var err error
	
	if roomNames[spaceSplit[0]] {
		return fmt.Errorf("ERROR: Duplicated room initialazation %d", i+1)
	}
	roomNames[spaceSplit[0]] = true

	if _, err = strconv.Atoi(spaceSplit[1]); err != nil {
		return fmt.Errorf("ERROR: invalid coordinates for room at line %d", i+1)
	}
	
	if _, err = strconv.Atoi(spaceSplit[2]); err != nil {
		return fmt.Errorf("ERROR: invalid coordinates for room at line %d", i+1)
	}

	if roomCordinations[strings.Join(spaceSplit[1:], " ")] {
		return fmt.Errorf("ERROR: duplicate coordinates for room at line %d", i+1)
	}

	roomCordinations[strings.Join(spaceSplit[1:], " ")] = true
	return nil
}

func fillRoomData(str string) error {
    split := strings.Split(str, "-")

    if split[0] == split[1] {
        return fmt.Errorf("room %s can't link to itself", split[0])
    }

    for _, link := range Rooms[split[0]].links {
        if link == split[1] {
            return fmt.Errorf("duplicate link between %s and %s", split[0], split[1])
        }
    }

    if (split[0] == start && split[1] == end) || (split[0] == end && split[1] == start) {
        validPaths = [][]string{{start, end}}
        return nil
    }

    roomA := Rooms[split[0]]
    roomA.links = append(roomA.links, split[1])
    Rooms[split[0]] = roomA

    roomB := Rooms[split[1]]
    roomB.links = append(roomB.links, split[0])
    Rooms[split[1]] = roomB

    return nil
}

func HandleExport(ants []int, turns int, shortestPathIndex int, originalData string) {

	fmt.Println(originalData+"\n")

    var result = make([]string, turns, turns)
    AntsMoved := 1

     for i,Ants := range ants {
		for j:= 0 ; j < Ants; j++{
			for k, room:= range allValidPaths[shortestPathIndex][i][1:] {
				if result[k+j] != "" {
					result[k+j] += " "
				}
				result[k+j] += "L"+strconv.Itoa(AntsMoved)+"-"+room
			}
			AntsMoved++
		}
    }
    fmt.Print(strings.Join(result, "\n"))
}