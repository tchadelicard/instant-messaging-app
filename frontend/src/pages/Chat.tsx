import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { FaSearch, FaPaperPlane } from "react-icons/fa";
import axiosInstance from "../api/axiosInstance";
import { User, Message } from "../types";

const Chat: React.FC = () => {
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [message, setMessage] = useState<string>("");
  const [messages, setMessages] = useState<Message[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [search, setSearch] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(true);

  const navigate = useNavigate();
  const currentUserId = parseInt(localStorage.getItem("user_id") || "0", 10);

  useEffect(() => {
    const validateToken = async () => {
      const token = localStorage.getItem("token");

      if (!token) {
        navigate("/login");
        return;
      }

      try {
        axiosInstance.defaults.headers.common[
          "Authorization"
        ] = `Bearer ${token}`;
        await axiosInstance.get("/users/self");
        fetchUsers();
        setLoading(false);
      } catch {
        localStorage.removeItem("token");
        localStorage.removeItem("user_id");
        localStorage.removeItem("username");
        navigate("/login");
      }
    };

    validateToken();
  }, [navigate]);

  const fetchUsers = async () => {
    try {
      const response = await axiosInstance.get<User[]>("/users");
      setUsers(response.data || []);
    } catch (err) {
      console.error("Failed to fetch users:", err);
    }
  };

  const fetchMessages = async (userId: number) => {
    try {
      const response = await axiosInstance.get<Message[]>(
        `/messages/${userId}`
      );
      setMessages(response.data || []);
    } catch (err) {
      console.error("Failed to fetch messages:", err);
    }
  };

  const handleSendMessage = async () => {
    if (message.trim() && selectedUser) {
      try {
        const newMessage = {
          content: message,
          receiver_id: selectedUser.id,
        };
        const response = await axiosInstance.post<Message>(
          `/messages/${selectedUser.id}`,
          newMessage
        );
        setMessages((prev) => [...prev, response.data]);
        setMessage("");
      } catch (err) {
        console.error("Failed to send message:", err);
      }
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <p>Loading...</p>
      </div>
    );
  }

  return (
    <div className="flex h-screen bg-gradient-to-r from-green-400 via-blue-500 to-purple-500">
      {/* Left column: User list */}
      <div className="w-1/4 bg-gray-100 bg-opacity-90 border-r border-gray-300 flex flex-col min-w-1/4 max-w-xs">
        <div className="p-4">
          <div className="relative">
            <input
              type="text"
              placeholder="Search users"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="w-full pl-10 pr-4 py-2 border rounded-full focus:outline-none focus:ring-2 focus:ring-indigo-300"
              aria-label="Search users"
            />
            <FaSearch className="absolute left-3 top-3 text-gray-400" />
          </div>
        </div>
        <div className="flex-1 overflow-y-auto">
          {users
            .filter((user) => user.id !== currentUserId) // Exclude current user
            .filter((user) =>
              user.username.toLowerCase().includes(search.toLowerCase())
            ) // Apply search filter
            .map((user) => (
              <div
                key={user.id}
                className={`flex items-center p-3 border-b border-gray-200 hover:bg-gray-200 cursor-pointer transition duration-150 ease-in-out ${
                  selectedUser?.id === user.id ? "bg-indigo-200" : ""
                }`}
                onClick={() => {
                  setSelectedUser(user);
                  fetchMessages(user.id); // Fetch messages here
                }}
              >
                <img
                  src={`https://picsum.photos/50/50?random=${user.id}`} // Replace with user.avatar if available
                  alt={user.username}
                  className="w-10 h-10 rounded-full mr-3"
                />
                <span className="font-medium text-gray-700">
                  {user.username}
                </span>
              </div>
            ))}
        </div>
      </div>

      {/* Right column: Chat conversation */}
      <div className="flex-1 flex flex-col bg-gray-50 bg-opacity-90">
        {selectedUser ? (
          <>
            <div className="bg-gray-100 bg-opacity-90 border-b border-gray-300 p-4 flex items-center">
              <img
                src={`https://picsum.photos/50/50?random=${selectedUser.id}`} // Replace with selectedUser.avatar if available
                alt={selectedUser.username}
                className="w-10 h-10 rounded-full mr-3"
              />
              <span className="font-medium text-gray-700">
                {selectedUser.username}
              </span>
            </div>
            <div className="flex-1 overflow-y-auto p-4 space-y-4">
              {messages.map((msg) => (
                <div
                  key={msg.id}
                  className={`flex ${
                    msg.sender_id === currentUserId
                      ? "justify-end"
                      : "justify-start"
                  }`}
                >
                  <div
                    className={`max-w-xs px-4 py-2 rounded-lg ${
                      msg.sender_id === currentUserId
                        ? "bg-indigo-500 text-white"
                        : "bg-gray-200"
                    }`}
                  >
                    {msg.content}
                  </div>
                </div>
              ))}
            </div>
            <div className="bg-gray-100 bg-opacity-90 border-t border-gray-300 p-4">
              <div className="flex items-center">
                <input
                  type="text"
                  value={message}
                  onChange={(e) => setMessage(e.target.value)}
                  placeholder="Type a message"
                  className="flex-1 border rounded-full px-4 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-300"
                  onKeyPress={(e) => e.key === "Enter" && handleSendMessage()}
                />
                <button
                  onClick={handleSendMessage}
                  className="ml-2 bg-indigo-500 text-white rounded-full p-2 hover:bg-indigo-600 focus:outline-none focus:ring-2 focus:ring-indigo-300"
                >
                  <FaPaperPlane />
                </button>
              </div>
            </div>
          </>
        ) : (
          <div className="flex-1 flex items-center justify-center text-gray-500">
            Select a user to start chatting
          </div>
        )}
      </div>
    </div>
  );
};

export default Chat;
