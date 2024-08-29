package parser

import (
	"fmt"
	"strconv"
)

func ExtractIBCTransferFromEventsFromJson(idx int, jsonData []byte) ([]IBCFromCosmWasm, error) {
	events, error := ParseEvents(jsonData)
	if error != nil {
		return nil, fmt.Errorf("failed to parse events: %w", error)
	}

	return ExtractIBCTransferFromEvents(idx, events)
}

func ExtractIBCTransferFromEvents(idx int, events []Event) ([]IBCFromCosmWasm, error) {

	var ibcTransfers []IBCFromCosmWasm

	containsIBCTransfer := false
	for _, event := range events {
		if event.Type == "ibc_transfer" {
			for _, attr := range event.Attributes {
				if attr.Key == "msg_index" && attr.Value == strconv.Itoa(idx) {
					containsIBCTransfer = true
				}
			}
		}
	}

	if !containsIBCTransfer {
		return ibcTransfers, nil
	}

	sendPacketEvent := Event{}

	for _, event := range events {
		if event.Type == "send_packet" {
			for _, attr := range event.Attributes {
				if attr.Key == "msg_index" && attr.Value == strconv.Itoa(idx) {
					sendPacketEvent = event
					break
				}
			}
		}
	}

	if sendPacketEvent.Type == "" {
		return ibcTransfers, nil
	}

	packetAttributes := PacketAttributes{}
	for _, attr := range sendPacketEvent.Attributes {
		switch attr.Key {
		case "connection_id":
			packetAttributes.ConnectionID = attr.Value
		case "packet_channel_ordering":
			packetAttributes.PacketChannelOrdering = attr.Value
		case "packet_connection":
			packetAttributes.PacketConnection = attr.Value
		case "packet_data":
			packetAttributes.PacketData = attr.Value
		case "packet_data_hex":
			packetAttributes.PacketDataHex = attr.Value
		case "packet_dst_channel":
			packetAttributes.PacketDstChannel = attr.Value
		case "packet_dst_port":
			packetAttributes.PacketDstPort = attr.Value
		case "packet_sequence":
			packetAttributes.PacketSequence = attr.Value
		case "packet_src_channel":
			packetAttributes.PacketSrcChannel = attr.Value
		case "packet_src_port":
			packetAttributes.PacketSrcPort = attr.Value
		case "packet_timeout_height":
			packetAttributes.PacketTimeoutHeight = attr.Value
		case "packet_timeout_timestamp":
			packetAttributes.PacketTimeoutTimestamp = attr.Value
		case "msg_index":
			packetAttributes.MsgIndex = attr.Value
		}
	}

	if packetAttributes.ConnectionID == "" {
		return ibcTransfers, nil
	}

	packetData, err := packetAttributes.toPacketData()
	if err != nil {
		return ibcTransfers, fmt.Errorf("failed to map packet data to PacketData: %w", err)
	}

	return []IBCFromCosmWasm{
		{
			Amount:             packetData.Amount,
			Denom:              packetData.Denom,
			Receiver:           packetData.Receiver,
			Sender:             packetData.Sender,
			SourceChannel:      packetAttributes.PacketSrcChannel,
			SourcePort:         packetAttributes.PacketSrcPort,
			DestinationChannel: packetAttributes.PacketDstChannel,
			DestinationPort:    packetAttributes.PacketDstPort,
		},
	}, nil

}
