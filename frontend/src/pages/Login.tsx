import React, { useState, useEffect } from "react";
import axiosInstance from "../api/axiosInstance";
import { useNavigate } from "react-router-dom";
import { User } from "../types";

const Login: React.FC = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const navigate = useNavigate();

  useEffect(() => {
    const validateToken = async () => {
      const token = localStorage.getItem("token");

      if (token) {
        try {
          // Validate token by fetching user data
          await axiosInstance.get<User>("/users/self", {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          });

          // If valid, redirect to chat
          navigate("/chat");
        } catch {
          // If token is invalid, clear it
          localStorage.removeItem("token");
          localStorage.removeItem("user_id");
          localStorage.removeItem("username");
        }
      }
    };

    validateToken();
  }, [navigate]);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      // Fetch and save token in localStorage
      const loginResponse = await axiosInstance.post("/login", {
        username,
        password,
      });
      localStorage.setItem("token", loginResponse.data.token);

      axiosInstance.defaults.headers.common[
        "Authorization"
      ] = `Bearer ${loginResponse.data.token}`;

      // Fetch user data and save in localStorage
      const userReponse = await axiosInstance.get<User>("/users/self");

      localStorage.setItem("user_id", userReponse.data.id.toString());
      localStorage.setItem("username", userReponse.data.username);

      // Redirect to the chat page
      navigate("/chat");
    } catch (err: any) {
      setError(err.response?.data?.error || "Something went wrong");
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
            />
          </div>
          <button
            type="submit"
            className="w-full bg-blue-500 text-white py-2 rounded-md hover:bg-blue-600"
          >
            Se connecter
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

export default Login;
