package logic

import (
	"bytes"
	"fmt"
	"log"

	"github.com/captncraig/agi/binutils"
)

type printer struct {
	w       *bytes.Buffer
	lineLen int
}

func (p *printer) f(format string, args ...interface{}) {
	n, _ := fmt.Fprintf(p.w, format, args...)
	p.lineLen += n
}
func (p *printer) p(args ...interface{}) {
	n, _ := fmt.Fprint(p.w, args...)
	p.lineLen += n
}
func (p *printer) nl() {
	fmt.Fprintln(p.w)
	p.lineLen = 0
}
func (p *printer) indent(n int) {
	fmt.Println("!", n)
	for i := 0; i < n; i++ {
		p.p("    ")
	}
}

type StatementNode interface {
	Addr() int
	Print(int, *printer)
	Size() int
}

// statement nodes
type StatementBase struct {
	Address int
}

func (s StatementBase) Addr() int {
	return s.Address
}

type ListNode []StatementNode

func (l ListNode) Addr() int {
	return l[0].Addr()
}
func (l ListNode) Print(indent int, p *printer) {
	for _, n := range l {
		n.Print(indent, p)
	}
}

func (l ListNode) Size() int {
	total := 0
	for _, n := range l {
		total += n.Size()
	}
	return total
}

type CommandNode struct {
	StatementBase
	Command *command
	Args    []byte
}

func (c *CommandNode) Print(indent int, p *printer) {
	//indent
	p.indent(indent)
	p.f("%s(", c.Command.Name)
	for i, arg := range c.Args {
		if i > 0 {
			p.p(", ")
		}
		at := c.Command.ArgTypes[i]
		p.f("%s%d", atPrefixes[at], arg)
	}
	p.p(") ~", c.Size())
	p.nl()
}

func (c *CommandNode) Size() int {
	return 1 + len(c.Args)
}

type IfNode struct {
	StatementBase
	Condition ExpressionNode
	Commands  ListNode
	Jump      uint16
}

func (i *IfNode) Size() int {
	return 2 + i.Condition.Size() + 2 + i.Commands.Size()
}

func (i *IfNode) Print(indent int, p *printer) {
	p.indent(indent)
	p.p("if (")
	i.Condition.Print(p)
	p.p(") { ~", i.Size())
	p.nl()
	i.Commands.Print(indent+1, p)
	p.indent(indent)
	p.p("}")
	p.nl()
}

type GotoNode struct {
	StatementBase
	Jump uint16
}

func (g *GotoNode) Size() int {
	return 3
}

func (g *GotoNode) Print(indent int, p *printer) {
	p.indent(indent)
	p.f("goto %d", int16(g.Jump))
	p.nl()
}

func Disassemble(dat []byte) string {
	ln := ListNode{}
	r := binutils.NewReader(dat)
	for r.HasMore() {
		st, err := parseStatement(r)
		if err != nil {
			buf := &bytes.Buffer{}
			p := &printer{w: buf}
			ln.Print(0, p)
			fmt.Println(buf.String())
			log.Fatal(err)
		}
		ln = append(ln, st)
	}
	buf := &bytes.Buffer{}
	p := &printer{w: buf}
	ln.Print(0, p)
	return buf.String()
}

func parseStatement(r *binutils.Reader) (StatementNode, error) {
	op := r.Peek()
	if op == 0xff {
		return parseIf(r)
	}
	if op == 0xfe {
		r.Take()
		gn := &GotoNode{}
		gn.Address = r.Addr()
		gn.Jump = r.Take16()
		return gn, nil
	}
	return parseCommand(r)
}

func parseCommand(r *binutils.Reader) (*CommandNode, error) {
	op := int(r.Take())
	if op >= len(AgiCommands) {
		return nil, fmt.Errorf("parseCommand: unknown op: 0x%x", op)
	}
	cmd := AgiCommands[op]
	cn := &CommandNode{
		Command: cmd,
		Args:    make([]byte, len(cmd.ArgTypes)),
	}
	for i := range cmd.ArgTypes {
		cn.Args[i] = r.Take()
	}
	return cn, nil
}

