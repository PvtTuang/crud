import React, { useState } from "react";
import axios from "axios";

const Login = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const handleLogin = async (e) => {
    e.preventDefault();
    try {
      const res = await axios.post(
        `${import.meta.env.VITE_AUTH_API_URL}/login`,
        { username, password }
      );

      // ✅ เก็บ token
      localStorage.setItem("token", res.data.token);

      // ✅ ถ้ามี user_id จาก backend ให้เก็บด้วย
      if (res.data.user_id) {
        localStorage.setItem("user_id", res.data.user_id);
      }

      alert("เข้าสู่ระบบสำเร็จ!");
      window.location.href = "/products";
    } catch (err) {
      alert("เข้าสู่ระบบไม่สำเร็จ: " + (err.response?.data?.error || err.message));
    }
  };

  return (
    <div style={{ padding: "2rem" }}>
      <h2>Login</h2>
      <form onSubmit={handleLogin}>
        <div>
          <input
            type="text"
            placeholder="Username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
        </div>
        <br />
        <div>
          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </div>
        <br />
        <button type="submit">Login</button>
      </form>
      <p>
        ยังไม่มีบัญชี? <a href="/register">สมัครสมาชิก</a>
      </p>
    </div>
  );
};

export default Login;
