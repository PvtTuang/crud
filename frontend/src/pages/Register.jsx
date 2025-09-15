import React, { useState } from "react";
import axios from "axios";

const Register = () => {
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleRegister = async (e) => {
    e.preventDefault();
    try {
      await axios.post(`${import.meta.env.VITE_AUTH_API_URL}/register`, {
        username,
        email,
        password,
      });

      alert("สมัครสมาชิกสำเร็จ! กรุณาเข้าสู่ระบบ");
      window.location.href = "/login";
    } catch (err) {
      // ปรับการแสดง error ให้ชัดเจน
      const errorMsg =
        err.response?.data?.error ||
        err.response?.data?.message ||
        err.message;
      alert("สมัครไม่สำเร็จ: " + JSON.stringify(errorMsg));
      console.error(err);
    }
  };

  return (
    <div style={{ padding: "2rem" }}>
      <h2>Register</h2>
      <form onSubmit={handleRegister}>
        <div>
          <input
            type="text"
            placeholder="Username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
          />
        </div>
        <br />
        <div>
          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>
        <br />
        <div>
          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        <br />
        <button type="submit">Register</button>
      </form>

      <p>
        มีบัญชีแล้ว? <a href="/login">เข้าสู่ระบบ</a>
      </p>
    </div>
  );
};

export default Register;
