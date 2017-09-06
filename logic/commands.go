package logic

type ArgType int

const (
	atVar ArgType = iota
	atNum
	atFlag
	atSObj
	atIObj
	atStr
	atMsg
	atWord
	atCtrl
)

// correspond to ArgTypes 1:1
var atPrefixes = []string{
	"v",
	"#",
	"f",
	"O",
	"I",
	"S",
	"M",
	"W",
	"C",
}

type command struct {
	Name     string
	ArgTypes []ArgType
}

var TestCommands = []*command{
	{Name: "false"},
	{Name: "equaln", ArgTypes: []ArgType{atVar, atNum}},
	{Name: "equalv", ArgTypes: []ArgType{atVar, atVar}},
	{Name: "lessn", ArgTypes: []ArgType{atVar, atNum}},
	{Name: "lessv", ArgTypes: []ArgType{atVar, atVar}},
	{Name: "greatern", ArgTypes: []ArgType{atVar, atNum}},
	{Name: "greaterv", ArgTypes: []ArgType{atVar, atVar}},
	{Name: "isset", ArgTypes: []ArgType{atFlag}},
	{Name: "issetv", ArgTypes: []ArgType{atVar}},
	{Name: "has", ArgTypes: []ArgType{atIObj}},
	{Name: "obj.in.room", ArgTypes: []ArgType{atIObj, atVar}},
	{Name: "posn", ArgTypes: []ArgType{atSObj, atNum, atNum, atNum, atNum}},
	{Name: "controller", ArgTypes: []ArgType{atCtrl}},
	{Name: "have.key"},
	{Name: "said"},
	{Name: "compare.strings", ArgTypes: []ArgType{atStr, atStr}},
	{Name: "obj.in.box", ArgTypes: []ArgType{atSObj, atNum, atNum, atNum, atNum}},
	{Name: "center.posn", ArgTypes: []ArgType{atSObj, atNum, atNum, atNum, atNum}},
	{Name: "right.posn", ArgTypes: []ArgType{atSObj, atNum, atNum, atNum, atNum}},
}

