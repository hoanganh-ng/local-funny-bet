import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUiStore = defineStore('ui', () => {
  const showSidebarOnMobile = ref(false)

  function toggleMobileNav() {
    showSidebarOnMobile.value = !showSidebarOnMobile.value
  }

  function closeMobileSidebar() {
    showSidebarOnMobile.value = false
  }

  return {
    showSidebarOnMobile,
    toggleMobileNav,
    closeMobileSidebar
  }
})
