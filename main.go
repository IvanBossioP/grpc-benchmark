package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"grpc-benchmark/connection"
	pb "grpc-benchmark/protobuf"
	"grpc-benchmark/util"
	"io"
	"log"
	"os"
	"sort"
	"sync"
	"time"
)

type BenchmarkConfig struct {
	Nodes             []string `json:"nodes"`
	BenchmarkDuration int      `json:"benchmarkDuration"`
	DetectionAddress  string   `json:"detectionAddress"`
}

type DetectionResult struct {
	Node      string
	Timestamp int
}

type BenchmarkResult struct {
	Node       string
	TimesFirst int
	WinRate    float64
}

func main() {

	file, err := os.Open("config.json")

	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	var config BenchmarkConfig
	if err := json.Unmarshal(byteValue, &config); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if len(config.Nodes) < 2 {
		log.Fatalf("There gotta be at least 2 nodes in the config.json")
	}

	if config.BenchmarkDuration <= 0 {
		log.Fatalf("Benchmark duration has to be positive and more than zero")
	}

	if !util.IsValidSolanaAddress(config.DetectionAddress) {
		log.Fatalf("Detection address is not a valid solana address")
	}

	ctx := context.Background()

	subscriptionRequest := pb.SubscribeRequest{
		Transactions: map[string]*pb.SubscribeRequestFilterTransactions{
			"filter": {
				Vote:           util.BoolPtr(false),
				Failed:         util.BoolPtr(false),
				AccountInclude: []string{config.DetectionAddress},
			},
		},
		Commitment: util.CommitmentPtr(pb.CommitmentLevel_PROCESSED),
	}

	updateMx := sync.Mutex{}

	detections := map[string][]DetectionResult{}

	results := map[string]int{}

	total := 0

	fmt.Printf("Starting benchmark, duration %d seconds...\n", config.BenchmarkDuration)

	for _, node := range config.Nodes {
		results[node] = 0
		client := connection.InitGeyserClient(node)

		go func(client pb.GeyserClient) {
			stream, err := client.Subscribe(ctx)
			if err != nil {
				log.Fatalf("Failed to subscribe on address %s, error: %v", node, err)
			}

			stream.Send(&subscriptionRequest)

			for {
				msg, err := stream.Recv()
				if err != nil {
					log.Fatalf("Error receiving message: %v", err)
				}

				if _, ok := msg.GetUpdateOneof().(*pb.SubscribeUpdate_Transaction); ok {

					tx := msg.GetTransaction()

					base64Sig := base64.RawStdEncoding.EncodeToString(tx.Transaction.Signature)

					result := DetectionResult{
						Timestamp: time.Now().Nanosecond(),
						Node:      node,
					}

					updateMx.Lock()
					if _, ok := detections[base64Sig]; !ok {
						detections[base64Sig] = []DetectionResult{result}
					} else {

						detections[base64Sig] = append(detections[base64Sig], result)
					}

					if len(detections[base64Sig]) == len(config.Nodes) {

						total++

						winner := detections[base64Sig][0]

						results[winner.Node]++

						delete(detections, base64Sig)
					}
					updateMx.Unlock()

				}

			}
		}(client)
	}

	for i := config.BenchmarkDuration; i > 0; i-- {

		fmt.Printf("\rTime left: %ds, %d transactions detected", i, total)
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("\rBenchmark completed, calcolating results...\n")

	if total == 0 {
		fmt.Printf("No transactions detected, consider increasing the benchmark duration or changing detection address")
		return
	}

	resultsArray := []BenchmarkResult{}

	for k, v := range results {

		resultsArray = append(resultsArray, BenchmarkResult{
			Node: k, TimesFirst: v, WinRate: (float64(v) / float64(total)) * 100,
		})

	}

	sort.Slice(resultsArray, func(i, j int) bool {
		return resultsArray[i].TimesFirst > resultsArray[j].TimesFirst
	})

	fmt.Printf("Detected %d transactions\n", total)

	fmt.Printf("The winner is %s\n", resultsArray[0].Node)
	fmt.Printf("Ranking:\n")

	for i, result := range resultsArray {
		fmt.Printf("%d) %s detected first %d transactions with a winrate of %.2f%%\n", i+1, result.Node, result.TimesFirst, result.WinRate)

	}

}
