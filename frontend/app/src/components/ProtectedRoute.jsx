import React, { useEffect, useState } from "react";
import { Navigate } from "react-router-dom";
import Cookies from "js-cookie";

const ProtectedRoute = ({ children }) => {
  const [isAuthenticated, setIsAuthenticated] = useState(null); // null = loading

  useEffect(() => {
    const checkToken = async () => {
      try {
        const response = await fetch("/api/auth/refresh", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                refresh_token: Cookies.get("refresh_token")
            }),
        });

        if (response.ok) {
            const data = await response.json();
            sessionStorage.setItem("access_token", data.access_token) // more security

            Cookies.set("refresh_token", data.refresh_token); // default 
            Cookies.set("token_type", data.token_type); // default

            setIsAuthenticated(true);
        } else {
            setIsAuthenticated(false);
        }
      }     catch (error) {
         setIsAuthenticated(false);
      }
    };

    checkToken();
  }, []);

  if (isAuthenticated === null) return <div>Loading...</div>; // можно заменить на спиннер

  return isAuthenticated ? children : <Navigate to="/login" replace />;
};

export default ProtectedRoute;
