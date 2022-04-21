//Arturs Lazebniks 201RDB223
package main

import (
	"fmt"
	"strconv"
	"math"
	"math/rand"
	"time"
	"strings"
	
	"net/http"
	"io"
	"encoding/json"
	_ "embed"
)

type (
	GameGraph[T comparable] struct {
		start *Node[T]
	}
	
	Node[T comparable] struct {
		currentState GameState[T]
		nexts []*Node[T]
		EvalValue EvaluationValue
	}
	
	EvaluationValue struct {
		typeOfValue EvaluationType
		value int
	}
	
	EvaluationType int
	MinimaxVal int
	PlayerOrder int
	Path[T comparable] []*Node[T]
	
	GameState[T comparable] struct {
		firstPlayerHand Hand[T]
		secondPlayerHand Hand[T]
		alterableObject T
	}
	
	GameType[T comparable] struct {
		firstPlayerHand Hand[T]
		secondPlayerHand Hand[T]
		firstPlayerWinCond func(x T) bool
		secondPlayerWinCond func(x T) bool
	}
	
	Hand[T comparable] map[string](Operation[T])
	
	Operation[T comparable] struct {
		title string
		action func(x T) T
		condition func(x T) bool
	}
	
	MyByte byte
	MyFloat float64
	MyUint uint
)

func (b MyByte) String() string {
	return fmt.Sprintf("%08b", b)
}

func (f MyFloat) String() string {
	return fmt.Sprintf("%.4f", f)
}

func (f MyUint) String() string {
	return fmt.Sprintf("%d", f)
}

type AnyGameType interface{
		GameTypeImplementation()
}
	
