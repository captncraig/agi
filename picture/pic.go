package picture

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// ImageData is a single frame. Indexed [X][Y].
// The screen mode used by the AGI games is the 320x200x16 standard EGA mode.
// However, all graphics is designed to be shown on a 160x200x16 mode.
// This was apparently the resolution that the original PCjr interpreter used.
// They stuck with it when they started supporting EGA and thus have a situation where each AGI pixel has a width of two normal 320x200 pixels.
type ImageData [160][200]byte

type pic struct {
	picture, priority           *ImageData
	pictureDraw, priorityDraw   bool
	pictureColor, priorityColor byte
	brush                       byte
	setPic, setPri              bool
}

const (
	black byte = iota
	blue
	green
	cyan
	red
	magenta
	brown
	lightGrey
	darkGrey
	lightBlue
	lightGreen
	lightCyan
	lightRed
	lightMagenta
	yellow
	white
)

type CommandDef struct {
	op   byte
	name string
	nArg byte // 255 is var
	f    func(*pic, []byte)
}

type Command struct {
	Args   []byte
	Op     byte
	OpName string
}

type Picture []*Command

var cmds = []*CommandDef{
	{
		op:   0xf0,
		name: "setPictureColor",
		nArg: 1,
		f: func(p *pic, args []byte) {
			p.pictureDraw = true
			if len(args) != 1 {
				log.Println("setPictureColor: not enough args")
				return
			}
			p.pictureColor = args[0]
		},
	},
	{
		op:   0xf1,
		name: "disablePicture",
		nArg: 0,
		f: func(p *pic, args []byte) {
			p.pictureDraw = false
		},
	},
	{
		op:   0xf2,
		name: "setPriorityColor",
		nArg: 1,
		f: func(p *pic, args []byte) {
			p.priorityDraw = true
			if len(args) != 1 {
				log.Println("setPriorityColor: not enough args")
				return
			}
			p.priorityColor = args[0]
		},
	},
	{
		op:   0xf3,
		name: "disablePriority",
		nArg: 0,
		f: func(p *pic, args []byte) {
			p.priorityDraw = false
		},
	},
	{
		op:   0xf4,
		name: "drawYCorner",
		nArg: 255,
		f: func(p *pic, args []byte) {
			if len(args) < 2 {
				log.Println("drawYCorner: not enough args")
				return
			}
			x, y := args[0], args[1]
			i := 2
			for {
				if i >= len(args) {
					break
				}
				y1 := args[i]
				i++
				p.line(x, y, x, y1)
				y = y1
				if i >= len(args) {
					break
				}
				x1 := args[i]
				i++
				p.line(x, y, x1, y)
				x = x1
			}
		},
	},
	{
		op:   0xf5,
		name: "drawXCorner",
		nArg: 255,
		f: func(p *pic, args []byte) {
			if len(args) < 2 {
				log.Println("drawXCorner: not enough args")
				return
			}
			x, y := args[0], args[1]
			i := 2
			for {
				if i >= len(args) {
					break
				}
				x1 := args[i]
				i++
				p.line(x, y, x1, y)
				x = x1
				if i >= len(args) {
					break
				}
				y1 := args[i]
				i++
				p.line(x, y, x, y1)
				y = y1
			}
		},
	},
	{
		op:   0xf6,
		name: "absoluteLine",
		nArg: 255,
		f: func(p *pic, args []byte) {
			if len(args) < 2 {
				log.Println("absoluteLine: not enough args")
				return
			}
			if len(args)%2 != 0 {
				log.Println("absoluteLine: odd number of args")
			}
			x, y := args[0], args[1]
			for i := 2; i < len(args)-1; i += 2 {
				x1, y1 := args[i], args[i+1]
				p.line(x, y, x1, y1)
				x, y = x1, y1
			}
		},
	},
	{
		op:   0xf7,
		name: "relativeLine",
		nArg: 255,
		f: func(p *pic, args []byte) {
			if len(args) < 2 {
				log.Println("relativeLine: not enough args")
				return
			}
			x, y := args[0], args[1]
			for i := 2; i < len(args); i++ {
				d := args[i]
				yOff := d & 0x07
				y1 := y + yOff
				if d&0x08 != 0 {
					y1 = y - yOff
				}
				xOff := (d >> 4) & 0x07
				x1 := x + xOff
				if d&0x80 != 0 {
					x1 = x - xOff
				}
				p.line(x, y, x1, y1)
				x, y = x1, y1
			}
		},
	},
	{
		op:   0xf8,
		name: "fill",
		nArg: 255,
		f: func(p *pic, args []byte) {
			if len(args)%2 != 0 {
				log.Println("fill: odd number of args")
			}
			for i := 0; i < len(args)-1; i += 2 {
				x, y := args[i], args[i+1]
				p.fill(x, y)
			}
		},
	},
	{
		op:   0xf9,
		name: "changePen",
		nArg: 1,
		f: func(p *pic, args []byte) {
			if len(args) != 1 {
				log.Println("changePen: wrong number of args")
				return
			}
			p.brush = args[0]
		},
	},
	{
		op:   0xfa,
		name: "plotPen",
		nArg: 255,
		f: func(p *pic, args []byte) {
			if p.brush&0x20 == 0 {
				//If the pen style is set to solid, then the arguments
				//are just a list of coordinates to be plotted.
				if len(args)%2 != 0 {
					log.Println("plotPen(solid): odd number of args")
				}
				for i := 0; i < len(args)-1; i += 2 {
					x, y := args[i], args[i+1]
					p.set(x, y)
				}
			} else {
				//If the pen style is set to splatter brush (texture),
				//then the arguments are in groups of three with the first
				//argument giving the texture number and the other two giving the coordinates.
				if len(args)%3 != 0 {
					log.Println("plotPen(texture): odd number of args")
				}
				log.Println("texture plot not implemented")
			}
		},
	},
}

