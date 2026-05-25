<template>
  <div>
    <t-card :title="project?.name || ''">
      <template #actions>
        <t-button @click="$router.push('/monitors')">{{ $t('nav.monitors') }}</t-button>
      </template>
      <t-row :gutter="[16, 16]">
        <t-col :span="3">
          <t-statistic :title="$t('project.monitorCount')" :value="dashboard?.total_monitors ?? 0" />
        </t-col>
        <t-col :span="3">
          <t-statistic :title="$t('project.uptime')" :value="dashboard?.uptime?.toFixed(1) ?? '0'" suffix="%" />
        </t-col>
        <t-col :span="3">
          <t-statistic :title="$t('project.avgLatency')" :value="dashboard?.avg_latency?.toFixed(1) ?? '0'" suffix="ms" />
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
          <t-button variant="text" @click="$router.push(`/monitors/${row.id}`)">{{ $t('common.edit') }}</t-button>
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