var GameModes map[string]AnyGameType = map[string]AnyGameType{
	"Natural Numbers": GameType[MyUint]{
		firstPlayerHand: Hand[MyUint]{
			"add 1":Operation[MyUint]{"add 1",func(x MyUint) MyUint{return x+1},func(x MyUint) bool{return true}},
			"multiply by 3":Operation[MyUint]{"multiply by 3",func(x MyUint) MyUint{return x*3},func(x MyUint) bool{return true}},
			"substract 1":Operation[MyUint]{"substract 1",func(x MyUint) MyUint{return x-1},func(x MyUint) bool{return x>1}},
			"divide by 2":Operation[MyUint]{"divide by 2",func(x MyUint) MyUint{return x/2},func(x MyUint) bool{return x%2==0}},
			"count bits":Operation[MyUint]{"count bits",func(x MyUint) MyUint{return MyUint(math.Log2(float64(x)))+1},func(x MyUint) bool{return x>0}},
			"reverse":Operation[MyUint]{"reverse",func(x MyUint) MyUint{return reverse(x)},func(x MyUint) bool{return len(fmt.Sprint(x))>1}}},
		secondPlayerHand: Hand[MyUint]{
			"add 1":Operation[MyUint]{"add 1",func(x MyUint) MyUint{return x+1},func(x MyUint) bool{return true}},
			"multiply by 3":Operation[MyUint]{"multiply by 3",func(x MyUint) MyUint{return x*3},func(x MyUint) bool{return true}},
			"substract 1":Operation[MyUint]{"substract 1",func(x MyUint) MyUint{return x-1},func(x MyUint) bool{return x>1}},
			"divide by 2":Operation[MyUint]{"divide by 2",func(x MyUint) MyUint{return x/2},func(x MyUint) bool{return x%2==0}},
			"count bits":Operation[MyUint]{"count bits",func(x MyUint) MyUint{return MyUint(math.Log2(float64(x)))+1},func(x MyUint) bool{return x>0}},
			"reverse":Operation[MyUint]{"reverse",func(x MyUint) MyUint{return reverse(x)},func(x MyUint) bool{return len(fmt.Sprint(x))>1}}},
		firstPlayerWinCond: func(x MyUint) bool{return x%2==0},
		secondPlayerWinCond: func(x MyUint) bool{return x%2==1}},
		"Byte": GameType[MyByte]{
		firstPlayerHand: Hand[MyByte]{
			"or 00001111":Operation[MyByte]{"or 00001111",func(x MyByte) MyByte{return x|0b00001111},func(x MyByte) bool{return true}},
			"or 11110000":Operation[MyByte]{"or 11110000",func(x MyByte) MyByte{return x|0b11110000},func(x MyByte) bool{return true}},
			"and 00001111":Operation[MyByte]{"and 00001111",func(x MyByte) MyByte{return x&0b00001111},func(x MyByte) bool{return true}},
			"and 11110000":Operation[MyByte]{"and 11110000",func(x MyByte) MyByte{return x&0b11110000},func(x MyByte) bool{return true}},
			"not":Operation[MyByte]{"not",func(x MyByte) MyByte{return ^x},func(x MyByte) bool{return strings.Count(fmt.Sprint(x),"0")!=strings.Count(fmt.Sprint(x),"1")}},
			"xor 01010101":Operation[MyByte]{"xor 01010101",func(x MyByte) MyByte{return x^0b01010101},func(x MyByte) bool{return strings.Count(fmt.Sprint(x),"0")==strings.Count(fmt.Sprint(x),"1")}}},
		secondPlayerHand: Hand[MyByte]{
			"or 00001111":Operation[MyByte]{"or 00001111",func(x MyByte) MyByte{return x|0b00001111},func(x MyByte) bool{return true}},
			"or 11110000":Operation[MyByte]{"or 11110000",func(x MyByte) MyByte{return x|0b11110000},func(x MyByte) bool{return true}},
			"and 00001111":Operation[MyByte]{"and 00001111",func(x MyByte) MyByte{return x&0b00001111},func(x MyByte) bool{return true}},
			"and 11110000":Operation[MyByte]{"and 11110000",func(x MyByte) MyByte{return x&0b11110000},func(x MyByte) bool{return true}},
			"not x":Operation[MyByte]{"not x",func(x MyByte) MyByte{return ^x},func(x MyByte) bool{return strings.Count(fmt.Sprint(x),"0")!=strings.Count(fmt.Sprint(x),"1")}},
			"xor 01010101":Operation[MyByte]{"xor 01010101",func(x MyByte) MyByte{return x^0b01010101},func(x MyByte) bool{return strings.Count(fmt.Sprint(x),"0")==strings.Count(fmt.Sprint(x),"1")}}},
		firstPlayerWinCond: func(x MyByte) bool{return strings.Count(fmt.Sprint(x),"0")>4},
		secondPlayerWinCond: func(x MyByte) bool{return strings.Count(fmt.Sprint(x),"1")>4}},
		"[0;1]": GameType[MyFloat]{
		firstPlayerHand: Hand[MyFloat]{
			"x^2":Operation[MyFloat]{"x^2",func(x MyFloat) MyFloat{return x*x},func(x MyFloat) bool{return x<0.5}},
			"sqtr(x)":Operation[MyFloat]{"sqtr(x)",func(x MyFloat) MyFloat{return MyFloat(math.Sqrt(float64(x)))},func(x MyFloat) bool{return x>0.5}},
			"1-x":Operation[MyFloat]{"1-x",func(x MyFloat) MyFloat{return 1-x},func(x MyFloat) bool{return true}},
			"x/2":Operation[MyFloat]{"x/2",func(x MyFloat) MyFloat{return x/2},func(x MyFloat) bool{return true}},
			"x*2":Operation[MyFloat]{"x*2",func(x MyFloat) MyFloat{return x*2},func(x MyFloat) bool{return x*2<=1}},
			"avg(1,x)":Operation[MyFloat]{"avg(1,x)",func(x MyFloat) MyFloat{return MyFloat((1+float64(x))/2)},func(x MyFloat) bool{return x<0.5}}},
		secondPlayerHand: Hand[MyFloat]{
			"x^2":Operation[MyFloat]{"x^2",func(x MyFloat) MyFloat{return x*x},func(x MyFloat) bool{return x<0.5}},
			"sqtr(x)":Operation[MyFloat]{"sqtr(x)",func(x MyFloat) MyFloat{return MyFloat(math.Sqrt(float64(x)))},func(x MyFloat) bool{return x>0.5}},
			"1-x":Operation[MyFloat]{"1-x",func(x MyFloat) MyFloat{return 1-x},func(x MyFloat) bool{return true}},
			"x/2":Operation[MyFloat]{"x/2",func(x MyFloat) MyFloat{return x/2},func(x MyFloat) bool{return true}},
			"x*2":Operation[MyFloat]{"x*2",func(x MyFloat) MyFloat{return x*2},func(x MyFloat) bool{return x*2<=1}},
			"avg(0,x)":Operation[MyFloat]{"avg(0,x)",func(x MyFloat) MyFloat{return MyFloat((0+float64(x))/2)},func(x MyFloat) bool{return x>0.5}}},
		firstPlayerWinCond: func(x MyFloat) bool{return x>0.5},
		secondPlayerWinCond: func(x MyFloat) bool{return x<0.5}}}

