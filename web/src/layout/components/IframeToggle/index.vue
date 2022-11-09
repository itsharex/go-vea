<template>
  <transition-group name="fade-transform" mode="out-in">
    <inner-link
      v-for="(item, index) in tagsView.iframeViews"
      :key="item.path"
      :iframeId="'iframe' + index"
      v-show="route.path === item.path"
      :src="item.meta?.link"
    ></inner-link>
  </transition-group>
</template>

<script lang="ts" setup>
import useStore from '@/store'
const { tagsView } = useStore()

const route = useRoute()
</script>
