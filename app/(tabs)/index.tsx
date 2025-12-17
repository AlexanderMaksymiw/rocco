import { View, Text } from "react-native";
import { Link } from 'expo-router';

export default function Home() {
  return (
    <View className="flex-1 items-center justify-center bg-white">
      <Text className="text-blue-500 text-xl font-bold">
        NativeWind is working!
      </Text>
            <Link href="/(auth)/sign-up">
        <Text style={{ marginTop: 20, color: 'blue' }}>Go to Sign Up</Text>
      </Link>
            <Link href="/(auth)/landing">
        <Text style={{ marginTop: 20, color: 'blue' }}>Go to landing</Text>
      </Link>
    </View>
  );
}
