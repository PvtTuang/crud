import React, { useEffect, useState } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";

const Cart = () => {
  const [cart, setCart] = useState(null);
  const [loading, setLoading] = useState(false);

  const token = localStorage.getItem("token");
  const userId = localStorage.getItem("user_id");
  const navigate = useNavigate();

  useEffect(() => {
    if (!token) {
      window.location.href = "/login";
      return;
    }
    fetchCart();
  }, [token]);

  // ดึงตะกร้า
  const fetchCart = async () => {
    try {
      setLoading(true);
      const res = await axios.get(
        `${import.meta.env.VITE_CRUD_API_URL}/carts/${userId}`,
        { headers: { Authorization: `Bearer ${token}` } }
      );

      // ✅ กรองสินค้าที่ถูกลบออก
      const filtered = {
        ...res.data,
        products: res.data.products.filter((p) => !p.product?.deleted_at),
      };
      setCart(filtered);
    } catch (err) {
      console.error(err);
      alert("โหลดตะกร้าไม่สำเร็จ");
    } finally {
      setLoading(false);
    }
  };

  // ลบสินค้าออกจาก cart
  const removeFromCart = async (productId) => {
    try {
      // optimistic update
      setCart((prev) => ({
        ...prev,
        products: prev.products.filter((p) => p.product.id !== productId),
      }));

      await axios.delete(
        `${import.meta.env.VITE_CRUD_API_URL}/carts/${userId}/products/${productId}`,
        { headers: { Authorization: `Bearer ${token}` } }
      );
    } catch (err) {
      console.error(err);
      alert("ลบสินค้าไม่สำเร็จ");
      fetchCart();
    }
  };

  // เคลียร์ตะกร้า
  const clearCart = async () => {
    if (!window.confirm("ต้องการล้างตะกร้าทั้งหมดหรือไม่?")) return;
    try {
      setCart({ ...cart, products: [] });
      await axios.delete(
        `${import.meta.env.VITE_CRUD_API_URL}/carts/${userId}/clear`,
        { headers: { Authorization: `Bearer ${token}` } }
      );
    } catch (err) {
      console.error(err);
      alert("ล้างตะกร้าไม่สำเร็จ");
      fetchCart();
    }
  };

  // Checkout
  const checkout = async () => {
    if (!window.confirm("ยืนยันการสั่งซื้อ?")) return;
    try {
      await axios.post(
        `${import.meta.env.VITE_CRUD_API_URL}/carts/${userId}/checkout`,
        {},
        { headers: { Authorization: `Bearer ${token}` } }
      );
      alert("✅ Checkout success");
      setCart({ ...cart, products: [] });
    } catch (err) {
      console.error(err);
      alert("Checkout ไม่สำเร็จ");
      fetchCart();
    }
  };

  if (loading) return <p>กำลังโหลด...</p>;
  if (!cart || !cart.products || cart.products.length === 0) {
    return (
      <div style={{ padding: "2rem" }}>
        <h2>Cart</h2>
        <p>ยังไม่มีสินค้าในตะกร้า</p>
        <button onClick={() => navigate("/products")}>
          ⬅️ กลับไปหน้าสินค้า
        </button>
      </div>
    );
  }

  const totalPrice = cart.products.reduce(
    (sum, p) => sum + (p.product?.price ?? 0) * p.quantity,
    0
  );

  return (
    <div style={{ padding: "2rem" }}>
      <h2>🛒 Cart</h2>
      <ul>
        {cart.products
          .filter((p) => !p.product?.deleted_at)
          .map((p) => (
            <li key={p.id}>
              {p.product?.name ?? p.product_id} — {p.quantity} ชิ้น
              <span style={{ marginLeft: "0.5rem", color: "gray" }}>
                ({p.product?.price ?? "?"} ฿)
              </span>
              <button
                onClick={() => removeFromCart(p.product.id)}
                style={{ marginLeft: "1rem", color: "red" }}
              >
                ลบ
              </button>
            </li>
          ))}
      </ul>
      <h3 style={{ marginTop: "1rem" }}>💰 รวมทั้งหมด: {totalPrice} ฿</h3>
      <div style={{ marginTop: "1rem" }}>
        <button onClick={checkout} style={{ marginRight: "1rem" }}>
          ✅ Checkout
        </button>
        <button
          onClick={clearCart}
          style={{ marginRight: "1rem", color: "red" }}
        >
          🗑️ Clear Cart
        </button>
        <button onClick={() => navigate("/products")}>
          ⬅️ กลับไปหน้าสินค้า
        </button>
      </div>
    </div>
  );
};

export default Cart;
