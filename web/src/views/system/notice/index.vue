<template>
  <div class="app-container">
    <el-form :model="queryParams" ref="queryFormRef" :inline="true" v-show="showSearch" label-width="68px">
      <el-form-item label="公告标题" prop="noticeTitle">
        <el-input v-model="queryParams.noticeTitle" placeholder="请输入公告标题" clearable @keyup.enter="handleQuery" />
      </el-form-item>
      <el-form-item label="操作人员" prop="createBy">
        <el-input v-model="queryParams.createBy" placeholder="请输入操作人员" clearable @keyup.enter="handleQuery" />
      </el-form-item>
      <el-form-item label="类型" prop="noticeType">
        <el-select v-model="queryParams.noticeType" placeholder="公告类型" clearable>
          <el-option v-for="dict in sys_notice_type" :key="dict.value" :label="dict.label" :value="dict.value" />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="handleQuery"> <ep:search /> 搜索 </el-button>
        <el-button @click="resetQuery"> <ep:refresh /> 重置 </el-button>
      </el-form-item>
    </el-form>

    <el-row :gutter="10" class="mb8">
      <el-col :span="1.5">
        <el-button type="primary" plain @click="handleAdd" v-hasPerms="['system:notice:add']"><ep:plus /> 新增</el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button type="danger" plain :disabled="multiple" @click="handleDelete" v-hasPerms="['system:notice:remove']"><ep:delete />删除</el-button>
      </el-col>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList"></right-toolbar>
    </el-row>

    <el-table border v-loading="loading" :data="noticeList" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55" align="center" />
      <el-table-column label="序号" align="center" prop="noticeId" width="100" />
      <el-table-column label="公告标题" align="center" prop="noticeTitle" :show-overflow-tooltip="true" />
      <el-table-column label="公告类型" align="center" prop="noticeType" width="100">
        <template #default="scope">
          <dict-tag :options="sys_notice_type" :value="scope.row.noticeType" />
        </template>
      </el-table-column>
      <el-table-column label="状态" align="center" prop="status" width="100">
        <template #default="scope">
          <dict-tag :options="sys_notice_status" :value="scope.row.status" />
        </template>
      </el-table-column>
      <el-table-column label="创建者" align="center" prop="createBy" width="100" />
      <el-table-column label="创建时间" align="center" prop="createTime" width="100">
        <template #default="scope">
          <span>{{ parseTime(scope.row.createTime, '{y}-{m}-{d}') }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" class-name="small-padding fixed-width">
        <template #default="scope">
          <el-button link type="primary" @click="handleUpdate(scope.row)" v-hasPerms="['system:notice:edit']"><ep:edit />修改</el-button>
          <el-button link type="danger" @click="handleDelete(scope.row)" v-hasPerms="['system:notice:remove']"><ep:delete />删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total > 0" :total="total" v-model:page="queryParams.pageNum" v-model:limit="queryParams.pageSize" @pagination="getList" />

    <!-- 添加或修改公告对话框 -->
    <el-dialog :title="dialog.title" v-model="dialog.visible" width="780px" append-to-body>
      <el-form ref="noticeFormRef" :model="form" :rules="rules" label-width="80px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="公告标题" prop="noticeTitle">
              <el-input v-model="form.noticeTitle" placeholder="请输入公告标题" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="公告类型" prop="noticeType">
              <el-select v-model="form.noticeType" placeholder="请选择">
                <el-option v-for="dict in sys_notice_type" :key="dict.value" :label="dict.label" :value="dict.value"></el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="状态">
              <el-radio-group v-model="form.status">
                <el-radio v-for="dict in sys_notice_status" :key="dict.value" :label="dict.value">{{ dict.label }}</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="内容">
              <el-input :rows="6" type="textarea" placeholder="请输入内容" v-model="form.noticeContent" />
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

<script lang="ts" setup name="Notice">
import { listNotice, getNotice, delNotice, addNotice, updateNotice } from '@/api/system/notice'
import useCurrentInstance from '@/hooks/useCurrentInstance'
import { NoticeFormData } from '@/types/api/notice'
import { Dialog } from '@/types/common'

const { proxy } = useCurrentInstance()
const { sys_notice_status, sys_notice_type } = proxy.useDict('sys_notice_status', 'sys_notice_type')

const noticeList = ref([])
const loading = ref(true)
const showSearch = ref(true)
const ids = ref([])
const single = ref(true)
const multiple = ref(true)
const total = ref(0)

const queryFormRef = ref<ElForm>(null)
const noticeFormRef = ref<ElForm>(null)

const data = reactive({
  dialog: {
    visible: false,
    title: ''
  } as Dialog,
  form: {
    status: '0'
  } as NoticeFormData,
  queryParams: {
    pageNum: 1,
    pageSize: 10,
    noticeTitle: undefined,
    createBy: undefined,
    status: undefined
  },
  rules: {
    noticeTitle: [{ required: true, message: '公告标题不能为空', trigger: 'blur' }],
    noticeType: [{ required: true, message: '公告类型不能为空', trigger: 'change' }]
  }
})

const { dialog, queryParams, form, rules } = toRefs(data)

/** 查询公告列表 */
function getList() {
  loading.value = true
  listNotice(queryParams.value).then(response => {
    noticeList.value = response.data.rows
    total.value = response.data.total
    loading.value = false
  })
}
/** 关闭弹窗 */
function closeDialog() {
  dialog.value.visible = false
  noticeFormRef.value?.resetFields()
  noticeFormRef.value?.clearValidate()
  form.value.noticeId = undefined
}
/** 搜索按钮操作 */
function handleQuery() {
  queryParams.value.pageNum = 1
  getList()
}
/** 重置按钮操作 */
function resetQuery() {
  queryFormRef.value.resetFields()
  handleQuery()
}
/** 多选框选中数据 */
function handleSelectionChange(selection:any) {
  ids.value = selection.map((item:any) => item.noticeId)
  single.value = selection.length != 1
  multiple.value = !selection.length
}
/** 新增按钮操作 */
function handleAdd() {
  dialog.value.visible = true
  dialog.value.title = '添加公告'
}
/**修改按钮操作 */
function handleUpdate(row: { [key: string]: any }) {
  const noticeId = row.noticeId || ids.value
  getNotice(noticeId).then(response => {
    dialog.value.visible = true
    dialog.value.title = '修改公告'
    nextTick(() => {
      form.value = response.data
    })
  })
}
/** 提交按钮 */
function submitForm() {
  noticeFormRef.value.validate((valid:any) => {
    if (valid) {
      if (form.value.noticeId != undefined) {
        updateNotice(form.value).then(() => {
          proxy.$modal.msgSuccess('修改成功')
          closeDialog()
          getList()
        })
      } else {
        addNotice(form.value).then(() => {
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
  let noticeIds = []
  if (row.noticeId !== undefined) {
    noticeIds.push(row.noticeId)
  } else {
    noticeIds = ids.value
  }
  proxy.$modal
    .confirm('是否确认删除公告编号为"' + noticeIds + '"的数据项？')
    .then(function () {
      return delNotice({ids: noticeIds})
    })
    .then(() => {
      getList()
      proxy.$modal.msgSuccess('删除成功')
    })
    .catch(() => {})
}

onMounted(() => {
  getList()
})
</script>
