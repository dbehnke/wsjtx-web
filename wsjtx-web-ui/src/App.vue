<script setup lang="ts">
import { onMounted } from 'vue'
import { useWsjtxStore } from '@/stores/wsjtx'
import ControlPanel from '@/components/ControlPanel.vue'
import BandActivity from '@/components/BandActivity.vue'

const store = useWsjtxStore()

onMounted(() => {
  store.connect()
})
</script>

<template>
  <div class="app-container">
    <header>
      <h1>WSJT-X Web Interface</h1>
    </header>

    <main>
      <ControlPanel />
      <BandActivity />
    </main>
  </div>
</template>

<style scoped>
.app-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

header {
  margin-bottom: 2rem;
}

h1 {
  color: #fff;
  font-weight: 300;
}
</style>

<style>
body {
  background: #121212;
  color: #e0e0e0;
  margin: 0;
  font-family: 'Inter', sans-serif;
}
</style>
