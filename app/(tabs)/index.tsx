import { View, Text } from "react-native";
import { Link } from "expo-router";
import { useEffect, useState } from "react";
import AsyncStorage from "@react-native-async-storage/async-storage";
import { router } from "expo-router";

export default function Home() {
  const [renderHome, setRenderHome] = useState(true);

  useEffect(() => {
    async function verifyUserAuthentication() {
      const token = await AsyncStorage.getItem("token");

      if (!token) {
        router.replace("/landing");
        return;
      }

      setRenderHome(false);
    }

    verifyUserAuthentication();
  }, []);

  if (renderHome) return null;

  return (
    <View className="flex-1 items-center justify-center bg-white">
      <Text className="text-blue-500 text-xl font-bold">
        NativeWind is working!
      </Text>
      <Link href="/(auth)/sign-up">
        <Text style={{ marginTop: 20, color: "blue" }}>Go to Sign Up</Text>
      </Link>
      <Link href="/(auth)/landing">
        <Text style={{ marginTop: 20, color: "blue" }}>Go to landing</Text>
      </Link>
    </View>
  );
}
