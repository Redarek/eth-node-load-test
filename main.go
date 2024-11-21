package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"sync"
	"time"
)

func main() {
	// Настройки
	nodeURL := "ws://websocket_endpoint" // Ваш WebSocket эндпоинт
	privateKeyHex := "from_private_key"  // Приватный ключ отправителя
	toAddressHex := "0x_to_public_key"   // Адрес получателя

	// Подключение к ноде
	client, err := ethclient.Dial(nodeURL)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum node: %v", err)
	}
	defer client.Close()

	// Преобразование приватного ключа
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Invalid private key: %v", err)
	}

	// Получение адреса отправителя
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKey)

	// Получение стартового nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}

	// Получение цены газа
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to get gas price: %v", err)
	}

	toAddress := common.HexToAddress(toAddressHex)
	amount := big.NewInt(1000000000000000) // Сумма в Wei (0.001 ETH)

	fmt.Println("Starting DoS test...")

	// Количество горутин
	numGoroutines := 10

	// Канал для передачи ошибок
	errCh := make(chan error, numGoroutines)

	// WaitGroup для ожидания завершения всех горутин
	var wg sync.WaitGroup

	// Запуск горутин
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(startNonce uint64) {
			defer wg.Done()
			sendTransactions(client, privateKey, fromAddress, toAddress, amount, gasPrice, startNonce, errCh)
		}(nonce)
		nonce += 1000 // Увеличиваем nonce для предотвращения конфликтов
	}

	// Закрытие канала ошибок после завершения горутин
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// Обработка ошибок
	for err := range errCh {
		if err != nil {
			log.Printf("Error during transaction sending: %v", err)
		}
	}
}
func sendTransactions(client *ethclient.Client, privateKey *ecdsa.PrivateKey, fromAddress, toAddress common.Address, amount, gasPrice *big.Int, startNonce uint64, errCh chan<- error) {
	// Счетчик запросов
	var requestCount int
	startTime := time.Now()

	for {
		// Создание транзакции
		tx := types.NewTransaction(startNonce, toAddress, amount, uint64(21000), gasPrice, nil)

		// Подписание транзакции
		chainID, err := client.NetworkID(context.Background())
		if err != nil {
			errCh <- fmt.Errorf("failed to get chain ID: %v", err)
			return
		}

		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
		if err != nil {
			errCh <- fmt.Errorf("failed to sign transaction: %v", err)
			return
		}

		// Отправка транзакции
		err = client.SendTransaction(context.Background(), signedTx)
		if err != nil {
			if err.Error() == "nonce too low" {
				startNonce++
				continue // Пропускаем и пробуем следующий nonce
			}
			errCh <- fmt.Errorf("failed to send transaction: %v", err)
		} else {
			requestCount++
			fmt.Printf("Transaction sent: %s\n", signedTx.Hash().Hex())
		}

		// Увеличение nonce
		startNonce++

		// Замер количества запросов в секунду
		elapsed := time.Since(startTime)
		if elapsed >= time.Second {
			fmt.Printf("Requests per second: %d\n", requestCount)
			requestCount = 0
			startTime = time.Now()
		}
	}
}
