import React, { useEffect } from "react";
import { useNavigate } from "react-router-dom";

const Home: React.FC = () => {
  const navigate = useNavigate();

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (token) {
      // Redirect to the chat page if the user is logged in
      navigate("/chat");
    }
  }, [navigate]);

  return (
    <div className="h-screen w-screen bg-gradient-to-r from-green-400 to-blue-500 flex items-center justify-center">
      <div className="w-full max-w-md bg-white p-6 rounded-lg shadow-lg">
        <h1 className="text-4xl font-bold text-gray-800 text-center mb-6">
          Bienvenue 👋
        </h1>
        <p className="text-gray-600 text-center mb-6">
          Connectez-vous pour discuter avec vos amis ou rejoignez un groupe.
        </p>
        <div className="flex justify-around">
          <button
            onClick={() => navigate("/login")}
            className="bg-blue-500 text-white py-2 px-4 rounded-md hover:bg-blue-600"
          >
            Connexion
          </button>
          <button
            onClick={() => navigate("/register")}
            className="bg-gray-300 text-gray-700 py-2 px-4 rounded-md hover:bg-gray-400"
          >
            Inscription
          </button>
        </div>
      </div>
    </div>
  );
};

export default Home;
