import {fileURLToPath} from "node:url"
import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'
import Components from 'unplugin-vue-components/vite'
import {AntDesignVueResolver} from 'unplugin-vue-components/resolvers'


export default defineConfig({
    base: '/static',
    server: {
        port: 31945,
    },
    plugins: [
        vue(),
        Components({
            dirs: ['src/components'],
            dts: 'src/types/components.d.ts',
            resolvers: [
                AntDesignVueResolver({
                    importStyle: false,
                }),
            ]
        }),
    ],
    resolve: {
        alias: {
            '@': fileURLToPath(new URL('./src', import.meta.url))
        }
    }
})
