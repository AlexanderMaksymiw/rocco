import { View, Text } from "react-native";
import { Link } from "expo-router";
import AuthBackground from "../../components/ui/AuthBackground";
import AppButton from "../../components/ui/AppButton";

export default function LandingScreen() {
  return (
    <AuthBackground>
      <View
        style={{
          flex: 1,
          justifyContent: "center",
          alignItems: "center",
          backgroundColor: "rgba(0,0,0,0.2)",
          paddingHorizontal: 24,
        }}
      >
        <View>
          <Text className="text-white font-black text-6xl pt-10">ROCCO</Text>
        </View>
        <View className="flex-1 justify-end items-center pb-15 gap-5">
          <Link href="/(auth)/sign-up" asChild>
            <AppButton>Join for Free</AppButton>
          </Link>
          <Link href="/(auth)/sign-in">
            <Text className="font-bold text-white">Log In</Text>
          </Link>
        </View>
      </View>
    </AuthBackground>
  );
}
