<template>
  <div class="dashboard-page">
    <t-row :gutter="[16, 16]">
      <t-col :span="3">
        <div class="stat-card">
          <div class="stat-icon stat-icon--total"><t-icon name="root-list" size="24px" /></div>
          <div class="stat-info">
            <div class="stat-title">{{ $t('dashboard.totalMonitors') }}</div>
            <div class="stat-value">{{ stats?.total_monitors ?? 0 }}</div>
          </div>
        </div>
      </t-col>
      <t-col :span="3">
        <div class="stat-card">
          <div class="stat-icon stat-icon--up"><t-icon name="check-circle" size="24px" /></div>
          <div class="stat-info">
            <div class="stat-title">{{ $t('dashboard.upMonitors') }}</div>
            <div class="stat-value stat-value--up">{{ stats?.up_monitors ?? 0 }}</div>
          </div>
        </div>
      </t-col>
      <t-col :span="3">
        <div class="stat-card">
          <div class="stat-icon stat-icon--down"><t-icon name="close-circle" size="24px" /></div>
          <div class="stat-info">
            <div class="stat-title">{{ $t('dashboard.downMonitors') }}</div>
            <div class="stat-value stat-value--down">{{ stats?.down_monitors ?? 0 }}</div>
          </div>
        </div>
      </t-col>
      <t-col :span="3">
        <div class="stat-card">
          <div class="stat-icon stat-icon--latency"><t-icon name="time" size="24px" /></div>
          <div class="stat-info">
            <div class="stat-title">{{ $t('dashboard.avgLatency') }}</div>
            <div class="stat-value">{{ stats?.avg_latency?.toFixed(1) ?? '0' }}<span class="stat-unit">ms</span></div>
          </div>
        </div>
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

const darkAxis = { axisLine: { lineStyle: { color: 'rgba(255,255,255,0.1)' } }, axisLabel: { color: 'rgba(255,255,255,0.4)' }, splitLine: { lineStyle: { color: 'rgba(255,255,255,0.04)' } } }

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
    latencyData.push({ name: m.name, type: 'line', smooth: true, data: latencies, lineStyle: { width: 2 }, symbol: 'none' })
    uptimeData.push({ name: m.name, type: 'line', smooth: true, data: uptimes, lineStyle: { width: 2 }, symbol: 'none' })
  }

  if (latencyChartRef.value) {
    latencyChart = echarts.init(latencyChartRef.value)
    latencyChart.setOption({
      tooltip: { trigger: 'axis', backgroundColor: '#1a2540', borderColor: 'rgba(255,255,255,0.06)', textStyle: { color: '#fff' } },
      legend: { textStyle: { color: 'rgba(255,255,255,0.5)' } },
      xAxis: { type: 'category', ...darkAxis },
      yAxis: { type: 'value', name: 'ms', nameTextStyle: { color: 'rgba(255,255,255,0.4)' }, ...darkAxis },
      series: latencyData,
      grid: { left: 50, right: 20, top: 40, bottom: 30 },
    })
  }

  if (uptimeChartRef.value) {
    uptimeChart = echarts.init(uptimeChartRef.value)
    uptimeChart.setOption({
      tooltip: { trigger: 'axis', backgroundColor: '#1a2540', borderColor: 'rgba(255,255,255,0.06)', textStyle: { color: '#fff' } },
      legend: { textStyle: { color: 'rgba(255,255,255,0.5)' } },
      xAxis: { type: 'category', ...darkAxis },
      yAxis: { type: 'value', name: '%', max: 100, nameTextStyle: { color: 'rgba(255,255,255,0.4)' }, ...darkAxis },
      series: uptimeData,
      grid: { left: 50, right: 20, top: 40, bottom: 30 },
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
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.stat-icon--total { background: rgba(0, 212, 255, 0.1); color: #00d4ff; }
.stat-icon--up { background: rgba(82, 196, 26, 0.1); color: #73d13d; }
.stat-icon--down { background: rgba(255, 77, 79, 0.1); color: #ff7875; }
.stat-icon--latency { background: rgba(123, 97, 255, 0.1); color: #b37feb; }

.stat-info {
  flex: 1;
  min-width: 0;
}
.stat-title {
  font-size: 12px;
  color: var(--upp-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 1px;
  margin-bottom: 4px;
}
.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: var(--upp-text-primary);
  line-height: 1;
}
.stat-value--up { color: #73d13d; }
.stat-value--down { color: #ff7875; }
.stat-unit {
  font-size: 14px;
  font-weight: 400;
  color: var(--upp-text-tertiary);
  margin-left: 2px;
}
</style>