type settings struct {
	Mode string `json:"mode"`
	Order string `json:"order"`
}

type buttonUse struct {
	Can bool `json:"can"`
	Val string `json:"val"`
}

type result struct {
	Html string `json:"html"`
	Op string `json:"op"`
}

var game bool = false

var opinput chan string = make(chan string,10)
var opoutput chan string = make(chan string,10)

const (
	none EvaluationType = 2<<iota
	strict
	
	minimize MinimaxVal = 3<<iota
	maximize
	
	firstP PlayerOrder = iota+1
	secondP
)

//go:embed CORS/index.html
var mainhtml string

func mainhtmlHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Access-Control-Allow-Origin","*")
	if req.Method == "GET" {
		io.WriteString(w, mainhtml)
	}else if req.Method == "POST" {
		data, err := io.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {return}
		var info settings
		json.Unmarshal(data, &info)
		var newhtml,op string
		switch info.Mode {
			case "Natural Numbers": newhtml,op=SetUpGame[MyUint](info)
			case "Byte": newhtml,op=SetUpGame[MyByte](info)
			case "[0;1]": newhtml,op=SetUpGame[MyFloat](info)
		}
		d, _ := json.Marshal(result{newhtml, op})
		io.WriteString(w, string(d))
	} else {w.WriteHeader(405)}
}

//go:embed CORS/js/script.js
var mainjs string

func mainscriptHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Access-Control-Allow-Origin","*")
	if req.Method == "GET" {io.WriteString(w, mainjs); game=false} else {w.WriteHeader(405)}
}

//go:embed CORS/game.html
var gamehtmltempl string

//go:embed CORS/js/game.js
var gamejs string

func gamescriptHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Access-Control-Allow-Origin","*")
	if req.Method == "GET" {io.WriteString(w, gamejs)} else {w.WriteHeader(405)}
}

func gameHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Access-Control-Allow-Origin","*")
	if req.Method == "POST" {
		data, err := io.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {return}
		inf:=string(data)
		var toSent buttonUse
		opinput<-inf
		v:=<-opoutput
		if v != "no" {
			toSent.Can=true
			toSent.Val=v
		} else {
			toSent.Can=false
		}
		str, _:=json.Marshal(toSent)
		io.WriteString(w, string(str))
	} else if req.Method == "GET" {
		out:=<-opoutput
		io.WriteString(w, out)
	} else {
		w.WriteHeader(405)
	}
}

func reverse(a MyUint) MyUint {
	var s string
	for _, ch := range fmt.Sprint(a) {
		s=string(ch)+s
	}
	i,_:=strconv.Atoi(s)
	return MyUint(i)
}

