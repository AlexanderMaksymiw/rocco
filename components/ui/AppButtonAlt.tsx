import { Pressable, Text } from "react-native";

type AppButtonProps = {
  children: React.ReactNode;
  onPress?: () => void;
  className?: string;
};

export default function AppButtonAlt({
  children,
  onPress,
  className = "",
}: AppButtonProps) {
  return (
    <Pressable
      onPress={onPress}
      className={`bg-gray-600 py-3 px-36 rounded-4xl items-center justify-center ${className}`}
    >
      <Text className="text-white text-center font-bold">{children}</Text>
    </Pressable>
  );
}
