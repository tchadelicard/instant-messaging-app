export interface User {
  id: number;
  username: string;
}

export interface Message {
  id: number;
  sender_id: number;
  receiver_id: number;
  content: string;
}
