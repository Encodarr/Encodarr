import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.tsx";
import "./index.css";
import { SSEProvider } from "./contexts/webSocketContext.tsx";
import { ThemeProvider } from "./contexts/themeContext.tsx";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <SSEProvider>
      <ThemeProvider>
        <App />
      </ThemeProvider>
    </SSEProvider>
  </React.StrictMode>
);
