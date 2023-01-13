<template>
  <div class="app-container">
    <el-form :model="queryParams" ref="queryFormRef" :inline="true" v-show="showSearch" label-width="68px">
      <el-form-item label="登录地址" prop="ipaddr">
        <el-input v-model="queryParams.ipaddr" placeholder="请输入登录地址" clearable style="width: 240px" @keyup.enter="handleQuery" />
      </el-form-item>
      <el-form-item label="用户名称" prop="username">
        <el-input v-model="queryParams.username" placeholder="请输入用户名称" clearable style="width: 240px" @keyup.enter="handleQuery" />
      </el-form-item>
      <el-form-item label="状态" prop="status">
        <el-select v-model="queryParams.status" placeholder="登录状态" clearable style="width: 240px">
          <el-option v-for="dict in sys_common_status" :key="dict.value" :label="dict.label" :value="dict.value" />
        </el-select>
      </el-form-item>
      <el-form-item label="登录时间" style="width: 308px">
        <el-date-picker
          v-model="dateRange"
          value-format="YYYY-MM-DD"
          type="daterange"
          range-separator="-"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
        ></el-date-picker>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="handleQuery"><ep:search /> 搜索</el-button>
        <el-button @click="resetQuery"><ep:refresh /> 重置</el-button>
      </el-form-item>
    </el-form>

    <el-row :gutter="10" class="mb8">
      <el-col :span="1.5">
        <el-button type="danger" plain :disabled="multiple" @click="handleDelete" v-hasPerms="['monitor:logininfor:remove']"><ep:delete /> 删除</el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button type="danger" plain @click="handleClean" v-hasPerms="['monitor:logininfor:remove']"><ep:delete />清空</el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button type="primary" plain :disabled="single" @click="handleUnlock" v-hasPerms="['monitor:logininfor:unlock']"><ep:unlock /> 解锁</el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button type="warning" plain @click="handleExport" v-hasPerms="['monitor:logininfor:export']"><ep:download /> 导出</el-button>
      </el-col>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList"></right-toolbar>
    </el-row>

    <el-table
      border
      ref="loginLogTableRef"
      v-loading="loading"
      :data="logininforList"
      @selection-change="handleSelectionChange"
      :default-sort="defaultSort"
      @sort-change="handleSortChange"
    >
      <el-table-column type="selection" width="55" align="center" />
      <el-table-column label="访问编号" align="center" prop="id" />
      <el-table-column
        label="用户名称"
        align="center"
        prop="username"
        :show-overflow-tooltip="true"
        sortable="custom"
        :sort-orders="['descending', 'ascending']"
      />
      <el-table-column label="地址" align="center" prop="ipaddr" :show-overflow-tooltip="true" />
      <el-table-column label="登录地点" align="center" prop="loginLocation" :show-overflow-tooltip="true" />
      <el-table-column label="操作系统" align="center" prop="os" :show-overflow-tooltip="true" />
      <el-table-column label="浏览器" align="center" prop="browser" :show-overflow-tooltip="true" />
      <el-table-column label="登录状态" align="center" prop="status">
        <template #default="scope">
          <dict-tag :options="sys_common_status" :value="scope.row.status" />
        </template>
      </el-table-column>
      <el-table-column label="描述" align="center" prop="msg" />
      <el-table-column label="访问时间" align="center" prop="loginTime" sortable="custom" :sort-orders="['descending', 'ascending']" width="180">
        <template #default="scope">
          <span>{{ parseTime(scope.row.loginTime) }}</span>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total > 0" :total="total" v-model:page="queryParams.pageNum" v-model:limit="queryParams.pageSize" @pagination="getList" />
  </div>
</template>

<script lang="ts" setup name="Logininfor">
import { list, delLogininfor, cleanLogininfor, unlockLogininfor } from '@/api/monitor/loginLog'
import useCurrentInstance from '@/hooks/useCurrentInstance'

const { proxy } = useCurrentInstance()
const { sys_common_status } = proxy.useDict('sys_common_status')

const logininforList = ref([])
const loading = ref(true)
const showSearch = ref(true)
const ids = ref([])
const single = ref(true)
const multiple = ref(true)
const selectName = ref('')
const total = ref(0)
const dateRange = ref([])
const defaultSort = ref({ prop: 'loginTime', order: 'descending' })

const queryFormRef = ref<ElForm>(null)
const loginLogTableRef = ref<ElTable>(null)

// 查询参数
const queryParams = ref({
  pageNum: 1,
  pageSize: 10,
  ipaddr: undefined,
  username: undefined,
  status: undefined,
  orderByColumn: undefined,
  isAsc: undefined
})

/** 查询登录日志列表 */
function getList() {
  loading.value = true
  list(proxy.addDateRange(queryParams.value, dateRange.value)).then(response => {
    logininforList.value = response.data.rows
    total.value = response.data.total
    loading.value = false
  })
}
/** 搜索按钮操作 */
function handleQuery() {
  queryParams.value.pageNum = 1
  getList()
}
/** 重置按钮操作 */
function resetQuery() {
  dateRange.value = []
  queryFormRef.value.resetFields()
  queryParams.value.pageNum = 1
  loginLogTableRef.value.sort(defaultSort.value.prop, defaultSort.value.order)
}
/** 多选框选中数据 */
function handleSelectionChange(selection:any) {
  ids.value = selection.map((item:any) => item.infoId)
  multiple.value = !selection.length
  single.value = selection.length != 1
  selectName.value = selection.map((item:any) => item.username)
}
/** 排序触发事件 */
function handleSortChange(column, prop, order) {
  queryParams.value.orderByColumn = column.prop
  queryParams.value.isAsc = column.order
  getList()
}
/** 删除按钮操作 */
function handleDelete(row: { [key: string]: any }) {
  const infoIds = row.infoId || ids.value
  proxy.$modal
    .confirm('是否确认删除访问编号为"' + infoIds + '"的数据项?')
    .then(function () {
      return delLogininfor(infoIds)
    })
    .then(() => {
      getList()
      proxy.$modal.msgSuccess('删除成功')
    })
    .catch(() => {})
}
/** 清空按钮操作 */
function handleClean() {
  proxy.$modal
    .confirm('是否确认清空所有登录日志数据项?')
    .then(function () {
      return cleanLogininfor()
    })
    .then(() => {
      getList()
      proxy.$modal.msgSuccess('清空成功')
    })
    .catch(() => {})
}
/** 解锁按钮操作 */
function handleUnlock() {
  const username = selectName.value
  proxy.$modal
    .confirm('是否确认解锁用户"' + username + '"数据项?')
    .then(function () {
      return unlockLogininfor(username)
    })
    .then(() => {
      proxy.$modal.msgSuccess('用户' + username + '解锁成功')
    })
    .catch(() => {})
}
/** 导出按钮操作 */
function handleExport() {
  proxy.download(
    'monitor/logininfor/export',
    {
      ...queryParams.value
    },
    `config_${new Date().getTime()}.xlsx`
  )
}

onMounted(() => {
  getList()
})
</script>
