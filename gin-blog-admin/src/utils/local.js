const CryptoSecret = '__SecretKey__'

/**
 * Store serialized data to LocalStorage
 * @param {string} key
 * @param {any} value Object needs to be serialized
 * @param {number} expire 过期时间(秒), 传 0 表示不过期
 */
export function setLocal(key, value, expire = 60 * 60 * 24 * 7) {
  const data = JSON.stringify({
    value,
    time: Date.now(),
    expire: expire ? Date.now() + expire * 1000 : null,
  })
  window.localStorage.setItem(key, encrypto(data)) // Encrypted storage
}

/**
 * Get data from LocalStorage, decrypt and deserialize, return based on expiration
 * @param {string} key
 */
export function getLocal(key) {
  const encryptedVal = window.localStorage.getItem(key)
  const val = encryptedVal ? decrypto(encryptedVal) : null
  // 解密或反序列化失败(数据被改坏、换了加密方式)时当作没有数据
  if (val) {
    const { value, expire } = val
    // Return if not expired
    if (!expire || expire > Date.now()) {
      return value
    }
  }
  // 过期或数据不可用则移除
  removeLocal(key)
  return null
}

export function removeLocal(key) {
  window.localStorage.removeItem(key)
}

/**
 * Encrypt data: Base64 encryption
 * btoa 只接受 Latin1, 中文密码等字符会直接抛错, 先转成百分号编码
 * @param {string} data
 */
function encrypto(data) {
  return btoa(encodeURIComponent(CryptoSecret + data))
}

/**
 * 解密数据: Base64 解密, 失败时返回 null 而不是抛错
 * @param {string} cipherText
 */
function decrypto(cipherText) {
  try {
    const originalText = decodeURIComponent(atob(cipherText)).replace(CryptoSecret, '')
    return JSON.parse(originalText)
  }
  catch {
    return null
  }
}
