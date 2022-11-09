export default [
  {
    url: '/blog/article/list',
    method: 'get',
    response: () => {
      return {
        code: 200,
        message: 'ok',
        data: ['tom', 'jerry']
      }
    }
  }
]
