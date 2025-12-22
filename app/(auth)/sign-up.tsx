import { View, Text, TextInput } from "react-native";
import { Link } from "expo-router";
import AppButtonIcon from "@/components/ui/AppButtonIcon";
import AppButtonAlt from "@/components/ui/AppButtonAlt";
import AppText from "@/components/ui/AppText";
import Heading from "@/components/ui/Appheading";
import * as Google from "expo-auth-session/providers/google";
import { useEffect } from "react";
import { GoogleIcon } from "@/components/ui/IconSymbol";
import { useTheme } from "@react-navigation/native";

export default function SignUpScreen() {
  const [request, response, promptAsync] = Google.useIdTokenAuthRequest({
    clientId: "<YOUR_IOS_CLIENT_ID>",
    androidClientId: "<YOUR_ANDROID_CLIENT_ID>",
  });

  useEffect(() => {
    if (response?.type === "success") {
      const { id_token } = response.params;
      // Send id_token to backend to create/sign-in user
    }
  }, [response]);

  const { colors } = useTheme();

  return (
    <View className="flex-1 px-6 gap-5 bg-neutral-900">
      <Heading>Create an Account</Heading>

      <View className="pt-5">
        <AppButtonIcon
          onPress={() => console.log("Sign in with Google")}
          icon={<GoogleIcon size={24} />}
        >
          Continue with Google
        </AppButtonIcon>
      </View>

      <View className="flex-row items-center my-3">
        <View className="flex-1 h-px bg-neutral-400" />
        <Text className="mx-3 text-white">Or</Text>
        <View className="flex-1 h-px bg-neutral-400" />
      </View>

      <TextInput
        className="border-b-gray-600 border rounded-xl text-gray-200"
        placeholder="Username"
        placeholderTextColor="#C7C9CE"
        style={{ height: 55, paddingVertical: 10, paddingHorizontal: 15 }}
      />
      <TextInput
        className="border-b-gray-600 border rounded-xl text-gray-200"
        placeholder="Email"
        placeholderTextColor="#C7C9CE"
        style={{ height: 55, paddingVertical: 10, paddingHorizontal: 15 }}
      />
      <TextInput
        className="border-b-gray-600 border rounded-xl text-gray-200"
        placeholder="Password"
        placeholderTextColor="#C7C9CE"
        secureTextEntry={true}
        style={{ height: 55, paddingVertical: 10, paddingHorizontal: 15 }}
      />
      <AppButtonAlt>Sign up</AppButtonAlt>
    </View>
  );
}
