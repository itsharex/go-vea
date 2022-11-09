const wsUrl = 'wss://baidu.com/ws' // websocket 默认连接地址
let websocket: any // 用于存储实例化后websocket
let needReconnect = false // 连接成功过后需要设置发送间隔
let websocketTimeout = 0 // 保存重连延迟函数
//设置心跳配置,定时发送信息，确保连接成功
const heartCheck = {
  timeout: 3000, // 定时发送socket
  timeoutSend: 0, // 发送socket延迟函数
  serverTimeoutNumber: 0, //延迟关闭连接
  reset() {
    this.clear()
    this.start()
  },
  start() {
    // socket连接发送
    const self = this
    this.timeoutSend = window.setTimeout(() => {
      if (websocket) {
        websocket.send(JSON.stringify({ action: 'ping', type: 'ping' }))
        self.serverTimeoutNumber = window.setTimeout(() => {
          if (websocket) {
            websocket.close()
          }
        }, self.timeout)
      }
    }, this.timeout)
  },
  clear() {
    //清除定时器
    clearTimeout(this.timeoutSend)
    clearTimeout(this.serverTimeoutNumber)
  }
}

/**
 * 连接websocket
 * @param url 连接地址
 * @returns
 */
export function openSocket(url: string = '') {
  if (!url) {
    url = wsUrl
  }
  const user_id = '' // 当前用户信息，后台配置
  const session_id = '' // 获取验证信息，后台配置
  if (!user_id || !session_id) {
    notification.error({ message: '重新登录' })
    return
  }
  const connectURL = `${url}?session_id=${session_id}&user_id=${user_id}`
  websocket = null
  needReconnect = true
  // 监听websocket
  watchSocket(connectURL)
}

/**
 * 连接socket并且监听socket
 * @param url
 */
function watchSocket(url: string) {
  const client = getWebsocket(url)

  // 消息接收
  client.onmessage = (data: any) => {
    heartCheck.reset()
    console.log('-onmessage-收到的消息为--', data)
  }

  // 连接已准备好
  client.onopen = (data: any) => {
    heartCheck.start()
    console.log('-onopen-连接已准备好--', data)
    if (websocketTimeout) {
      console.log('--websocketTimeout--', websocketTimeout)
      clearTimeout(websocketTimeout)
    }
  }

  // 关闭连接
  client.onclose = (data: any) => {
    console.log('-onclose-关闭连接--', data)
    reconnect()
  }

  // 错误处理
  client.onerror = (data: any) => {
    console.log('-onerror-错误处理--', data)
  }
}

/**
 * 创建socket实例
 * @param url
 */
function getWebsocket(url: string) {
  if (!websocket) {
    websocket = new WebSocket(url)
  }
  return websocket
}

/**
 * 重新连接
 */
function reconnect() {
  if (needReconnect) {
    websocketTimeout = window.setTimeout(() => {
      console.log('重新连接')
      openSocket()
    }, 1000)
  }
}

/**
 * 发送数据
 * @param data
 * @returns
 */
function sendMessage(data) {
  // 发送数据
  if (websocket) {
    websocket.send(JSON.stringify({ action: 'ping', type: 'ping', data }))
  }
}

/**
 * 处理收到的数据
 * @param data
 * @returns
 */
function receiveMessage() {
  // 处理内容
}
