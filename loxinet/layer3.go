/*
 * Copyright (c) 2022 NetLOX Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at:
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package loxinet

import (
	"errors"
	"fmt"
	"net"

	cmn "github.com/loxilb-io/loxilb/common"
	tk "github.com/loxilb-io/loxilib"
)

// constants
const (
	L3ErrBase = iota - RtErrBase - 1000
	L3AddrErr
	L3ObjErr
)

// IfaKey - key to find a ifa entry
type IfaKey struct {
	Obj string
}

// IfaEnt - the ifa-entry
type IfaEnt struct {
	IfaAddr   net.IP
	IfaNet    net.IPNet
	Secondary bool
}

// Ifa  - a singe ifa can contain multiple ifas
type Ifa struct {
	Key  IfaKey
	Zone *Zone
	Sync DpStatusT
	Ifas []*IfaEnt
}

// L3H - context container
type L3H struct {
	IfaMap map[IfaKey]*Ifa
	Zone   *Zone
}

// L3Init - Initialize the layer3 subsystem
func L3Init(zone *Zone) *L3H {
	var nl3h = new(L3H)
	nl3h.IfaMap = make(map[IfaKey]*Ifa)
	nl3h.Zone = zone

	return nl3h
}

// IfaAdd - Adds an interface IP address (primary or secondary) and associate it with Obj
// Obj can be anything but usually it is the name of a valid interface
func (l3 *L3H) IfaAdd(Obj string, Cidr string) (int, error) {
	var sec bool = false
	addr, network, err := net.ParseCIDR(Cidr)
	if err != nil {
		return L3AddrErr, errors.New("ip address parse error")
	}

	ifObjID := -1
	pObj := l3.Zone.Ports.PortFindByName(Obj)
	if pObj != nil {
		ifObjID = pObj.SInfo.OsID
	}

	key := IfaKey{Obj}
	ifa := l3.IfaMap[key]

	if ifa == nil {
		ifa = new(Ifa)
		ifaEnt := new(IfaEnt)
		ifa.Key.Obj = Obj
		ifa.Zone = l3.Zone
		ifaEnt.IfaAddr = addr
		ifaEnt.IfaNet = *network
		ifa.Ifas = append(ifa.Ifas, ifaEnt)
		l3.IfaMap[key] = ifa

		// ifa needs related self-routes
		ra := RtAttr{0, 0, false, ifObjID}
		_, err = mh.zr.Rt.RtAdd(*network, RootZone, ra, nil)
		if err != nil {
			tk.LogIt(tk.LogDebug, "ifa add - %s:%s self-rt error\n", addr.String(), Obj)
			return L3AddrErr, errors.New("self-route add error")
		}

		ifa.DP(DpCreate)

		return 0, nil
	}

	for _, ifaEnt := range ifa.Ifas {
		if ifaEnt.IfaAddr.Equal(addr) {
			tk.LogIt(tk.LogDebug, "ifa add - exists %s:%s\n", addr.String(), Obj)
			return L3AddrErr, errors.New("ip address exists")
		}

		// if network part of an added ifa is equal to previously
		// existing ifa, then it is considered a secondary ifa
		if ifaEnt.IfaNet.IP.Equal(network.IP) {
			pfxSz1, _ := ifaEnt.IfaNet.Mask.Size()
			pfxSz2, _ := network.Mask.Size()

			if pfxSz1 == pfxSz2 {
				sec = true
			}
		}
	}

	ifaEnt := new(IfaEnt)
	ifa.Key.Obj = Obj
	ifaEnt.IfaAddr = addr
	ifaEnt.IfaNet = *network
	ifaEnt.Secondary = sec
	ifa.Ifas = append(ifa.Ifas, ifaEnt)

	// ifa needs to related self-routes
	// FIXME - Code duplication with primary address route above
	ra := RtAttr{0, 0, false, ifObjID}
	_, err = mh.zr.Rt.RtAdd(*network, RootZone, ra, nil)
	if err != nil {
		tk.LogIt(tk.LogDebug, " - %s:%s self-rt error\n", addr.String(), Obj)
		return L3AddrErr, errors.New("self-route add error")
	}

	ifa.DP(DpCreate)

	tk.LogIt(tk.LogDebug, "ifa added %s:%s\n", addr.String(), Obj)

	return 0, nil
}

// IfaDelete - Deletes an interface IP address (primary or secondary) and de-associate from Obj
// Obj can be anything but usually it is the name of a valid interface
func (l3 *L3H) IfaDelete(Obj string, Cidr string) (int, error) {
	var found bool = false
	addr, network, err := net.ParseCIDR(Cidr)
	if err != nil {
		tk.LogIt(tk.LogError, "ifa delete - malformed %s:%s\n", addr.String(), Obj)
		return L3AddrErr, errors.New("ip address parse error")
	}

	key := IfaKey{Obj}
	ifa := l3.IfaMap[key]

	if ifa == nil {
		tk.LogIt(tk.LogError, "ifa delete - no such %s:%s\n", addr.String(), Obj)
		return L3AddrErr, errors.New("no such ip address")
	}

	for index, ifaEnt := range ifa.Ifas {
		if ifaEnt.IfaAddr.Equal(addr) {

			if ifaEnt.IfaNet.IP.Equal(network.IP) {
				pfxSz1, _ := ifaEnt.IfaNet.Mask.Size()
				pfxSz2, _ := network.Mask.Size()

				if pfxSz1 == pfxSz2 {
					ifa.Ifas = append(ifa.Ifas[:index], ifa.Ifas[index+1:]...)
					found = true
				}
			}
		}
	}

	if found == true {
		// delete self-routes related to this ifa
		_, err = mh.zr.Rt.RtDelete(*network, RootZone)
		if err != nil {
			tk.LogIt(tk.LogError, "ifa delete %s:%s self-rt error\n", addr.String(), Obj)
			// Continue after logging error because there is noway to fallback
		}
		if len(ifa.Ifas) == 0 {
			delete(l3.IfaMap, ifa.Key)

			ifa.DP(DpRemove)

			tk.LogIt(tk.LogDebug, "ifa deleted %s:%s\n", addr.String(), Obj)
		}
		return 0, nil
	}

	tk.LogIt(tk.LogDebug, "ifa delete - no such %s:%s\n", addr.String(), Obj)
	return L3AddrErr, errors.New("no such ifa")
}

// IfaDeleteAll - Deletes all interface IP address (primary or secondary) and de-associate from Obj
// Obj can be anything but usually it is the name of a valid interface
func (l3 *L3H) IfaDeleteAll(Obj string) (int, error) {

	key := IfaKey{Obj}
	ifa := l3.IfaMap[key]

	if ifa != nil {
		for _, ifaEnt := range ifa.Ifas {
			ones, _ := ifaEnt.IfaNet.Mask.Size()
			cidr := fmt.Sprintf("%s/%d", ifaEnt.IfaAddr.String(), ones)
			l3.IfaDelete(Obj, cidr)
		}
		ifa.Ifas = nil
		delete(l3.IfaMap, key)
	}

	return 0, nil
}

// IfaSelect - Given any ip address, select optimal ip address from Obj's ifa list
// This is useful to determine source ip address when sending traffic
// to the given ip address
func (l3 *L3H) IfaSelect(Obj string, addr net.IP, findAny bool) (int, net.IP) {

	key := IfaKey{Obj}
	ifa := l3.IfaMap[key]

	if ifa == nil {
		return L3ObjErr, net.IPv4(0, 0, 0, 0)
	}

	for _, ifaEnt := range ifa.Ifas {
		if ifaEnt.Secondary == true {
			continue
		}

		if tk.IsNetIPv6(addr.String()) && tk.IsNetIPv4(ifaEnt.IfaNet.IP.String()) {
			continue
		}

		if ifaEnt.IfaNet.Contains(addr) {
			return 0, ifaEnt.IfaAddr
		}
	}

	if findAny == false {
		return L3AddrErr, net.IPv4(0, 0, 0, 0)
	}

	// Select first IP
	if len(ifa.Ifas) > 0 {
		return 0, ifa.Ifas[0].IfaAddr
	}

	return L3AddrErr, net.IPv4(0, 0, 0, 0)
}

// IfaSelectAny - Given any dest ip address, select optimal interface source ip address
// This is useful to determine source ip address when sending traffic to the given ip address
func (l3 *L3H) IfaSelectAny(addr net.IP, findAny bool) (int, net.IP) {
	var err int
	var tDat tk.TrieData
	var firstIP *net.IP

	v6 := false
	IfObj := ""
	firstIP = nil

	if tk.IsNetIPv4(addr.String()) {
		err, _, tDat = l3.Zone.Rt.Trie4.FindTrie(addr.String())
	} else {
		v6 = true
		err, _, tDat = l3.Zone.Rt.Trie6.FindTrie(addr.String())
	}
	if err == 0 {
		switch rtn := tDat.(type) {
		case *Neigh:
			if rtn != nil {
				IfObj = rtn.OifPort.Name
			}
		case *int:
			p := l3.Zone.Ports.PortFindByOSID(*rtn)
			if p != nil {
				IfObj = p.Name
			}
		default:
			break
		}
	}

	if IfObj != "" {
		return l3.IfaSelect(IfObj, addr, findAny)
	}

	for _, ifa := range l3.IfaMap {
		if ifa.Key.Obj == "lo" {
			continue
		}

		for _, ifaEnt := range ifa.Ifas {
			if ifaEnt.Secondary == true {
				continue
			}

			if v6 && tk.IsNetIPv4(ifaEnt.IfaNet.IP.String()) {
				continue
			}

			if ifaEnt.IfaNet.Contains(addr) {
				return 0, ifaEnt.IfaAddr
			}

			if firstIP == nil {
				firstIP = &ifaEnt.IfaAddr
			}
		}
	}

	if findAny == false {
		return L3AddrErr, net.IPv4(0, 0, 0, 0)
	}

	if firstIP != nil {
		return 0, *firstIP
	}

	return L3AddrErr, net.IPv4(0, 0, 0, 0)
}

// Ifa2String - Format an ifa to a string
func Ifa2String(ifa *Ifa, it IterIntf) {
	var str string
	for _, ifaEnt := range ifa.Ifas {
		var flagStr string
		if ifaEnt.Secondary {
			flagStr = "Secondary"
		} else {
			flagStr = "Primary"
		}
		plen, _ := ifaEnt.IfaNet.Mask.Size()
		str = fmt.Sprintf("%s/%d - %s", ifaEnt.IfaAddr.String(), plen, flagStr)
	}

	it.NodeWalker(str)
}

// Ifas2String - Format all ifas to string
func (l3 *L3H) Ifas2String(it IterIntf) error {
	for _, ifa := range l3.IfaMap {
		Ifa2String(ifa, it)
	}
	return nil
}

// IfaMkString - Given an ifa return its string representation
func IfaMkString(ifa *Ifa) string {
	var str string
	for _, ifaEnt := range ifa.Ifas {
		var flagStr string
		if ifaEnt.Secondary {
			flagStr = "Secondary"
		} else {
			flagStr = "Primary"
		}
		plen, _ := ifaEnt.IfaNet.Mask.Size()
		str = fmt.Sprintf("%s/%d - %s", ifaEnt.IfaAddr.String(), plen, flagStr)
	}

	return str
}

// IfObjMkString - given an ifa object, get all its member ifa's string rep
func (l3 *L3H) IfObjMkString(obj string) string {
	key := IfaKey{obj}
	ifa := l3.IfaMap[key]
	if ifa != nil {
		return IfaMkString(ifa)
	}

	return ""
}

// IfaGet - Get All of the IPv4Address in the Ifa
func (l3 *L3H) IfaGet() []cmn.IpAddrGet {
	var ret []cmn.IpAddrGet
	for ifName, ifa := range l3.IfaMap {
		var tmpIPa cmn.IpAddrGet
		tmpIPa.Dev = ifName.Obj
		tmpIPa.Sync = cmn.DpStatusT(ifa.Sync)
		for _, ip := range ifa.Ifas {
			o, _ := ip.IfaNet.Mask.Size()
			tmpIPa.IP = append(tmpIPa.IP, fmt.Sprintf("%s/%d", ip.IfaAddr.String(), o))
		}
		ret = append(ret, tmpIPa)
	}
	return ret
}

// DP - Sync state of L3 entities to data-path
func (ifa *Ifa) DP(work DpWorkT) int {
	port := ifa.Zone.Ports.PortFindByName(ifa.Key.Obj)

	if port == nil {
		return -1
	}

	// In case of remove request, we need to make sure
	// there are no other port IFAs with similar l2 address
	if work == DpRemove {
		for _, ent := range ifa.Zone.L3.IfaMap {
			if ifa.Zone.Ports.PortL2AddrMatch(ent.Key.Obj, port) == true {
				return 0
			}
		}
	}

	rmWq := new(RouterMacDpWorkQ)
	rmWq.Work = work
	rmWq.Status = &ifa.Sync

	for i := 0; i < 6; i++ {
		rmWq.L2Addr[i] = uint8(port.HInfo.MacAddr[i])
	}

	rmWq.PortNum = port.PortNo

	mh.dp.ToDpCh <- rmWq

	if port.SInfo.PortType&cmn.PortVxlanBr == cmn.PortVxlanBr {
		rmWq := new(RouterMacDpWorkQ)
		rmWq.Work = work
		rmWq.Status = &ifa.Sync

		if port.SInfo.PortReal == nil {
			return 0
		}

		up := port.SInfo.PortReal

		for i := 0; i < 6; i++ {
			rmWq.L2Addr[i] = uint8(up.HInfo.MacAddr[i])
		}

		rmWq.PortNum = up.PortNo
		rmWq.TunID = port.HInfo.TunID
		rmWq.TunType = DpTunVxlan
		rmWq.BD = port.L2.Vid

		mh.dp.ToDpCh <- rmWq

	}

	return 0
}
