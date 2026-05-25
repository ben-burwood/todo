import { defineConfig, minimal2023Preset } from '@vite-pwa/assets-generator/config'

export default defineConfig({
  headLinkOptions: { preset: '2023' },
  preset: {
    ...minimal2023Preset,
    transparent: {
      sizes: minimal2023Preset.transparent.sizes,
    },
  },
  images: ['public/pwa-icon.svg'],
})
