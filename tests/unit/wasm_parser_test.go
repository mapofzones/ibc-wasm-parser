package tests

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	cometbftabci "github.com/cometbft/cometbft/abci/types"
	parser "github.com/mapofzones/ibc-wasm-parser/pkg"
	tendermintabci "github.com/tendermint/tendermint/abci/types"
)

func loadEventsFromJson(fileName string) ([]cometbftabci.Event, []tendermintabci.Event) {
	filePath := fmt.Sprintf("../fixtures/%s", fileName)
	data, err := os.ReadFile(filePath)
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
	files := []string{"events_without_msg_index.json", "events_with_msg_index.json"}

	for idx, file := range files {
		cometEvents, tendermintEvents := loadEventsFromJson(file)
		cometEventsJson, _ := json.MarshalIndent(cometEvents, "", "  ")
		doExtractIBCTransferFromEvents(t, idx-1, cometEventsJson)
		tendermintEventsJson, _ := json.MarshalIndent(tendermintEvents, "", "  ")
		doExtractIBCTransferFromEvents(t, idx-1, tendermintEventsJson)
	}

}

func doExtractIBCTransferFromEvents(t *testing.T, idx int, jsonData []byte) {
	ibcTransfers, err := parser.ExtractIBCTransferFromEventsFromJson(idx, jsonData)

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
