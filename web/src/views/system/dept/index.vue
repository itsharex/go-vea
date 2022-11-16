<template>
  <div class="app-container">
    <el-form :model="queryParams" ref="queryFormRef" :inline="true" v-show="showSearch">
      <el-form-item label="部门名称" prop="deptName">
        <el-input v-model="queryParams.deptName" placeholder="请输入部门名称" clearable @keyup.enter="handleQuery" />
      </el-form-item>
      <el-form-item label="状态" prop="status">
        <el-select v-model="queryParams.status" placeholder="部门状态" clearable>
          <el-option v-for="dict in sys_normal_disable" :key="dict.value" :label="dict.label" :value="dict.value" />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="handleQuery"> <ep:search /> 搜索 </el-button>
        <el-button @click="resetQuery"> <ep:refresh /> 重置 </el-button>
      </el-form-item>
    </el-form>

    <el-row :gutter="10" class="mb8">
      <el-col :span="1.5">
        <el-button type="primary" plain @click="handleAdd" v-hasPerms="['system:dept:add']"> <ep:plus /> 新增 </el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button type="info" plain @click="toggleExpandAll"> <ep:sort /> 展开/折叠 </el-button>
      </el-col>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList"></right-toolbar>
    </el-row>

    <el-table
      border
      v-if="refreshTable"
      v-loading="loading"
      :data="deptList"
      row-key="deptId"
      :default-expand-all="isExpandAll"
      :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
    >
      <el-table-column prop="deptName" label="部门名称" width="260"></el-table-column>
      <el-table-column prop="orderNum" label="排序" width="200"></el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="scope">
          <dict-tag :options="sys_normal_disable" :value="scope.row.status" />
        </template>
      </el-table-column>
      <el-table-column label="创建时间" align="center" prop="createTime" width="200">
        <template #default="scope">
          <span>{{ parseTime(scope.row.createTime) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" class-name="small-padding fixed-width">
        <template #default="scope">
          <el-button link type="primary" @click="handleUpdate(scope.row)" v-hasPerms="['system:dept:edit']"> <ep:edit /> 修改 </el-button>
          <el-button link type="primary" @click="handleAdd(scope.row)" v-hasPerms="['system:dept:add']"> <ep:plus /> 新增 </el-button>
          <el-button link type="danger" v-if="scope.row.parentId != 0" @click="handleDelete(scope.row)" v-hasPerms="['system:dept:remove']">
            <ep:delete /> 删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 添加或修改部门对话框 -->
    <el-dialog :title="dialog.title" v-model="dialog.visible" width="600px" @close="closeDialog" append-to-body>
      <el-form ref="deptFormRef" :model="form" :rules="rules" label-width="80px">
        <el-row>
          <el-col :span="24" v-if="form.parentId !== 0">
            <el-form-item label="上级部门" prop="parentId">
              <el-tree-select
                style="width: 100%"
                v-model="form.parentId"
                :data="deptOptions"
                :props="{ value: 'deptId', label: 'deptName', children: 'children' }"
                value-key="deptId"
                placeholder="选择上级部门"
                check-strictly
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="部门名称" prop="deptName">
              <el-input v-model="form.deptName" placeholder="请输入部门名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="显示排序" prop="orderNum">
              <el-input-number style="width: 100%" v-model="form.orderNum" controls-position="right" :min="0" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="负责人" prop="leader">
              <el-input v-model="form.leader" placeholder="请输入负责人" maxlength="20" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="联系电话" prop="phone">
              <el-input v-model="form.phone" placeholder="请输入联系电话" maxlength="11" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="form.email" placeholder="请输入邮箱" maxlength="50" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="部门状态">
              <el-radio-group v-model="form.status">
                <el-radio v-for="dict in sys_normal_disable" :key="dict.value" :label="dict.value">{{ dict.label }}</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
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

<script lang="ts" setup name="Dept">
import { listDeptTree, getDept, delDept, addDept, updateDept, listDeptExcludeChild } from '@/api/system/dept'
import useCurrentInstance from '@/hooks/useCurrentInstance'
import { DeptFormData } from '@/types/api/dept'
import { Dialog } from '@/types/common'

const { proxy } = useCurrentInstance()
const { sys_normal_disable } = proxy.useDict('sys_normal_disable')

const deptList = ref([])
const loading = ref(true)
const showSearch = ref(true)
const deptOptions = ref([])
const isExpandAll = ref(true)
const refreshTable = ref(true)

const queryFormRef = ref<ElForm>(null)
const deptFormRef = ref<ElForm>(null)

const data = reactive({
  dialog: {
    visible: false,
    title: ''
  } as Dialog,
  form: {
    orderNum: 0,
    status: '0'
  } as DeptFormData,
  queryParams: {
    deptName: undefined,
    status: undefined
  },
  rules: {
    parentId: [{ required: true, message: '上级部门不能为空', trigger: 'blur' }],
    deptName: [{ required: true, message: '部门名称不能为空', trigger: 'blur' }],
    orderNum: [{ required: true, message: '显示排序不能为空', trigger: 'blur' }],
    email: [{ type: 'email', message: '请输入正确的邮箱地址', trigger: ['blur', 'change'] }],
    phone: [{ pattern: /^1[3|4|5|6|7|8|9][0-9]\d{8}$/, message: '请输入正确的手机号码', trigger: 'blur' }]
  }
})

const { dialog, queryParams, form, rules } = toRefs(data)

/** 查询部门列表 */
function getList() {
  loading.value = true
  listDeptTree(queryParams.value).then(response => {
    deptList.value = response.data
    loading.value = false
  })
}
/** 关闭弹窗 */
function closeDialog() {
  dialog.value.visible =  false
  deptFormRef.value?.resetFields()
  deptFormRef.value?.clearValidate()
}
/** 搜索按钮操作 */
function handleQuery() {
  getList()
}
/** 重置按钮操作 */
function resetQuery() {
  queryFormRef.value?.resetFields()
  handleQuery()
}
/** 新增按钮操作 */
function handleAdd(row: { [key: string]: any }) {
  listDeptTree(queryParams.value).then(response => {
    deptOptions.value = response.data
  })
  if (row != undefined) {
    form.value.parentId = row.deptId
  }
  dialog.value.visible = true
  dialog.value.title = '添加部门'
}
/** 展开/折叠操作 */
function toggleExpandAll() {
  refreshTable.value = false
  isExpandAll.value = !isExpandAll.value
  nextTick(() => {
    refreshTable.value = true
  })
}
/** 修改按钮操作 */
function handleUpdate(row: { [key: string]: any }) {
  listDeptExcludeChild({ExcludeDeptId: row.deptId}).then(response => {
    deptOptions.value = response.data
  })
  getDept(row.deptId).then(response => {
    form.value = response.data
    dialog.value.visible = true
    dialog.value.title = '修改部门'
  })
}
/** 提交按钮 */
function submitForm() {
  deptFormRef.value.validate((valid:any) => {
    if (valid) {
      if (form.value.deptId != undefined) {
        updateDept(form.value).then(() => {
          proxy.$modal.msgSuccess('修改成功')
          closeDialog()
          getList()
        })
      } else {
        addDept(form.value).then(() => {
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
  proxy.$modal
    .confirm('是否确认删除名称为"' + row.deptName + '"的数据项?')
    .then(function () {
      return delDept(row.deptId)
    })
    .then(() => {
      getList()
      proxy.$modal.msgSuccess('删除成功')
    })
    .catch(() => {})
}

onMounted(() =>{
  getList()
})
</script>
