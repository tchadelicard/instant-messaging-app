import React, { useState, useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import { FaSearch, FaPaperPlane } from "react-icons/fa";
import { User, Message } from "../types";

const Chat: React.FC = () => {
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [message, setMessage] = useState<string>("");
  const [messages, setMessages] = useState<Message[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [search, setSearch] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(true);
  const [authenticated, setAuthenticated] = useState<boolean>(false);

  const navigate = useNavigate();
  const currentUserId = parseInt(localStorage.getItem("user_id") || "0", 10);
  const token = localStorage.getItem("token");
  const socketRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    if (!token) {
      console.error("No token found. Redirecting to login.");
      navigate("/login");
      return;
    }

    const connectWebSocket = () => {
      const ws = new WebSocket("ws://localhost:8080/ws/auth");
      socketRef.current = ws;

      ws.onopen = () => {
        console.log("WebSocket connected. Sending token...");
        ws.send(JSON.stringify({ type: "auth", token: token }));
      };

      ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        console.log("WebSocket message received:", data);

        switch (data.type) {
          case "auth":
            if (data.success) {
              console.log("WebSocket authenticated successfully.");
              setAuthenticated(true);
              fetchUsers(); // Fetch users after authentication
              setLoading(false);
            } else {
              console.error("WebSocket authentication failed:", data.message);
              localStorage.clear();
              navigate("/login");
            }
            break;

          case "get_users_response":
            console.log("Users received:", data.data.users);
            setUsers(data.data.users || []);
            restoreSelectedUser(data.data.users); // Restore the previously selected user
            break;

          case "get_messages_response":
            console.log("Message received:", data.data.messages);
            setMessages((prev) => [...prev, ...data.data.messages]);
            break;

          case "error":
            console.error("Error from WebSocket:", data.message);
            break;

          default:
            console.log(data);
            console.warn("Unknown message type:", data.type);
        }
      };

      ws.onclose = () => {
        console.log("WebSocket disconnected.");
        setAuthenticated(false);
      };

      ws.onerror = (event) => {
        console.error("WebSocket encountered an error:", event);
      };
    };

    connectWebSocket();

    return () => {
      console.log("Cleaning up WebSocket connection...");
      socketRef.current?.close();
    };
  }, [token, navigate]);

  const fetchUsers = () => {
    console.log("Fetching users via WebSocket...");
    socketRef.current?.send(
      JSON.stringify({
        type: "getUsers",
      })
    );
  };

  const restoreSelectedUser = (users: User[]) => {
    const storedUserId = localStorage.getItem("selected_user_id");
    if (storedUserId) {
      const restoredUser = users.find(
        (user) => user.id === parseInt(storedUserId, 10)
      );
      if (restoredUser) {
        setSelectedUser(restoredUser);
        fetchMessages(restoredUser.id);
      }
    }
  };

  const fetchMessages = (userId: number) => {
    console.log(`Fetching messages for user ${userId} via WebSocket...`);
    socketRef.current?.send(
      JSON.stringify({
        type: "getMessages",
        receiver_id: userId,
      })
    );
  };

  const handleSendMessage = () => {
    if (socketRef.current && message.trim() && selectedUser) {
      const newMessage = {
        type: "sendMessage",
        content: message,
        receiver_id: selectedUser.id,
      };
      console.log("Sending message:", newMessage);
      socketRef.current.send(JSON.stringify(newMessage));
      setMessage("");
    }
  };

  const handleUserSelection = (user: User) => {
    setSelectedUser(user);
    setMessages([]); // Clear messages when switching users
    localStorage.setItem("selected_user_id", user.id.toString()); // Save selected user in localStorage
    fetchMessages(user.id);
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
                onClick={() => handleUserSelection(user)}
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
