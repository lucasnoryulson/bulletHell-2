// Construido como parte da disciplina: Sistemas Distribuidos - PUCRS - Escola Politecnica
// Professor: Fernando Dotti  (https://fldotti.github.io/)
// Implementação do algoritmo de Exclusão Mútua Distribuída (Ricart-Agrawala)

package DIMEX

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/brunobaa/bullethell/PP2PLink"
)

// Tipos de estado

type State int

const (
	noMX State = iota
	wantMX
	inMX
)

type dmxReq int

const (
	ENTER dmxReq = iota
	EXIT
)

type dmxResp struct{}

// Estrutura do módulo DIMEX

type DIMEX_Module struct {
	Req       chan dmxReq
	Ind       chan dmxResp
	processes []string
	id        int
	st        State
	waiting   []bool
	lcl       int
	reqTs     int
	nbrResps  int
	dbg       bool
	mutex     sync.Mutex

	Pp2plink *PP2PLink.PP2PLink
}

// Inicialização

func NewDIMEX(_addresses []string, _id int, _dbg bool) *DIMEX_Module {
	p2p := PP2PLink.NewPP2PLink(_addresses[_id], _dbg)

	dmx := &DIMEX_Module{
		Req:       make(chan dmxReq, 1),
		Ind:       make(chan dmxResp, 1),
		processes: _addresses,
		id:        _id,
		st:        noMX,
		waiting:   make([]bool, len(_addresses)),
		lcl:       0,
		reqTs:     0,
		nbrResps:  0,
		dbg:       _dbg,
		Pp2plink:  p2p,
	}

	dmx.Start()
	dmx.outDbg("Init DIMEX!")
	return dmx
}

// Execução paralela

func (module *DIMEX_Module) Start() {
	go func() {
		for {
			select {
			case dmxR := <-module.Req:
				if dmxR == ENTER {
					module.outDbg("app pede mx")
					module.handleUponReqEntry()
				} else if dmxR == EXIT {
					module.outDbg("app libera mx")
					module.handleUponReqExit()
				}
			case msgOutro := <-module.Pp2plink.Ind:
				if strings.Contains(msgOutro.Message, "respOk") {
					module.outDbg("<<<---- recebe respOk de outro")
					module.handleUponDeliverRespOk(msgOutro)
				} else if strings.Contains(msgOutro.Message, "reqEntry") {
					module.outDbg("<<<---- recebe reqEntry de outro")
					module.handleUponDeliverReqEntry(msgOutro)
				}
			}
		}
	}()
}

// ENTRADA NA SESSÃO CRÍTICA

func (module *DIMEX_Module) handleUponReqEntry() {
	module.mutex.Lock()
	defer module.mutex.Unlock()
	module.lcl++
	module.reqTs = module.lcl
	module.nbrResps = 0
	module.st = wantMX

	for i, addr := range module.processes {
		if i != module.id {
			msg := fmt.Sprintf("reqEntry,%d,%d", module.id, module.reqTs)
			module.sendToLink(addr, msg, "  ")
		}
	}
}

// SAÍDA DA SESSÃO CRÍTICA

func (module *DIMEX_Module) handleUponReqExit() {
	module.mutex.Lock()
	defer module.mutex.Unlock()
	for i, waiting := range module.waiting {
		if waiting {
			addr := module.processes[i]
			msg := "respOk"
			module.sendToLink(addr, msg, "  ")
			module.waiting[i] = false
		}
	}
	module.st = noMX
}

// RECEBEU respOk

func (module *DIMEX_Module) handleUponDeliverRespOk(msgOutro PP2PLink.PP2PLink_Ind_Message) {
	module.mutex.Lock()
	defer module.mutex.Unlock()
	module.nbrResps++
	if module.nbrResps == len(module.processes)-1 {
		module.st = inMX
		module.Ind <- dmxResp{}
	}
}

// RECEBEU reqEntry

func (module *DIMEX_Module) handleUponDeliverReqEntry(msgOutro PP2PLink.PP2PLink_Ind_Message) {
	module.mutex.Lock()
	defer module.mutex.Unlock()
	fields := strings.Split(msgOutro.Message, ",")
	rid := atoi(fields[1])
	rts := atoi(fields[2])
	module.lcl = max(module.lcl, rts) + 1

	sendOk := false
	if module.st == noMX {
		sendOk = true
	} else if module.st == wantMX {
		if before(module.reqTs, module.id, rts, rid) {
			sendOk = true
		}
	}

	if sendOk {
		module.sendToLink(module.processes[rid], "respOk", "  ")
	} else {
		for i, addr := range module.processes {
			if strings.Contains(msgOutro.From, addr) {
				module.waiting[i] = true
				break
			}
		}

	}
}

func before(ts1, id1, ts2, id2 int) bool {
	if ts1 < ts2 {
		return true
	} else if ts1 == ts2 {
		return id1 < id2
	}
	return false
}

func atoi(s string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return n
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (module *DIMEX_Module) sendToLink(address string, content string, space string) {
	module.outDbg(space + " ---->>>>   to: " + address + "     msg: " + content)
	module.Pp2plink.Req <- PP2PLink.PP2PLink_Req_Message{
		To:      address,
		Message: content,
	}
}

func (module *DIMEX_Module) outDbg(s string) {
	if module.dbg {
		fmt.Println(". . . . . . . . . . . . [ DIMEX : " + s + " ]")
	}
}
