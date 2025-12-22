import { API_URL } from "../../config";
import AsyncStorage from "@react-native-async-storage/async-storage";

export async function loginUser(email, password) {
  try {
    const response = await fetch(`${API_URL}/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    });

    if (!response.ok) throw new Error(`Status: ${response.status}`);

    const data = await response.json();
    console.log(data.token);
    return data.token;
  } catch (err) {
    console.error("Login failed:", err.message);
  }
}

export async function createUser(username, email, password) {
  try {
    const response = await fetch(`${API_URL}/signup`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, email, password }),
    });

    if (!response.ok) throw new Error(`Status: ${response.status}`);

    const data = await response.json();
    await AsyncStorage.setItem("userToken", data.token);
    return data.token;
  } catch (err) {
    console.error("Failed to create account:", err.message);
    return null;
  }
}
