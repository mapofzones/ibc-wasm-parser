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

	var sendPacketEvents []Event = make([]Event, 0)

	for _, event := range events {
		if event.Type == "send_packet" {
			if idx < 0 {
				sendPacketEvents = append(sendPacketEvents, event)
				continue
			} else {
				for _, attr := range event.Attributes {
					fmt.Println("INDEX=", idx)
					if attr.Key == "msg_index" && attr.Value == strconv.Itoa(idx) {
						sendPacketEvents = append(sendPacketEvents, event)
					}
				}
			}
		}
	}

	if len(sendPacketEvents) == 0 {
		return ibcTransfers, nil
	}

	for _, sendPacketEvent := range sendPacketEvents {
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
			continue
		}

		packetData, err := packetAttributes.toPacketData()
		if err != nil {
			continue
		}

		ibcTransfers = append(ibcTransfers, IBCFromCosmWasm{
			Amount:             packetData.Amount,
			Denom:              packetData.Denom,
			Receiver:           packetData.Receiver,
			Sender:             packetData.Sender,
			SourceChannel:      packetAttributes.PacketSrcChannel,
			SourcePort:         packetAttributes.PacketSrcPort,
			DestinationChannel: packetAttributes.PacketDstChannel,
			DestinationPort:    packetAttributes.PacketDstPort,
		})
	}

	return ibcTransfers, nil
}
