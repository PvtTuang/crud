import React, { useEffect, useState } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom"; // ✅ ใช้ navigate

const History = () => {
  const [histories, setHistories] = useState([]);
  const [loading, setLoading] = useState(false);

  const token = localStorage.getItem("token");
  const userId = localStorage.getItem("user_id");
  const navigate = useNavigate(); // ✅ init navigate

  useEffect(() => {
    if (!token) {
      window.location.href = "/login";
      return;
    }
    fetchHistory();
  }, [token]);

  // ดึงประวัติการสั่งซื้อ
  const fetchHistory = async () => {
    try {
      setLoading(true);
      const res = await axios.get(
        `${import.meta.env.VITE_CRUD_API_URL}/history/${userId}`,
        { headers: { Authorization: `Bearer ${token}` } }
      );
      setHistories(res.data);
    } catch (err) {
      console.error(err);
      alert("โหลดประวัติการสั่งซื้อไม่สำเร็จ");
    } finally {
      setLoading(false);
    }
  };

  if (loading) return <p>กำลังโหลด...</p>;

  if (!histories || histories.length === 0) {
    return (
      <div style={{ padding: "2rem" }}>
        <h2>📜 Purchase History</h2>
        <p>ยังไม่มีประวัติการสั่งซื้อ</p>
        {/* ✅ ปุ่มกลับไป Products */}
        <button onClick={() => navigate("/products")}>⬅️ กลับไปหน้าสินค้า</button>
      </div>
    );
  }

  return (
    <div style={{ padding: "2rem" }}>
      <h2>📜 Purchase History</h2>
      {histories.map((h) => (
        <div
          key={h.id}
          style={{
            border: "1px solid #ccc",
            borderRadius: "8px",
            padding: "1rem",
            marginBottom: "1rem",
          }}
        >
          <h4>🕒 {new Date(h.created_at).toLocaleString()}</h4>
          <ul>
            {h.items?.map((item) => (
              <li key={item.id}>
                {item.product?.name ?? item.product_id} — {item.quantity} ชิ้น
                <span style={{ marginLeft: "0.5rem", color: "gray" }}>
                  ({item.product?.price ?? "?"} ฿)
                </span>
              </li>
            ))}
          </ul>
        </div>
      ))}

      {/* ✅ ปุ่มกลับไป Products */}
      <button
        onClick={() => navigate("/products")}
        style={{ marginTop: "1rem" }}
      >
        ⬅️ กลับไปหน้าสินค้า
      </button>
    </div>
  );
};

export default History;
