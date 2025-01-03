import React, { useState } from "react";

const ChatBox: React.FC = () => {
  const [messages, setMessages] = useState<
    { content: string; sender: string }[]
  >([]);
  const [newMessage, setNewMessage] = useState<string>("");

  const sendMessage = () => {
    if (newMessage.trim() !== "") {
      setMessages([...messages, { content: newMessage, sender: "Moi" }]);
      setNewMessage("");
    }
  };

  return (
    <div className="flex flex-col h-full bg-gray-100 rounded-lg shadow-inner p-4">
      {/* Messages */}
      <div className="flex-1 overflow-y-auto mb-4 space-y-3">
        {messages.map((msg, index) => (
          <div
            key={index}
            className={`flex ${
              msg.sender === "Moi" ? "justify-end" : "justify-start"
            }`}
          >
            <div
              className={`p-3 rounded-lg text-white ${
                msg.sender === "Moi" ? "bg-blue-500" : "bg-gray-400"
              }`}
            >
              {msg.content}
            </div>
          </div>
        ))}
      </div>

      {/* Input */}
      <div className="flex items-center space-x-3">
        <input
          type="text"
          value={newMessage}
          onChange={(e) => setNewMessage(e.target.value)}
          className="flex-1 p-3 border border-gray-300 rounded-lg"
          placeholder="Écrivez un message..."
        />
        <button
          onClick={sendMessage}
          className="bg-blue-500 text-white px-4 py-2 rounded-lg hover:bg-blue-600 transition"
        >
          Envoyer
        </button>
      </div>
    </div>
  );
};

export default ChatBox;
