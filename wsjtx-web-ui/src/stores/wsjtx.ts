import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

export const useWsjtxStore = defineStore('wsjtx', () => {
    const isConnected = ref(false)
    const messages = ref<any[]>([])
    const status = ref<any>(null)

    let ws: WebSocket | null = null

    const audioQueue = ref<Float32Array[]>([])
    const audioContext = new (window.AudioContext || (window as any).webkitAudioContext)()

    const audioDevices = ref<{ id: string, name: string }[]>([])

    function connect() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
        ws = new WebSocket(`${protocol}//${window.location.host}/ws`)
        ws.binaryType = 'arraybuffer'

        ws.onopen = () => {
            isConnected.value = true
            console.log('WebSocket connected')
            // Fetch devices on connect
            fetchAudioDevices()
        }

        ws.onmessage = (event) => {
            if (event.data instanceof ArrayBuffer) {
                // Audio data (Int16 PCM)
                const int16Data = new Int16Array(event.data)
                const float32Data = new Float32Array(int16Data.length)
                for (let i = 0; i < int16Data.length; i++) {
                    const sample = int16Data[i]
                    if (sample !== undefined) {
                        float32Data[i] = sample / 32768.0
                    }
                }
                audioQueue.value.push(float32Data)
                if (audioQueue.value.length > 50) {
                    audioQueue.value.shift() // Keep queue size manageable
                }
            } else {
                try {
                    const data = JSON.parse(event.data)
                    handleMessage(data)
                } catch (e) {
                    console.error('Error parsing message:', e)
                }
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
        } else if (msg.type === "audioDevices") {
            audioDevices.value = msg.data
        }
    }

    function sendMessage(type: string, data: any) {
        if (ws && ws.readyState === WebSocket.OPEN) {
            ws.send(JSON.stringify({ type, data }))
        } else {
            console.error('WebSocket not connected')
        }
    }

    function fetchAudioDevices() {
        sendMessage('getAudioDevices', {})
    }

    function setAudioDevice(id: string) {
        sendMessage('setAudioDevice', { id })
    }

    return { isConnected, messages, status, connect, sendMessage, audioQueue, audioContext, audioDevices, fetchAudioDevices, setAudioDevice }
})
