import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import axiosInstance from "../api/axiosInstance";

const Login: React.FC = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    const token = localStorage.getItem("token");

    if (token) {
      const ws = new WebSocket("ws://localhost:8080/ws/auth");

      ws.onopen = () => {
        console.log("WebSocket connection opened for authenticated user.");
        ws.send(JSON.stringify({ type: "auth", token }));
      };

      ws.onmessage = (event) => {
        const data = JSON.parse(event.data);

        if (data.type === "auth" && data.success) {
          console.log("Authenticated successfully.");
          ws.send(JSON.stringify({ type: "getSelf" }));
        } else if (data.type === "get_self_response") {
          const { id, username } = data.data.user;

          // Save the user ID and username
          localStorage.setItem("user_id", id.toString());
          localStorage.setItem("username", username);

          navigate("/chat");
        } else if (data.type === "error") {
          console.error("Error from WebSocket:", data.message);
          localStorage.removeItem("token");
          localStorage.removeItem("user_id");
          localStorage.removeItem("username");
        }
      };

      ws.onerror = () => {
        console.error("WebSocket error occurred.");
        localStorage.removeItem("token");
        localStorage.removeItem("user_id");
        localStorage.removeItem("username");
        navigate("/login");
      };

      ws.onclose = () => {
        console.log("WebSocket connection closed.");
      };

      return () => {
        ws.close();
      };
    }
  }, [navigate]);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      // Call the login API to get the UUID
      const response = await axiosInstance.post("/login", {
        username,
        password,
      });

      const { uuid } = response.data;

      // Establish WebSocket connection using the UUID
      const ws = new WebSocket(`ws://localhost:8080/ws/${uuid}`);

      ws.onopen = () => {
        console.log("WebSocket connection opened for UUID:", uuid);
      };

      ws.onmessage = async (event) => {
        const data = JSON.parse(event.data);

        if (data.success) {
          // Save the token
          localStorage.setItem("token", data.token);

          // Reconnect to the authenticated WebSocket
          const authWs = new WebSocket("ws://localhost:8080/ws/auth");

          authWs.onopen = () => {
            console.log("Authenticated WebSocket connection opened.");
            authWs.send(JSON.stringify({ type: "auth", token: data.token }));
          };

          authWs.onmessage = (authEvent) => {
            const authData = JSON.parse(authEvent.data);

            if (authData.type === "auth" && authData.success) {
              console.log("Authenticated successfully.");
              authWs.send(JSON.stringify({ type: "getSelf" }));
            } else if (authData.type === "get_self_response") {
              console.log("Received user data:", authData.data);
              const { id, username } = authData.data.user;

              // Save user ID and username
              localStorage.setItem("user_id", id.toString());
              localStorage.setItem("username", username);

              navigate("/chat");
            } else if (authData.type === "error") {
              setError(authData.message || "Failed to authenticate.");
            }
          };

          authWs.onerror = () => {
            setError("Error during WebSocket authentication.");
          };

          authWs.onclose = () => {
            console.log("Authenticated WebSocket connection closed.");
          };
        } else if (!data.success) {
          setError(data.message || "Login failed.");
        } else {
          console.log("Unexpected response from WebSocket:", data);
          setError("Unexpected response from WebSocket.");
        }

        ws.close();
        setLoading(false);
      };

      ws.onerror = () => {
        setError("WebSocket error occurred.");
        setLoading(false);
      };

      ws.onclose = () => {
        console.log("WebSocket connection closed.");
      };
    } catch (err: any) {
      setError(err.response?.data?.error || "Something went wrong.");
      setLoading(false);
    }
  };

  return (
    <div className="h-screen w-screen bg-gradient-to-r from-blue-500 to-purple-600 flex items-center justify-center">
      <div className="w-full max-w-md bg-white p-6 rounded-lg shadow-lg">
        <h2 className="text-3xl font-bold text-gray-800 text-center mb-6">
          Connexion
        </h2>
        {error && (
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative mb-4">
            <strong className="font-bold">Erreur: </strong>
            <span className="block sm:inline">{error}</span>
          </div>
        )}
        <form onSubmit={handleLogin}>
          <div className="mb-4">
            <label htmlFor="username" className="block text-gray-700">
              Nom d'utilisateur
            </label>
            <input
              type="text"
              id="username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="w-full p-2 border border-gray-300 rounded-md"
              placeholder="Votre nom d'utilisateur"
              disabled={loading}
            />
          </div>
          <div className="mb-4">
            <label htmlFor="password" className="block text-gray-700">
              Mot de passe
            </label>
            <input
              type="password"
              id="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full p-2 border border-gray-300 rounded-md"
              placeholder="Votre mot de passe"
              disabled={loading}
            />
          </div>
          <button
            type="submit"
            className="w-full bg-blue-500 text-white py-2 rounded-md hover:bg-blue-600"
            disabled={loading}
          >
            {loading ? "Connexion..." : "Se connecter"}
          </button>
        </form>
        <button
          onClick={() => navigate("/")}
          className="mt-4 w-full bg-gray-300 text-gray-700 py-2 rounded-md hover:bg-gray-400"
          disabled={loading}
        >
          Retour à l'accueil
        </button>
      </div>
    </div>
  );
};

export default Login;
