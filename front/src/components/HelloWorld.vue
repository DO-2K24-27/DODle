<script setup lang="ts">
defineProps<{
  msg: string
}>()
</script>

<template>
  <div class="hello-world">
    <h1>{{ msg }}</h1>
  </div>
</template>

<style scoped>
.hello-world {
  margin: 1rem;
  padding: 1rem;
}
</style>