func main() {
	rand.Seed(int64(time.Now().Nanosecond()))
	http.HandleFunc("/", mainhtmlHandler)
	http.HandleFunc("/js/script.js", mainscriptHandler)
	http.HandleFunc("/js/game.js", gamescriptHandler)
	http.HandleFunc("/game/", gameHandler)
	
	err := http.ListenAndServe(":8080", nil)
	panic(err)
}

func SetUpGame[T comparable](s settings) (string,string) {
	var GM GameType[T]
	var Graph GameGraph[T]
	GM=GameModes[s.Mode].(GameType[T])
	v:=GenRand[T]()
	Graph=GM.GenerateGameGraph(v)
	var PlayerIs PlayerOrder
	if s.Order=="first" {
		PlayerIs = firstP
	} else if s.Order=="second" {
		PlayerIs = secondP
	} else {
		panic("undefined player")
	}
	var BotIs PlayerOrder
	BotIs=PlayerIs.Next()
	Graph.Evaluate(BotIs, GM)
	paths:=Graph.start.FindPath(Graph.start.EvalValue.value)
	
	var firstPlOps, secondPlOps, BotsHand, PlayerHand []string
	for k, _ := range GM.firstPlayerHand {
		firstPlOps = append(firstPlOps, k)
	}
	for k, _ := range GM.secondPlayerHand {
		secondPlOps = append(secondPlOps, k)
	}
	if BotIs == firstP {
		BotsHand, PlayerHand = firstPlOps, secondPlOps
	} else if BotIs == secondP {
		BotsHand, PlayerHand = secondPlOps, firstPlOps
	} else {
		panic("undefined player")
	}
	var BotsHandMap, PlayerHandMap map[string]string = make(map[string]string, 6),make(map[string]string, 6)
	BotsHandMap["botop1"]=BotsHand[0];PlayerHandMap["op1"]=PlayerHand[0]
	BotsHandMap["botop2"]=BotsHand[1];PlayerHandMap["op2"]=PlayerHand[1]
	BotsHandMap["botop3"]=BotsHand[2];PlayerHandMap["op3"]=PlayerHand[2]
	BotsHandMap["botop4"]=BotsHand[3];PlayerHandMap["op4"]=PlayerHand[3]
	BotsHandMap["botop5"]=BotsHand[4];PlayerHandMap["op5"]=PlayerHand[4]
	BotsHandMap["botop6"]=BotsHand[5];PlayerHandMap["op6"]=PlayerHand[5]
	var playerHand Hand[T]
	var botHand Hand[T]
	if PlayerIs==firstP {
		playerHand=GM.firstPlayerHand
		botHand=GM.secondPlayerHand
	} else {
		playerHand=GM.secondPlayerHand
		botHand=GM.firstPlayerHand
	}
	op:=""
	if BotIs==firstP {
		ri:=rand.Intn(len(paths))
		op=paths[ri][0]
		v=GM.firstPlayerHand[op].action(v)
		botHand.Remove(botHand[op])
		paths=Update(paths, op)
		Graph.start.Move(op)
		for k, v := range BotsHandMap {
			if v==op {
				op=k
				break
			}
		}
	}
	game=true
	go PlayGame[T](GM, Graph, paths, BotsHandMap, PlayerHandMap, PlayerIs, botHand, playerHand, v)
	var plWinC string
	switch any(v).(type) {
		case MyUint: 
			if PlayerIs == firstP {
				plWinC="x%2==0"
			} else {
				plWinC="x%2==1"
			}
		case MyByte:
			if PlayerIs == firstP {
				plWinC="count of '0' > count of '1'"
			} else {
				plWinC="count of '1' > count of '0'"
			}
		case MyFloat:
			if PlayerIs == firstP {
				plWinC="x>0.5"
			} else {
				plWinC="x<0.5"
			}
	}
	return GenGameHTML(BotsHand, PlayerHand, fmt.Sprint(v), plWinC),op
}

