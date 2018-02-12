package asteroid_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/angadn/asteroid"
)

func TestSIPHeaderWatch(t *testing.T) {
	passed := false
	w := asteroid.NewSIPHeaderWatch([]string{"X-Twilio-CallSid", "X-Twilio-RecordingSid"}, func(
		callID string, headers map[string]string,
	) {
		if callID == "3ca63b4114c4730415f57b1b217d040e@35.197.101.20:5060" &&
			headers["X-Twilio-CallSid"] == "CA23a4978e378035d1389c0838183e47d4" &&
			headers["X-Twilio-RecordingSid"] == "RE3f460611a574b29e327b2a8acbff28d0" {
			passed = true
		} else {
			passed = false // If any message comes, it should contain all our headers!
		}
	})

	w.SetReader(bytes.NewReader([]byte(asteriskLogs)))
	w.Start()
	time.Sleep(100 * time.Millisecond)
	w.Stop()
	if !passed {
		t.Fail()
	}
}

func TestSIPHeaderWatchForReason(t *testing.T) {
	passed := false
	w := asteroid.NewSIPHeaderWatch([]string{"Reason"}, func(
		callID string, headers map[string]string,
	) {
		if callID == "3ca63b4114c4730415f57b1b217d040e@35.197.101.20:5060" {
			passed = true
		} else {
			passed = false // If any message comes, it should contain all our headers!
		}
	})

	w.SetReader(bytes.NewReader([]byte(asteriskLogs)))
	w.Start()
	time.Sleep(100 * time.Millisecond)
	w.Stop()
	if !passed {
		t.Fail()
	}
}

func TestSIPDestructionWatch(t *testing.T) {
	passed := false
	w := asteroid.NewSIPDestructionWatch(func(callID string) {
		if callID == "3ca63b4114c4730415f57b1b217d040e@35.197.101.20:5060" {
			passed = true
		} else {
			passed = false
		}
	})

	w.SetReader(bytes.NewReader([]byte(asteriskLogs)))
	w.Start()
	time.Sleep(100 * time.Millisecond)
	w.Stop()
	if !passed {
		t.Fail()
	}
}

/************************************************
*					WARNING						*
* These logs are imperfect. They don't contain  *
* terminal-markup. 								*
* eg. "[OKReally destroying SIP dialog..."		*
* If a test-case succeeds but deployment fails, *
* this is probably why!							*
*************************************************/

const asteriskLogs = `Asterisk 11.13.1~dfsg-2+deb8u4, Copyright (C) 1999 - 2013 Digium, Inc. and others.
Created by Mark Spencer <markster@digium.com>
Asterisk comes with ABSOLUTELY NO WARRANTY; type 'core show warranty' for details.
This is free software, with components licensed under the GNU General Public
License version 2 and other licenses; you are welcome to redistribute it under
certain conditions. Type 'core show license' for details.
=========================================================================
Connected to Asterisk 11.13.1~dfsg-2+deb8u4 currently running on gke-asterisk-test-default-pool-3afa9d9d-3lnn (pid = 10)
<SIP/twilio0-00000002>AGI Rx << EXEC BACKGROUND /tmp/162657525
<SIP/twilio0-00000002>AGI Tx >> 200 result=0
gke-asterisk-test-default-pool-3afa9d9d-3lnn*CLI>
gke-asterisk-test-default-pool-3afa9d9d-3lnn*CLI>

<--- SIP read from UDP:54.172.60.1:5060 --->
BYE sip:<<REMOVED>> SIP/2.0
CSeq: 1 BYE
From: <sip:<<REMOVED>>>;tag=09977975_6772d868_99437e75-7a54-4d55-9dfc-d895b4dba97c
To: <sip:<<REMOVED>>>;tag=as289b691a
Call-ID: 3ca63b4114c4730415f57b1b217d040e@35.197.101.20:5060
Max-Forwards: 68
Via: SIP/2.0/UDP <<REMOVED>>
Via: SIP/2.0/UDP <<REMOVED>>
Reason: Q.850;cause=16;text="Normal call clearing"
User-Agent: Twilio Gateway
X-Twilio-CallSid: CA23a4978e378035d1389c0838183e47d4
X-Twilio-RecordingSid: RE3f460611a574b29e327b2a8acbff28d0
X-Twilio-RecordingDuration: 25
Content-Length: 0

<------------->
--- (14 headers 0 lines) ---
Sending to 54.172.60.1:5060 (no NAT)
Scheduling destruction of SIP dialog '3ca63b4114c4730415f57b1b217d040e@35.197.101.20:5060' in 32000 ms (Method: BYE)

<--- Transmitting (no NAT) to 54.172.60.1:5060 --->
SIP/2.0 200 OK
Via: SIP/2.0/UDP 54.172.60.1:5060;branch=z9hG4bK9d39.72521cd1.0;received=54.172.60.1
Via: SIP/2.0/UDP 172.18.5.109:5060;rport=5060;received=34.228.39.101;branch=z9hG4bK99437e75-7a54-4d55-9dfc-d895b4dba97c_6772d868_342-4206590906768630031
From: <sip:<<REMOVED>>>;tag=09977975_6772d868_99437e75-7a54-4d55-9dfc-d895b4dba97c
To: <sip:<<REMOVED>>>;tag=as289b691a
Call-ID: 3ca63b4114c4730415f57b1b217d040e@35.197.101.20:5060
CSeq: 1 BYE
Server: Asterisk PBX 11.13.1~dfsg-2+deb8u4
Allow: INVITE, ACK, CANCEL, OPTIONS, BYE, REFER, SUBSCRIBE, NOTIFY, INFO, PUBLISH, MESSAGE
Supported: replaces, timer
Content-Length: 0


<------------>
[Jan 31 07:49:38] NOTICE[28951]: pbx_spool.c:402 attempt_thread: Call completed to SIP/twilio0/<<REMOVED>>
[0KReally destroying SIP dialog '3ca63b4114c4730415f57b1b217d040e@35.197.101.20:5060' Method: BYE
Really destroying SIP dialog '3ca63b4114c4730415f57b1b217d040e@35.197.101.20:5060' Method: BYE`
