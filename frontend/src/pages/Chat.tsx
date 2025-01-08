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
  //const [authenticated, setAuthenticated] = useState<boolean>(false);

  const navigate = useNavigate();
  const currentUserId = parseInt(localStorage.getItem("user_id") || "0", 10);
  const token = localStorage.getItem("token");
  const socketRef = useRef<WebSocket | null>(null);
  const chatContainerRef = useRef<HTMLDivElement | null>(null); // Add ref for chat container

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
              //setAuthenticated(true);
              fetchUsers(); // Fetch users after authentication
              setLoading(false);
            } else {
              console.error("WebSocket authentication failed:", data.message);
              handleLogout();
            }
            break;
          case "get_users_response":
            setUsers(data.data.users || []);
            restoreSelectedUser(data.data.users);
            break;
          case "get_messages_response":
            setMessages((prev) => [...prev, ...data.data.messages]);
            break;
          case "send_message_response":
            setMessages((prev) => [...prev, data.data.message]);
            break;
          case "error":
            console.error("Error from WebSocket:", data.message);
            break;

          default:
            console.warn("Unknown message type:", data.type);
        }
      };

      ws.onclose = () => {
        console.log("WebSocket disconnected.");
        //setAuthenticated(false);
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

  useEffect(() => {
    if (chatContainerRef.current) {
      chatContainerRef.current.scrollTop =
        chatContainerRef.current.scrollHeight;
    }
  }, [messages]);

  const fetchUsers = () => {
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
      socketRef.current.send(JSON.stringify(newMessage));
      setMessage("");
    }
  };

  const handleUserSelection = (user: User) => {
    setSelectedUser(user);
    setMessages([]);
    localStorage.setItem("selected_user_id", user.id.toString());
    fetchMessages(user.id);
  };

  const handleLogout = () => {
    console.log("Logging out...");
    localStorage.clear(); // Clear user session data
    socketRef.current?.close(); // Close WebSocket connection
    navigate("/login"); // Redirect to login page
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
            .filter((user) => user.id !== currentUserId)
            .filter((user) =>
              user.username.toLowerCase().includes(search.toLowerCase())
            )
            .map((user) => (
              <div
                key={user.id}
                className={`flex items-center p-3 border-b border-gray-200 hover:bg-gray-200 cursor-pointer transition duration-150 ease-in-out ${
                  selectedUser?.id === user.id ? "bg-indigo-200" : ""
                }`}
                onClick={() => handleUserSelection(user)}
              >
                <img
                  src={`https://picsum.photos/50/50?random=${user.id}`}
                  alt={user.username}
                  className="w-10 h-10 rounded-full mr-3"
                />
                <span className="font-medium text-gray-700">
                  {user.username}
                </span>
              </div>
            ))}
        </div>
        <div className="p-4 border-t border-gray-300">
          <button
            onClick={handleLogout}
            className="w-full py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition duration-150"
          >
            Logout
          </button>
        </div>
      </div>

      {/* Right column: Chat conversation */}
      <div className="flex-1 flex flex-col bg-gray-50 bg-opacity-90">
        {selectedUser ? (
          <>
            <div className="bg-gray-100 bg-opacity-90 border-b border-gray-300 p-4 flex items-center">
              <img
                src={`https://picsum.photos/50/50?random=${selectedUser.id}`}
                alt={selectedUser.username}
                className="w-10 h-10 rounded-full mr-3"
              />
              <span className="font-medium text-gray-700">
                {selectedUser.username}
              </span>
            </div>
            <div
              ref={chatContainerRef}
              className="flex-1 overflow-y-auto p-4 space-y-4"
            >
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
