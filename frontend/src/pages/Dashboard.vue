<template>
  <div class="dashboard-page">
    <t-row :gutter="[16, 16]">
      <t-col :span="3">
        <t-card>
          <t-statistic :title="$t('dashboard.totalMonitors')" :value="stats?.total_monitors ?? 0" />
        </t-card>
      </t-col>
      <t-col :span="3">
        <t-card>
          <t-statistic :title="$t('dashboard.upMonitors')" :value="stats?.up_monitors ?? 0" color="green" />
        </t-card>
      </t-col>
      <t-col :span="3">
        <t-card>
          <t-statistic :title="$t('dashboard.downMonitors')" :value="stats?.down_monitors ?? 0" color="red" />
        </t-card>
      </t-col>
      <t-col :span="3">
        <t-card>
          <t-statistic :title="$t('dashboard.avgLatency')" :value="stats?.avg_latency?.toFixed(1) ?? '0'" suffix="ms" />
        </t-card>
      </t-col>
    </t-row>

    <t-row :gutter="[16, 16]" style="margin-top: 16px">
      <t-col :span="6">
        <t-card :title="$t('dashboard.latencyTrend')">
          <div ref="latencyChartRef" style="height: 300px"></div>
        </t-card>
      </t-col>
      <t-col :span="6">
        <t-card :title="$t('dashboard.uptimeTrend')">
          <div ref="uptimeChartRef" style="height: 300px"></div>
        </t-card>
      </t-col>
    </t-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useDashboardStore } from '../stores/dashboard'
import { getMonitors, getMonitorResults } from '../api/monitor'
import * as echarts from 'echarts'

const dashboardStore = useDashboardStore()
const stats = ref(dashboardStore.stats)

const latencyChartRef = ref<HTMLElement>()
const uptimeChartRef = ref<HTMLElement>()
let latencyChart: echarts.ECharts | null = null
let uptimeChart: echarts.ECharts | null = null

async function loadData() {
  await dashboardStore.fetchStats()
  stats.value = dashboardStore.stats

  const monitors = await getMonitors()
  const latencyData: any[] = []
  const uptimeData: any[] = []

  for (const m of monitors.slice(0, 10)) {
    const results = await getMonitorResults(m.id, 50)
    const latencies = results.map((r: any) => r.latency).reverse()
    const uptimes = results.map((r: any) => (r.status === 'up' ? 100 : 0)).reverse()
    latencyData.push({ name: m.name, type: 'line', data: latencies })
    uptimeData.push({ name: m.name, type: 'line', data: uptimes })
  }

  if (latencyChartRef.value) {
    latencyChart = echarts.init(latencyChartRef.value)
    latencyChart.setOption({
      tooltip: { trigger: 'axis' },
      legend: {},
      xAxis: { type: 'category' },
      yAxis: { type: 'value', name: 'ms' },
      series: latencyData,
    })
  }

  if (uptimeChartRef.value) {
    uptimeChart = echarts.init(uptimeChartRef.value)
    uptimeChart.setOption({
      tooltip: { trigger: 'axis' },
      legend: {},
      xAxis: { type: 'category' },
      yAxis: { type: 'value', name: '%', max: 100 },
      series: uptimeData,
    })
  }
}

onMounted(() => {
  loadData()
  window.addEventListener('resize', () => {
    latencyChart?.resize()
    uptimeChart?.resize()
  })
})

onUnmounted(() => {
  latencyChart?.dispose()
  uptimeChart?.dispose()
})
</script>

<style scoped>
.dashboard-page {
  max-width: 1400px;
}
</style>
