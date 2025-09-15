import React, { useState, useEffect } from "react";
import axios from "axios";

const Products = () => {
  const [products, setProducts] = useState([]);
  const [name, setName] = useState("");
  const [price, setPrice] = useState("");
  const [loading, setLoading] = useState(false);
  const token = localStorage.getItem("token"); // JWT จาก localStorage

  // ตรวจสอบ token และโหลดสินค้า
  useEffect(() => {
    if (!token) {
      window.location.href = "/login";
      return;
    }
    fetchProducts();
  }, [token]);

  // ฟังก์ชันดึงสินค้า
  const fetchProducts = async () => {
    try {
      setLoading(true);
      const res = await axios.get(`${import.meta.env.VITE_CRUD_API_URL}/products`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      console.log("Backend response:", res.data);

      // map backend field ให้ตรง frontend
      const mappedProducts = res.data.map((p) => ({
        id: p.ID ?? p.id,        // GORM struct ใช้ ID
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

  // สร้างสินค้าใหม่
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
      fetchProducts(); // refresh list
    } catch (err) {
      alert("สร้างสินค้าไม่สำเร็จ");
      console.error(err);
    }
  };

  // ลบสินค้า
  const handleDelete = async (id) => {
    if (!id) {
      alert("ID สินค้าไม่ถูกต้อง");
      return;
    }
    if (!window.confirm("คุณต้องการลบสินค้านี้หรือไม่?")) return;

    try {
      await axios.delete(`${import.meta.env.VITE_CRUD_API_URL}/products/${id}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      fetchProducts();
    } catch (err) {
      alert("ลบสินค้าไม่สำเร็จ");
      console.error(err);
    }
  };

  return (
    <div style={{ padding: "2rem" }}>
      <h2>Products</h2>

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
              {p.name} - {p.price}฿
              <button
                onClick={() => handleDelete(p.id)}
                style={{ marginLeft: "1rem", color: "red" }}
              >
                Delete
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export default Products;
