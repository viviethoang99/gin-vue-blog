import { h } from 'vue'
import { Icon } from '@iconify/vue'
import { NIcon } from 'naive-ui'
import dayjs from 'dayjs'

export * from './http'
export * from './local'
export * from './naiveTool'

// Relative image path => full image URL, used for local file uploads
// If it contains http, it's a web image resource
// Otherwise it's an image on the server; prepend the server URL
const SERVER_URL = import.meta.env.VITE_SERVER_URL
export function convertImgUrl(imgUrl) {
  if (!imgUrl) {
    return 'http://dummyimage.com/400x400'
  }
  // Web resource
  if (imgUrl.startsWith('http')) {
    return imgUrl
  }
  return `${SERVER_URL}/${imgUrl}`
}

/**
 * Format date
 */
export function formatDate(date = undefined, format = 'YYYY-MM-DD') {
  return dayjs(date).format(format)
}

/**
 * Render icon with NIcon
 */
export function renderIcon(icon, props = { size: 12 }) {
  return () => h(NIcon, props, { default: () => h(Icon, { icon }) })
}

// Frontend file export: pass file content and name
export function downloadFile(content, fileName) {
  const aEle = document.createElement('a') // Create download link
  aEle.download = fileName // Set download filename
  aEle.style.display = 'none'// Hidden downloadable link
  // Convert string content to blob URL
  const blob = new Blob([content])
  aEle.href = URL.createObjectURL(blob)
  // Bind click event
  document.body.appendChild(aEle)
  aEle.click()
  // Then remove
  document.body.removeChild(aEle)
}
