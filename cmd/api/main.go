package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/timkippdev/cryptochain/pkg/blockchain"
	"github.com/timkippdev/cryptochain/pkg/pubsub"
	"github.com/timkippdev/cryptochain/pkg/wallet"
)

const (
	rootNodePort = 3000
)

var (
	domain          string
	port            int
	redisAddress    string
	rootNodeAddress string
	serverAddress   string
)

func init() {
	domain = "localhost"

	if os.Getenv("REDIS_ADDRESS") != "" {
		redisAddress = os.Getenv("REDIS_ADDRESS")
	} else {
		redisAddress = fmt.Sprintf("%s:6379", domain)
	}

	port = 3000
	if os.Getenv("IS_PEER") != "" {
		rand.Seed(time.Now().Unix())
		port = port + rand.Intn(1000)
	}

	serverAddress = fmt.Sprintf("%s:%d", domain, port)
	rootNodeAddress = fmt.Sprintf("http://%s:%d", domain, rootNodePort)
}

func main() {
	bc := blockchain.NewBlockchain()
	walletInstance := wallet.NewWallet()
	tp := wallet.NewTransactionPool()
	ps := pubsub.NewPubSub(redisAddress, bc)
	defer ps.Close()

	router := mux.NewRouter()
	router.HandleFunc("/health-check", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	router.HandleFunc("/api/blocks", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(bc.GetChain())
	})

	router.HandleFunc("/api/mine", func(w http.ResponseWriter, r *http.Request) {
		b := bc.AddBlock("some-data")
		ps.BroadcastChain()

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(b)
	}).Methods("POST")

	router.HandleFunc("/api/transaction", func(w http.ResponseWriter, r *http.Request) {
		t, err := walletInstance.CreateTransaction(&ecdsa.PublicKey{}, 75)
		if err != nil {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}

		tp.AddTransaction(t)

		fmt.Printf("%+v\n", tp)
		fmt.Printf("%+v\n", t)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(201)
		err = json.NewEncoder(w).Encode(map[string]*wallet.Transaction{
			"data": t,
		})
		if err != nil {
			fmt.Println(err)
		}
	}).Methods("POST")

	srv := &http.Server{
		Handler:      router,
		Addr:         serverAddress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if port != rootNodePort {
		fmt.Println("Syncing Chain...")
		syncInitialChain(bc)
		fmt.Printf("Synced Chain Successfully (%d blocks)\n", len(bc.GetChain()))
		fmt.Println("-------")
	}

	fmt.Println("Redis Address:", redisAddress)
	fmt.Println("Server Address:", serverAddress)

	log.Fatal(srv.ListenAndServe())
}

func syncInitialChain(bc *blockchain.Blockchain) {
	res, err := http.Get(fmt.Sprintf("%s/api/blocks", rootNodeAddress))
	if err != nil {
		panic(err)
	}
	if res.StatusCode != http.StatusOK {
		panic(errors.New("unable to retrieve root node chain"))
	}

	var chain []*blockchain.Block
	err = json.NewDecoder(res.Body).Decode(&chain)
	if err != nil {
		panic(errors.New("unable to decode root node chain"))
	}

	bc.ReplaceChain(chain)
}
