import { ImageBackground, View } from 'react-native';


const bgImage = require('../../assets/images/rocco-2.jpg');

export default function AuthBackground({ children }: { children?: React.ReactNode }) {
  return (
    <ImageBackground
      source={bgImage}
      style={{ flex: 1 }}
      resizeMode="cover"
    >
      <View style={{ flex: 1 }}>{children}</View>
    </ImageBackground>
  );
}
