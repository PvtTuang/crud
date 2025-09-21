import React, { useState, useEffect } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";

const Products = () => {
  const [products, setProducts] = useState([]);
  const [name, setName] = useState("");
  const [price, setPrice] = useState("");
  const [loading, setLoading] = useState(false);

  const [editId, setEditId] = useState(null);
  const [editName, setEditName] = useState("");
  const [editPrice, setEditPrice] = useState("");

  const token = localStorage.getItem("token");
  const userId = localStorage.getItem("user_id"); // 👉 เก็บ user id ไว้หลัง login
  const navigate = useNavigate();

  useEffect(() => {
    if (!token) {
      window.location.href = "/login";
      return;
    }
    fetchProducts();
  }, [token]);

  // ดึงสินค้า
  const fetchProducts = async () => {
    try {
      setLoading(true);
      const res = await axios.get(
        `${import.meta.env.VITE_CRUD_API_URL}/products`,
        {
          headers: { Authorization: `Bearer ${token}` },
        }
      );

      const mappedProducts = res.data.map((p) => ({
        id: p.ID ?? p.id,
        name: p.Name ?? p.name,
        price: p.Price ?? p.price,
      }));
      setProducts(mappedProducts);
    } catch (err) {
      alert("โหลดสินค้าไม่สำเร็จ");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  // เพิ่มสินค้าใหม่
  const handleCreate = async (e) => {
    e.preventDefault();
    if (!name || !price) {
      alert("กรุณากรอกชื่อและราคาสินค้า");
      return;
    }
    try {
      await axios.post(
        `${import.meta.env.VITE_CRUD_API_URL}/products`,
        { name, price: parseFloat(price) },
        { headers: { Authorization: `Bearer ${token}` } }
      );
      setName("");
      setPrice("");
      fetchProducts();
    } catch (err) {
      alert("สร้างสินค้าไม่สำเร็จ");
      console.error(err);
    }
  };

  // ลบสินค้า
  const handleDelete = async (id) => {
    if (!window.confirm("คุณต้องการลบสินค้านี้หรือไม่?")) return;

    try {
      await axios.delete(
        `${import.meta.env.VITE_CRUD_API_URL}/products/${id}`,
        { headers: { Authorization: `Bearer ${token}` } }
      );
      fetchProducts();
    } catch (err) {
      alert("ลบสินค้าไม่สำเร็จ");
      console.error(err);
    }
  };

  // เพิ่มสินค้าไป Cart
  const handleAddToCart = async (id) => {
    try {
      await axios.post(
        `${import.meta.env.VITE_CRUD_API_URL}/carts/${userId}/products`,
        { product_id: id, quantity: 1 },
        { headers: { Authorization: `Bearer ${token}` } }
      );
      alert("เพิ่มสินค้าไปที่ตะกร้าแล้ว ✅");
    } catch (err) {
      alert("เพิ่มสินค้าไป Cart ไม่สำเร็จ");
      console.error(err);
    }
  };

  // เริ่มแก้ไข
  const startEdit = (product) => {
    setEditId(product.id);
    setEditName(product.name);
    setEditPrice(product.price);
  };

  // ยกเลิกแก้ไข
  const cancelEdit = () => {
    setEditId(null);
    setEditName("");
    setEditPrice("");
  };

  // บันทึกการแก้ไข
  const handleUpdate = async (id) => {
    try {
      await axios.put(
        `${import.meta.env.VITE_CRUD_API_URL}/products/${id}`,
        { name: editName, price: parseFloat(editPrice) },
        { headers: { Authorization: `Bearer ${token}` } }
      );
      cancelEdit();
      fetchProducts();
    } catch (err) {
      alert("แก้ไขสินค้าไม่สำเร็จ");
      console.error(err);
    }
  };

  return (
    <div style={{ padding: "2rem" }}>
      <h2>Products</h2>

      {/* ปุ่มไป Cart และ History */}
      <div style={{ marginBottom: "1rem" }}>
        <button onClick={() => navigate("/cart")} style={{ marginRight: "1rem" }}>
          🛒 ไปที่ Cart
        </button>
        <button onClick={() => navigate("/history")}>
          📜 ดูประวัติการสั่งซื้อ
        </button>
      </div>

      {/* Form เพิ่มสินค้า */}
      <form onSubmit={handleCreate} style={{ marginBottom: "2rem" }}>
        <input
          type="text"
          placeholder="Name"
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
        <input
          type="number"
          placeholder="Price"
          value={price}
          onChange={(e) => setPrice(e.target.value)}
        />
        <button type="submit">Add Product</button>
      </form>

      <hr />

      {/* List สินค้า */}
      {loading ? (
        <p>กำลังโหลดสินค้า...</p>
      ) : products.length === 0 ? (
        <p>ยังไม่มีสินค้า</p>
      ) : (
        <ul>
          {products.map((p) => (
            <li key={p.id} style={{ marginBottom: "0.5rem" }}>
              {editId === p.id ? (
                <>
                  <input
                    type="text"
                    value={editName}
                    onChange={(e) => setEditName(e.target.value)}
                  />
                  <input
                    type="number"
                    value={editPrice}
                    onChange={(e) => setEditPrice(e.target.value)}
                  />
                  <button onClick={() => handleUpdate(p.id)}>Save</button>
                  <button
                    onClick={cancelEdit}
                    style={{ marginLeft: "0.5rem" }}
                  >
                    Cancel
                  </button>
                </>
              ) : (
                <>
                  {p.name} - {p.price}฿
                  <button
                    onClick={() => startEdit(p)}
                    style={{ marginLeft: "1rem" }}
                  >
                    Edit
                  </button>
                  <button
                    onClick={() => handleDelete(p.id)}
                    style={{ marginLeft: "0.5rem", color: "red" }}
                  >
                    Delete
                  </button>
                  <button
                    onClick={() => handleAddToCart(p.id)}
                    style={{ marginLeft: "0.5rem", color: "green" }}
                  >
                    Add to Cart
                  </button>
                </>
              )}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export default Products;