func (c *Command) String() string {
	args := []string{}
	for _, a := range c.Args {
		args = append(args, fmt.Sprintf("0x%x", a))
	}
	return fmt.Sprintf("%s(0x%x) [%s]", c.OpName, c.Op, strings.Join(args, " "))
}

var ops = reverseCmds()

func reverseCmds() map[byte]*CommandDef {
	o := make(map[byte]*CommandDef, len(cmds))
	for _, cmd := range cmds {
		o[cmd.op] = cmd
	}
	return o
}

func (p Picture) Render() (*ImageData, *ImageData) {
	return p.renderByStep(nil, nil)
}

func (p Picture) RenderToFile(dir string, idx int, prio bool) error {
	a, b := p.Render()
	i := a
	dotPart := "."
	if prio {
		dotPart = ".priority."
		i = b
	}
	fName := filepath.Join(dir, fmt.Sprintf("%d%spng", idx, dotPart))
	img := i.Image()
	f, err := os.Create(fName)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func (p Picture) RenderGifs(dir string, idx int) {
	aFile := filepath.Join(dir, fmt.Sprintf("%d.gif", idx))
	bFile := filepath.Join(dir, fmt.Sprintf("%d.priority.gif", idx))
	if _, err := os.Stat(aFile); !os.IsNotExist(err) {
		if _, err := os.Stat(bFile); !os.IsNotExist(err) {
			return
		}
	}
	priImgs := []ImageData{}
	picImgs := []ImageData{}
	ca := make(chan *ImageData)
	cb := make(chan *ImageData)
	go func() {
		for pi := range ca {
			picImgs = append(picImgs, *pi)
		}
	}()
	go func() {
		for pi := range cb {
			priImgs = append(priImgs, *pi)
		}
	}()
	p.renderByStep(ca, cb)
	close(ca)
	close(cb)
	renderGif(aFile, picImgs)
	renderGif(bFile, priImgs)
}

func renderGif(fname string, imgs []ImageData) {
	g := &gif.GIF{}
	for _, img := range imgs {
		i2 := img.Image()
		g.Image = append(g.Image, i2)
		g.Delay = append(g.Delay, 25)
		g.Config = image.Config{
			ColorModel: i2.ColorModel(),
			Height:     i2.Bounds().Dy(),
			Width:      i2.Bounds().Dx(),
		}
	}
	f, err := os.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err = gif.EncodeAll(f, g); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Wrote", fname)
}

func (p Picture) renderByStep(pCh chan<- *ImageData, prCh chan<- *ImageData) (*ImageData, *ImageData) {
	pic := &pic{
		picture:       &ImageData{},
		priority:      &ImageData{},
		priorityColor: 4,
	}
	for x := 0; x < 160; x++ {
		for y := 0; y < 200; y++ {
			pic.picture[x][y] = white
			pic.priority[x][y] = red
		}
	}
	if pCh != nil {
		pCh <- pic.picture
	}
	if prCh != nil {
		prCh <- pic.priority
	}
	for _, cmd := range p {
		op := ops[cmd.Op]
		if op == nil {
			panic(fmt.Sprintf("Non existant OP 0x%x", cmd.Op))
		}
		op.f(pic, cmd.Args)
		if pCh != nil && pic.setPic {
			pCh <- pic.picture
		}
		if prCh != nil && pic.setPri {
			prCh <- pic.priority
		}
		pic.setPic = false
		pic.setPri = false
	}
	return pic.picture, pic.priority
}

func Decode(dat []byte) (Picture, error) {
	var cmds = Picture{}
	for i := 0; i < len(dat); {
		opCode := dat[i]
		if opCode == 0xff {
			if i != len(dat)-1 {
				return nil, fmt.Errorf("Extra data remaining after EOF command")
			}
			break
		}
		i++
		cDef := ops[opCode]
		if cDef == nil {
			log.Fatalf("Opcode 0x%x not defined", opCode)
		}
		cmd := &Command{
			Op:     opCode,
			OpName: cDef.name,
		}
		for i < len(dat) {
			if cDef.nArg != 255 && len(cmd.Args) == int(cDef.nArg) {
				//satisfied
				break
			}
			arg := dat[i]
			if cDef.nArg == 255 && arg >= 0xf0 {
				break
			}
			i++
			cmd.Args = append(cmd.Args, arg)
		}
		cmds = append(cmds, cmd)
	}
	return cmds, nil
}

func (p *pic) fill(x, y byte) {
	//fmt.Println("FILL", x, y)
	q := []byte{x, y}
	ok := func(x, y byte) bool {
		if x >= 160 || y >= 200 {
			return false
		}
		if !p.pictureDraw && !p.priorityDraw {
			return false
		}
		if (p.pictureDraw && p.pictureColor == white) || (p.priorityDraw && p.priorityColor == red) {
			return false
		}
		prioOK := p.priority[x][y] == red
		if p.priorityDraw && !p.pictureDraw {
			return prioOK
		}
		picOK := p.picture[x][y] == white
		if !p.priorityDraw && p.pictureDraw {
			return picOK
		}
		// TODO: is this correct? Should I check both boundaries?
		if picOK != prioOK {
			//log.Fatal("HMMMMMM")
		}
		return picOK

	}
	enq := func(x, y byte) {
		if ok(x, y) {
			q = append(q, x)
			q = append(q, y)
		}
	}
	for len(q) > 0 {
		x = q[0]
		y = q[1]
		q = q[2:]
		if !ok(x, y) {
			continue
		}
		p.set(x, y)
		enq(x, y-1)
		enq(x-1, y)
		enq(x, y+1)
		enq(x+1, y)
	}
}

func (p *pic) set(x, y byte) {
	if x >= 160 || y >= 200 {
		log.Println("Draw out of bounds", x, y)
		return
	}
	if p.pictureDraw {
		p.picture[x][y] = p.pictureColor
		p.setPic = true
	}
	if p.priorityDraw {
		p.priority[x][y] = p.priorityColor
		p.setPri = true
	}
}

// ported from code in http://www.agidev.com/articles/agispec/agispecs-7.html#ss7.4
func (p *pic) line(bx1, by1, bx2, by2 byte) {
	x1, y1, x2, y2 := float64(bx1), float64(by1), float64(bx2), float64(by2)
	height := y2 - y1
	width := x2 - x1
	x, y := x1, y1
	// addx and addy initialized as if they are minor dimension
	addX := height
	if height != 0 {
		addX = width / math.Abs(height)
	}
	addY := width
	if width != 0 {
		addY = height / math.Abs(width)
	}
	round := func(n float64, step float64) byte {
		if step < 0 {
			if n-math.Floor(n) <= .501 {
				return byte(math.Floor(n))
			}
			return byte(math.Ceil(n))
		}
		if n-math.Floor(n) < .499 {
			return byte(math.Floor(n))
		}
		return byte(math.Ceil(n))
	}
	pset := func(x, y float64) {
		p.set(round(x, addX), round(y, addY))
	}
	var done func() bool
	if math.Abs(width) > math.Abs(height) {
		// more x steps than y
		addX = 0
		if width != 0 {
			addX = width / math.Abs(width)
		}
		done = func() bool {
			return x == x2
		}
	} else {
		// more y steps than x
		addY = 0
		if height != 0 {
			addY = height / math.Abs(height)
		}
		done = func() bool {
			return y == y2
		}
	}
	for !done() {
		pset(x, y)
		x += addX
		y += addY
	}
	pset(x2, y2)
}

const a uint8 = 255

var pal = color.Palette{
	color.RGBA{0x00, 0x00, 0x00, a},
	color.RGBA{0x00, 0x00, 0x2A, a},
	color.RGBA{0x00, 0x2A, 0x00, a},
	color.RGBA{0x00, 0x2A, 0x2A, a},
	color.RGBA{0x2A, 0x00, 0x00, a},
	color.RGBA{0x2A, 0x00, 0x2A, a},
	color.RGBA{0x2A, 0x15, 0x00, a},
	color.RGBA{0x2A, 0x2A, 0x2A, a},
	color.RGBA{0x15, 0x15, 0x15, a},
	color.RGBA{0x15, 0x15, 0x3F, a},
	color.RGBA{0x15, 0x3F, 0x15, a},
	color.RGBA{0x15, 0x3F, 0x3F, a},
	color.RGBA{0x3F, 0x15, 0x15, a},
	color.RGBA{0x3F, 0x15, 0x3F, a},
	color.RGBA{0x3F, 0x3F, 0x15, a},
	color.RGBA{0x3F, 0x3F, 0x3F, a},
}

func init() {
	scale := func(b byte) byte {
		// x / 0x3f == x2 / 0xff
		return byte(uint16(b) * 0xff / 0x3f)
	}
	for i, cl := range pal {
		c := cl.(color.RGBA)
		pal[i] = color.RGBA{scale(c.R), scale(c.G), scale(c.B), a}
	}
}

func (i *ImageData) Image() *image.Paletted {
	var img = image.NewPaletted(image.Rect(0, 0, 320, 200), pal)
	for x := 0; x < 160; x++ {
		for y := 0; y < 200; y++ {
			ix := x * 2
			c := i[x][y]
			img.Set(ix, y, pal[c])
			img.Set(ix+1, y, pal[c])
		}
	}
	return img
}
