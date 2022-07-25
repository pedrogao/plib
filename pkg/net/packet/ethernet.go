package packet

// This is a generated file! Please edit source .ksy file and use kaitai-struct-compiler to rebuild

import (
	"bytes"

	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
)

/**
 * Ethernet frame is a OSI data link layer (layer 2) protocol data unit
 * for Ethernet networks. In practice, many other networks and/or
 * in-file dumps adopted the same format for encapsulation purposes.
 * @see <a href="https://ieeexplore.ieee.org/document/7428776">Source</a>
 */

type EthernetFrame_EtherTypeEnum int

const (
	EthernetFrame_EtherTypeEnum__Ipv4          EthernetFrame_EtherTypeEnum = 2048
	EthernetFrame_EtherTypeEnum__X75Internet   EthernetFrame_EtherTypeEnum = 2049
	EthernetFrame_EtherTypeEnum__NbsInternet   EthernetFrame_EtherTypeEnum = 2050
	EthernetFrame_EtherTypeEnum__EcmaInternet  EthernetFrame_EtherTypeEnum = 2051
	EthernetFrame_EtherTypeEnum__Chaosnet      EthernetFrame_EtherTypeEnum = 2052
	EthernetFrame_EtherTypeEnum__X25Level3     EthernetFrame_EtherTypeEnum = 2053
	EthernetFrame_EtherTypeEnum__Arp           EthernetFrame_EtherTypeEnum = 2054
	EthernetFrame_EtherTypeEnum__Ieee8021qTpid EthernetFrame_EtherTypeEnum = 33024
	EthernetFrame_EtherTypeEnum__Ipv6          EthernetFrame_EtherTypeEnum = 34525
)

type EthernetFrame struct {
	DstMac       []byte
	SrcMac       []byte
	EtherType1   EthernetFrame_EtherTypeEnum
	Tci          *EthernetFrame_TagControlInfo
	EtherType2   EthernetFrame_EtherTypeEnum
	Body         any
	_io          *kaitai.Stream
	_root        *EthernetFrame
	_parent      any
	_raw_Body    []byte
	_f_etherType bool
	etherType    EthernetFrame_EtherTypeEnum
}

func NewEthernetFrame() *EthernetFrame {
	return &EthernetFrame{}
}

func (this *EthernetFrame) Read(io *kaitai.Stream, parent any, root *EthernetFrame) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp1, err := this._io.ReadBytes(int(6))
	if err != nil {
		return err
	}
	tmp1 = tmp1
	this.DstMac = tmp1
	tmp2, err := this._io.ReadBytes(int(6))
	if err != nil {
		return err
	}
	tmp2 = tmp2
	this.SrcMac = tmp2
	tmp3, err := this._io.ReadU2be()
	if err != nil {
		return err
	}
	this.EtherType1 = EthernetFrame_EtherTypeEnum(tmp3)
	if this.EtherType1 == EthernetFrame_EtherTypeEnum__Ieee8021qTpid {
		tmp4 := NewEthernetFrame_TagControlInfo()
		err = tmp4.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Tci = tmp4
	}
	if this.EtherType1 == EthernetFrame_EtherTypeEnum__Ieee8021qTpid {
		tmp5, err := this._io.ReadU2be()
		if err != nil {
			return err
		}
		this.EtherType2 = EthernetFrame_EtherTypeEnum(tmp5)
	}
	tmp6, err := this.EtherType()
	if err != nil {
		return err
	}
	switch tmp6 {
	case EthernetFrame_EtherTypeEnum__Ipv4:
		tmp7, err := this._io.ReadBytesFull()
		if err != nil {
			return err
		}
		tmp7 = tmp7
		this._raw_Body = tmp7
		_io__raw_Body := kaitai.NewStream(bytes.NewReader(this._raw_Body))
		tmp8 := NewIpv4Packet()
		err = tmp8.Read(_io__raw_Body, this, nil)
		if err != nil {
			return err
		}
		this.Body = tmp8
	case EthernetFrame_EtherTypeEnum__Ipv6:
		tmp9, err := this._io.ReadBytesFull()
		if err != nil {
			return err
		}
		tmp9 = tmp9
		this._raw_Body = tmp9
		_io__raw_Body := kaitai.NewStream(bytes.NewReader(this._raw_Body))
		tmp10 := NewIpv6Packet()
		err = tmp10.Read(_io__raw_Body, this, nil)
		if err != nil {
			return err
		}
		this.Body = tmp10
	default:
		tmp11, err := this._io.ReadBytesFull()
		if err != nil {
			return err
		}
		tmp11 = tmp11
		this._raw_Body = tmp11
	}
	return err
}

/**
 * Ether type can be specied in several places in the frame. If
 * first location bears special marker (0x8100), then it is not the
 * real ether frame yet, an additional payload (`tci`) is expected
 * and real ether type is upcoming next.
 */
func (this *EthernetFrame) EtherType() (v EthernetFrame_EtherTypeEnum, err error) {
	if this._f_etherType {
		return this.etherType, nil
	}
	var tmp12 EthernetFrame_EtherTypeEnum
	if this.EtherType1 == EthernetFrame_EtherTypeEnum__Ieee8021qTpid {
		tmp12 = this.EtherType2
	} else {
		tmp12 = this.EtherType1
	}
	this.etherType = EthernetFrame_EtherTypeEnum(tmp12)
	this._f_etherType = true
	return this.etherType, nil
}

/**
 * Destination MAC address
 */

/**
 * Source MAC address
 */

/**
 * Either ether type or TPID if it is a IEEE 802.1Q frame
 */

/**
 * Tag Control Information (TCI) is an extension of IEEE 802.1Q to
 * support VLANs on normal IEEE 802.3 Ethernet network.
 */
type EthernetFrame_TagControlInfo struct {
	Priority     uint64
	DropEligible bool
	VlanId       uint64
	_io          *kaitai.Stream
	_root        *EthernetFrame
	_parent      *EthernetFrame
}

func NewEthernetFrame_TagControlInfo() *EthernetFrame_TagControlInfo {
	return &EthernetFrame_TagControlInfo{}
}

func (this *EthernetFrame_TagControlInfo) Read(io *kaitai.Stream, parent *EthernetFrame, root *EthernetFrame) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp13, err := this._io.ReadBitsIntBe(3)
	if err != nil {
		return err
	}
	this.Priority = tmp13
	tmp14, err := this._io.ReadBitsIntBe(1)
	if err != nil {
		return err
	}
	this.DropEligible = tmp14 != 0
	tmp15, err := this._io.ReadBitsIntBe(12)
	if err != nil {
		return err
	}
	this.VlanId = tmp15
	return err
}

/**
 * Priority Code Point (PCP) is used to specify priority for
 * different kinds of traffic.
 */

/**
 * Drop Eligible Indicator (DEI) specifies if frame is eligible
 * to dropping while congestion is detected for certain classes
 * of traffic.
 */

/**
 * VLAN Identifier (VID) specifies which VLAN this frame
 * belongs to.
 */
