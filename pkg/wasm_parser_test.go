package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	abcitypes "github.com/cometbft/cometbft/abci/types"
)

func loadEventsFromJson() []abcitypes.Event {
	data, err := os.ReadFile("events.json")
	if err != nil {
		fmt.Println("Error reading json file:", err)
		return nil
	}

	var events []abcitypes.Event
	err = json.Unmarshal(data, &events)
	if err != nil {
		fmt.Println("Error unmarshalling json data:", err)
		return nil
	}

	return events

}

func TestExtractIBCTransferFromEvents(t *testing.T) {
	events := loadEventsFromJson()

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
