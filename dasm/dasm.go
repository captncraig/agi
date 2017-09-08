package dasm

import (
	"bytes"
	"fmt"
	"log"

	"github.com/captncraig/agi/logic"
)

type Stmt interface{}

type StmtList []Stmt

type CommandStmt struct {
	Addr   uint16
	OpCode byte
	OpName string
	Args   []byte
}

type IfStmt struct {
	Addr      uint16
	Condition AndExpr
	Code      StmtList
	Rel       uint16 // else jump
	Abs       uint16 // for convenience
}

type GotoStmt struct {
	Addr uint16
	Rel  int16
	Abs  uint16
}

type Expr interface{}

type AndExpr []Expr

type OrExpr []Expr

type TestExpr struct {
	OpCode byte
	OpName string
	Args   []uint16
}

type NotExpr struct {
	inner Expr
}

type Program struct {
	list StmtList
}

type reader struct {
	dat []byte
	idx uint16
}

func (r *reader) peek() byte {
	return r.dat[r.idx]
}
func (r *reader) take() byte {
	d := r.dat[r.idx]
	r.idx++
	return d
}
func (r *reader) has() bool {
	return r.idx < uint16(len(r.dat))
}

//Decompile performs a first pass, simple instruction decoding
func Decompile(instructions []byte) *Program {
	r := &reader{instructions, 0}
	l := StmtList{}
	for r.has() {
		l = append(l, r.parseStmt())
	}
	return &Program{list: l}
}

func (r *reader) parseStmt() Stmt {
	a := r.idx
	d := r.take()
	if d == 0xff {
		return r.parseIf()
	}
	if d == 0xfe {
		rel := uint16(r.take()) | uint16(r.take())<<8
		return &GotoStmt{
			Addr: a,
			Rel:  int16(rel),
			Abs:  uint16(int16(rel) + int16(r.idx)),
		}
	}
	if int(d) >= len(logic.AgiCommands) {
		log.Fatalf("Unknown AGI Command: 0x%02x", d)
	}
	cmd := logic.AgiCommands[d]
	cn := &CommandStmt{
		Addr:   a,
		OpCode: d,
		OpName: cmd.Name,
		Args:   make([]byte, len(cmd.ArgTypes)),
	}
	for i := range cmd.ArgTypes {
		cn.Args[i] = r.take()
	}
	return cn
}
func (r *reader) parseIf() *IfStmt {
	i := &IfStmt{
		Addr:      r.idx - 1,
		Condition: r.parseAnd(),
		Rel:       uint16(r.take()) | uint16(r.take())<<8,
	}
	i.Abs = i.Rel + r.idx
	for r.idx < i.Abs {
		i.Code = append(i.Code, r.parseStmt())
	}
	return i
}

func (r *reader) parseAnd() AndExpr {
	a := AndExpr{}
	for {
		d := r.peek()
		if d == 0xff {
			r.take()
			break
		}
		a = append(a, r.parseExpr())
	}
	return a
}

func (r *reader) parseExpr() Expr {
	d := r.take()
	if d == 0xfd {
		return &NotExpr{
			inner: r.parseExpr(),
		}
	}
	if d == 0xfc {
		return r.parseOr()
	}
	if int(d) >= len(logic.TestCommands) {
		log.Fatalf("Unknown Test Command: 0x%02x", d)
	}
	cmd := logic.TestCommands[d]
	te := TestExpr{
		OpCode: d,
		OpName: cmd.Name,
	}
	if cmd.Name == "said" {
		//special case
		n := int(r.take())
		te.Args = make([]uint16, n)
		for i := 0; i < n; i++ {
			te.Args[i] = uint16(r.take()) | uint16(r.take())<<8
		}
	} else {
		te.Args = make([]uint16, len(cmd.ArgTypes))
		for i := range cmd.ArgTypes {
			te.Args[i] = uint16(r.take())
		}
	}
	return te
}

func (r *reader) parseOr() OrExpr {
	o := OrExpr{}
	for {
		d := r.peek()
		if d == 0xfc {
			r.take()
			break
		}
		o = append(o, r.parseExpr())
	}
	return o
}

func (p *Program) String() string {
	buf := &bytes.Buffer{}
	printStmt(p.list, buf, 0)
	return buf.String()
}

func printAddr(buf *bytes.Buffer, addr uint16, indent int) {
	fmt.Fprintf(buf, "0x%04x: ", addr)
	ind(buf, indent)
}

func ind(buf *bytes.Buffer, indent int) {
	for i := 0; i < indent; i++ {
		buf.WriteString("    ")
	}
}

func printStmt(s Stmt, buf *bytes.Buffer, indent int) {
	switch t := s.(type) {
	case StmtList:
		for _, st := range t {
			printStmt(st, buf, indent)
		}
	case *IfStmt:
		printAddr(buf, t.Addr, indent)
		buf.WriteString("if(")
		buf.WriteString(") {")
		buf.WriteString("\n")
		for _, st := range t.Code {
			printStmt(st, buf, indent+1)
		}
		buf.WriteString("        ")
		ind(buf, indent)
		buf.WriteString("}\n")
	case *CommandStmt:
		printAddr(buf, t.Addr, indent)
		buf.WriteString(t.OpName)
		buf.WriteString("\n")
	case *GotoStmt:
		printAddr(buf, t.Addr, indent)
		buf.WriteString("goto")
		buf.WriteString("\n")
	default:
		log.Fatalf("no print for %T", s)
	}
}

func printExpr(e Expr, buf *bytes.Buffer) {

}
