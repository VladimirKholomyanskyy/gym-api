import { createSystem, defaultConfig, defineConfig } from "@chakra-ui/react";

const config = defineConfig({
  theme: {
    tokens: {
      colors: {
        magenta: {
          50: { value: "#fce4f1" }, // lightest
          100: { value: "#f8b1d4" },
          200: { value: "#f57dbe" },
          300: { value: "#f24cb5" }, // you can adjust these shades
          400: { value: "#e1008d" }, // main magenta color (you can adjust this)
          500: { value: "#c10075" }, // darker magenta
          600: { value: "#9c005f" },
          700: { value: "#75004a" },
          800: { value: "#4f0034" },
          900: { value: "#2a001f" }, // darkest
        },
        neon: {
          100: { value: "#ccfff9" },
          200: { value: "#99fff3" },
          300: { value: "#66ffee" },
          400: { value: "#00fff0" }, // Main neon cyan
          500: { value: "#00ccc0" },
          600: { value: "#009990" },
          700: { value: "#006660" },
          800: { value: "#003330" },
          900: { value: "#001a18" },
        },
      },
      fonts: {
        heading: { value: "Audiowide" },
        body: { value: "Audiowide" },
      },
    },
    semanticTokens: {
      colors: {
        bg: {
          primary: { value: "#1A202C" },
          subtle: { value: "#121212" },
        },
        fg: {
          light: { value: "#F7FAFC" },
          dark: { value: "#2D3748" },
        },
      },
    },
  },
  globalCss: {
    body: {
      backgroundColor: "bg.subtle",
      color: "fg.light",
    },
  },
});

export const system = createSystem(defaultConfig, config);
