<script setup lang="ts">
import { useWsjtxStore } from '@/stores/wsjtx'
import { storeToRefs } from 'pinia'

const store = useWsjtxStore()
const { messages } = storeToRefs(store)

function formatTime(ms: number) {
  // ms is milliseconds since midnight
  const date = new Date(ms)
  // Since it's since midnight, we can just format it as UTC time if we treat it as a duration,
  // or just do simple math.
  const seconds = Math.floor((ms / 1000) % 60)
  const minutes = Math.floor((ms / (1000 * 60)) % 60)
  const hours = Math.floor((ms / (1000 * 60 * 60)) % 24)
  
  return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`
}

function handleReply(msg: any) {
  // Construct a reply message based on the decode
  // For now, we just send the decode info back as a "reply" command
  // The backend will wrap it in a ReplyMessage
  store.sendMessage('reply', {
    Time: msg.Time,
    SNR: msg.SNR,
    DeltaTime: msg.DeltaTime,
    DeltaFrequency: msg.DeltaFrequency,
    Mode: msg.Mode,
    Message: msg.Message,
    LowConfidence: msg.LowConfidence,
    Modifiers: 0
  })
}
</script>

<template>
  <div class="band-activity">
    <h2>Band Activity</h2>
    <table>
      <thead>
        <tr>
          <th>Time</th>
          <th>dB</th>
          <th>DT</th>
          <th>Freq</th>
          <th>Mode</th>
          <th>Message</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(msg, index) in messages" :key="index" @dblclick="handleReply(msg)" class="decode-row">
          <td>{{ formatTime(msg.Time) }}</td>
          <td>{{ msg.SNR }}</td>
          <td>{{ msg.DeltaTime.toFixed(1) }}</td>
          <td>{{ msg.DeltaFrequency }}</td>
          <td>{{ msg.Mode }}</td>
          <td class="message">{{ msg.Message }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.band-activity {
  background: #1e1e1e;
  color: #d4d4d4;
  padding: 1rem;
  border-radius: 8px;
  overflow-y: auto;
  max-height: 500px;
}

table {
  width: 100%;
  border-collapse: collapse;
  font-family: monospace;
}

th, td {
  padding: 4px 8px;
  text-align: left;
}

th {
  border-bottom: 1px solid #444;
}

tr:nth-child(even) {
  background: #252526;
}

.message {
  color: #9cdcfe;
}

.decode-row {
  cursor: pointer;
}

.decode-row:hover {
  background: #2a2d2e !important;
}
</style>
