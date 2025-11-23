import React, { useEffect, useState } from "react";

export default function Items({ token }) {
  const [items, setItems] = useState([]);
  const [cartCount, setCartCount] = useState(0);

  useEffect(() => {
    loadItems();
    loadCartCount();
  }, []);

  async function loadItems() {
    const res = await fetch("http://localhost:8080/items");
    const data = await res.json();
    setItems(Array.isArray(data) ? data : []);
  }

  async function loadCartCount() {
    const res = await fetch("http://localhost:8080/carts");
    const allCarts = await res.json();

    if (allCarts.length > 0 && allCarts[0].items) {
      setCartCount(allCarts[0].items.length);
    }
  }

  async function addToCart(itemId) {
    const res = await fetch("http://localhost:8080/carts/", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + token,
      },
      body: JSON.stringify({ item_id: itemId, quantity: 1 }),
    });

    if (res.ok) {
      alert("Item added to cart");
      loadCartCount(); // Update badge
    }
  }

  async function checkout() {
    const res = await fetch("http://localhost:8080/carts");
    const carts = await res.json();
    const cart = carts[0];

    if (!cart || cart.items.length === 0) {
      alert("Cart is empty!");
      return;
    }

    const res2 = await fetch("http://localhost:8080/orders/", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + token,
      },
      body: JSON.stringify({ cart_id: cart.id }),
    });

    if (res2.ok) {
      alert("Order successful");
      setCartCount(0);
    }
  }

  async function showCart() {
    const res = await fetch("http://localhost:8080/carts");
    const carts = await res.json();
    alert(JSON.stringify(carts, null, 2));
  }

  async function showOrders() {
    const res = await fetch("http://localhost:8080/orders");
    const orders = await res.json();
    alert(JSON.stringify(orders.map((o) => o.id), null, 2));
  }

  return (
    <div>
      <div className="items-header">
        <div className="items-header-title">Items</div>

        <div className="items-toolbar">

          {/* CHECKOUT BTN */}
          <button className="btn btn-primary" onClick={checkout}>
            Checkout
          </button>

          {/* CART BTN WITH BADGE */}
          <button className="btn btn-secondary" onClick={showCart} style={{ position: "relative" }}>
            Cart
            {cartCount > 0 && (
              <span
                style={{
                  background: "#ef4444",
                  color: "white",
                  borderRadius: "999px",
                  fontSize: "0.7rem",
                  padding: "2px 6px",
                  position: "absolute",
                  top: "-6px",
                  right: "-6px",
                  boxShadow: "0 0 8px rgba(239, 68, 68, 0.6)",
                }}
              >
                {cartCount}
              </span>
            )}
          </button>

          {/* ORDER HISTORY BTN */}
          <button className="btn btn-secondary" onClick={showOrders}>
            Order History
          </button>
        </div>
      </div>

      {/* ITEMS GRID */}
      <div className="items-wrapper">
        <div className="items-grid">
          {items.map((item) => (
            <div key={item.id} className="item-card">
              <div>
                <div className="item-name">{item.name}</div>
                <div className="item-desc">{item.description}</div>
              </div>

              <div className="item-footer">
                <span className="item-price">₹{item.price}</span>
                <button className="btn btn-primary" onClick={() => addToCart(item.id)}>
                  Add to Cart
                </button>
              </div>
            </div>
          ))}
        </div>

        {items.length === 0 && <p className="items-empty">No items available.</p>}
      </div>
    </div>
  );
}
