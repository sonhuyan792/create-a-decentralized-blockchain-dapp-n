// 5yk4_create_a_decent.go
// A decentralized blockchain dApp notifier built with Go

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Config holds the configuration for the dApp notifier
type Config struct {
	EthereumNodeURL string `json:"ethereum_node_url"`
	ContractAddress string `json:"contract_address"`
	 ABI            string `json:"abi"`
	NotifierURL     string `json:"notifier_url"`
}

// NewConfig returns a new Config instance
func NewConfig() *Config {
	return &Config{}
}

// LoadConfig loads the configuration from a JSON file
func LoadConfig(filename string) (*Config, error) {
 configFile, err := os.Open(filename)
 if err != nil {
  return nil, err
 }
 defer configFile.Close()
 var config Config
 err = json.NewDecoder(configFile).Decode(&config)
 return &config, err
}

// Notifier struct holds the configuration and a pointer to the Ethereum client
type Notifier struct {
	config *Config
	client *ethclient.Client
}

// NewNotifier returns a new Notifier instance
func NewNotifier(config *Config) (*Notifier, error) {
 client, err := ethclient.Dial(config.EthereumNodeURL)
 if err != nil {
  return nil, err
 }
 return &Notifier{config: config, client: client}, nil
}

// Start starts the notifier
func (n *Notifier) Start() error {
 log.Println("Notifier started")
 for {
  // Filter for new blocks
  filterer, err := n.client.NewBlockFilter()
  if err != nil {
   return err
  }
  defer filterer.Uninstall()
  for {
   select {
   case block := <-filterer.Chan():
    // Process new block
    n.processBlock(block)
   }
  }
 }
}

// processBlock process a new block
func (n *Notifier) processBlock(block *types.Block) {
 // Get the contract instance
 contractAddress := common.HexToAddress(n.config.ContractAddress)
 instance, err := myContract.NewMyContract(contractAddress, n.client)
 if err != nil {
  log.Println(err)
  return
 }
 // Call the contract to get the events
 events, err := instance.GetEvents(&bind.CallOpts{}, block.Number.Int64())
 if err != nil {
  log.Println(err)
  return
 }
 // Send notifications for each event
 for _, event := range events {
  n.sendNotification(event)
 }
}

// sendNotification sends a notification to the notifier URL
func (n *Notifier) sendNotification(event *MyContractMyEvent) {
 // Create a notification message
 message := fmt.Sprintf("Event triggered: %s, Args: %v", event.Raw.Name, event.Raw.Data)
 // Send the notification
 resp, err := http.Get(n.config.NotifierURL + "?message=" + url.QueryEscape(message))
 if err != nil {
  log.Println(err)
  return
 }
 defer resp.Body.Close()
 if resp.StatusCode != http.StatusOK {
  log.Println("Error sending notification:", resp.Status)
 }
}

func main() {
 config, err := LoadConfig("config.json")
 if err != nil {
  log.Fatal(err)
 }
 notifier, err := NewNotifier(config)
 if err != nil {
  log.Fatal(err)
 }
 err = notifier.Start()
 if err != nil {
  log.Fatal(err)
 }
}