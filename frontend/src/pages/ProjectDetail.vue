<template>
  <div>
    <t-card :title="project?.name || ''">
      <template #actions>
        <t-button @click="$router.push('/monitors')">{{ $t('nav.monitors') }}</t-button>
      </template>
      <t-row :gutter="[16, 16]">
        <t-col :span="3">
          <div class="stat-card">
            <div class="stat-icon stat-icon--total"><t-icon name="root-list" size="24px" /></div>
            <div class="stat-info">
              <div class="stat-title">{{ $t('project.monitorCount') }}</div>
              <div class="stat-value">{{ dashboard?.total_monitors ?? 0 }}</div>
            </div>
          </div>
        </t-col>
        <t-col :span="3">
          <div class="stat-card">
            <div class="stat-icon stat-icon--uptime"><t-icon name="chart" size="24px" /></div>
            <div class="stat-info">
              <div class="stat-title">{{ $t('project.uptime') }}</div>
              <div class="stat-value">{{ dashboard?.uptime?.toFixed(1) ?? '0' }}<span class="stat-unit">%</span></div>
            </div>
          </div>
        </t-col>
        <t-col :span="3">
          <div class="stat-card">
            <div class="stat-icon stat-icon--latency"><t-icon name="time" size="24px" /></div>
            <div class="stat-info">
              <div class="stat-title">{{ $t('project.avgLatency') }}</div>
              <div class="stat-value">{{ dashboard?.avg_latency?.toFixed(1) ?? '0' }}<span class="stat-unit">ms</span></div>
            </div>
          </div>
        </t-col>
      </t-row>
    </t-card>

    <t-card style="margin-top: 16px" :title="$t('nav.monitors')">
      <t-table :data="monitors" :columns="columns" row-key="id">
        <template #status="{ row }">
          <t-tag :theme="row.enabled ? 'success' : 'default'">
            {{ row.enabled ? $t('monitor.enabled') : $t('monitor.disabled') }}
          </t-tag>
        </template>
        <template #operation="{ row }">
          <t-button variant="text" @click="$router.push(`/monitors/${row.id}`)">{{ $t('common.view') }}</t-button>
        </template>
      </t-table>
    </t-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getProject } from '../api/project'
import { getMonitorsByProject } from '../api/monitor'
import { getProjectDashboard } from '../api/dashboard'
import { useI18n } from 'vue-i18n'

const route = useRoute()
const { t } = useI18n()

const project = ref<any>(null)
const monitors = ref<any[]>([])
const dashboard = ref<any>(null)

const columns = [
  { colKey: 'id', title: 'ID', width: 60 },
  { colKey: 'name', title: t('monitor.name') },
  { colKey: 'type', title: t('monitor.type') },
  { colKey: 'target', title: t('monitor.target') },
  { colKey: 'status', title: t('monitor.status') },
  { colKey: 'operation', title: t('common.operation'), width: 100 },
]

onMounted(async () => {
  const id = Number(route.params.id)
  project.value = await getProject(id)
  monitors.value = await getMonitorsByProject(id)
  dashboard.value = await getProjectDashboard(id)
})
</script>

<style scoped>
.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px 24px;
  background: var(--upp-bg-card);
  border: 1px solid var(--upp-border);
  border-radius: 12px;
  backdrop-filter: blur(12px);
  transition: all 0.3s;
}
.stat-card:hover {
  border-color: var(--upp-border-hover);
  box-shadow: 0 0 20px rgba(0, 180, 255, 0.06);
}
.stat-icon {
  width: 48px; height: 48px; border-radius: 12px;
  display: flex; align-items: center; justify-content: center; flex-shrink: 0;
}
.stat-icon--total { background: rgba(0, 212, 255, 0.1); color: #00d4ff; }
.stat-icon--uptime { background: rgba(82, 196, 26, 0.1); color: #73d13d; }
.stat-icon--latency { background: rgba(123, 97, 255, 0.1); color: #b37feb; }
.stat-info { flex: 1; min-width: 0; }
.stat-title {
  font-size: 12px; color: var(--upp-text-tertiary);
  text-transform: uppercase; letter-spacing: 1px; margin-bottom: 4px;
}
.stat-value {
  font-size: 28px; font-weight: 700; color: var(--upp-text-primary); line-height: 1;
}
.stat-unit {
  font-size: 14px; font-weight: 400; color: var(--upp-text-tertiary); margin-left: 2px;
}
</style>
