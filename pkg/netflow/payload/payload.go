// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package payload

// FlowPayload contains network devices flows
type FlowPayload struct {
	// Keys
	SrcAddr        string `json:"src_addr"`
	SrcPort        uint32 `json:"src_port"`
	DstAddr        string `json:"dst_addr"`
	DstPort        uint32 `json:"dst_port"`
	Proto          uint32 `json:"proto"`
	Tos            uint32 `json:"Tos"`
	InputInterface uint32 `json:"input_interface"`

	// Non-Keys
	ReceivedTimestamp uint64 `json:"received_timestamp"`
	StartTimestamp    uint64 `json:"start_timestamp"`
	EndTimestamp      uint64 `json:"end_timestamp"`
	SamplerAddr       string `json:"sampler_addr"`
	FlowType          string `json:"flow_type"`
	OutputInterface   uint32 `json:"output_interface"`
	Direction         uint32 `json:"direction"`
	Bytes             uint64 `json:"bytes"`
	Packets           uint64 `json:"packets"`
}
