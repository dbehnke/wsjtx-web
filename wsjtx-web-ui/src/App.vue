<script setup lang="ts">
import { onMounted } from 'vue'
import { useWsjtxStore } from '@/stores/wsjtx'
import HomeView from '@/views/HomeView.vue'

const store = useWsjtxStore()

onMounted(() => {
  store.connect()
})
</script>

<template>
  <div class="min-h-screen bg-gray-900 text-gray-200 font-sans p-8">
    <header class="mb-8">
      <h1 class="text-white font-light text-3xl">WSJT-X Web Interface</h1>
    </header>

    <HomeView />
  </div>
</template>