func GenGameHTML(firstPlOps, secondPlOps []string, initVal, plWinC string) string {
	return fmt.Sprintf(gamehtmltempl,
		firstPlOps[0],firstPlOps[1],firstPlOps[2],firstPlOps[3],firstPlOps[4],firstPlOps[5],
		initVal, plWinC,
		secondPlOps[0],secondPlOps[1],secondPlOps[2],secondPlOps[3],secondPlOps[4],secondPlOps[5])
}

func PlayGame[T comparable](GM GameType[T], Graph GameGraph[T], paths [][]string, BotsHandMap, PlayerHandMap map[string]string, PlayerIs PlayerOrder, botHand, playerHand Hand[T], v T) {
	var currPlayer = PlayerIs
	CurrNode:=Graph.start
	endGame:=false
	for game {
		if endGame {
			game=false
			break
		}
		if PlayerIs == currPlayer {
			op:=playerHand[PlayerHandMap[<-opinput]]
			if op.condition(v) {
				v=op.action(v)
				playerHand.Remove(op)
				paths=Update(paths, op.title)
				endGame=!CurrNode.Move(op.title)
				if !endGame && len(paths)==0 {
					paths=CurrNode.FindPath(CurrNode.EvalValue.value)
				}
				opoutput<-fmt.Sprint(v)
				currPlayer=currPlayer.Next()
				if !CheckHand[T](botHand, v) {
					endGame=true
					game=false
					break
				}
			} else {
				opoutput<-"no"
				if endGame {
					game=false
					break
				}
			}
		} else {
			if endGame {
				game=false
				break
			}
			if len(paths)==0 {
				paths=CurrNode.FindPath(CurrNode.EvalValue.value)
			}
			if len(paths)==0 {
				endGame=true
				game=false
				break
			}
			ri:=rand.Intn(len(paths))
			op:=paths[ri][0]
			if _, ok := botHand[op]; !ok {
				endGame=true
				game=false
				break
			}
			v=botHand[op].action(v)
			botHand.Remove(botHand[op])
			paths=Update(paths, op)
			endGame=!CurrNode.Move(op)
			if !endGame && len(paths)==0 {
				paths=CurrNode.FindPath(CurrNode.EvalValue.value)
			}
			
			var tagid string
			for k, v := range BotsHandMap {
				if v==op {
					tagid=k
					break
				}
			}
			if endGame {
				game=false
				break
			} else {
				opoutput<-fmt.Sprintf(`{"op":"%s","val":"%s"}`, tagid, v)
			}
			currPlayer=currPlayer.Next()
			if !CheckHand[T](playerHand, v) {
				endGame=true
				game=false
				break
			}
		}
	}
	fpwon:=GM.firstPlayerWinCond(v)
	spwon:=GM.secondPlayerWinCond(v)
	if !fpwon && !spwon {
		opoutput<-"draw"
	}else if (PlayerIs == firstP && fpwon) || (PlayerIs == secondP && spwon) {
		opoutput<-"win"
	} else {
		opoutput<-"lose"
	}
}

func CheckHand[T comparable](h Hand[T], v T) bool {
	var a bool = false
	for _, op := range h {
		a=a||op.condition(v)
	}
	return a
}

func (n *Node[T]) Move(op string) bool {
	if len(n.nexts)==0 {return false}
	for _, next := range n.nexts {
		fphd:=Diff[T](n.currentState.firstPlayerHand,next.currentState.firstPlayerHand)
		sphd:=Diff[T](n.currentState.secondPlayerHand,next.currentState.secondPlayerHand)
		if len(fphd) > 0 {
			if fphd[0]==op {
				*n=*next
				if len(n.nexts)==0 {return false}
				return true
			} 
		} else if len(sphd) > 0 {
			if sphd[0]==op {
				*n=*next
				if len(n.nexts)==0 {return false}
				return true			
			}
		}
	}
	return true
}

