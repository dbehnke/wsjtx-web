<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { useWsjtxStore } from '@/stores/wsjtx'
import { storeToRefs } from 'pinia'

const store = useWsjtxStore()
const { audioQueue } = storeToRefs(store)

const canvasRef = ref<HTMLCanvasElement | null>(null)
let animationId: number
let fft: AnalyserNode
let audioCtx: AudioContext

// Configuration
const FFT_SIZE = 2048
const SAMPLE_RATE = 12000 // Must match backend

onMounted(() => {
  if (!store.audioContext) return
  audioCtx = store.audioContext
  
  // Create Analyser
  fft = audioCtx.createAnalyser()
  fft.fftSize = FFT_SIZE
  fft.smoothingTimeConstant = 0.0 // No smoothing for waterfall
  
  // Start animation loop
  // animate() // using setInterval now
})

onUnmounted(() => {
  cancelAnimationFrame(animationId)
})

// Process audio queue
watch(audioQueue, (newQueue) => {
  if (newQueue.length === 0) return
  
  // Take the oldest buffer
  const buffer = newQueue[0]
  if (!buffer) return

  // In a real implementation, we would schedule this buffer to play/process smoothly.
  // For visualization, we can just feed it to the FFT.
  // However, Web Audio API works by connecting nodes.
  // We need to create a buffer source, fill it, and connect to analyser.
  
  // Create AudioBuffer
  const audioBuffer = audioCtx.createBuffer(1, buffer.length, SAMPLE_RATE)
  audioBuffer.copyToChannel(buffer as any, 0)
  
  const source = audioCtx.createBufferSource()
  source.buffer = audioBuffer
  source.connect(fft)
  source.start()
  
  // Remove from queue (this is a simplification, race conditions possible if watch triggers too fast)
  // Better to have a loop that consumes the queue.
}, { deep: true })

// Actually, watching the queue is not efficient for streaming.
// Let's use a processing loop.

function animate() {
  if (!canvasRef.value || !fft) {
    animationId = requestAnimationFrame(animate)
    return
  }

  const canvas = canvasRef.value
  const ctx = canvas.getContext('2d')
  if (!ctx) return

  // 1. Get FFT data
  const bufferLength = fft.frequencyBinCount
  const dataArray = new Uint8Array(bufferLength)
  fft.getByteFrequencyData(dataArray)

  // 2. Shift existing canvas content down
  // We want the newest line at the top (or bottom? WSJT-X usually scrolls down, newest at top)
  // Let's put newest at top.
  ctx.drawImage(canvas, 0, 1)

  // 3. Draw new line at top (y=0)
  const width = canvas.width
  const imageData = ctx.createImageData(width, 1)
  const data = imageData.data

  // Map FFT bins to pixels
  // We want to show 0 to MAX_FREQ (e.g., 4000 Hz)
  const MAX_FREQ = 4000
  // The FFT runs at the AudioContext's sample rate, not the source sample rate.
  const nyquist = audioCtx.sampleRate / 2
  const binsToDisplay = Math.floor((MAX_FREQ / nyquist) * bufferLength)

  for (let x = 0; x < width; x++) {
    // Map x (0..width) to binIndex (0..binsToDisplay)
    const binIndex = Math.floor((x / width) * binsToDisplay)
    const value = dataArray[binIndex] || 0

    // WSJT-X Style Color Map
    let r = 0, g = 0, b = 0
    
    // Noise floor subtraction (Gain control)
    // Subtract 30 from value to make background darker
    let v = value - 30
    if (v < 0) v = 0
    
    if (v < 50) {
        b = 20 + (v / 50) * 235
    } else if (v < 100) {
        const p = (v - 50) / 50
        r = p * 255
        b = 255 - (p * 255)
    } else if (v < 150) {
        const p = (v - 100) / 50
        r = 255
        g = p * 255
    } else {
        const p = Math.min(1, (v - 150) / 105)
        r = 255
        g = 255
        b = p * 255
    }

    const i = x * 4
    data[i] = r     // R
    data[i + 1] = g // G
    data[i + 2] = b // B
    data[i + 3] = 255 // Alpha
  }

  // Check for time interval (15s)
  const now = new Date()
  const seconds = now.getSeconds()
  
  if (seconds % 15 === 0) {
      // Draw green line
      for (let i = 0; i < data.length; i += 4) {
          data[i] = 0
          data[i + 1] = 255
          data[i + 2] = 0
      }
  }

  ctx.putImageData(imageData, 0, 0)
  
  if (seconds % 15 === 0) {
      // Draw timestamp
      ctx.fillStyle = 'white'
      ctx.font = '10px sans-serif'
      const timeStr = now.toISOString().substr(11, 5) // HH:MM
      ctx.fillText(timeStr, 5, 10)
  }
}

// We need a better way to feed the FFT than the watcher.
// Let's use a ScriptProcessor or AudioWorklet, or just schedule sources.
// For simplicity in this MVP, let's just create a loop that checks the queue.

setInterval(() => {
    if (store.audioQueue.length > 0) {
        const buffer = store.audioQueue.shift()
        if (buffer && audioCtx) {
             const audioBuffer = audioCtx.createBuffer(1, buffer.length, SAMPLE_RATE)
             // TS is very strict about ArrayBuffer vs SharedArrayBuffer here.
             // We know it's a standard Float32Array from the store.
             audioBuffer.copyToChannel(buffer as any, 0)
             
             const source = audioCtx.createBufferSource()
             source.buffer = audioBuffer
             source.connect(fft)
             // We don't connect to destination (speakers) to avoid feedback/noise
             source.start()
        }
    }
}, 100) // Check every 100ms. This is crude but might work for visualization.

// Draw loop (slower rate)
setInterval(() => {
    if (canvasRef.value && fft) {
        animate()
    }
}, 500) // Update every 500ms

</script>

<template>
  <div class="flex flex-col gap-1 flex-1">
      <div class="relative h-5 bg-gray-900 text-gray-400 text-xs rounded">
          <span class="absolute transform -translate-x-1/2" style="left: 0%; transform: translateX(0); left: 4px !important;">0</span>
          <span class="absolute transform -translate-x-1/2" style="left: 25%">1000</span>
          <span class="absolute transform -translate-x-1/2" style="left: 50%">2000</span>
          <span class="absolute transform -translate-x-1/2" style="left: 75%">3000</span>
          <span class="absolute transform -translate-x-1/2" style="left: 98%">4000</span>
      </div>
      <div class="w-full bg-black border border-gray-700 rounded overflow-hidden">
        <canvas ref="canvasRef" width="1024" height="400" class="w-full h-[400px] block"></canvas>
      </div>
  </div>
</template>
