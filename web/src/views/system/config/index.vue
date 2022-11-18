<template>
  <div class="app-container">
    <el-form :model="queryParams" ref="queryFormRef" :inline="true" v-show="showSearch" label-width="68px">
      <el-form-item label="参数名称" prop="configName">
        <el-input v-model="queryParams.configName" placeholder="请输入参数名称" clearable style="width: 240px" @keyup.enter="handleQuery" />
      </el-form-item>
      <el-form-item label="参数键名" prop="configKey">
        <el-input v-model="queryParams.configKey" placeholder="请输入参数键名" clearable style="width: 240px" @keyup.enter="handleQuery" />
      </el-form-item>
      <el-form-item label="系统内置" prop="configType">
        <el-select v-model="queryParams.configType" placeholder="系统内置" clearable>
          <el-option v-for="dict in sys_yes_no" :key="dict.value" :label="dict.label" :value="dict.value" />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="handleQuery"> <ep:search /> 搜索 </el-button>
        <el-button @click="resetQuery"> <ep:refresh /> 重置 </el-button>
      </el-form-item>
    </el-form>

    <el-row :gutter="10" class="mb8">
      <el-col :span="1.5">
        <el-button type="primary" plain @click="handleAdd" v-hasPerms="['system:config:add']"> <ep:plus /> 新增 </el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button type="danger" plain :disabled="multiple" @click="handleDelete" v-hasPerms="['system:config:remove']"> <ep:delete /> 删除 </el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button type="warning" plain @click="handleExport" v-hasPerms="['system:config:export']"> <ep:download /> 导出 </el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button type="danger" plain @click="handleRefreshCache" v-hasPerms="['system:config:remove']"> <ep:refresh />刷新缓存 </el-button>
      </el-col>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList"></right-toolbar>
    </el-row>

    <el-table border v-loading="loading" :data="configList" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55" align="center" />
      <el-table-column label="参数主键" align="center" prop="configId" />
      <el-table-column label="参数名称" align="center" prop="configName" :show-overflow-tooltip="true" />
      <el-table-column label="参数键名" align="center" prop="configKey" :show-overflow-tooltip="true" />
      <el-table-column label="参数键值" align="center" prop="configValue" />
      <el-table-column label="系统内置" align="center" prop="configType">
        <template #default="scope">
          <dict-tag :options="sys_yes_no" :value="scope.row.configType" />
        </template>
      </el-table-column>
      <el-table-column label="备注" align="center" prop="remark" :show-overflow-tooltip="true" />
      <el-table-column label="创建时间" align="center" prop="createTime" width="180">
        <template #default="scope">
          <span>{{ parseTime(scope.row.createTime) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" width="150" class-name="small-padding fixed-width">
        <template #default="scope">
          <el-button link type="primary" @click="handleUpdate(scope.row)" v-hasPerms="['system:config:edit']"> <ep:edit /> 修改 </el-button>
          <el-button link type="danger" @click="handleDelete(scope.row)" v-hasPerms="['system:config:remove']"> <ep:delete /> 删除 </el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total > 0" :total="total" v-model:page="queryParams.pageNum" v-model:limit="queryParams.pageSize" @pagination="getList" />

    <!-- 添加或修改参数配置对话框 -->
    <el-dialog :title="dialog.title" v-model="dialog.visible" width="500px" @close="closeDialog" append-to-body>
      <el-form ref="configFormRef" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="参数名称" prop="configName">
          <el-input v-model="form.configName" placeholder="请输入参数名称" />
        </el-form-item>
        <el-form-item label="参数键名" prop="configKey">
          <el-input v-model="form.configKey" placeholder="请输入参数键名" />
        </el-form-item>
        <el-form-item label="参数键值" prop="configValue">
          <el-input v-model="form.configValue" placeholder="请输入参数键值" />
        </el-form-item>
        <el-form-item label="系统内置" prop="configType">
          <el-radio-group v-model="form.configType">
            <el-radio v-for="dict in sys_yes_no" :key="dict.value" :label="dict.value">
              {{ dict.label }}
            </el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="form.remark" type="textarea" placeholder="请输入内容" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button type="primary" @click="submitForm">确 定</el-button>
          <el-button @click="closeDialog">取 消</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup name="Config">
import { listConfig, getConfig, delConfig, addConfig, updateConfig, refreshCache } from '@/api/system/config'
import useCurrentInstance from '@/hooks/useCurrentInstance'
import { ConfigFormData } from '@/types/api/config'
import { Dialog } from '@/types/common'

const { proxy } = useCurrentInstance()
const { sys_yes_no } = proxy.useDict('sys_yes_no')

const configList = ref([])
const loading = ref(true)
const showSearch = ref(true)
const ids = ref([])
const single = ref(true)
const multiple = ref(true)
const total = ref(0)
const dateRange = ref([])

const queryFormRef = ref<ElForm>(null)
const configFormRef = ref<ElForm>(null)

const data = reactive({
  dialog: {
    visible: false,
    title: ''
  } as Dialog,
  form: {
    configType: 'Y'
  } as  ConfigFormData,
  queryParams: {
    pageNum: 1,
    pageSize: 10,
    configName: undefined,
    configKey: undefined,
    configType: undefined
  },
  rules: {
    configName: [{ required: true, message: '参数名称不能为空', trigger: 'blur' }],
    configKey: [{ required: true, message: '参数键名不能为空', trigger: 'blur' }],
    configValue: [{ required: true, message: '参数键值不能为空', trigger: 'blur' }]
  }
})

const { dialog, queryParams, form, rules } = toRefs(data)

/** 查询参数列表 */
function getList() {
  loading.value = true
  listConfig(queryParams.value).then(response => {
    configList.value = response.data.rows
    total.value = response.data.total
    loading.value = false
  })
}
/** 关闭弹窗 */
function closeDialog() {
  dialog.value.visible = false
  configFormRef.value?.resetFields()
  configFormRef.value?.clearValidate()
  form.value.configId = undefined
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
  handleQuery()
}
/** 多选框选中数据 */
function handleSelectionChange(selection:any) {
  ids.value = selection.map((item:any) => item.configId)
  single.value = selection.length != 1
  multiple.value = !selection.length
}
/** 新增按钮操作 */
function handleAdd() {
  dialog.value.visible = true
  dialog.value.title = '添加参数'
}
/** 修改按钮操作 */
function handleUpdate(row: { [key: string]: any }) {
  const configId = row.configId || ids.value
  getConfig(configId).then(response => {
    dialog.value.visible = true
    dialog.value.title = '修改参数'
    nextTick(() => {
      form.value = response.data
    })
  })
}
/** 提交按钮 */
function submitForm() {
  configFormRef.value.validate((valid:any) => {
    if (valid) {
      if (form.value.configId != undefined) {
        updateConfig(form.value).then(() => {
          proxy.$modal.msgSuccess('修改成功')
          closeDialog()
          getList()
        })
      } else {
        addConfig(form.value).then(() => {
          proxy.$modal.msgSuccess('新增成功')
          closeDialog()
          getList()
        })
      }
    }
  })
}
/** 删除按钮操作 */
function handleDelete(row: { [key: string]: any }) {
  let configIds = []
  if (row.configId !== undefined) {
    configIds.push(row.configId)
  } else {
    configIds = ids.value
  }
  proxy.$modal
    .confirm('是否确认删除参数编号为"' + configIds + '"的数据项？')
    .then(() => {
      return delConfig({ ids: configIds })
    })
    .then(() => {
      getList()
      proxy.$modal.msgSuccess('删除成功')
    })
    .catch(() => {})
}
/** 导出按钮操作 */
function handleExport() {
  proxy.download(
    'system/config/export',
    {
      ...queryParams.value
    },
    `config_${new Date().getTime()}.xlsx`
  )
}
/** 刷新缓存按钮操作 */
function handleRefreshCache() {
  refreshCache().then(() => {
    proxy.$modal.msgSuccess('刷新缓存成功')
  })
}

onMounted(() => {
  getList()
})
</script>
