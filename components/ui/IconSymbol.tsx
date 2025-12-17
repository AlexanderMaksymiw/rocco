import MaterialIcons from "@expo/vector-icons/MaterialIcons";
import { Image, ImageStyle } from "react-native";

type IconProps = {
  size?: number;
  color?: string;
  style?: ImageStyle;
};

export const GarageIcon = ({ size, color }: IconProps) => (
  <MaterialIcons name="garage" size={size ?? 24} color={color ?? "black"} />
);

export const CarIcon = ({ size, color }: IconProps) => (
  <MaterialIcons name="drive-eta" size={size ?? 24} color={color ?? "black"} />
);

export const SettingsIcon = ({ size, color }: IconProps) => (
  <MaterialIcons name="settings" size={size ?? 24} color={color ?? "black"} />
);

export const MaintenanceIcon = ({ size, color }: IconProps) => (
  <MaterialIcons name="add-task" size={size ?? 24} color={color ?? "black"} />
);

export const ReminderIcon = ({ size, color }: IconProps) => (
  <MaterialIcons
    name="notifications"
    size={size ?? 24}
    color={color ?? "black"}
  />
);

export const GoogleIcon = ({ size = 24 }: { size?: number }) => (
  <Image
    source={require("../../assets/icons/google-icon.png")}
    style={{ width: size, height: size }}
    resizeMode="contain"
  />
);
