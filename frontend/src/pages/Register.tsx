import React, { useState, useEffect } from "react";
import axiosInstance from "../api/axiosInstance";
import { useNavigate } from "react-router-dom";

const Register: React.FC = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState("");
  const navigate = useNavigate();

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (token) {
      navigate("/chat");
    }
  }, [navigate]);

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setSuccess("");
    setLoading(true);

    try {
      // Send registration request
      const response = await axiosInstance.post("/register", {
        username,
        password,
      });

      const { uuid } = response.data;

      // Open WebSocket connection
      const ws = new WebSocket(`ws://localhost:8080/ws/${uuid}`);
      ws.onopen = () => console.log("WebSocket connection established");
      ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        console.log(data);
        setSuccess(data.message);
        setLoading(false);
        ws.close();

        // Redirect to login page after successful registration
        if (data.message === "Registration successful") {
          setTimeout(() => navigate("/login"), 2000);
        }
      };
      ws.onerror = () => {
        setError("WebSocket error occurred");
        setLoading(false);
      };
      ws.onclose = () => console.log("WebSocket connection closed");
    } catch (err: any) {
      setError(err.response?.data?.error || "Something went wrong");
      setLoading(false);
    }
  };

  return (
    <div className="h-screen w-screen bg-gradient-to-r from-purple-600 to-blue-500 flex items-center justify-center">
      <div className="w-full max-w-md bg-white p-6 rounded-lg shadow-lg">
        <h2 className="text-3xl font-bold text-gray-800 text-center mb-6">
          Inscription
        </h2>
        {error && (
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative mb-4">
            <strong className="font-bold">Erreur: </strong>
            <span className="block sm:inline">{error}</span>
          </div>
        )}
        {success && (
          <div className="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded relative mb-4">
            <strong className="font-bold">Succès: </strong>
            <span className="block sm:inline">{success}</span>
          </div>
        )}
        <form onSubmit={handleRegister}>
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
              placeholder="Choisissez un nom d'utilisateur"
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
              placeholder="Choisissez un mot de passe"
            />
          </div>
          <button
            type="submit"
            className={`w-full bg-blue-500 text-white py-2 rounded-md hover:bg-blue-600 ${
              loading ? "opacity-50 cursor-not-allowed" : ""
            }`}
            disabled={loading}
          >
            {loading ? "En cours..." : "S'inscrire"}
          </button>
        </form>
        <button
          onClick={() => navigate("/")}
          className="mt-4 w-full bg-gray-300 text-gray-700 py-2 rounded-md hover:bg-gray-400"
        >
          Retour à l'accueil
        </button>
      </div>
    </div>
  );
};

export default Register;
