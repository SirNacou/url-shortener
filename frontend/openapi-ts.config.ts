import { defineConfig } from "@hey-api/openapi-ts"

export default defineConfig({
  input: './src/api/openapi.json',
  output: './src/api/gen',
  plugins: [
    '@hey-api/typescript',
    { name: '@hey-api/sdk', validator: true, transformer: true },
    { name: '@hey-api/client-ky', runtimeConfigPath: './src/hey-api.ts' },
    {
      name: 'zod',
      dates: {
        offset: true
      },
      types: {
        input: true
      },
      responses: {
        types: {
          output: true
        }
      }
    },
    { name: '@tanstack/react-query' }
  ]
})