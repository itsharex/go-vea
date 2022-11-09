<template>
  <div
    :class="{ 'has-logo': showLogo }"
    :style="{
      backgroundColor: sideTheme === 'theme-dark' ? variables.menuBackground : variables.menuLightBackground
    }"
  >
    <logo v-if="showLogo" :collapse="isCollapse" />
    <el-scrollbar :class="sideTheme" wrap-class="scrollbar-wrapper">
      <el-menu
        :default-active="activeMenu"
        :collapse="isCollapse"
        :background-color="sideTheme === 'theme-dark' ? variables.menuBackground : variables.menuLightBackground"
        :text-color="sideTheme === 'theme-dark' ? variables.menuColor : variables.menuLightColor"
        :unique-opened="true"
        :active-text-color="theme"
        :collapse-transition="false"
        mode="vertical"
      >
        <sidebar-item v-for="(route, index) in sidebarRouters" :key="route.path + index" :item="route" :base-path="route.path" />
      </el-menu>
    </el-scrollbar>
  </div>
</template>

<script lang="ts" setup>
import variables from '@/assets/styles/variables.module.scss'
import useStore from '@/store'

const { app, permission, setting } = useStore()
const route = useRoute()

const sidebarRouters = computed(() => permission.sidebarRouters)
const showLogo = computed(() => setting.sidebarLogo)
const sideTheme = computed(() => setting.sideTheme)
const theme = computed(() => setting.theme)
const isCollapse = computed(() => !app.sidebar.opened)

const activeMenu = computed(() => {
  const { meta, path } = route
  // if set path, the sidebar will highlight the path you set
  if (meta.activeMenu) {
    return meta.activeMenu
  }
  return path
})
</script>
