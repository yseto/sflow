package sflow

import (
	"bytes"
	"fmt"
	"github.com/fstelzer/sflow/flow_records"
	"net"
	"reflect"
	"testing"
)

func TestEncodeDecodeExtendedGatewayFlowRecord(t *testing.T) {
	rec := flow_records.ExtendedGatewayFlow{
		NextHopType:          2,
		NextHop:              net.ParseIP("2001:0db8:ac10:fe01::"), //IPv4 fails with the DeepEqual
		As:                   1234,
		SrcAs:                4321,
		SrcPeerAs:            5678,
		DstAsPathSegmentsLen: 1,
		DstAsPathSegments: []flow_records.ExtendedGatewayFlowASPathSegment{{
			SegType: 1,
			SegLen:  3,
			Seg:     []uint32{1234, 4321, 65535},
		}},
		CommunitiesLen: 3,
		Communities:    []uint32{1, 18, 42011},
		LocalPref:      255,
	}

	b := &bytes.Buffer{}

	err := rec.Encode(b)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Received binary results: %+#v\n", b)

	// Skip the header section. It's 8 bytes.
	var headerBytes [8]byte

	_, err = b.Read(headerBytes[:])
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := flow_records.DecodeExtendedGatewayFlow(b)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Got     : %+#v\n", decoded)
	fmt.Printf("Expected: %+#v\n", rec)
	if !reflect.DeepEqual(rec, decoded) {
		t.Errorf("expected\n%+#v\n, got\n%+#v", rec, decoded)
	}
}

func TestEncodeDecodeRawPacketFlowRecord(t *testing.T) {
	rec := flow_records.RawPacketFlow{
		Protocol:    1,
		FrameLength: 318,
		Stripped:    4,
		HeaderSize:  128,
		Header: []byte{0x00, 0xD0, 0x01, 0xFF, 0x58,
			0x00, 0x00, 0x16, 0x3C, 0xC2, 0xA9, 0xAB,
			0x08, 0x00, 0x45, 0x00, 0x01, 0x2C, 0x00,
			0x00, 0x40, 0x00, 0x40, 0x11, 0xD1, 0x58,
			0xC7, 0x3A, 0xA1, 0x96, 0xC5, 0xA1, 0x39,
			0xF6, 0xC8, 0xD5, 0x26, 0x00, 0x01, 0x18,
			0xA6, 0x17, 0x64, 0x31, 0x3A, 0x72, 0x64,
			0x32, 0x3A, 0x69, 0x64, 0x32, 0x30, 0x3A,
			0x6B, 0x96, 0x8B, 0xCA, 0x4A, 0xC0, 0xB5,
			0xCF, 0x10, 0x3A, 0xD6, 0xBF, 0x8D, 0xD7,
			0x34, 0x01, 0x46, 0x51, 0xB7, 0xFA, 0x35,
			0x3A, 0x6E, 0x6F, 0x64, 0x65, 0x73, 0x32,
			0x30, 0x38, 0x3A, 0x61, 0x4A, 0xB8, 0x64,
			0x54, 0xEE, 0x85, 0x5F, 0x13, 0x9A, 0x20,
			0x96, 0xE9, 0x83, 0xFF, 0xCF, 0xF4, 0xD0,
			0xC5, 0xA5, 0xDE, 0x67, 0x0A, 0x8F, 0xDB,
			0x1D, 0x61, 0x4A, 0x78, 0x12, 0x83, 0x31,
			0xA3, 0x77, 0x86, 0x68, 0x5E, 0x1C, 0x24,
			0xCE, 0x33, 0x19, 0xDE,
		},
	}

	b := &bytes.Buffer{}

	err := rec.Encode(b)
	if err != nil {
		t.Fatal(err)
	}

	// Skip the header section. It's 8 bytes.
	var headerBytes [8]byte

	_, err = b.Read(headerBytes[:])
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := flow_records.DecodeRawPacketFlow(b)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(rec, decoded) {
		t.Errorf("expected\n%+#v\n, got\n%+#v", rec, decoded)
	}
}
