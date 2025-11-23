import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

export const useWsjtxStore = defineStore('wsjtx', () => {
    const isConnected = ref(false)
    const messages = ref<any[]>([])
    const status = ref<any>(null)

    let ws: WebSocket | null = null

    function connect() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
        ws = new WebSocket(`${protocol}//${window.location.host}/ws`)

        ws.onopen = () => {
            isConnected.value = true
            console.log('WebSocket connected')
        }

        ws.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data)
                handleMessage(data)
            } catch (e) {
                console.error('Error parsing message:', e)
            }
        }

        ws.onclose = () => {
            isConnected.value = false
            console.log('WebSocket disconnected')
            setTimeout(connect, 1000) // Reconnect
        }

        ws.onerror = (error) => {
            console.error('WebSocket error:', error)
        }
    }

    function handleMessage(msg: any) {
        // msg.type is string "2" for Decode, "1" for Status, etc.
        if (msg.type === "2") { // Decode
            messages.value.unshift(msg.data)
            if (messages.value.length > 100) {
                messages.value.pop()
            }
        } else if (msg.type === "1") { // Status
            status.value = msg.data
        }
    }

    function sendMessage(type: string, data: any) {
        if (ws && ws.readyState === WebSocket.OPEN) {
            ws.send(JSON.stringify({ type, data }))
        } else {
            console.error('WebSocket not connected')
        }
    }

    return { isConnected, messages, status, connect, sendMessage }
})
