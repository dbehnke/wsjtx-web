<script setup lang="ts">
import { useWsjtxStore } from '@/stores/wsjtx'
import { storeToRefs } from 'pinia'
import { ref, watch } from 'vue'

const store = useWsjtxStore()
const { status, isConnected, audioDevices } = storeToRefs(store)

const selectedAudioDevice = ref("")

watch(audioDevices, (devices) => {
  if (devices.length > 0 && !selectedAudioDevice.value) {
    // Select first device by default if nothing selected
    // Or maybe we should ask backend for current device?
    // For now, let's just default to empty (default device) or first one.
    // Actually, backend starts with default.
  }
})

function haltTx() {
  store.sendMessage('halt', { AutoTxOnly: false })
}

function enableTx() {
  // "Enable Tx" usually means setting "Auto Tx" on?
  // Or just ensuring TxEnabled is true?
  // The "Halt Tx" message stops transmission.
  // To "Enable Tx", we might need to send a configuration change or just rely on the user clicking "Enable Tx" in WSJT-X if we can't control it fully.
  // But wait, the protocol has "Halt Tx". Does it have "Enable Tx"?
  // The "Halt Tx" message has "AutoTxOnly" flag.
  // If we want to re-enable, maybe we just don't halt?
  // Actually, "Enable Tx" button in WSJT-X toggles the "Enable Tx" checkbox.
  // Placeholder for Enable Tx
  console.log("Enable Tx clicked")
  // store.sendMessage('enableTx', {})
}

function changeAudioDevice() {
  store.setAudioDevice(selectedAudioDevice.value)
}
</script>

<template>
  <div class="bg-slate-800 text-slate-100 p-4 rounded-lg flex flex-col gap-4">
    <div class="p-2 text-center rounded font-bold" :class="isConnected ? 'bg-green-600' : 'bg-red-600'">
      {{ isConnected ? 'Connected' : 'Disconnected' }}
    </div>

    <div class="grid grid-cols-2 gap-2" v-if="status">
      <div class="flex flex-col bg-slate-700 p-2 rounded">
        <label class="text-xs text-slate-400">Dial Freq</label>
        <span class="font-bold">{{ (status.DialFrequency / 1000000).toFixed(6) }} MHz</span>
      </div>
      <div class="flex flex-col bg-slate-700 p-2 rounded">
        <label class="text-xs text-slate-400">Mode</label>
        <span class="font-bold">{{ status.Mode }}</span>
      </div>
      <div class="flex flex-col bg-slate-700 p-2 rounded">
        <label class="text-xs text-slate-400">DX Call</label>
        <span class="font-bold">{{ status.DxCall }}</span>
      </div>
      <div class="flex flex-col bg-slate-700 p-2 rounded">
        <label class="text-xs text-slate-400">Report</label>
        <span class="font-bold">{{ status.Report }}</span>
      </div>
      <div class="flex flex-col bg-slate-700 p-2 rounded">
        <label class="text-xs text-slate-400">Tx Enabled</label>
        <span class="font-bold" :class="{ 'text-red-500': status.TxEnabled }">{{ status.TxEnabled ? 'YES' : 'NO' }}</span>
      </div>
      <div class="flex flex-col bg-slate-700 p-2 rounded">
        <label class="text-xs text-slate-400">Transmitting</label>
        <span class="font-bold" :class="{ 'text-red-500': status.Transmitting }">{{ status.Transmitting ? 'YES' : 'NO' }}</span>
      </div>
    </div>

    <div class="flex flex-col gap-2">
      <label>Audio Input:</label>
      <select v-model="selectedAudioDevice" @change="changeAudioDevice" class="p-2 rounded border-none bg-slate-700 text-white">
        <option value="">Default Device</option>
        <option v-for="device in audioDevices" :key="device.id" :value="device.id">
          {{ device.name }}
        </option>
      </select>
    </div>

    <div class="flex gap-2">
      <button class="flex-1 p-3 border-none rounded cursor-pointer font-bold transition-opacity hover:opacity-90 bg-red-600 text-white" @click="haltTx">Halt Tx</button>
      <button class="flex-1 p-3 border-none rounded cursor-pointer font-bold transition-opacity hover:opacity-90 bg-orange-500 text-white" @click="enableTx">Enable Tx</button>
    </div>
  </div>
</template>
