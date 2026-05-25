<template>
  <div>
    <t-card :title="monitor?.name || ''">
      <template #actions>
        <t-button @click="$router.push('/monitors')">{{ $t('nav.monitors') }}</t-button>
      </template>
      <t-row :gutter="[16, 16]">
        <t-col :span="3">
          <t-statistic :title="$t('monitor.status')" :value="latestResult?.status === 'up' ? $t('monitor.up') : $t('monitor.down')" />
        </t-col>
        <t-col :span="3">
          <t-statistic :title="$t('monitor.latency')" :value="latestResult?.latency?.toFixed(1) ?? '0'" suffix="ms" />
        </t-col>
        <t-col :span="3">
          <t-statistic :title="$t('monitor.uptime')" :value="stats?.uptime?.toFixed(1) ?? '0'" suffix="%" />
        </t-col>
        <t-col :span="3">
          <t-statistic :title="$t('monitor.avgLatency')" :value="stats?.avg_latency?.toFixed(1) ?? '0'" suffix="ms" />
        </t-col>
      </t-row>
    </t-card>

    <t-card style="margin-top: 16px" :title="$t('dashboard.latencyTrend')">
      <div ref="chartRef" style="height: 350px"></div>
    </t-card>

    <t-card style="margin-top: 16px" :title="$t('monitor.history')">
      <t-table :data="results" :columns="resultColumns" row-key="id" :loading="loading" />
    </t-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { getMonitor, getMonitorResults, getMonitorStats, type MonitorResult, type MonitorStats } from '../api/monitor'
import { useI18n } from 'vue-i18n'
import * as echarts from 'echarts'

const route = useRoute()
const { t } = useI18n()

const monitor = ref<any>(null)
const results = ref<MonitorResult[]>([])
const stats = ref<MonitorStats | null>(null)
const latestResult = ref<MonitorResult | null>(null)
const loading = ref(false)
const chartRef = ref<HTMLElement>()
let chart: echarts.ECharts | null = null

const resultColumns = [
  { colKey: 'id', title: 'ID', width: 60 },
  { colKey: 'status', title: t('monitor.status') },
  { colKey: 'latency', title: t('monitor.latency') },
  { colKey: 'message', title: 'Message' },
  { colKey: 'created_at', title: t('monitor.createdAt'), cell: (_h: any, { row }: any) => formatTime(row.created_at) },
]

function formatTime(val: string) {
  if (!val) return ''
  const d = new Date(val)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}:${pad(d.getMonth() + 1)}:${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}

onMounted(async () => {
  const id = Number(route.params.id)
  loading.value = true
  try {
    monitor.value = await getMonitor(id)
    results.value = await getMonitorResults(id, 100)
    stats.value = await getMonitorStats(id)
    latestResult.value = results.value.length > 0 ? results.value[0] : null

    // Render chart
    if (chartRef.value) {
      chart = echarts.init(chartRef.value)
      const reversed = [...results.value].reverse()
      chart.setOption({
        tooltip: { trigger: 'axis' },
        xAxis: { type: 'category', data: reversed.map((r) => r.created_at) },
        yAxis: { type: 'value', name: 'ms' },
        series: [
          {
            type: 'line',
            data: reversed.map((r) => r.latency),
            areaStyle: { opacity: 0.3 },
            itemStyle: { color: '#1890ff' },
          },
        ],
      })
    }
  } finally {
    loading.value = false
  }

  window.addEventListener('resize', onResize)
})

function onResize() {
  chart?.resize()
}

onUnmounted(() => {
  chart?.dispose()
  window.removeEventListener('resize', onResize)
})
</script>
