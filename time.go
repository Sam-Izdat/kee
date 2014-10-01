// Copyright 2014 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kee

import (
	"encoding/binary"
	"sync"
	"time"
)

// A Time represents a time as the number of 100's of nanoseconds since 15 Oct
// 1582.
type Time int64

const (
	lillian    = 2299160          // Julian day of 15 Oct 1582
	unix       = 2440587          // Julian day of 1 Jan 1970
	epoch      = unix - lillian   // Days between epochs
	g1582      = epoch * 86400    // seconds between epochs
	g1582ns100 = g1582 * 10000000 // 100s of a nanoseconds between epochs
)

var (
	mu        sync.Mutex
	lasttime  uint64 // last time we returned
	clockSeq uint16 // clock sequence for this run

	timeNow = time.Now // for testing
)

// UnixTime converts t the number of seconds and nanoseconds using the Unix
// epoch of 1 Jan 1970.
func (t Time) UnixTime() (sec, nsec int64) {
	sec = int64(t - g1582ns100)
	nsec = (sec % 10000000) * 100
	sec /= 10000000
	return sec, nsec
}

// GetTime returns the current Time (100s of nanoseconds since 15 Oct 1582) and
// adjusts the clock sequence as needed.  An error is returned if the current
// time cannot be determined.
func GetTime() (Time, error) {
	defer mu.Unlock()
	mu.Lock()
	return getTime()
}

func getTime() (Time, error) {
	t := timeNow()

	// If we don't have a clock sequence already, set one.
	if clockSeq == 0 {
		setClockSequence(-1)
	}
	now := uint64(t.UnixNano()/100) + g1582ns100

	// If time has gone backwards with this clock sequence then we
	// increment the clock sequence
	if now <= lasttime {
		clockSeq = ((clockSeq + 1) & 0x3fff) | 0x8000
	}
	lasttime = now
	return Time(now), nil
}

// ClockSequence returns the current clock sequence, generating one if not
// already set.  The clock sequence is only used for Version 1 UUIDs.
//
// The UUID package does not use global static storage for the clock sequence or
// the last time a UUID was generated.  Unless SetClockSequence a new random
// clock sequence is generated the first time a clock sequence is requested by
// ClockSequence, GetTime, or NewUUID.  (section 4.2.1.1) sequence is generated
// for
func ClockSequence() int {
	defer mu.Unlock()
	mu.Lock()
	return clockSequence()
}

func clockSequence() int {
	if clockSeq == 0 {
		setClockSequence(-1)
	}
	return int(clockSeq & 0x3fff)
}

// SetClockSequence sets the clock sequence to the lower 14 bits of seq.  Setting to
// -1 causes a new sequence to be generated.
func SetClockSequence(seq int) {
	defer mu.Unlock()
	mu.Lock()
	setClockSequence(seq)
}

func setClockSequence(seq int) {
	if seq == -1 {
		var b [2]byte
		randomBits(b[:]) // clock sequence
		seq = int(b[0])<<8 | int(b[1])
	}
	oldSeq := clockSeq
	clockSeq = uint16(seq&0x3fff) | 0x8000 // Set our variant
	if oldSeq != clockSeq {
		lasttime = 0
	}
}

// Time returns the time in 100s of nanoseconds since 15 Oct 1582 encoded in
// UUID.  It returns false if UUID is not valid.  The time is only well defined
// for version 1 and 2 UUIDs.
func (id KUUID) Time() (Time, bool) {
	bytes := id.slc
	if len(bytes) != 16 {
		return 0, false
	}
	time := int64(binary.BigEndian.Uint32(bytes[0:4]))
	time |= int64(binary.BigEndian.Uint16(bytes[4:6])) << 32
	time |= int64(binary.BigEndian.Uint16(bytes[6:8])&0xfff) << 48
	return Time(time), true
}

// ClockSequence returns the clock sequence encoded in UUID.  It returns false
// if UUID is not valid.  The clock sequence is only well defined for version 1
// and 2 UUIDs.
func (id KUUID) ClockSequence() (int, bool) {	
	bytes := id.slc
	if len(bytes) != 16 {
		return 0, false
	}
	return int(binary.BigEndian.Uint16(bytes[8:10])) & 0x3fff, true
}
