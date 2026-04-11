<template>
  <div class="mcp-setting">
    <div class="setting-row">
      <div class="setting-info">
        <div class="info-label">MCP 服务</div>
        <div class="info-desc">Model Context Protocol TCP 服务，用于 AI 工具集成</div>
      </div>
      <a-switch v-model:checked="switchLoading" :loading="loading" @change="handleToggle"/>
    </div>

    <div class="setting-row">
      <div class="setting-info">
        <div class="info-label">运行状态</div>
        <div class="info-desc">当前 MCP 服务连接状态</div>
      </div>
      <div class="status-indicator">
        <a-badge :status="mcpStore.isRunning ? 'success' : 'default'" :text="mcpStore.isRunning ? '运行中' : '已停止'"/>
      </div>
    </div>

    <div class="setting-row">
      <div class="setting-info">
        <div class="info-label">端口</div>
        <div class="info-desc">MCP TCP 服务监听端口</div>
      </div>
      <div class="port-display">{{ mcpStore.port }}</div>
    </div>

    <!-- Claude Code 配置示例 -->
    <div v-if="mcpStore.isRunning" class="config-section">
      <div class="config-title">Claude Code 配置</div>
      <div class="config-desc">在 <code>~/.claude/settings.json</code> 中添加以下配置：</div>
      <pre class="config-code"><code>{{ configExample }}</code></pre>
      <div class="config-tip">
        配置后重启 Claude Code，然后可以使用 MCP 工具查询账本和交易记录。
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {ref, computed, onMounted} from 'vue';
import {useMcpStore} from '@/stores/mcpStore';

const mcpStore = useMcpStore();
const loading = ref(false);
const switchLoading = ref(false);

const configExample = computed(() => {
  return JSON.stringify({
    "mcpServers": {
      "billadm": {
        "type": "tcp",
        "host": "127.0.0.1",
        "port": mcpStore.port
      }
    }
  }, null, 2);
});

onMounted(async () => {
  await mcpStore.refreshStatus();
  switchLoading.value = mcpStore.isRunning;
});

const handleToggle = async (checked: boolean | string | number) => {
  loading.value = true;
  try {
    if (checked) {
      await mcpStore.start();
    } else {
      await mcpStore.stop();
    }
  } catch (error) {
    // revert switch state on error
    switchLoading.value = mcpStore.isRunning;
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.mcp-setting {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.setting-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
}

.setting-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-label {
  font-size: 14px;
  color: var(--billadm-color-text-major);
}

.info-desc {
  font-size: 12px;
  color: var(--billadm-color-text-minor);
}

.status-indicator {
  display: flex;
  align-items: center;
}

.port-display {
  font-size: 14px;
  font-family: monospace;
  padding: 6px 12px;
  background-color: var(--billadm-color-minor-background);
  border-radius: 4px;
}

.config-section {
  margin-top: 12px;
  padding: 16px;
  background-color: var(--billadm-color-minor-background);
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.config-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--billadm-color-text-major);
}

.config-desc {
  font-size: 12px;
  color: var(--billadm-color-text-minor);
}

.config-desc code {
  font-family: monospace;
  background-color: var(--billadm-color-icon-hover-bg);
  padding: 2px 4px;
  border-radius: 3px;
}

.config-code {
  background-color: #1e1e1e;
  color: #d4d4d4;
  padding: 12px;
  border-radius: 6px;
  font-size: 12px;
  font-family: 'Consolas', 'Monaco', monospace;
  overflow-x: auto;
  margin: 0;
}

.config-code code {
  white-space: pre;
}

.config-tip {
  font-size: 12px;
  color: var(--billadm-color-text-minor);
}
</style>