func Update(p [][]string, op string) [][]string {
	for i:=0; i <len(p); {
		path:=p[i]
		if path[0]!=op {
			if i < len(p)-1 {
				p=append(p[:i],p[i+1:]...)
			} else {
				p=p[:i]
			}
		} else {
			i++
		}
	}
	
	for i:=0; i < len(p); {
		if len(p[i])>1 {
			p[i]=p[i][1:]
			i++
		} else {
			if i < len(p)-1 {
				p=append(p[:i],p[i+1:]...)
			} else {
				p=p[:i]
			}
		}
	}
	
	return p
}

func GenRand[T comparable]() T {
	var filter T
	var val any = filter
	
	switch val.(type) {
		case MyUint:
			val = MyUint(rand.Intn(10)+1) //[1;10]
		case MyFloat:
			val = MyFloat(rand.Float64()/2+0.25) //[0.25;0.75)
		case MyByte:
			var indexes [8]int = [8]int{0,1,2,3,4,5,6,7}
			rand.Shuffle(8, func (i,j int) {indexes[i],indexes[j]=indexes[j],indexes[i]})
			var n byte
			for i := 0; i < 4; i++ {
				n|=1<<indexes[i]
			}
			val=MyByte(n) //randomly located 4 '1' and 4 '0'
	}
	
	return val.(T)
}

func (GT GameType[T]) GameTypeImplementation() {}

func (GType GameType[T]) GenerateGameGraph(initVal T) GameGraph[T] {
	var SG GameGraph[T]
	SG.start = new(Node[T])
	SG.start.EvalValue=EvaluationValue{none,0}
	SG.start.currentState = GameState[T]{firstPlayerHand:GType.firstPlayerHand, secondPlayerHand:GType.secondPlayerHand, alterableObject:initVal}
	SG.OptimizedExpand(firstP, []*Node[T]{SG.start})
	return SG
}

func (SG GameGraph[T]) OptimizedExpand(currentPlayer PlayerOrder, currentLayer []*Node[T]) {
	var newLayer []*Node[T]
	for _, n := range currentLayer {
		var nexts []*Node[T]
		var currentPlayerHand Hand[T]
		obj := n.currentState.alterableObject
		if currentPlayer == firstP {
			currentPlayerHand = MapCopy[T](n.currentState.firstPlayerHand)
		} else if currentPlayer == secondP {
			currentPlayerHand = MapCopy[T](n.currentState.secondPlayerHand)
		} else {panic("undefined player")}
		
		for _, op := range currentPlayerHand {
			if op.condition(obj) {
				var newState GameState[T]
				newState.alterableObject = op.action(obj)
				newPlayerHand := MapCopy[T](currentPlayerHand)
				newPlayerHand.Remove(op)
				if currentPlayer == firstP {
					newState.firstPlayerHand = newPlayerHand
					newState.secondPlayerHand = n.currentState.secondPlayerHand
				} else if currentPlayer == secondP {
					newState.secondPlayerHand = newPlayerHand
					newState.firstPlayerHand = n.currentState.firstPlayerHand
				} else {panic("undefined player")}
				if found, node := newState.FindIn(newLayer); !found {
					node=&Node[T]{newState, make([]*Node[T],0), EvaluationValue{none, 0}}
					nexts = append(nexts, node)
					newLayer = append(newLayer, node)
				} else {
					nexts = append(nexts, node)
				}
			}
		}
		n.nexts = nexts
	}
	
	if len(newLayer)>0 {
		SG.OptimizedExpand(currentPlayer.Next(), newLayer)
	}
}

func MapCopy[T comparable](h Hand[T]) Hand[T] {
	newh:=make(Hand[T], len(h))
	for k, v := range h {
		newh[k]=v
	}
	return newh
}

func (SG GameGraph[T]) Evaluate(BotsOrder PlayerOrder, GType GameType[T]) {
	var minmax MinimaxVal
	if BotsOrder == firstP {
		minmax = maximize
	} else if BotsOrder == secondP {
		minmax = minimize
	} else {
		panic("undefined player")
	}
	SG.start.Evaluate(minmax, GType, BotsOrder)
}

