/*
 * Copyright(C) 2015 Simon Schmidt
 * 
 * This Source Code Form is subject to the terms of the
 * Mozilla Public License, v. 2.0. If a copy of the MPL
 * was not distributed with this file, You can obtain one at
 * http://mozilla.org/MPL/2.0/.
 */

/*
 This Package models POSIX-ACLs including their representation as Xattrs.
 */
package posix_acl

import "bytes"
import "encoding/binary"
import "fmt"

const ACL_USER_OWNER  = 0x0001
const ACL_USER        = 0x0002
const ACL_GROUP_OWNER = 0x0004
const ACL_GROUP       = 0x0008
const ACL_MASK        = 0x0010
const ACL_OTHERS      = 0x0020

type AclSID uint64
func (a *AclSID) SetUid(uid uint32) {
	*a = AclSID(uid)|(ACL_USER<<32)
}
func (a *AclSID) SetGid(gid uint32) {
	*a = AclSID(gid)|(ACL_GROUP<<32)
}
// One of ACL_* (Except ACL_ACCESS/ACL_DEFAULTS)
func (a *AclSID) SetType(tp int) {
	*a = AclSID(tp)<<32
}
// One of ACL_* (Except ACL_ACCESS/ACL_DEFAULTS)
func (a AclSID) GetType() int {
	return int(a>>32)
}
func (a AclSID) GetID() uint32 {
	return uint32(a&0xffffffff)
}
func (a AclSID) String() string {
	switch(a>>32){
	case ACL_USER_OWNER:return "u::"
	case ACL_USER:return fmt.Sprintf("u:%v:",a.GetID())
	case ACL_GROUP_OWNER:return "g::"
	case ACL_GROUP:return fmt.Sprintf("g:%v:",a.GetID())
	case ACL_MASK:return "m::"
	case ACL_OTHERS:return "o::"
	}
	return "?:"
}

type AclElement struct{
	AclSID
	Perm uint16
}
func (a AclElement) String() string {
	str := ""
	if (a.Perm&4)!=0 { str+="r" }
	if (a.Perm&2)!=0 { str+="w" }
	if (a.Perm&1)!=0 { str+="x" }
	return fmt.Sprintf("%v%v",a.AclSID,str)
}

type Acl struct{
	Version uint32
	List []AclElement
}
func (a *Acl)Decode(xattr []byte) {
	var elem AclElement
	ae := new(aclElem)
	nr := bytes.NewReader(xattr)
	e := binary.Read(nr,binary.LittleEndian,&a.Version)
	if e!=nil { a.Version=0; return }
	if len(a.List)>0 {
		a.List=a.List[:0]
	}
	for binary.Read(nr,binary.LittleEndian,ae)==nil {
		elem.AclSID = (AclSID(ae.Tag)<<32) | AclSID(ae.Id)
		elem.Perm = ae.Perm
		a.List = append(a.List,elem)
	}
}
func (a *Acl)Encode() []byte {
	buf := new(bytes.Buffer)
	ae := new(aclElem)
	binary.Write(buf,binary.LittleEndian,&a.Version)
	for _,elem := range a.List {
		ae.Tag = uint16(elem.GetType())
		ae.Perm = elem.Perm
		ae.Id = elem.GetID()
		binary.Write(buf,binary.LittleEndian,ae)
	}
	return buf.Bytes()
}

type aclElem struct{
	Tag uint16
	Perm uint16
	Id uint32
}



