import { Pressable, Text, View } from "react-native";

type AppButtonProps = {
  children: React.ReactNode;
  onPress?: () => void;
  className?: string;
  icon?: React.ReactNode;
};

export default function AppButtonAlt({
  children,
  onPress,
  className = "",
  icon,
}: AppButtonProps) {
  return (
    <Pressable
      onPress={onPress}
      className={`flex-row border-2 border-gray-600 py-4 px-15 gap-5 rounded-4xl  ${className}`}
    >
      {icon && <View className="">{icon}</View>}
      <Text className="text-white font-bold">{children}</Text>
    </Pressable>
  );
}
