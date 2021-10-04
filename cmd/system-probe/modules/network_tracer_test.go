// +build linux windows

package modules

import (
	"net/http/httptest"
	"testing"

	"github.com/DataDog/datadog-agent/pkg/network"
	"github.com/DataDog/datadog-agent/pkg/network/encoding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"inet.af/netaddr"
)

func TestDecode(t *testing.T) {
	rec := httptest.NewRecorder()

	in := &network.Connections{
		BufferedData: network.BufferedData{
			Conns: []network.ConnectionStats{
				{
					Source:               netaddr.MustParseIP("10.1.1.1"),
					Dest:                 netaddr.MustParseIP("10.2.2.2"),
					MonotonicSentBytes:   1,
					LastSentBytes:        2,
					MonotonicRecvBytes:   100,
					LastRecvBytes:        101,
					LastUpdateEpoch:      50,
					MonotonicRetransmits: 201,
					LastRetransmits:      201,
					Pid:                  6000,
					NetNS:                7,
					SPort:                1000,
					DPort:                9000,
					IPTranslation: &network.IPTranslation{
						ReplSrcIP:   netaddr.MustParseIP("20.1.1.1"),
						ReplDstIP:   netaddr.MustParseIP("20.1.1.1"),
						ReplSrcPort: 40,
						ReplDstPort: 70,
					},

					Type:      network.UDP,
					Direction: network.LOCAL,
				},
			},
		},
	}

	marshaller := encoding.GetMarshaler(encoding.ContentTypeJSON)
	expected, err := marshaller.Marshal(in)
	require.NoError(t, err)

	writeConnections(rec, marshaller, in)

	rec.Flush()
	out := rec.Body.Bytes()
	assert.Equal(t, expected, out)

}
