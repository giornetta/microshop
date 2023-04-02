package events

const ProductTopic Topic = "Products"

func init() {
	registerEvent[ProductCreated](ProductCreatedType)
	registerEvent[ProductUpdated](ProductUpdatedType)
	registerEvent[ProductDeleted](ProductDeletedType)
}

const (
	ProductCreatedType Type = "Product.Created"
	ProductUpdatedType Type = "Product.Updated"
	ProductDeletedType Type = "Product.Deleted"
)

type ProductEvent struct {
	ProductId string `json:"product_id"`
}

func (e ProductEvent) Key() Key { return Key(e.ProductId) }

func (e ProductEvent) Topic() Topic { return ProductTopic }

type ProductCreated struct {
	ProductEvent
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Amount      int     `json:"amount"`
}

func (ProductCreated) Type() Type { return ProductCreatedType }

type ProductUpdated struct {
	ProductEvent
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Amount      int     `json:"amount"`
}

func (ProductUpdated) Type() Type { return ProductUpdatedType }

type ProductDeleted struct {
	ProductEvent
}

func (ProductDeleted) Type() Type { return ProductDeletedType }
