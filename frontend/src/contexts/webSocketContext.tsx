import React, { useState, useEffect, useRef, createContext } from "react";
import { EventSourcePolyfill } from "event-source-polyfill";

export type SSEContextType = {
  data: any;
  eventSource: EventSourcePolyfill | null;
};

export const SSEContext = createContext<SSEContextType | undefined>(undefined);

interface SSEProviderProps {
  children: React.ReactNode;
}

export const SSEProvider: React.FC<SSEProviderProps> = ({ children }) => {
  const [data, setData] = useState({});
  const eventSource = useRef<EventSourcePolyfill | null>(null);
  const reconnectInterval = useRef<number | null>(null);

  const connectSSE = () => {
    const token = localStorage.getItem("token");
    if (!token) {
      console.error("No token found in localStorage");
      return;
    }

    if (eventSource.current) {
      eventSource.current.close();
    }

    eventSource.current = new EventSourcePolyfill(
      `http://${window.location.hostname}:7889/api/events/stream`,
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
        withCredentials: true,
      }
    );

    // Handle different event types
    const eventTypes = [
      "settings",
      "system",
      "profiles",
      "containers",
      "codecs",
      "encoders",
      "queue",
      "series",
      "movies",
      "history",
      "logs",
    ];

    eventTypes.forEach((eventType) => {
      eventSource.current?.addEventListener(
        eventType,
        (event: MessageEvent) => {
          const newData = JSON.parse(event.data);
          setData((prevData) => ({ ...prevData, [eventType]: newData }));
        }
      );
    });

    eventSource.current.onopen = () => {
      if (reconnectInterval.current) {
        clearInterval(reconnectInterval.current);
        reconnectInterval.current = null;
      }
    };

    eventSource.current.onerror = (error) => {
      console.error("SSE error:", error);
      eventSource.current?.close();
      if (!reconnectInterval.current) {
        reconnectInterval.current = window.setInterval(connectSSE, 5000);
      }
    };
  };

  useEffect(() => {
    connectSSE();

    return () => {
      if (eventSource.current) {
        eventSource.current.close();
      }
      if (reconnectInterval.current) {
        clearInterval(reconnectInterval.current);
      }
    };
  }, []);

  const value: SSEContextType = {
    data,
    eventSource: eventSource.current,
  };

  return <SSEContext.Provider value={value}>{children}</SSEContext.Provider>;
};
