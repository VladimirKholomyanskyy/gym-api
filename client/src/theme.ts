import { createSystem, defaultConfig, defineConfig } from "@chakra-ui/react";

const config = defineConfig({
  theme: {
    tokens: {
      fonts: {
        heading: { value: "Audiowide" },
        body: { value: "Audiowide" },
      },
    },
  },
});

export const system = createSystem(defaultConfig, config);
