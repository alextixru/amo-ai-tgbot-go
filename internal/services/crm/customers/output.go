package customers

// CustomerOutput покупатель в читаемом формате (с именами вместо числовых ID)
type CustomerOutput struct {
	ID                  int                    `json:"id"`
	Name                string                 `json:"name,omitempty"`
	NextPrice           int                    `json:"next_price,omitempty"`
	NextDate            string                 `json:"next_date,omitempty"`
	StatusName          string                 `json:"status_name,omitempty"`
	Periodicity         int                    `json:"periodicity,omitempty"`
	ResponsibleUserName string                 `json:"responsible_user_name,omitempty"`
	CreatedByName       string                 `json:"created_by_name,omitempty"`
	UpdatedByName       string                 `json:"updated_by_name,omitempty"`
	CreatedAt           string                 `json:"created_at,omitempty"`
	UpdatedAt           string                 `json:"updated_at,omitempty"`
	IsDeleted           bool                   `json:"is_deleted,omitempty"`
	Ltv                 int                    `json:"ltv,omitempty"`
	PurchasesCount      int                    `json:"purchases_count,omitempty"`
	AverageCheck        int                    `json:"average_check,omitempty"`
	Tags                []string               `json:"tags,omitempty"`
	Segments            []CustomerSegmentBrief `json:"segments,omitempty"`
	Contacts            []CustomerEntityBrief  `json:"contacts,omitempty"`
	Companies           []CustomerEntityBrief  `json:"companies,omitempty"`
}

// CustomerSegmentBrief краткая информация о сегменте
type CustomerSegmentBrief struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

// CustomerEntityBrief краткая информация о связанной сущности
type CustomerEntityBrief struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

// BonusPointsResult результат операции с бонусными баллами
type BonusPointsResult struct {
	// Balance текущий баланс после операции
	Balance int `json:"balance"`
	// Operation тип операции: "earn" или "redeem"
	Operation string `json:"operation"`
	// Points количество баллов в операции
	Points int `json:"points"`
}

// CustomersListOutput список покупателей с пагинацией
type CustomersListOutput struct {
	Customers []*CustomerOutput `json:"customers"`
	HasMore   bool              `json:"has_more"`
}

// TransactionsListOutput список транзакций с пагинацией
type TransactionsListOutput struct {
	Transactions []TransactionOutput `json:"transactions"`
	HasMore      bool                `json:"has_more"`
}

// TransactionOutput транзакция в читаемом формате
type TransactionOutput struct {
	ID        int    `json:"id"`
	Price     int    `json:"price,omitempty"`
	Comment   string `json:"comment,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

// SegmentsListOutput список сегментов с пагинацией
type SegmentsListOutput struct {
	Segments []*SegmentOutput `json:"segments"`
	HasMore  bool             `json:"has_more"`
}

// SegmentOutput сегмент в читаемом формате
type SegmentOutput struct {
	ID             int    `json:"id"`
	Name           string `json:"name,omitempty"`
	Color          string `json:"color,omitempty"`
	CustomersCount int    `json:"customers_count,omitempty"`
}
