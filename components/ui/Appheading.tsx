import { Text, TextProps } from "react-native";

type HeadingProps = TextProps & {
  className?: string;
};

export default function Heading({ className = "", ...props }: HeadingProps) {
  return (
    <Text
      className={`text-heading font-bold text-white text-4xl ${className}`}
      {...props}
    />
  );
}
