package sector

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type XY struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Connector struct {
	Map map[string]Sector `json:"sector"`
}

type Sector struct {
	X    int    `json:"x"`
	Y    int    `json:"y"`
	Name string `json:"sector"`
	Abb  string `json:"abb"`
	File string `json:"filepath"`
}

func NewConnector() *Connector {
	c := Connector{}
	c.Map = make(map[string]Sector)
	return &c
}

func (c *Connector) Add(X, Y int, Name, File string) {
	xy := XY{X, Y}
	for len(Name) < 4 {
		Name += " "
	}
	abb := strings.Join(strings.Split(Name, "")[:4], "")
	c.Map[fmt.Sprintf("{%v,%v}", xy.X, xy.Y)] = Sector{
		X:    X,
		Y:    Y,
		Name: Name,
		Abb:  abb,
		File: File,
	}
}

var ConnectorPath = `c:\Users\pemaltynov\travellermap\res\merge\compiledSectorCoordinates.json`

func (c *Connector) Save() error {

	bt, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ConnectorPath, bt, 0755)
}

func LoadConnector() (*Connector, error) {
	bt, err := os.ReadFile(ConnectorPath)
	if err != nil {
		return nil, err
	}
	c := NewConnector()
	err = json.Unmarshal(bt, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func Name(sx, sy int) string {

	return "Undefined"
}

func Abb(name string) string {
	letters := strings.Split(name, "")
	return strings.Join(letters[:4], "")
}

func Hex(lx, ly int) string {
	w := fmt.Sprintf("%v", ly)
	for len(w) < 2 {
		w = "0" + w
	}
	h := fmt.Sprintf("%v", lx)
	for len(h) < 2 {
		h = "0" + h
	}
	return h + w
}

func ParseLocalFromHex(hex string) (int, int) {
	lx, ly := -1, -1
	if hex == "" {
		fmt.Println("no hex")
		return lx, ly
	}
	if len(hex) == 3 {
		hex = "0" + hex
	}
	err := fmt.Errorf("not parsed")
	pts := strings.Split(hex, "")
	lx, err = strconv.Atoi(pts[0] + pts[1])
	if err != nil {
		fmt.Println(pts)
		return -1, -1
	}
	ly, err = strconv.Atoi(pts[2] + pts[3])
	if err != nil {
		fmt.Println(pts)
		return -1, -1
	}
	return lx, ly
}
