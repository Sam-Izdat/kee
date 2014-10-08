// Copyright 2011 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kee

import "net"

var node struct {
	interfaces []net.Interface // cached list of interfaces
	ifname     string          // name of interface being used
	nodeID     []byte          // hardware for version 1 UUIDs
}

// NodeInterface returns the name of the interface from which the NodeID was
// derived.  The interface "user" is returned if the NodeID was set by
// SetNodeID.
func (_ UUIDCtrl) NodeInterface() string {
	return node.ifname
}

// SetNodeInterface selects the hardware address to be used for Version 1 UUIDs.
// If name is "" then the first usable interface found will be used or a random
// Node ID will be generated.  If a named interface cannot be found then false
// is returned.
//
// SetNodeInterface never fails when name is "".
func (c UUIDCtrl) SetNodeInterface(name string) bool {
	if node.interfaces == nil {
		var err error
		node.interfaces, err = net.Interfaces()
		if err != nil && name != "" {
			return false
		}
	}

	for _, ifs := range node.interfaces {
		if len(ifs.HardwareAddr) >= 6 && (name == "" || name == ifs.Name) {
			if c.setNodeID(ifs.HardwareAddr) {
				node.ifname = ifs.Name
				return true
			}
		}
	}

	// We found no interfaces with a valid hardware address.  If name
	// does not specify a specific interface generate a random Node ID
	// (section 4.1.6)
	if name == "" {
		if node.nodeID == nil {
			node.nodeID = make([]byte, 6)
		}
		randomBits(node.nodeID)
		return true
	}
	return false
}

// NodeID returns a slice of a copy of the current Node ID, setting the Node ID
// if not already set.
func (c UUIDCtrl) NodeID() []byte {
	if node.nodeID == nil {
		c.SetNodeInterface("")
	}
	nid := make([]byte, 6)
	copy(nid, node.nodeID)
	return nid
}

// SetNodeID sets the Node ID to be used for Version 1 UUIDs.  The first 6 bytes
// of id are used.  If id is less than 6 bytes then false is returned and the
// Node ID is not set.
func (c UUIDCtrl) SetNodeID(id []byte) bool {
	if c.setNodeID(id) {
		node.ifname = "user"
		return true
	}
	return false
}

func (_ UUIDCtrl) setNodeID(id []byte) bool {
	if len(id) < 6 {
		return false
	}
	if node.nodeID == nil {
		node.nodeID = make([]byte, 6)
	}
	copy(node.nodeID, id)
	return true
}

// NodeID returns the 6 byte node id encoded in UUID.  It returns nil if UUID is
// not valid.  The NodeID is only well defined for version 1 and 2 UUIDs.
func (id KUUID) NodeID() []byte {
	bytes := id.slc
	if len(bytes) != 16 {
		return nil
	}
	node := make([]byte, 6)
	copy(node, bytes[10:])
	return node
}