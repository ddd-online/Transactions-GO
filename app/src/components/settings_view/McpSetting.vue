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

    <!-- Claude Code 配置 -->
    <div v-if="mcpStore.isRunning" class="config-section">
      <div class="config-title">Claude Code 配置</div>
      <div class="config-desc">在终端运行以下命令添加 MCP 服务器：</div>
      <pre class="config-code"><code>{{ configCommand }}</code></pre>

      <div class="config-title" style="margin-top: 16px;">可用工具</div>
      <div class="tools-list">
        <div class="tool-item">
          <div class="tool-name">query_ledgers</div>
          <div class="tool-desc">查询所有账本列表，返回账本 ID、名称和描述</div>
          <div class="tool-example">
            <div class="example-label">输出样例</div>
            <pre class="example-output"><code>[a1b2c3d4] 主账本 - 我的日常开支记录
[e5f6g7h8] 生意账本</code></pre>
          </div>
        </div>
        <div class="tool-item">
          <div class="tool-name">query_transactions</div>
          <div class="tool-desc">按条件查询交易记录</div>
          <div class="tool-params">
            <code>ledger_id</code> (必填) - 账本ID<br>
            <code>time_range</code> - 时间戳范围 [start, end]，单位为秒级时间戳<br>
            <code>transaction_type</code> - expense/income/transfer，收入不含转账<br>
            <code>category</code> - 分类名称<br>
            <code>tags</code> - 标签列表<br>
            <code>description</code> - 备注包含的字符<br>
            <code>offset</code> - 分页偏移<br>
            <code>limit</code> - 每页数量，最大100
          </div>
          <div class="tool-example">
            <div class="example-label">输出样例</div>
            <pre class="example-output"><code>共 2 条记录

2026-04-10 14:30:00 | expense | 金额: 58.50 | 分类: 餐饮 | 备注: 午餐 | 标签: 工作餐
2026-04-10 09:15:00 | income | 金额: 15000.00 | 分类: 工资 | 备注: 月薪 | 标签: 固定收入 [离群值]</code></pre>
          </div>
        </div>
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

const configCommand = computed(() => {
  return `claude mcp add --transport http transactions http://127.0.0.1:${mcpStore.port}/mcp`;
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
  gap: var(--billadm-space-lg);
}

.setting-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--billadm-space-xl);
  padding: var(--billadm-space-md) 0;
  border-bottom: 1px solid var(--billadm-color-divider);
}

.setting-row:last-child {
  border-bottom: none;
}

.setting-info {
  display: flex;
  flex-direction: column;
  gap: var(--billadm-space-xs);
}

.info-label {
  font-size: var(--billadm-size-text-body);
  font-weight: 500;
  color: var(--billadm-color-text-major);
}

.info-desc {
  font-size: var(--billadm-size-text-caption);
  color: var(--billadm-color-text-secondary);
}

.status-indicator {
  display: flex;
  align-items: center;
}

.port-display {
  font-size: var(--billadm-size-text-body);
  font-family: monospace;
  padding: var(--billadm-space-xs) var(--billadm-space-md);
  background-color: var(--billadm-color-minor-background);
  border-radius: var(--billadm-radius-md);
}

.config-section {
  margin-top: var(--billadm-space-md);
  padding: var(--billadm-space-lg);
  background-color: var(--billadm-color-minor-background);
  border-radius: var(--billadm-radius-lg);
  display: flex;
  flex-direction: column;
  gap: var(--billadm-space-md);
}

.config-title {
  font-size: var(--billadm-size-text-body);
  font-weight: 600;
  color: var(--billadm-color-text-major);
}

.config-desc {
  font-size: var(--billadm-size-text-caption);
  color: var(--billadm-color-text-secondary);
}

.config-desc code {
  font-family: monospace;
  background-color: var(--billadm-color-hover-bg);
  padding: 2px var(--billadm-space-xs);
  border-radius: var(--billadm-radius-sm);
}

.config-code {
  background-color: #1e1e1e;
  color: #d4d4d4;
  padding: var(--billadm-space-md);
  border-radius: var(--billadm-radius-md);
  font-size: var(--billadm-size-text-caption);
  font-family: 'Consolas', 'Monaco', monospace;
  overflow-x: auto;
  margin: 0;
}

.config-code code {
  white-space: pre;
}

.config-tip {
  font-size: var(--billadm-size-text-caption);
  color: var(--billadm-color-text-secondary);
}

.tools-list {
  display: flex;
  flex-direction: column;
  gap: var(--billadm-space-md);
}

.tool-item {
  padding: var(--billadm-space-md);
  background-color: var(--billadm-color-hover-bg);
  border-radius: var(--billadm-radius-md);
}

.tool-name {
  font-family: monospace;
  font-size: var(--billadm-size-text-body-sm);
  font-weight: 500;
  color: var(--billadm-color-primary);
  margin-bottom: var(--billadm-space-xs);
}

.tool-desc {
  font-size: var(--billadm-size-text-caption);
  color: var(--billadm-color-text-secondary);
}

.tool-params {
  font-size: var(--billadm-size-text-small);
  color: var(--billadm-color-text-secondary);
  margin-top: var(--billadm-space-sm);
  line-height: 1.6;
}

.tool-params code {
  font-family: monospace;
  background-color: rgba(0, 0, 0, 0.08);
  padding: 1px var(--billadm-space-xs);
  border-radius: var(--billadm-radius-sm);
}

.tool-example {
  margin-top: var(--billadm-space-sm);
  padding-top: var(--billadm-space-sm);
  border-top: 1px dashed var(--billadm-color-divider);
}

.example-label {
  font-size: var(--billadm-size-text-small);
  color: var(--billadm-color-text-secondary);
  margin-bottom: var(--billadm-space-xs);
}

.example-output {
  background-color: #1e1e1e;
  color: #d4d4d4;
  padding: var(--billadm-space-sm);
  border-radius: var(--billadm-radius-md);
  font-size: var(--billadm-size-text-small);
  font-family: monospace;
  overflow-x: auto;
  margin: 0;
  line-height: 1.5;
}

.example-output code {
  white-space: pre;
}
</style>
