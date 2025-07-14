<template>
  <div class="game-view">
    <GuessGame />
  </div>
</template>

<script setup lang="ts">
import GuessGame from '@/components/GuessGame.vue';
</script>

<style scoped>
.game-view {
  min-height: 100vh;
  width: 100%;
  background: #0e0e0f;
  padding: 0;
}
</style>
