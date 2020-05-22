package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/zserge/lorca"
)

const (
	Height = 600
	Width  = 800
)

var (
	ui           lorca.UI
	lastDrawer   worldDrawer
	disableEvent bool
	playerName   string
	controller   *Controller
)

func fuck(msg string)               { ui.Eval(fmt.Sprintf(`console.log("%s")`, msg)) }
func renderField(html string)       { ui.Eval(fmt.Sprintf("renderField(`%s`);", html)) }
func renderUpperPlayer(html string) { ui.Eval(fmt.Sprintf("renderUpperPlayer(`%s`);", html)) }
func renderUpperTime(html string)   { ui.Eval(fmt.Sprintf("renderUpperTime(`%s`);", html)) }
func renderTurn(html string)        { ui.Eval(fmt.Sprintf("renderTurn(`%s`);", html)) }
func renderName(html string)        { ui.Eval(fmt.Sprintf("renderName(`%s`);", html)) }
func renderLowerPlayer(html string) { ui.Eval(fmt.Sprintf("renderLowerPlayer(`%s`);", html)) }
func renderLowerTime(html string)   { ui.Eval(fmt.Sprintf("renderLowerTime(`%s`);", html)) }
func width() int                    { return ui.Eval(`width()`).Int() }
func height() int                   { return ui.Eval(`height()`).Int() }

func render(drawer worldDrawer) {
	log.Println("Start rendering...")
	field := drawer.DrawField()
	renderField(field)

	upperPlayer, upperTime := drawer.genUpperSpec()
	lowerPlayer, lowerTime := drawer.genLowerSpec()
	renderUpperPlayer(upperPlayer)
	renderLowerPlayer(lowerPlayer)
	renderUpperTime(upperTime)
	renderLowerTime(lowerTime)

	turn := drawer.genTurn()
	renderTurn(turn)

	name := drawer.genPlayerName()
	renderName(name)

	lastDrawer = drawer
}

func main() {

	initGUI()
	defer ui.Close()
	log.Println("Start GUI...")

	time.Sleep(time.Millisecond * 500)

	disableEvent = true
	controller, err := NewController("world", "paper-soccer", "verdict")
	if err != nil {
		panic(err)
	}
	controller.Run()

	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-ui.Done():
	}

	log.Println("exiting...")
}

func initGUI() {
	var err error
	args := []string{}
	if runtime.GOOS == "linux" {
		args = append(args, "--class=Lorca")
	}

	ui, err = lorca.New("", "", Width, Height, args...)
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	go http.Serve(ln, http.FileServer(FS))
	err = ui.Load(fmt.Sprintf("http://%s", ln.Addr()))
	if err != nil {
		panic(err)
	}

	w, _ := NewWorld("wait", "wait", 11, 15)
	drawer := worldDrawer{
		World:  w,
		width:  width(),
		height: height(),
	}

	lastDrawer = drawer
	render(drawer)
}
