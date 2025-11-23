package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ========== helpers ==========

func findUserByUsernameAndPassword(username, password string) *User {
	for i := range DB.Users {
		if DB.Users[i].Username == username && DB.Users[i].Password == password {
			return &DB.Users[i]
		}
	}
	return nil
}

func findUserByToken(token string) *User {
	for i := range DB.Users {
		if DB.Users[i].Token == token {
			return &DB.Users[i]
		}
	}
	return nil
}

func findCartByUserID(userID uint) *Cart {
	for i := range DB.Carts {
		if DB.Carts[i].UserID == userID {
			return &DB.Carts[i]
		}
	}
	return nil
}

func findCartByIDForUser(cartID, userID uint) *Cart {
	for i := range DB.Carts {
		if DB.Carts[i].ID == cartID && DB.Carts[i].UserID == userID {
			return &DB.Carts[i]
		}
	}
	return nil
}

// ========== USERS ==========

func CreateUser(c *gin.Context) {
	var in struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := User{
		ID:       nextUserID,
		Username: in.Username,
		Password: in.Password,
		Token:    "",
	}
	nextUserID++
	DB.Users = append(DB.Users, user)
	SaveDB()
	c.JSON(http.StatusCreated, gin.H{"id": user.ID, "username": user.Username})
}

func ListUsers(c *gin.Context) {
	c.JSON(http.StatusOK, DB.Users)
}

func LoginUser(c *gin.Context) {
	var in struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := findUserByUsernameAndPassword(in.Username, in.Password)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username/password"})
		return
	}

	// single token per user
	token := uuid.NewString()
	user.Token = token
	SaveDB()

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// ========== ITEMS ==========

func CreateItem(c *gin.Context) {
	var in struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Price       int    `json:"price"`
	}
	if err := c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item := Item{
		ID:          nextItemID,
		Name:        in.Name,
		Description: in.Description,
		Price:       in.Price,
	}
	nextItemID++
	DB.Items = append(DB.Items, item)
	SaveDB()

	c.JSON(http.StatusCreated, item)
}

func ListItems(c *gin.Context) {
	c.JSON(http.StatusOK, DB.Items)
}

// ========== CARTS ==========

func CreateOrAddToCart(c *gin.Context) {
	u := c.MustGet("user").(User)

	var in struct {
		ItemID   uint `json:"item_id"`
		Quantity int  `json:"quantity"`
	}
	if err := c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if in.Quantity <= 0 {
		in.Quantity = 1
	}

	// find or create cart
	cart := findCartByUserID(u.ID)
	if cart == nil {
		newCart := Cart{
			ID:     nextCartID,
			UserID: u.ID,
			Items:  []CartItem{},
		}
		nextCartID++
		DB.Carts = append(DB.Carts, newCart)
		cart = &DB.Carts[len(DB.Carts)-1]
	}

	// find existing cart item
	var existing *CartItem
	for i := range cart.Items {
		if cart.Items[i].ItemID == in.ItemID {
			existing = &cart.Items[i]
			break
		}
	}

	if existing == nil {
		ci := CartItem{
			ID:       nextCartItemID,
			CartID:   cart.ID,
			ItemID:   in.ItemID,
			Quantity: in.Quantity,
		}
		nextCartItemID++
		cart.Items = append(cart.Items, ci)
	} else {
		existing.Quantity += in.Quantity
	}

	SaveDB()
	c.JSON(http.StatusOK, cart)
}

func ListCarts(c *gin.Context) {
	c.JSON(http.StatusOK, DB.Carts)
}

// ========== ORDERS ==========

func CreateOrderFromCart(c *gin.Context) {
	u := c.MustGet("user").(User)

	var in struct {
		CartID uint `json:"cart_id"`
	}
	if err := c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart := findCartByIDForUser(in.CartID, u.ID)
	if cart == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cart not found for user"})
		return
	}
	if len(cart.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cart is empty"})
		return
	}

	order := Order{
		ID:        nextOrderID,
		UserID:    u.ID,
		CreatedAt: time.Now().Format(time.RFC3339),
		Items:     []OrderItem{},
	}
	nextOrderID++

	for _, ci := range cart.Items {
		oi := OrderItem{
			ID:       nextOrderItemID,
			OrderID:  order.ID,
			ItemID:   ci.ItemID,
			Quantity: ci.Quantity,
		}
		nextOrderItemID++
		order.Items = append(order.Items, oi)
	}

	DB.Orders = append(DB.Orders, order)

	// clear cart items
	cart.Items = []CartItem{}
	SaveDB()

	c.JSON(http.StatusCreated, order)
}

func ListOrders(c *gin.Context) {
	c.JSON(http.StatusOK, DB.Orders)
}
