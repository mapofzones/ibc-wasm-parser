package parser

import (
	"encoding/base64"
	"fmt"
	"strconv"
)

func ExtractIBCTransferFromEventsFromJson(idx int, jsonData []byte, decodeKeys bool, decodeValues bool) ([]IBCFromCosmWasm, error) {
	events, error := ParseEvents(jsonData)
	if error != nil {
		return nil, fmt.Errorf("failed to parse events: %w", error)
	}

	return ExtractIBCTransferFromEvents(idx, events, decodeKeys, decodeValues)
}

func ExtractIBCTransferFromEvents(idx int, events []Event, decodeKeys bool, decodeValues bool) ([]IBCFromCosmWasm, error) {

	var ibcTransfers []IBCFromCosmWasm

	var sendPacketEvents []Event = make([]Event, 0)

	for _, event := range events {
		if event.Type == "send_packet" {
			if idx < 0 {
				sendPacketEvents = append(sendPacketEvents, event)
				continue
			} else {
				for _, attr := range event.Attributes {
					attrKey := attr.Key

					if decodeKeys {
						attrKey = decodeBase64String(attr.Key)
					}

					if attrKey == "msg_index" {
						attrValue := attr.Value
						if decodeValues {
							attrValue = decodeBase64String(attr.Value)
						}
						if attrValue == strconv.Itoa(idx) {
							sendPacketEvents = append(sendPacketEvents, event)
						}

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
			attrKey := attr.Key
			attrValue := attr.Value

			if decodeKeys {
				attrKey = decodeBase64String(attr.Key)
			}
			if decodeValues {
				attrValue = decodeBase64String(attr.Value)
			}

			switch attrKey {
			case "connection_id":
				packetAttributes.ConnectionID = attrValue
			case "packet_channel_ordering":
				packetAttributes.PacketChannelOrdering = attrValue
			case "packet_connection":
				packetAttributes.PacketConnection = attrValue
			case "packet_data":
				packetAttributes.PacketData = attrValue
			case "packet_data_hex":
				packetAttributes.PacketDataHex = attrValue
			case "packet_dst_channel":
				packetAttributes.PacketDstChannel = attrValue
			case "packet_dst_port":
				packetAttributes.PacketDstPort = attrValue
			case "packet_sequence":
				packetAttributes.PacketSequence = attrValue
			case "packet_src_channel":
				packetAttributes.PacketSrcChannel = attrValue
			case "packet_src_port":
				packetAttributes.PacketSrcPort = attrValue
			case "packet_timeout_height":
				packetAttributes.PacketTimeoutHeight = attrValue
			case "packet_timeout_timestamp":
				packetAttributes.PacketTimeoutTimestamp = attrValue
			case "msg_index":
				packetAttributes.MsgIndex = attrValue
			}
		}

		if packetAttributes.PacketSrcChannel == "" || packetAttributes.PacketDstChannel == "" {
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

// if it is base64 decode it and return it
// if not just return the plain string
func decodeBase64String(str string) string {
	decoded, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return str
	}
	return string(decoded)
}
