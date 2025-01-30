package Utils

import (
	"fmt"
	"lem-in/GlobVar"
	"os"
	"strconv"
	"strings"
)

// HandleExport prints the input data and simulates the movement of ants along the shortest path.
func HandleExport(ants []int, turns int, shortestPathIndex int, originalData string) {
    fmt.Println(originalData+"\n")

    var result = []string{}
    AntsMoved := 1

    // Simulate ant movement along the shortest path
     for i, Ants := range ants {
        for j:= 0 ; j < Ants; j++{
            for k, room:= range GlobVar.AllValidPaths[shortestPathIndex][i][1:] {
                if k+j > len(result)-1 {
                    result = append(result, "")
                }
                if result[k+j] != "" {
                    result[k+j] += " "
                }
                result[k+j] += "L"+strconv.Itoa(AntsMoved)+"-"+room
            }
            AntsMoved++
        }
    }
    // Print the movement of ants
    fmt.Print(strings.Join(result, "\n"))
}


// ParsingData parses the input data and initializes global variables.
// It validates the number of ants, rooms, and links.
func ParsingData(str string) error {
    var err error
    roomCordinations := make(map[string]bool)

    split := strings.Split(str, "\n")
	roomNames := make(map[string]bool)

	GlobVar.AntsNumber, err = strconv.Atoi(split[0])
	if  err != nil || GlobVar.AntsNumber <= 0 {
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

                if GlobVar.Start != "" {
                    return fmt.Errorf("ERROR: multiple ##start found at line %d", i+1)
                }
				spacesplit := strings.Split(split[i+1], " ")
				
				err := checkIsValideRoomInitialisation(roomNames, roomCordinations,spacesplit,i+1)
				if err != nil {
					return err
				}
                GlobVar.Start = spacesplit[0]
                i++
                
            } else if  split[i] == "##end" {
                if i+1 >= len(split) {
                    return fmt.Errorf("ERROR: end-flag trailling in the end at line %d", i+1)
                }
                
				if GlobVar.End != "" {
                    return fmt.Errorf("ERROR: multiple ##end found at line %d", i+1)
                }
				spacesplit := strings.Split(split[i+1], " ")

                GlobVar.End = spacesplit[0]    
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

	if (GlobVar.Start == "" || GlobVar.End == "") {
		return fmt.Errorf("ERROR: missing start or end room;")
	} else if GlobVar.Start == GlobVar.End {
		return fmt.Errorf("ERROR: start is end;")
	}
    return nil
}


// checkIsValideRoomInitialisation validates room initialization.
// It ensures that room names and coordinates are unique and valid.
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

// fillRoomData stores links between rooms in the global Rooms map.
// It ensures that links are valid and not duplicated.
func fillRoomData(str string) error {
    split := strings.Split(str, "-")

    if split[0] == split[1] {
        return fmt.Errorf("room %s can't link to itself", split[0])
    }

    for _, link := range GlobVar.Rooms[split[0]].Links {
        if link == split[1] {
            return fmt.Errorf("duplicate link between %s and %s", split[0], split[1])
        }
    }

    if (split[0] == GlobVar.Start && split[1] == GlobVar.End) || (split[0] == GlobVar.End && split[1] == GlobVar.Start) {
        GlobVar.ValidPaths = [][]string{{GlobVar.Start, GlobVar.End}}
        return nil
    }

    roomA := GlobVar.Rooms[split[0]]
    roomA.Links = append(roomA.Links, split[1])
    GlobVar.Rooms[split[0]] = roomA

    roomB := GlobVar.Rooms[split[1]]
    roomB.Links = append(roomB.Links, split[0])
    GlobVar.Rooms[split[1]] = roomB

    return nil
}