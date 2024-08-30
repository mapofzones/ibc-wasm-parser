package parser

import (
	"encoding/json"
	"fmt"
)

type IBCFromCosmWasm struct {
	Receiver           string
	Sender             string
	Denom              string
	Amount             string
	SourceChannel      string
	SourcePort         string
	DestinationChannel string
	DestinationPort    string
}

type PacketData struct {
	Denom    string `json:"denom"`
	Amount   string `json:"amount"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
}

type PacketAttributes struct {
	ConnectionID           string `json:"connection_id"`
	PacketChannelOrdering  string `json:"packet_channel_ordering"`
	PacketConnection       string `json:"packet_connection"`
	PacketData             string `json:"packet_data"`
	PacketDataHex          string `json:"packet_data_hex"`
	PacketDstChannel       string `json:"packet_dst_channel"`
	PacketDstPort          string `json:"packet_dst_port"`
	PacketSequence         string `json:"packet_sequence"`
	PacketSrcChannel       string `json:"packet_src_channel"`
	PacketSrcPort          string `json:"packet_src_port"`
	PacketTimeoutHeight    string `json:"packet_timeout_height"`
	PacketTimeoutTimestamp string `json:"packet_timeout_timestamp"`
	MsgIndex               string `json:"msg_index"`
}

func (pa PacketAttributes) toPacketData() (PacketData, error) {
	packetData := PacketData{}
	err := json.Unmarshal([]byte(pa.PacketData), &packetData)
	if err != nil {
		return PacketData{}, fmt.Errorf("failed to unmarshal packet data: %w", err)
	}

	// if amount or denom is empty, return an error
	if packetData.Amount == "" || packetData.Denom == "" {
		return PacketData{}, fmt.Errorf("amount or denom is empty")
	}

	// if sender or receiver is empty, return an error
	if packetData.Sender == "" || packetData.Receiver == "" {
		return PacketData{}, fmt.Errorf("sender or receiver is empty")
	}

	return packetData, nil
}

type EventAttribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Index bool   `json:"index"`
}

// Event represents an event with a type and a list of attributes
type Event struct {
	Type       string           `json:"type"`
	Attributes []EventAttribute `json:"attributes"`
}

func ParseEvents(jsonData []byte) ([]Event, error) {
	var events []Event
	err := json.Unmarshal(jsonData, &events)
	if err != nil {
		return nil, err
	}
	return events, nil
}
