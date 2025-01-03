import axios from "axios";

const BASE_URL = "http://localhost:8080/api"; // Replace with your backend URL

const axiosInstance = axios.create({
  baseURL: BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

// Add a response interceptor
axiosInstance.interceptors.response.use(
  (response) => {
    // If the response is successful, just return it
    return response;
  },
  (error) => {
    // Check for 401 Unauthorized
    if (error.response && error.response.status === 401) {
      // Remove user data from localStorage
      localStorage.removeItem("token");
      localStorage.removeItem("user_id");

      // Get the current route
      const currentPath = window.location.pathname;

      // Only redirect if the user is not on the login or register page
      if (currentPath !== "/login" && currentPath !== "/register") {
        window.location.href = "/login";
      }
    }

    // Return the error so it can be handled locally as well
    return Promise.reject(error);
  }
);

export default axiosInstance;