func (n *Node[T]) Evaluate(layer MinimaxVal, GType GameType[T], BotsOrder PlayerOrder) int {
	if len(n.nexts) == 0 {
		v:=n.currentState.alterableObject
		fpw:=GType.firstPlayerWinCond(v)
		spw:=GType.secondPlayerWinCond(v)
		if !fpw && !spw {
			n.EvalValue=EvaluationValue{strict, 0}
		} else if BotsOrder == firstP && fpw {
			n.EvalValue=EvaluationValue{strict, 1}
		} else {
			n.EvalValue=EvaluationValue{strict, -1}
		}
	}
	
	for _, next := range n.nexts {
		v:=next.Evaluate(layer.Next(), GType, BotsOrder)
		if n.EvalValue.typeOfValue == none {
			n.EvalValue = EvaluationValue{strict, v}
		} else {
			if layer == minimize {
				n.EvalValue.value = min(n.EvalValue.value, v)
			} else if layer == maximize {
				n.EvalValue.value = max(n.EvalValue.value, v)
			} else {
				panic("undefined minimax value")
			}
		}
	}
	
	return n.EvalValue.value
}

func (n *Node[T]) FindPath(val int) [][]string {
	var paths [][]string
	if len(n.nexts) == 0 {
		return paths
	}
	for _, next := range n.nexts {
		if next.EvalValue.value == val && next.EvalValue.typeOfValue==strict {
			ops:=Diff[T](n.currentState.firstPlayerHand, next.currentState.firstPlayerHand)
			if len(ops)==0 {
				ops=Diff[T](n.currentState.secondPlayerHand, next.currentState.secondPlayerHand)
			}
			if len(ops)==0 {
				continue
			}
			k:=ops[0]
			newPaths:=next.FindPath(val)
			for i:=range newPaths {
				newPaths[i]=append([]string{k},newPaths[i]...)
			}
			if len(newPaths)==0 {
				newPaths=[][]string{[]string{k}}
			}
			paths = append(paths, newPaths...)
		}
	}
	return paths
}

func Diff[T comparable](a, b Hand[T]) []string {
	var res []string
	for k, _ := range a {
		if _, ok := b[k]; !ok {
			res = append(res,k)
		}
	}
	return res
}

func min(a, b int) int {
	if a<b {return a} else {return b}
}

func max(a, b int) int {
	if a>b {return a} else {return b}
}

func (m MinimaxVal) Next() MinimaxVal {
	if m == minimize {
		return maximize
	} else if m == maximize {
		return minimize
	} else {
		panic("undefined minimax value")
		return 0
	}
}

func (p PlayerOrder) Next() PlayerOrder {
	//p%2+1
	if p == firstP {
		return secondP
	} else if p == secondP {
		return firstP
	} else {
		panic("undefined player")
		return 0
	}
}

func (h *Hand[T]) Remove(op Operation[T]) {
	delete(*h,op.title)
}

func (o Operation[T]) Equal(op Operation[T]) bool {
	return o.title==op.title
}

func (h1 Hand[T]) Equal(h2 Hand[T]) bool {
	if len(h1) != len(h2) {return false} 
	for k,_:=range h1 {
		if _, ok := h2[k]; !ok {return false}
		if !h1[k].Equal(h2[k]) {
			return false
		}
	}
	return true
}

func (s GameState[T]) Equal(st GameState[T]) bool {
	return s.firstPlayerHand.Equal(st.firstPlayerHand) && s.secondPlayerHand.Equal(st.secondPlayerHand) && s.alterableObject==st.alterableObject
}

func (st GameState[T]) FindIn(src []*Node[T]) (bool, *Node[T]) {
	for _, n := range src {
		if st.Equal(n.currentState) {
			return true, n
		}
	}
	return false, nil
}
