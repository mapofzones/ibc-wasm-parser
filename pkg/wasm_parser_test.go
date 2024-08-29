package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	cometbftabci "github.com/cometbft/cometbft/abci/types"
	tendermintabci "github.com/tendermint/tendermint/abci/types"
)

func loadEventsFromJson() ([]cometbftabci.Event, []tendermintabci.Event) {
	data, err := os.ReadFile("events.json")
	if err != nil {
		fmt.Println("Error reading json file:", err)
		return nil, nil
	}

	var cometEvents []cometbftabci.Event
	err = json.Unmarshal(data, &cometEvents)
	if err != nil {
		fmt.Println("Error unmarshalling json data:", err)
		return nil, nil
	}

	var tendermintEvents []tendermintabci.Event
	err = json.Unmarshal(data, &tendermintEvents)
	if err != nil {
		fmt.Println("Error unmarshalling json data for tendermint events:", err)
		return nil, nil
	}

	return cometEvents, tendermintEvents

}

func TestExtractIBCTransferFromEvents(t *testing.T) {
	cometEvents, tendermintEvents := loadEventsFromJson()

	doExtractIBCTransferFromEvents(t, cometEvents)

	// convert tendermint cometEvents to comet events
	cometEventsConverted := ConvertEventsToCmt(tendermintEvents)
	doExtractIBCTransferFromEvents(t, cometEventsConverted)

}

func doExtractIBCTransferFromEvents(t *testing.T, events []cometbftabci.Event) {
	ibcTransfers, err := ExtractIBCTransferFromEvents(0, events)

	if err != nil {
		t.Errorf("Error extracting ibc transfers: %v", err)
	}

	if len(ibcTransfers) != 1 {
		t.Errorf("Expected 1 ibc transfer, got %d", len(ibcTransfers))
	}

	ibcTransfer := ibcTransfers[0]
	if ibcTransfer.Receiver != "stars1ewwylz4klmn0jl95fjxjl956728mgszrch0saf" {
		t.Errorf("Expected receiver to be stars1ewwylz4klmn0jl95fjxjl956728mgszrch0saf, got %s", ibcTransfer.Receiver)
	}

	// this returns the contract that performed the transfer, but then should be
	// replaced later by the actual sender
	if ibcTransfer.Sender != "osmo1z24llw8lyafczgpza7qzdpmx273c4zxwjgl6qpu8ly34g3g9jd2q6vnd3t" {
		t.Errorf("Expected sender to be osmo1z24llw8lyafczgpza7qzdpmx273c4zxwjgl6qpu8ly34g3g9jd2q6vnd3t, got %s", ibcTransfer.Sender)
	}

	if ibcTransfer.Denom != "transfer/channel-75/ustars" {
		t.Errorf("Expected denom to be transfer/channel-75/ustars, got %s", ibcTransfer.Denom)
	}

	if ibcTransfer.Amount != "1569150" {
		t.Errorf("Expected amount to be 1569150, got %s", ibcTransfer.Amount)
	}

	if ibcTransfer.SourceChannel != "channel-75" {
		t.Errorf("Expected source channel to be channel-75, got %s", ibcTransfer.SourceChannel)
	}

	if ibcTransfer.SourcePort != "transfer" {
		t.Errorf("Expected source port to be transfer, got %s", ibcTransfer.SourcePort)
	}

	if ibcTransfer.DestinationChannel != "channel-0" {
		t.Errorf("Expected destination channel to be channel-0, got %s", ibcTransfer.DestinationChannel)
	}

	if ibcTransfer.DestinationPort != "transfer" {
		t.Errorf("Expected destination port to be transfer, got %s", ibcTransfer.DestinationPort)
	}

}

func TestConvertEventsToCmt(t *testing.T) {
	cometEvents, tendermintEvents := loadEventsFromJson()

	convertedEvents := ConvertEventsToCmt(tendermintEvents)

	if len(convertedEvents) != len(cometEvents) {
		t.Errorf("Expected %d converted events, got %d", len(cometEvents), len(convertedEvents))
	}

	for i, event := range convertedEvents {
		if event.Type != cometEvents[i].Type {
			t.Errorf("Expected event type to be %s, got %s", cometEvents[i].Type, event.Type)
		}

		if len(event.Attributes) != len(cometEvents[i].Attributes) {
			t.Errorf("Expected %d attributes, got %d", len(cometEvents[i].Attributes), len(event.Attributes))
		}

		for j, attr := range event.Attributes {
			if attr.Key != cometEvents[i].Attributes[j].Key {
				t.Errorf("Expected attribute key to be %s, got %s", cometEvents[i].Attributes[j].Key, attr.Key)
			}

			if attr.Value != cometEvents[i].Attributes[j].Value {
				t.Errorf("Expected attribute value to be %s, got %s", cometEvents[i].Attributes[j].Value, attr.Value)
			}
		}
	}
}
