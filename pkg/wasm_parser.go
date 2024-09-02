package parser

import (
	"encoding/base64"
	"fmt"
	"strconv"
)

func ExtractIBCTransferFromEventsFromJson(idx int, jsonData []byte, decodeKeys, decodeValues bool) ([]IBCFromCosmWasm, error) {
	events, err := ParseEvents(jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse events: %w", err)
	}

	return ExtractIBCTransferFromEvents(idx, events, decodeKeys, decodeValues)
}

func ExtractIBCTransferFromEvents(idx int, events []Event, decodeKeys, decodeValues bool) ([]IBCFromCosmWasm, error) {
	var ibcTransfers []IBCFromCosmWasm

	sendPacketEvents := filterSendPacketEvents(idx, events, decodeKeys, decodeValues)

	for _, event := range sendPacketEvents {
		packetAttributes := extractPacketAttributes(event.Attributes, decodeKeys, decodeValues)

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

func filterSendPacketEvents(idx int, events []Event, decodeKeys, decodeValues bool) []Event {
	var sendPacketEvents []Event

	for _, event := range events {
		if event.Type != "send_packet" {
			continue
		}

		if idx < 0 {
			sendPacketEvents = append(sendPacketEvents, event)
			continue
		}

		for _, attr := range event.Attributes {
			attrKey := decodeIfNeeded(attr.Key, decodeKeys)

			if attrKey == "msg_index" {
				attrValue := decodeIfNeeded(attr.Value, decodeValues)
				if attrValue == strconv.Itoa(idx) {
					sendPacketEvents = append(sendPacketEvents, event)
					break
				}
			}
		}
	}

	return sendPacketEvents
}

func extractPacketAttributes(attributes []EventAttribute, decodeKeys, decodeValues bool) PacketAttributes {
	packetAttributes := PacketAttributes{}

	for _, attr := range attributes {
		attrKey := decodeIfNeeded(attr.Key, decodeKeys)
		attrValue := decodeIfNeeded(attr.Value, decodeValues)

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

	return packetAttributes
}

func decodeIfNeeded(str string, decode bool) string {
	if decode {
		return decodeBase64String(str)
	}
	return str
}

// decodeBase64String attempts to decode a string from base64; returns the original string if decoding fails.
func decodeBase64String(str string) string {
	decoded, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return str
	}
	return string(decoded)
}