var AgiCommands = []*command{
	{Name: "return"},
	{Name: "increment", ArgTypes: []ArgType{atVar}},
	{Name: "decrement", ArgTypes: []ArgType{atVar}},
	{Name: "assignn", ArgTypes: []ArgType{atVar, atNum}},
	{Name: "assignv", ArgTypes: []ArgType{atVar, atVar}},
	{Name: "addn", ArgTypes: []ArgType{atVar, atNum}},
	{Name: "addv", ArgTypes: []ArgType{atVar, atVar}},
	{Name: "subn", ArgTypes: []ArgType{atVar, atNum}},
	{Name: "subv", ArgTypes: []ArgType{atVar, atVar}},
	{Name: "lindirectv", ArgTypes: []ArgType{atVar, atVar}},
	{Name: "rindirect", ArgTypes: []ArgType{atVar, atVar}},
	{Name: "lindirectn", ArgTypes: []ArgType{atVar, atNum}},
	{Name: "set", ArgTypes: []ArgType{atFlag}},
	{Name: "reset", ArgTypes: []ArgType{atFlag}},
	{Name: "toggle", ArgTypes: []ArgType{atFlag}},
	{Name: "set.v", ArgTypes: []ArgType{atVar}},
	{Name: "reset.v", ArgTypes: []ArgType{atVar}},
	{Name: "toggle.v", ArgTypes: []ArgType{atVar}},
	{Name: "new.room", ArgTypes: []ArgType{atNum}},
	{Name: "new.room.v", ArgTypes: []ArgType{atVar}},
	{Name: "load.logics", ArgTypes: []ArgType{atNum}},
	{Name: "load.logics.v", ArgTypes: []ArgType{atVar}},
	{Name: "call", ArgTypes: []ArgType{atNum}},
	{Name: "call.v", ArgTypes: []ArgType{atVar}},
	{Name: "load.pic", ArgTypes: []ArgType{atVar}},
	{Name: "draw.pic", ArgTypes: []ArgType{atVar}},
	{Name: "show.pic"},
	{Name: "discard.pic", ArgTypes: []ArgType{atVar}},
	{Name: "overlay.pic", ArgTypes: []ArgType{atVar}},
	{Name: "show.pri.screen"},
	{Name: "load.view", ArgTypes: []ArgType{atNum}},
	{Name: "load.view.v", ArgTypes: []ArgType{atVar}},
	{Name: "discard.view", ArgTypes: []ArgType{atNum}},
	{Name: "animate.obj", ArgTypes: []ArgType{atSObj}},
	{Name: "unanimate.all"},
	{Name: "draw", ArgTypes: []ArgType{atSObj}},
	{Name: "erase", ArgTypes: []ArgType{atSObj}},
	{Name: "position", ArgTypes: []ArgType{atSObj, atNum, atNum}},
	{Name: "position.v", ArgTypes: []ArgType{atSObj, atVar, atVar}},
	{Name: "get.posn", ArgTypes: []ArgType{atSObj, atVar, atVar}},
	{Name: "reposition", ArgTypes: []ArgType{atSObj, atVar, atVar}},
	{Name: "set.view", ArgTypes: []ArgType{atSObj, atNum}},
	{Name: "set.view.v", ArgTypes: []ArgType{atSObj, atVar}},
	{Name: "set.loop", ArgTypes: []ArgType{atSObj, atNum}},
	{Name: "set.loop.v", ArgTypes: []ArgType{atSObj, atVar}},
	{Name: "fix.loop", ArgTypes: []ArgType{atSObj}},
	{Name: "release.loop", ArgTypes: []ArgType{atSObj}},
	{Name: "set.cel", ArgTypes: []ArgType{atSObj, atNum}},
	{Name: "set.cel.v", ArgTypes: []ArgType{atSObj, atVar}},
	{Name: "last.cel", ArgTypes: []ArgType{atSObj, atVar}},
	{Name: "current.cel", ArgTypes: []ArgType{atSObj, atVar}},
	{Name: "current.loop", ArgTypes: []ArgType{atSObj, atVar}},
	{Name: "current.view", ArgTypes: []ArgType{atSObj, atVar}},
	{Name: "number.of.loops", ArgTypes: []ArgType{atSObj, atVar}},
	{Name: "set.priority", ArgTypes: []ArgType{atSObj, atNum}},
	{Name: "set.priority.v", ArgTypes: []ArgType{atSObj, atVar}},
	{Name: "release.priority", ArgTypes: []ArgType{atSObj}},
	{Name: "get.priority", ArgTypes: []ArgType{atSObj, atVar}},
	{Name: "stop.update", ArgTypes: []ArgType{atSObj}},
	{Name: "start.update", ArgTypes: []ArgType{atSObj}},
	{Name: "force.update", ArgTypes: []ArgType{atSObj}},
	{Name: "ignore.horizon", ArgTypes: []ArgType{atSObj}},
	{Name: "observe.horizon", ArgTypes: []ArgType{atSObj}},
	{Name: "set.horizon", ArgTypes: []ArgType{atNum}},
	{Name: "object.on.water", ArgTypes: []ArgType{atSObj}},
	{Name: "object.on.land", ArgTypes: []ArgType{atSObj}},
	{Name: "object.on.anything", ArgTypes: []ArgType{atSObj}},
	{Name: "ignore.objs", ArgTypes: []ArgType{atSObj}},
	{Name: "observe.objs", ArgTypes: []ArgType{atSObj}},
	{Name: "distance", ArgTypes: []ArgType{atSObj, atSObj, atVar}},
	{Name: "stop.cycling", ArgTypes: []ArgType{atSObj}},
	{Name: "start.cycling", ArgTypes: []ArgType{atSObj}},
	{Name: "normal.cycle", ArgTypes: []ArgType{atSObj}},
	{Name: "end.of.loop", ArgTypes: []ArgType{atSObj, atFlag}},
	{Name: "reverse.cycle", ArgTypes: []ArgType{atSObj}},
	{Name: "reverse.loop", ArgTypes: []ArgType{atSObj, atFlag}},
	{Name: "cycle.time", ArgTypes: []ArgType{atSObj, atVar}},
	{Name: "stop.motion", ArgTypes: []ArgType{atSObj}},
	{Name: "start.motion", ArgTypes: []ArgType{atSObj}},
	{Name: "step.size", ArgTypes: []ArgType{atSObj, atVar}},
	{Name: "step.time", ArgTypes: []ArgType{atSObj, atVar}},
	{Name: "move.obj", ArgTypes: []ArgType{atSObj, atNum, atNum, atNum, atFlag}},
	{Name: "move.obj.v", ArgTypes: []ArgType{atSObj, atVar, atVar, atNum, atFlag}},
	{Name: "follow.ego", ArgTypes: []ArgType{atSObj, atNum, atFlag}},
	{Name: "wander", ArgTypes: []ArgType{atSObj}},
	{Name: "normal.motion", ArgTypes: []ArgType{atSObj}},
	{Name: "set.dir", ArgTypes: []ArgType{atSObj, atVar}},
	{Name: "get.dir", ArgTypes: []ArgType{atSObj, atVar}},
	{Name: "ignore.blocks", ArgTypes: []ArgType{atSObj}},
	{Name: "observe.blocks", ArgTypes: []ArgType{atSObj}},
	{Name: "block", ArgTypes: []ArgType{atNum, atNum, atNum, atNum}},
	{Name: "unblock"},
	{Name: "get", ArgTypes: []ArgType{atIObj}},
	{Name: "get.v", ArgTypes: []ArgType{atVar}},
	{Name: "drop", ArgTypes: []ArgType{atIObj}},
	{Name: "put", ArgTypes: []ArgType{atIObj, atVar}},
	{Name: "put.v", ArgTypes: []ArgType{atVar, atVar}},
	{Name: "get.room.v", ArgTypes: []ArgType{atVar, atVar}},
	{Name: "load.sound", ArgTypes: []ArgType{atNum}},
	{Name: "sound", ArgTypes: []ArgType{atNum, atFlag}},
	{Name: "stop.sound"},
	{Name: "print", ArgTypes: []ArgType{atMsg}},
	{Name: "print.v", ArgTypes: []ArgType{atVar}},
	{Name: "display", ArgTypes: []ArgType{atNum, atNum, atMsg}},
	{Name: "display.v", ArgTypes: []ArgType{atVar, atVar, atVar}},
	{Name: "clear.lines", ArgTypes: []ArgType{atNum, atNum, atNum}},
	{Name: "text.screen"},
	{Name: "graphics"},
	{Name: "set.cursor.char", ArgTypes: []ArgType{atMsg}},
	{Name: "set.text.attribute", ArgTypes: []ArgType{atNum, atNum}},
	{Name: "shake.screen", ArgTypes: []ArgType{atNum}},
	{Name: "configure.screen", ArgTypes: []ArgType{atNum, atNum, atNum}},
	{Name: "status.line.on"},
	{Name: "status.line.off"},
	{Name: "set.string", ArgTypes: []ArgType{atStr, atMsg}},
	{Name: "get.string", ArgTypes: []ArgType{atStr, atMsg, atNum, atNum, atNum}},
	{Name: "word.to.string", ArgTypes: []ArgType{atWord, atStr}},
	{Name: "parse", ArgTypes: []ArgType{atStr}},
	{Name: "get.num", ArgTypes: []ArgType{atMsg, atVar}},
	{Name: "prevent.input"},
	{Name: "accept.input"},
	{Name: "set.key", ArgTypes: []ArgType{atNum, atNum, atCtrl}},
	{Name: "add.to.pic", ArgTypes: []ArgType{atNum, atNum, atNum, atNum, atNum, atNum, atNum}},
	{Name: "add.to.pic.v", ArgTypes: []ArgType{atVar, atVar, atVar, atVar, atVar, atVar, atVar}},
	{Name: "status"},
	{Name: "save.game"},
	{Name: "restore.game"},
	{Name: "init.disk"},
	{Name: "restart.game"},
	{Name: "show.obj", ArgTypes: []ArgType{atNum}},
	{Name: "random", ArgTypes: []ArgType{atNum, atNum, atVar}},
	{Name: "program.control"},
	{Name: "player.control"},
	{Name: "obj.status.v", ArgTypes: []ArgType{atVar}},
	{Name: "quit", ArgTypes: []ArgType{atNum}},
	{Name: "show.mem"},
	{Name: "pause"},
	{Name: "echo.line"},
	{Name: "cancel.line"},
	{Name: "init.joy"},
	{Name: "toggle.monitor"},
	{Name: "version"},
	{Name: "script.size", ArgTypes: []ArgType{atNum}},
	{Name: "set.game.id", ArgTypes: []ArgType{atMsg}},
	{Name: "log", ArgTypes: []ArgType{atMsg}},
	{Name: "set.scan.start"},
	{Name: "reset.scan.start"},
	{Name: "reposition.to", ArgTypes: []ArgType{atSObj, atNum, atNum}},
	{Name: "reposition.to.v", ArgTypes: []ArgType{atSObj, atVar, atVar}},
	{Name: "trace.on"},
	{Name: "trace.info", ArgTypes: []ArgType{atNum, atNum, atNum}},
	{Name: "print.at", ArgTypes: []ArgType{atMsg, atNum, atNum, atNum}},
	{Name: "print.at.v", ArgTypes: []ArgType{atMsg, atVar, atVar, atVar}},
	{Name: "discard.view.v", ArgTypes: []ArgType{atVar}},
	{Name: "clear.text.rect", ArgTypes: []ArgType{atNum, atNum, atNum, atNum, atNum}},
	{Name: "set.upper.left"},
	{Name: "set.menu", ArgTypes: []ArgType{atMsg}},
	{Name: "set.menu.item", ArgTypes: []ArgType{atMsg, atCtrl}},
	{Name: "submit.menu"},
	{Name: "enable.item", ArgTypes: []ArgType{atCtrl}},
	{Name: "disable.item", ArgTypes: []ArgType{atCtrl}},
	{Name: "menu.input"},
	{Name: "show.obj.v", ArgTypes: []ArgType{atVar}},
	{Name: "open.dialogue"},
	{Name: "close.dialogue"},
	{Name: "mul.n", ArgTypes: []ArgType{atVar, atNum}},
	{Name: "mul.v", ArgTypes: []ArgType{atVar, atVar}},
	{Name: "div.n", ArgTypes: []ArgType{atVar, atNum}},
	{Name: "div.v", ArgTypes: []ArgType{atVar, atVar}},
	{Name: "close.window"},
	{Name: "unknown170"},
	{Name: "unknown171"},
	{Name: "unknown172"},
	{Name: "unknown173"},
	{Name: "unknown174"},
	{Name: "unknown175"},
	{Name: "unknown176"},
	{Name: "unknown177"},
	{Name: "unknown178"},
	{Name: "unknown179"},
	{Name: "unknown180"},
	{Name: "unknown181"},
}
