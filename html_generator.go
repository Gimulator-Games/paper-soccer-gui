package main

import (
	"fmt"
	"strconv"
	"time"
)

type worldDrawer struct {
	World
	width, height int
	grid          [][]Position
}

func (w *worldDrawer) DrawField() string {
	var (
		html    = ""
		delta   = min(w.width/(w.Width+1), w.height/(w.Height+1))
		marginx = (w.width - delta*(w.Width-1)) / 2
		marginy = (w.height - delta*(w.Height-1)) / 2
	)
	w.grid = make([][]Position, w.Width)
	for i := 0; i < w.Width; i++ {
		w.grid[i] = make([]Position, w.Height)
	}

	for x := 0; x < w.Width; x++ {
		for y := 0; y < w.Height; y++ {
			xx := marginx + x*delta
			yy := marginy + y*delta
			w.grid[x][w.Height-y-1] = Position{X: xx, Y: yy}
		}
	}

	for _, move := range w.FilledMoves {
		color := "yellow"
		html += newLine(w.grid[move.From.X][move.From.Y].X, w.grid[move.From.X][move.From.Y].Y, w.grid[move.To.X][move.To.Y].X, w.grid[move.To.X][move.To.Y].Y, color)
	}

	for _, move := range w.Moves {
		name := move.Player.Name
		color := "yellow"
		if name == w.Player1.Name {
			color = "red"
		} else if name == w.Player2.Name {
			color = "blue"
		}
		html += newLine(w.grid[move.From.X][move.From.Y].X, w.grid[move.From.X][move.From.Y].Y, w.grid[move.To.X][move.To.Y].X, w.grid[move.To.X][move.To.Y].Y, color)
	}

	for x := 0; x < w.Width; x++ {
		for y := 0; y < w.Height; y++ {
			html += newClickable(
				fmt.Sprintf("click(%d, %d)", x+1, y+1),
				newCircle(w.grid[x][y].X, w.grid[x][y].Y, 5, "yellow"),
			)
		}
	}
	html += newCircle(w.grid[w.BallPos.X][w.BallPos.Y].X, w.grid[w.BallPos.X][w.BallPos.Y].Y, 10, "blue")

	return html
}

func (w worldDrawer) genUpperSpec() (string, string) {
	if w.Player1.Side == topSide {
		return w.Player1.Name, timeConverter(100)
	}
	if w.Player2.Side == topSide {
		return w.Player2.Name, timeConverter(100)
	}
	return "No Player", "00:00"
}

func (w worldDrawer) genLowerSpec() (string, string) {
	if w.Player1.Side == downSide {
		return w.Player1.Name, timeConverter(100)
	}
	if w.Player2.Side == downSide {
		return w.Player2.Name, timeConverter(100)
	}
	return "No Player", "00:00"
}

func (w worldDrawer) genTurn() string {
	if w.Player1.Name == w.Turn {
		return "Turn: " + w.Turn + " ---> red"
	} else if w.Player2.Name == w.Turn {
		return "Turn: " + w.Turn + " ---> blue"
	} else {
		return "Turn: -" 
	}
}

func (w worldDrawer) genPlayerName() string {
	if playerName == "" {
		return "Watcher"
	}
	return playerName
}

func timeConverter(duration int64) string {
	d := time.Duration(duration)
	min := int(d.Minutes())
	sec := int(d.Seconds()) - min*60

	return strconv.Itoa(min) + ":" + strconv.Itoa(sec)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func newClickable(jsFunction, html string) string {
	return fmt.Sprintf(`<a href="#" onclick="%s">%s</a>`, jsFunction, html)
}

func newCircle(cx, cy, r int, col string) string {
	return fmt.Sprintf(`<circle cx="%d" cy="%d" r="%d" stroke="black" stroke-width="1" fill="%s" />`, cx, cy, r, col)
}

func newLine(x1, y1, x2, y2 int, col string) string {
	return fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="2" />`, x1, y1, x2, y2, col)
}
