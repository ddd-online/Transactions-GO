<template>
  <div class="about-setting">
    <div class="about-header">
      <div class="app-logo">
        <svg viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg">
          <defs>
            <filter id="shadow" x="-20%" y="-20%" width="140%" height="140%">
              <feGaussianBlur stdDeviation="16" result="blur"/>
              <feOffset dx="0" dy="16"/>
              <feComponentTransfer><feFuncA type="linear" slope="0.15"/></feComponentTransfer>
              <feMerge><feMergeNode in="shadow"/><feMergeNode in="SourceGraphic"/></feMerge>
            </filter>
          </defs>
          <rect x="0" y="0" width="1024" height="1024" rx="200" fill="#f5f5f7"/>
          <g transform="translate(160, 160)" filter="url(#shadow)">
            <rect x="32" y="100" width="640" height="112" rx="24" fill="#52c41a"/>
            <rect x="32" y="260" width="640" height="112" rx="24" fill="#f5222d"/>
            <rect x="32" y="420" width="640" height="112" rx="24" fill="#1677ff"/>
            <g transform="translate(-5.8716431)">
              <polyline points="120,680 260,600 400,680 540,560 680,640" fill="none" stroke="#8e8e93" stroke-width="20" stroke-linecap="round" stroke-linejoin="round" transform="translate(-42.12836,1.1234229)"/>
              <circle cx="77.871643" cy="681.12341" r="18" fill="#8e8e93"/>
              <circle cx="217.87164" cy="601.12341" r="18" fill="#8e8e93"/>
              <circle cx="357.87164" cy="681.12341" r="18" fill="#8e8e93"/>
              <circle cx="497.87164" cy="561.12341" r="18" fill="#8e8e93"/>
              <circle cx="637.87164" cy="641.12341" r="18" fill="#8e8e93"/>
            </g>
          </g>
        </svg>
      </div>
      <h2 class="app-name">Transactions</h2>
      <p class="app-version">版本 1.0.0</p>
    </div>

    <div class="about-description">
      <p>简洁高效的个人财务管理工具</p>
      <p>多账本 · 分类标签 · 数据分析</p>
    </div>

    <div class="about-links">
      <div class="link-item">
        <span class="link-label">技术栈</span>
        <span class="link-value">Electron + Vue.js + Go</span>
      </div>
      <div class="link-item">
        <span class="link-label">构建时间</span>
        <span class="link-value">{{ buildTime }}</span>
      </div>
    </div>

    <div class="about-copyright">
      <p>© 2024 Transactions. All rights reserved.</p>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {ref} from 'vue';

const buildTime = ref('2024-04-20');
</script>

<style scoped>
.about-setting {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--billadm-space-xl);
  padding: var(--billadm-space-xl) 0;
}

.about-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--billadm-space-md);
}

.app-logo {
  width: 96px;
  height: 96px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.app-logo svg {
  width: 96px;
  height: 96px;
}

.app-name {
  font-size: 24px;
  font-weight: 600;
  color: var(--billadm-color-text-major);
  margin: 0;
}

.app-version {
  font-size: var(--billadm-size-text-body);
  color: var(--billadm-color-text-secondary);
  margin: 0;
}

.about-description {
  text-align: center;
  color: var(--billadm-color-text-secondary);
  font-size: var(--billadm-size-text-body);
  line-height: 1.8;
}

.about-links {
  display: flex;
  flex-direction: column;
  gap: var(--billadm-space-md);
  width: 100%;
  max-width: 320px;
  padding: var(--billadm-space-lg);
  background: var(--billadm-color-minor-background);
  border-radius: var(--billadm-radius-md);
}

.link-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.link-label {
  font-size: var(--billadm-size-text-body);
  color: var(--billadm-color-text-secondary);
}

.link-value {
  font-size: var(--billadm-size-text-body);
  color: var(--billadm-color-text-major);
  font-weight: 500;
}

.about-copyright {
  margin-top: auto;
  padding-top: var(--billadm-space-xl);
  text-align: center;
  color: var(--billadm-color-text-secondary);
  font-size: var(--billadm-size-text-caption);
}
</style>
