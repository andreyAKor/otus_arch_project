package blockchains

import "context"

type Blockchainer interface {
	// Расчет стоимости транзакции
	CalcTransaction(ctx context.Context, value string) (TransactionPrice, error)

	// Шлет транзакцию
	Send(
		ctx context.Context,
		walletFrom string, // Кошелек с которого отправляем
		walletTo string, // Кошелек на который отправляем
		value string, // Переводимое значение
	) (string, error)

	// Возвращает общую сумму баланса обменника в блокчейне
	SummaryBalance(ctx context.Context) (string, error)
}

type TransactionPrice struct {
	Value     string
	Fee       string
	AvgTxSize int
}
