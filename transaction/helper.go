package transaction

const (
	TransactionCollection     = "Transaction"
	OperationTypeCollection   = "OperationType"
	InvalidAccountError       = "Account is invalid"
	InvalidOperationTypeError = "OperationType is invalid"
)

type Helper struct{}

type IHelper interface {
	TransformAmount(amount float64, isCredit bool) float64
}

func NewTransactionHelper() *Helper {
	return &Helper{}
}

func (helper *Helper) TransformAmount(amount float64, isCredit bool) float64 {
	if !isCredit {
		return amount * (-1)
	}
	return amount
}
