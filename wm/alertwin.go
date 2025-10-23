package wm

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

var font xproto.Font
var win xproto.Window
var gc xproto.Gcontext
var screen *xproto.ScreenInfo

func createErrWindow(conn *xgb.Conn) {
	// Create a window
	border := screen.BlackPixel
	colrep, err := xproto.AllocColor(conn, screen.DefaultColormap, 255*257, 0*257, 0*257).Reply()
	if err == nil {
		border = colrep.Pixel
	} else {
		log.Println("couldn't make border colour", err)
	}
	log.Println(colrep.Pixel, colrep.Pixel, colrep.Red, colrep.Green, colrep.Blue)
	win, _ = xproto.NewWindowId(conn)
	xproto.CreateWindow(conn,
		screen.RootDepth,
		win,
		screen.Root,
		10, 10, screen.WidthInPixels-20, 75, 3,
		xproto.WindowClassInputOutput,
		screen.RootVisual,
		xproto.CwBackPixel|xproto.CwEventMask|xproto.CwBorderPixel,
		[]uint32{
			screen.BlackPixel,
			border,
			xproto.EventMaskExposure,
		},
	)

	xproto.ChangeWindowAttributes(conn, win, xproto.CwOverrideRedirect, []uint32{1})

	// Set the window title
	xproto.ChangeProperty(conn, xproto.PropModeReplace,
		win, xproto.AtomWmName, xproto.AtomString, 8,
		uint32(len("Config Error")), []byte("Config Error"),
	)

	xproto.MapWindow(conn, win)

	// Load a core X font (e.g., "fixed")
	font, _ = xproto.NewFontId(conn)
	xproto.OpenFont(conn, font, uint16(len("fixed")), "fixed")

	// Create a GC (graphics context) that uses this font
	gc, _ = xproto.NewGcontextId(conn)
	xproto.CreateGC(conn, gc, xproto.Drawable(win),
		xproto.GcForeground|xproto.GcFont,
		[]uint32{screen.WhitePixel, uint32(font)},
	)
}

func errwinclose(conn *xgb.Conn) {
	// Clean up
	xproto.DestroyWindow(conn, win)
	xproto.CloseFont(conn, font)
}

func (wm *WindowManager) errwin(msg string) {
	conn := wm.conn
	screen = wm.screen
	createErrWindow(conn)
	wm.pointerToWindow(wm.root)
	log.Println("ERRWIN MESSAGE:", msg, byte(len(msg)), strings.Count(msg, "\n"))
	lines := strings.Split(msg, "\n")
	var offset int16 = 0
	for i := range lines {
		err := xproto.ImageText8Checked(conn, byte(len(lines[i])), xproto.Drawable(win), gc, 20, 20+offset, lines[i]).Check()
		if err != nil {
			log.Println("err showing text", err.Error())
		}
		offset += 20
	}
}

func internAtom(conn *xgb.Conn, name string) (xproto.Atom, error) {
	reply, err := xproto.InternAtom(conn, true, uint16(len(name)), name).Reply()
	if err != nil {
		return 0, err
	}
	return reply.Atom, nil
}

func parseConfigError(err string) string {
	// get line num that error is referencing
	sb := strings.Index(err, "[")
	se := strings.Index(err, "]")

	str := err[sb+1 : se]
	ln := strings.Split(str, ":")[0]
	num, _ := strconv.Atoi(ln)

	// the error gives a list of lines before and after the actual error line so we cut out only the line we need and the arrow below that points to where it is on the line
	lines := strings.Split(err, "\n")
	numlines := lines[1 : len(lines)-1]
	final := lines[0] + "   "
	found := false
	for i := range numlines {
		if strings.Contains(numlines[i], fmt.Sprint(num, " |")) {
			final += numlines[i] + "\n"
			found = true
		} else if found {
			final += strings.Repeat(" ", 25+len(lines[0])) + numlines[i]
			found = false
		}
	}

	return final
}
