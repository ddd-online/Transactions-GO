import {defineStore} from 'pinia'
import {ref} from 'vue'
import {startMcpServer, stopMcpServer, getMcpStatus, type McpStatus} from "@/backend/api/mcp"
import NotificationUtil from "@/backend/notification"

export const useMcpStore = defineStore('mcp', () => {
    const isRunning = ref(false)
    const port = ref(31944)

    const refreshStatus = async () => {
        try {
            const status = await getMcpStatus()
            isRunning.value = status.running
            port.value = status.port
        } catch (error) {
            console.warn('Failed to get MCP status:', error)
        }
    }

    const start = async () => {
        try {
            const status = await startMcpServer()
            isRunning.value = status.running
            port.value = status.port
            NotificationUtil.success('MCP 服务已启动', `端口: ${status.port}`)
        } catch (error) {
            NotificationUtil.error('启动 MCP 服务失败', `${error}`)
            throw error
        }
    }

    const stop = async () => {
        try {
            const status = await stopMcpServer()
            isRunning.value = status.running
            NotificationUtil.success('MCP 服务已停止')
        } catch (error) {
            NotificationUtil.error('停止 MCP 服务失败', `${error}`)
            throw error
        }
    }

    return {
        isRunning,
        port,
        refreshStatus,
        start,
        stop,
    }
})