package main

// NOTE: No GORM tags; simple JSON models.

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"` // plain for demo (DO NOT do this in real prod)
	Token    string `json:"token"`
}

type Item struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

type CartItem struct {
	ID       uint `json:"id"`
	CartID   uint `json:"cart_id"`
	ItemID   uint `json:"item_id"`
	Quantity int  `json:"quantity"`
}

type Cart struct {
	ID     uint       `json:"id"`
	UserID uint       `json:"user_id"`
	Items  []CartItem `json:"items"`
}

type OrderItem struct {
	ID       uint `json:"id"`
	OrderID  uint `json:"order_id"`
	ItemID   uint `json:"item_id"`
	Quantity int  `json:"quantity"`
}

type Order struct {
	ID        uint        `json:"id"`
	UserID    uint        `json:"user_id"`
	CreatedAt string      `json:"created_at"`
	Items     []OrderItem `json:"items"`
}
