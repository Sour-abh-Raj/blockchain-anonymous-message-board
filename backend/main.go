package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/rs/cors"
)

type BlockChain struct {
    blocks []*Block
    mu     sync.Mutex
}

type Block struct {
    Hash     []byte `json:"hash"`
    Data     string `json:"data"`
    PrevHash []byte `json:"prevHash"`
    Time     string `json:"time"`
}

func (b *Block) DeriveHash() {
    info := []byte(b.Data + b.Time)
    hash := sha256.Sum256(info)
    b.Hash = hash[:]
}

func CreateBlock(data string, prevHash []byte) *Block {
    block := &Block{[]byte{}, data, prevHash, time.Now().Format(time.RFC3339)}
    block.DeriveHash()
    return block
}

func (chain *BlockChain) AddBlock(data string) {
    chain.mu.Lock()
    defer chain.mu.Unlock()

    prevBlock := chain.blocks[len(chain.blocks)-1]
    newBlock := CreateBlock(data, prevBlock.Hash)
    chain.blocks = append(chain.blocks, newBlock)
}

func Genesis() *Block {
    return CreateBlock("Genesis", []byte{})
}

func InitBlockChain() *BlockChain {
    return &BlockChain{blocks: []*Block{Genesis()}}
}

func SaveBlockchainToFile(blockchain *BlockChain, filename string) error {
    data, err := json.Marshal(blockchain.blocks)
    if err != nil {
        return err
    }
    return ioutil.WriteFile(filename, data, 0644)
}

func LoadBlockchainFromFile(filename string) (*BlockChain, error) {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    var blocks []*Block
    if err := json.Unmarshal(data, &blocks); err != nil {
        return nil, err
    }
    return &BlockChain{blocks: blocks}, nil
}

// Function to periodically save the blockchain
func PeriodicSave(blockchain *BlockChain, filename string, interval time.Duration) {
    for range time.Tick(interval) {
        if err := SaveBlockchainToFile(blockchain, filename); err != nil {
            fmt.Println("Error saving blockchain:", err)
        }
    }
}

func main() {
	var blockchain *BlockChain
    // Check if a saved blockchain file exists
    if _, err := os.Stat("blockchain.json"); err == nil {
        var err error
        blockchain, err = LoadBlockchainFromFile("blockchain.json")
        if err != nil {
            fmt.Println("Error loading blockchain:", err)
            return
        }
        fmt.Println("Loaded existing blockchain.")
    } else {
        // Initialize a new blockchain if no saved file exists
        blockchain = InitBlockChain()
        fmt.Println("Initialized new blockchain.")
    }
	
    // Start the periodic saving goroutine
    saveInterval := 10 * time.Minute
    go PeriodicSave(blockchain, "blockchain.json", saveInterval)
    // Create a new CORS handler with appropriate options
    corsHandler := cors.New(cors.Options{
        AllowedOrigins: []string{"*"}, // Adjust this as needed for your CORS policy
    })
    // API endpoint to get all messages
    http.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        blocks := blockchain.blocks
        response, _ := json.Marshal(blocks)

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write(response)

    })
    // Serve static files
    // http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("static"))))
    // http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// API endpoint to add a new message
    http.HandleFunc("/addBlock", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        var data struct {
            Data string `json:"data"`
        }

        if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        blockchain.AddBlock(data.Data)
        block := blockchain.blocks[len(blockchain.blocks)-1]
        response, _ := json.Marshal(block)

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        w.Write(response)
    })

    // Start the server
    fmt.Println("Listening on :8080")
    http.ListenAndServe(":8080", corsHandler.Handler(http.DefaultServeMux))
}