func parseIf(r *binutils.Reader) (*IfNode, error) {
	i := &IfNode{}
	i.Address = r.Addr()
	r.Take() // skip ff
	cond, err := parseAnd(r)
	if err != nil {
		return i, err
	}
	i.Condition = cond
	i.Jump = r.Take16()
	start := r.Addr()
	for r.Addr()-start < int(i.Jump) {
		st, err := parseStatement(r)
		if err != nil {
			return i, err
		}
		i.Commands = append(i.Commands, st)
	}
	return i, nil
}

// Expression Nodes
type ExpressionNode interface {
	Print(p *printer)
	Size() int
}
type AndNode []ExpressionNode

func (a AndNode) Print(p *printer) {
	for i, n := range a {
		if i > 0 {
			p.f(" && ")
		}
		n.Print(p)
	}
}
func (a AndNode) Size() int {
	total := 0
	for _, n := range a {
		total += n.Size()
	}
	return total
}

type OrNode []ExpressionNode

func (o OrNode) Size() int {
	total := 2
	for _, n := range o {
		total += n.Size()
	}
	return total
}

func (o OrNode) Print(p *printer) {
	p.p("(")
	for i, n := range o {
		if i > 0 {
			p.f(" || ")
		}
		n.Print(p)
	}
	p.p(")")
}

type ExprCommand struct {
	Command *command
	Args    []uint16
}

func (ec ExprCommand) Print(p *printer) {
	p.f("%s(", ec.Command.Name)
	for i, arg := range ec.Args {
		if i > 0 {
			p.p(", ")
		}
		//default for said
		at := atWord
		if i < len(ec.Command.ArgTypes) {
			at = ec.Command.ArgTypes[i]
		}
		p.f("%s%d", atPrefixes[at], arg)
	}
	p.p(")")
}

func (ec ExprCommand) Size() int {
	if ec.Command.Name == "said" {
		return 2 + 2*len(ec.Args)
	}
	return 1 + len(ec.Args)
}

type NotNode struct {
	Inner ExpressionNode
}

func (n NotNode) Print(p *printer) {
	p.p("!")
	n.Inner.Print(p)
}

func (n NotNode) Size() int {
	return 1 + n.Inner.Size()
}

// parse N expressions and combine them. Go until hit an 0xFF
func parseAnd(r *binutils.Reader) (AndNode, error) {
	n := AndNode{}
	for {
		op := r.Peek()
		if op == 0xff {
			r.Take()
			break
		}
		exp, err := parseExpr(r)
		if err != nil {
			return n, err
		}
		n = append(n, exp)
	}
	return n, nil
}

// parse a single expression
func parseExpr(r *binutils.Reader) (ExpressionNode, error) {
	op := int(r.Take())
	if op == 0xfd {
		inner, err := parseExpr(r)
		if err != nil {
			return nil, err
		}
		return NotNode{inner}, nil
	}
	if op == 0xfc {
		return parseOr(r)
	}
	if op >= len(TestCommands) {
		return nil, fmt.Errorf("parseExpr: unknown op: 0x%x", op)
	}
	cmd := TestCommands[op]
	ec := &ExprCommand{
		Command: cmd,
	}
	if cmd.Name == "said" {
		//special case
		n := int(r.Take())
		ec.Args = make([]uint16, n)
		for i := 0; i < n; i++ {
			ec.Args[i] = r.Take16()
		}
	} else {
		ec.Args = make([]uint16, len(cmd.ArgTypes))
		for i := range cmd.ArgTypes {
			ec.Args[i] = uint16(r.Take())
		}
	}
	return ec, nil
}

func parseOr(r *binutils.Reader) (OrNode, error) {
	m := OrNode{}
	for {
		op := r.Peek()
		if op == 0xfc {
			r.Take()
			return m, nil
		}
		inner, err := parseExpr(r)
		if err != nil {
			return nil, err
		}
		m = append(m, inner)
	}
}
