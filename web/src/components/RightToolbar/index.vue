<template>
  <div class="top-right-btn" :style="style">
    <el-row>
      <el-button-group>
        <el-tooltip class="item" effect="dark" :content="showSearch ? '隐藏搜索' : '显示搜索'" placement="top" v-if="search">
          <el-button @click="toggleSearch()">
            <ep:search />
          </el-button>
        </el-tooltip>
        <el-tooltip class="item" effect="dark" content="刷新" placement="top">
          <el-button @click="refresh()">
            <ep:refresh />
          </el-button>
        </el-tooltip>
        <el-tooltip class="item" effect="dark" content="显隐列" placement="top" v-if="columns">
          <el-button @click="showColumn()">
            <ep:menu />
          </el-button>
        </el-tooltip>
      </el-button-group>
    </el-row>
    <el-dialog :title="title" v-model="open" append-to-body>
      <div class="dialog-transfer">
        <el-transfer :titles="['显示', '隐藏']" v-model="value" :data="columns" @change="dataChange"></el-transfer>
      </div>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
const props = defineProps({
  showSearch: {
    type: Boolean,
    default: true
  },
  columns: {
    type: Array
  },
  search: {
    type: Boolean,
    default: true
  },
  gutter: {
    type: Number,
    default: 10
  }
})

const emits = defineEmits(['update:showSearch', 'queryTable'])

// 显隐数据
const value = ref([])
// 弹出层标题
const title = ref('显示/隐藏')
// 是否显示弹出层
const open = ref(false)

const style = computed(() => {
  const ret = {}
  if (props.gutter) {
    ret.marginRight = `${props.gutter / 2}px`
  }
  return ret
})

// 搜索
function toggleSearch() {
  emits('update:showSearch', !props.showSearch)
}

// 刷新
function refresh() {
  emits('queryTable')
}

// 右侧列表元素变化
function dataChange(data) {
  for (let item in props.columns) {
    const key = props.columns[item].key
    props.columns[item].visible = !data.includes(key)
  }
}

// 打开显隐列dialog
function showColumn() {
  open.value = true
}

// 显隐列初始默认隐藏列
for (let item in props.columns) {
  if (props.columns[item].visible === false) {
    value.value.push(parseInt(item))
  }
}
</script>

<style lang="scss" scoped>
:deep(.el-transfer__button) {
  border-radius: 50%;
  display: block;
  margin-left: 0px;
}
:deep(.el-transfer__button:first-child) {
  margin-bottom: 10px;
}

.dialog-transfer {
  text-align: center;
}
</style>
