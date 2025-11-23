<script setup lang="ts">
import ControlPanel from '@/components/ControlPanel.vue'
import BandActivity from '@/components/BandActivity.vue'
import WaterfallDisplay from '@/components/WaterfallDisplay.vue'
</script>

<template>
  <main class="grid grid-cols-1 lg:grid-cols-[1fr_350px] gap-4 h-[calc(100vh-8rem)]">
    <div class="overflow-hidden flex flex-col bg-gray-800 rounded-lg p-4">
      <BandActivity />
    </div>
    <div class="flex flex-col gap-4">
      <ControlPanel />
      <WaterfallDisplay />
    </div>
  </main>
</template>
