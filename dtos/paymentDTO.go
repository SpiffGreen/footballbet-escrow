package dtos

type WebhookQuery struct {
	Trxref    string `form:"trxref"`
	Reference string `form:"reference"`
}

type InitializePayment struct {
	Email  string
	Amount string
}

type SetBankDetails struct {
	Name          string
	AccountNumber string
	BankCode      string
}

type WithdrawFunds struct {
	Amount uint64
}
