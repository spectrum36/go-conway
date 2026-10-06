package main

import (
	"fmt"
	"math/rand/v2"
	"time"
	"os"
	"os/exec"
	flag "github.com/spf13/pflag"
	"strings"
	//"github.com/inancgumus/screen"
)
type Cell struct {
	State bool
	R, G, B uint8
}

func mod(a int, b int) int { return (a%b + b) % b }

func runCmd(n string) error{
	cmd := exec.Command(n)
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	return err
}

func sameColor(a, b Cell) bool {
	return a.R == b.R && a.G == b.G && a.B == b.B
}

func randCell(color [][]uint8) Cell {
	var out Cell
	out.State = true
	c := rand.IntN(len(color))
	out.R = color[c][0]
	out.G = color[c][1]
	out.B = color[c][2]
	return out
} 

func birth(a, b, c Cell, mut int, color [][]uint8) Cell {
	var out Cell
	switch {
	case sameColor(a, b), sameColor(a, c):
		out = a
	case sameColor(b, c):
		out = b
	default:
		out = [3]Cell{a,b,c}[rand.IntN(3)]
	}
	out.State = true
	if rand.IntN(100) < mut {
		c := rand.IntN(len(color))
		out.R = color[c][0]
		out.G = color[c][1]
		out.B = color[c][2]
	}
	return out
}

func conway(a [][]Cell, dx int, dy int, death int, life int, mut int, color [][]uint8) [][]Cell{
	wrap := len(os.Args) > 1 && os.Args[1] == "wrap"
	b := make([][]Cell, dy)
	for i := range b {
		b[i] = make([]Cell, dx)
	}
	for i := range a {	
		for p := range a[i] {
			cur := a[i][p]
			var cells [8]Cell
			n := 0
			same := 0
			for q := range 3 {
				for r	:= range 3 {
					if q == 1 && r == 1 {
						continue
					}
					y := i - 1 + q
					x := p - 1 + r
					if wrap {
						y = mod(y, dy)
						x = mod(x, dx)
					} else if y < 0 || y >= dy || x < 0 || x >= dx {
						continue
					}
					if !a[y][x].State {
						continue
					}
					cells[n] = a[y][x]
					n++
					if sameColor(cur, a[y][x]) {
						same++
					}
				}
			}

			if cur.State {
				if (same == 2 || same == 3) && rand.IntN(100) >= death {
					if rand.IntN(100) < mut {
						b[i][p] = randCell(color)
					} else {
						b[i][p] = cur
					}
				}
			} else if n == 3 {
				b[i][p] = birth(cells[0], cells[1], cells[2], mut, color)
			} else if rand.IntN(10000) < life {
				b[i][p] = randCell(color)
			}
		}
	}
	return b
}

func renderLine(a [][]Cell) (string, bool) {
	alive := false
	var out strings.Builder
	for i := range a {
		for p := range a[i] {
			if a[i][p].State {
				if !alive {
					alive = true
				}
				fmt.Fprintf(&out, "\033[38;2;%d;%d;%dm██", a[i][p].R, a[i][p].G, a[i][p].B)	
			} else {
				out.WriteString("\033[0m  ")
			}
		}
		out.WriteString("\033[0m\n")
	}
	return out.String(), alive
}

func colorCheck(a []uint8, b []uint8) bool {
	return a[0] == b[0] && a[1] == b[1] && a[2] == b[2]
}

func main () {
	prob := flag.Int("prob", 10, "probability of cell spawning during initialization")
	tm := flag.Int("time", 250, "number of milliseconds to wait between generations")
	dth := flag.Int("death", 0, "probability of cell death within generation")
	spc := flag.Int("species", 1, "number of species to choose from, 5 max")
	life := flag.Int("life", 0, "probability of life generating randomly, in .01% steps")
	mut := flag.Int("mut", 0, "probability of changing species")
	flag.Parse()
	if *prob > 100 || *prob < 0 {
		fmt.Println("enter generation probability between 0 and 100")
		return
	} else if *dth > 100 || *dth < 0 {
		fmt.Println("enter death probability between 0 and 100")
		return
	}
	//set up color choices
	color := make([][]uint8, *spc)
	count := 0
	for i := range color {	
		for {
			same := false
			temp := make([]uint8, 3)
			for p := range 3 {
				temp[p] = uint8(rand.IntN(256)) 
			}
			for p := range count {
				if colorCheck(temp, color[p]){
					same = true
				}
			}
			if !same {
				color[i] = make([]uint8, 3)
				color[i] = temp
				break
			}
		}
		count++
	}

	dx := 106
	dy := 61
	a := make([][]Cell, dy)
	for i := range a {
		a[i] = make([]Cell, dx)
		for p := range a[i] {
			if rand.IntN(100) <= *prob {
				if *spc > 1 {
					c := rand.IntN(*spc)
					a[i][p].R = color[c][0]
					a[i][p].G = color[c][1]
					a[i][p].B = color[c][2]
				} else {
					a[i][p].R = 255
					a[i][p].G = 255
					a[i][p].B = 255
				}
				a[i][p].State = true
			} else {
				a[i][p].State = false
			}	
		} 
	}
	alive := true	
	out := ""
	for alive {
		out, alive = renderLine(a)
		a = conway(a, dx, dy, *dth, *life, *mut, color)
		err := runCmd("reset")
		if err != nil {
			panic(err)
		}
		fmt.Print(out)
		time.Sleep(time.Duration(*tm) * time.Millisecond)
	}
}
