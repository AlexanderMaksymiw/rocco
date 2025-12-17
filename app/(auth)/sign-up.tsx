import { View, Text, TextInput } from "react-native";
import { Link } from "expo-router";
import AppButtonIcon from "@/components/ui/AppButtonIcon";
import AppButtonAlt from "@/components/ui/AppButtonAlt";
import AppText from "@/components/ui/AppText";
import Heading from "@/components/ui/Appheading";
import * as Google from "expo-auth-session/providers/google";
import { useEffect } from "react";
import { GoogleIcon } from "@/components/ui/IconSymbol";

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

  return (
    <View className="flex-1 px-6 gap-5">
      <Heading>Create an Account</Heading>

      <View className="pt-5">
        <AppButtonIcon
          onPress={() => console.log("Sign in with Google")}
          icon={<GoogleIcon size={24} />}
        >
          Continue with Google
        </AppButtonIcon>
      </View>

      <TextInput
        className="border-b-gray-600 border-2 rounded-xl text-gray-200"
        placeholder="Username"
        placeholderTextColor="#999"
        style={{ height: 55, paddingVertical: 10, paddingHorizontal: 15 }}
      />
      <TextInput
        className="border-b-gray-600 border-2 rounded-xl text-gray-200"
        placeholder="Email"
        placeholderTextColor="#999"
        style={{ height: 55, paddingVertical: 10, paddingHorizontal: 15 }}
      />
      <TextInput
        className="border-b-gray-600 border-2 rounded-xl text-gray-200"
        placeholder="Password"
        placeholderTextColor="#999"
        secureTextEntry={true}
        style={{ height: 55, paddingVertical: 10, paddingHorizontal: 15 }}
      />
      <AppButtonAlt>Sign up</AppButtonAlt>
    </View>
  );
}
