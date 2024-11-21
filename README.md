# Load Testing Tool for Ethereum Private Node

This project is a **Go-based load testing tool** designed to stress-test a private Ethereum node using **transactions**. It leverages **goroutines** for parallel execution and logs **requests per second (RPS)** to monitor performance.

---

## Features

- **Concurrent Execution**: Utilizes multiple goroutines to send transactions in parallel, maximizing node load.
- **Efficient Nonce Handling**: Each goroutine operates within its unique `nonce` range to prevent conflicts between transactions.
- **RPS Monitoring**: Logs requests per second for each goroutine, helping analyze system performance.
- **Customizable Settings**: Easily configurable WebSocket endpoint, private key, recipient address, and transaction parameters.

---

## How It Works

1. **Setup**: Connects to an Ethereum node using a WebSocket endpoint.
2. **Nonce Management**: Each goroutine starts with a unique `nonce` range and increments it independently to avoid collisions.
3. **Transaction Signing and Sending**:
   - Creates a raw Ethereum transaction.
   - Signs the transaction with the sender's private key.
   - Sends the transaction to the node.
4. **Performance Logging**: Each goroutine logs the number of requests it processes per second.

---

## Prerequisites

- Go installed on your system (1.18 or later).
- Access to a private Ethereum node with WebSocket enabled.
- A funded Ethereum account (private key).
- A recipient Ethereum address.

---

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/your-repo/ethereum-load-tester.git
   cd ethereum-load-tester
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

---

## Configuration

Edit the following parameters in the `main.go` file:
- **Node WebSocket URL**:
  ```go
  nodeURL := "ws://<NODE_IP>:<PORT>"
  ```
- **Sender's Private Key**:
  ```go
  privateKeyHex := "<YOUR_PRIVATE_KEY>"
  ```
- **Recipient Address**:
  ```go
  toAddressHex := "0x<RECIPIENT_ADDRESS>"
  ```
- **Transaction Amount** (in Wei):
  ```go
  amount := big.NewInt(1000000000000000) // 0.001 ETH
  ```

---

## Running the Tool

To run the load tester, execute:

```bash
go run main.go
```

---

## Output

1. **Transaction Logs**:
   - Logs each successfully sent transaction:
     ```plaintext
     Transaction sent: 0x123abc...
     ```
2. **Requests Per Second (RPS)**:
   - Logs the number of transactions processed per second by each goroutine:
     ```plaintext
     Requests per second: 150
     ```

---

## Code Highlights

1. **Goroutine Usage**:
   - Multiple goroutines (`numGoroutines`) are used to send transactions in parallel:
     ```go
     for i := 0; i < numGoroutines; i++ {
         go sendTransactions(...)
     }
     ```

2. **Nonce Management**:
   - Each goroutine starts with a unique `nonce` range, avoiding collisions:
     ```go
     nonce += 1000 // Increment to separate ranges
     ```

3. **RPS Logging**:
   - Logs RPS for each goroutine to analyze transaction throughput:
     ```go
     fmt.Printf("Requests per second: %d\n", requestCount)
     ```

---

## Notes

- Ensure the sender's account has sufficient funds for both the transaction amounts and gas fees.
- Test only on **private/test networks** to avoid unintended consequences.
- You can increase the number of goroutines (`numGoroutines`) to apply more load.

---

## License

This project is licensed under the [MIT License](LICENSE).

---

Feel free to reach out for questions or suggestions!
