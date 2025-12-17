/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./app/**/*.{js,ts,tsx}",
    "./components/**/*.{js,ts,tsx}",
  ],
  presets: [require("nativewind/preset")],
  theme: {
    extend: {
      colors: {
        primary: "#0a7ea4",
        backgroundLight: "#fff",
        backgroundDark: "#151718",
        textLight: "#11181C",
        textDark: "#ECEDEE",
      },
      fontSize: {
        heading: "28px",
        body: "16px",
      },
    },
  },
  plugins: [],
};
