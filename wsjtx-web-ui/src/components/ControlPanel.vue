<script setup lang="ts">
import { useWsjtxStore } from '@/stores/wsjtx'
import { storeToRefs } from 'pinia'

const store = useWsjtxStore()
const { status, isConnected } = storeToRefs(store)

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
  // Is there a message to set "Enable Tx"?
  // MsgSwitchConfiguration? MsgConfigure?
  // Let's look at protocol.go.
  // There is no explicit "Enable Tx" message.
  // However, we can send a "Reply" which might enable Tx if Auto Tx is on?
  // Or maybe we can't enable Tx remotely easily without a specific message.
  // Let's just implement Halt for now, and maybe "Enable Tx" is just a placeholder or we use a workaround if found.
  // For now, let's just log it or send a "halt" with AutoTxOnly=false (which stops it).
  // Actually, let's just implement Halt.
  console.log("Enable Tx not fully implemented yet")
}
</script>

<template>
  <div class="control-panel">
    <div class="header">
      <h2>Control Panel</h2>
      <span class="status-indicator" :class="{ connected: isConnected }">
        {{ isConnected ? 'Connected' : 'Disconnected' }}
      </span>
    </div>

    <div v-if="status" class="status-grid">
      <div class="stat-item">
        <label>Dial Freq</label>
        <span>{{ status.DialFrequency }} Hz</span>
      </div>
      <div class="stat-item">
        <label>Mode</label>
        <span>{{ status.Mode }}</span>
      </div>
      <div class="stat-item">
        <label>DX Call</label>
        <span>{{ status.DXCall || '--' }}</span>
      </div>
      <div class="stat-item">
        <label>Report</label>
        <span>{{ status.Report || '--' }}</span>
      </div>
      <div class="stat-item">
        <label>Tx Enabled</label>
        <span :class="{ active: status.TxEnabled }">{{ status.TxEnabled ? 'ON' : 'OFF' }}</span>
      </div>
      <div class="stat-item">
        <label>Transmitting</label>
        <span :class="{ active: status.Transmitting }">{{ status.Transmitting ? 'YES' : 'NO' }}</span>
      </div>
    </div>
    <div v-else class="no-status">
      Waiting for status...
    </div>

    <div class="controls">
      <button @click="haltTx" class="btn btn-danger">Halt Tx</button>
      <button @click="enableTx" class="btn btn-success">Enable Tx</button>
    </div>
  </div>
</template>

<style scoped>
.control-panel {
  background: #252526;
  color: #d4d4d4;
  padding: 1rem;
  border-radius: 8px;
  margin-bottom: 1rem;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.status-indicator {
  padding: 4px 8px;
  border-radius: 4px;
  background: #c72e2e;
  font-size: 0.8rem;
}

.status-indicator.connected {
  background: #379c37;
}

.status-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 1rem;
}

.stat-item {
  display: flex;
  flex-direction: column;
  background: #1e1e1e;
  padding: 0.5rem;
  border-radius: 4px;
}

.stat-item label {
  font-size: 0.8rem;
  color: #888;
  margin-bottom: 0.25rem;
}

.stat-item span {
  font-family: monospace;
  font-size: 1.1rem;
}

.active {
  color: #379c37;
  font-weight: bold;
}

.no-status {
  text-align: center;
  color: #888;
  padding: 2rem;
}

.controls {
  margin-top: 1rem;
  display: flex;
  gap: 1rem;
}

.btn {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: bold;
  color: white;
}

.btn-danger {
  background: #c72e2e;
}

.btn-danger:hover {
  background: #a82525;
}

.btn-success {
  background: #379c37;
}

.btn-success:hover {
  background: #2d802d;
}
</style>
