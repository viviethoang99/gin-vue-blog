const CryptoSecret = '__SecretKey__'

/**
 * Store serialized data to LocalStorage
 * @param {string} key
 * @param {any} value Object needs to be serialized
 * @param {number} expire
 */
export function setLocal(key, value, expire = 60 * 60 * 24 * 7) {
  const data = JSON.stringify({
    value,
    time: Date.now(),
    expire: expire ? new Date().getTime() + expire * 1000 : null,
  })
  window.localStorage.setItem(key, encrypto(data)) // Encrypted storage
}

/**
 * Get data from LocalStorage, decrypt and deserialize, return based on expiration
 * @param {string} key
 */
export function getLocal(key) {
  const encryptedVal = window.localStorage.getItem(key)
  if (encryptedVal) {
    const val = decrypto(encryptedVal) // Decrypt
    const { value, expire } = JSON.parse(val)
    // Return if not expired
    if (!expire || expire > new Date().getTime()) {
      return value
    }
  }
  // Remove if expired
  removeLocal(key)
  return null
}

export function removeLocal(key) {
  window.localStorage.removeItem(key)
}

export function clearLocal() {
  window.localStorage.clear()
}

/**
 * Encrypt data: Base64 encryption
 * @param {any} data - Data
 */
function encrypto(data) {
  const newData = JSON.stringify(data)
  const encryptedData = btoa(CryptoSecret + newData)
  return encryptedData
}

/**
 * Decrypt data: Base64 decryption
 * @param {string} cipherText - Cipher text
 */
function decrypto(cipherText) {
  const decryptedData = atob(cipherText)
  const originalText = decryptedData.replace(CryptoSecret, '')
  try {
    const parsedData = JSON.parse(originalText)
    return parsedData
  }
  catch (error) {
    return null
  }
}
