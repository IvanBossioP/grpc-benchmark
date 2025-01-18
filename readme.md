# Benchmark Tool for Solana GRPC Yellowstone Transaction Detection

This tool benchmarks multiple Solana nodes GRPC geyser by detecting transactions and calculating which node is the fastest in detecting transactions for a specific address.

## Setup

### Prerequisites
1. **Go Environment**:
   - Ensure you have Go 1.22 or later installed.
   - Download and install from [https://go.dev/dl/](https://go.dev/dl/).

2. **Dependencies**:
   - The program requires some libraries:
     - `grpc`
     - `protobuf`
   - Run the following command to install dependencies:
     ```bash
     go mod tidy
     ```

3. **Configuration File**:
   - Update the `config.json` file in the root directory using the following format:
     ```json
     {
       "nodes": [
         "http://grpc.node1.com:10000",
         "http://grpc.node2.com:10000"
       ],
       "benchmarkDuration": 30,
       "detectionAddress": "675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8"
     }
     ```
     - **nodes**: An array of Solana node addresses.
     - **benchmarkDuration**: Duration of the benchmark in seconds, defaults to 30 seconds.
     - **detectionAddress**: A valid Solana address to detect transactions, default one is Raydium V4 program.

---

## Running the Program

1. Build the program:
   - **Linux/MacOS**:
     ```bash
     go build -o benchmark
     ```
   - **Windows**:
     ```bash
     go build -o benchmark.exe
     ```

2. Run the program:
   - **Linux/MacOS**:
     ```bash
     ./benchmark
     ```
   - **Windows**:
     ```cmd
     benchmark.exe
     ```

## Example Output

```plaintext
Starting benchmark, duration 10 seconds...
Time left: 10s, 15 transactions detected
Time left: 9s, 124 transactions detected
...
Benchmark completed, calculating results...
Detected 2623 transactions
The winner is http://grpc.node1.com:10000
With 1671 transactions detected first and a winrate of 63.71%
```

---

## Possible Errors

1. **"Failed to open file"**:
   - Ensure the `config.json` file exists in the same directory as the program.

2. **"Failed to unmarshal JSON"**:
   - Verify the `config.json` file is correctly formatted.

3. **"There gotta be at least 2 nodes in the config.json"**:
   - Provide at least two Solana nodes in the `nodes` array.

5. **"Detection address is not a valid Solana address"**:
   - Provide a valid Solana address in the `detectionAddress` field.

## Notes
- This tool is designed to benchmark Solana nodes. Ensure the nodes provided in `config.json` are accessible and functional, nodes often require ip whitelisting.
- The benchmark detects new transactions on the "processed" level ignoring vote and failed transactions

## Limitations
- The benchmark program doesn't currently support x-token based GRPC geysers, feel free to make a pull request if you want to add support
- The benchmark program doesn't currently reconnect in case of disconnection, that would be useful in case geysers tested struggle to maintain a stable connection, again feel free make a pull request if you want to add support