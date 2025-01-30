package GlobVar

// Room represents a room in the network.
// It contains links to other rooms, a flag to indicate if it has been checked, and the previous room in the path.
type Room struct {
	Links        []string
	IsChecked    bool
	BeforeInPath string
}

// Global variables used throughout the program.
var (
	AntsNumber    int
	OriginalRooms = make(map[string]Room)
	Rooms         = make(map[string]Room)
	Start         string
	End           string
	ValidPaths    [][]string
	AllValidPaths [][][]string
)