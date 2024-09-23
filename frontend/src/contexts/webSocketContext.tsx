import React, { useState, useEffect, useRef, createContext } from "react";

export type WebSocketContextType = {
  data: any;
  ws: WebSocket | null;
};

export const WebSocketContext = createContext<WebSocketContextType | undefined>(
  undefined
);

interface WebSocketProviderProps {
  children: React.ReactNode;
}

export const WebSocketProvider: React.FC<WebSocketProviderProps> = ({
  children,
}) => {
  const [data, setData] = useState({});
  const ws = useRef<WebSocket | null>(null);
  const reconnectInterval = useRef<number | null>(null);

  const connectWebSocket = () => {
    const token = localStorage.getItem("token");
    if (!token) {
      console.error("No token found in localStorage");
      return;
    }

    ws.current = new WebSocket(
      `ws://${window.location.hostname}:7889/api/ws?token=${token}`
    );

    ws.current.onopen = () => {
      console.log("WebSocket connection opened");
      if (reconnectInterval.current) {
        clearInterval(reconnectInterval.current);
        reconnectInterval.current = null;
      }
    };

    ws.current.onmessage = (e) => {
      const newData = JSON.parse(e.data);
      setData((prevData) => ({ ...prevData, ...newData }));
    };

    ws.current.onerror = (error) => {
      console.error("WebSocket error:", error);
    };

    ws.current.onclose = (event) => {
      console.log("WebSocket connection closed", event);
      if (!reconnectInterval.current) {
        reconnectInterval.current = window.setInterval(connectWebSocket, 5000);
      }
    };
  };

  useEffect(() => {
    connectWebSocket();

    return () => {
      if (ws.current) {
        ws.current.close();
      }
      if (reconnectInterval.current) {
        clearInterval(reconnectInterval.current);
      }
    };
  }, []);

  const value: WebSocketContextType = {
    data,
    ws: ws.current,
  };

  return (
    <WebSocketContext.Provider value={value}>
      {children}
    </WebSocketContext.Provider>
  );
};