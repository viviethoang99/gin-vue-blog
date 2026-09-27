// Login type options
export const loginTypeOptions = [
  { label: 'Email', value: 1 },
  { label: 'QQ', value: 2 },
  { label: 'Weibo', value: 3 },
]

export const loginTypeMap = {
  1: { name: 'Email', tag: 'success' },
  2: { name: 'QQ', tag: 'info' },
  3: { name: 'Weibo', tag: 'warning' },
}

// Article type options
export const articleTypeOptions = [
  { label: 'Original', value: 1 },
  { label: 'Repost', value: 2 },
  { label: 'Translation', value: 3 },
]

export const articleTypeMap = {
  1: { name: 'Original', tag: 'error' },
  2: { name: 'Repost', tag: 'success' },
  3: { name: 'Translation', tag: 'warning' },
}

// Comment type options
export const commentTypeOptions = [
  { label: 'Article', value: 1 },
  { label: 'Friend Link', value: 2 },
  { label: 'Talk', value: 3 },
]

export const commentTypeMap = {
  1: { name: 'Article', tag: 'info' },
  2: { name: 'Friend Link', tag: 'warning' },
  3: { name: 'Talk', tag: 'error' },
}
