<template>
  <div class="bottom-bar">
    <div class="bottom-bar-left"></div>
    <div class="bottom-bar-right"></div>
  </div>
</template>

<style scoped>
.bottom-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 100%;
  width: 100%;
  padding: 0 12px;
}

.bottom-bar > * {
  -webkit-app-region: no-drag;
}

.bottom-bar-left {
  display: flex;
  align-items: center;
}

.bottom-bar-right {
  display: flex;
  align-items: center;
  gap: 4px;
}
</style>
