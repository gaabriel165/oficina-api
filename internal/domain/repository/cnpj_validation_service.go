package repository

type CNPJInfo struct {
	CNPJ              string
	CompanyName       string
	TradeName         string
	SituacaoCadastral int
	IsActive          bool
}

type CNPJValidationService interface {
	Fetch(cnpj string) (*CNPJInfo, error)
}
