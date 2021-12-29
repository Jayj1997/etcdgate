/*
 * @Author       : jayj
 * @Date         : 2021-12-15 16:43:58
 * @Description  : request middleware
 * @LastEditors  : jayj
 * @LastEditTime : 2021-12-15 17:10:51
 */

import axios from 'axios'

const req = axios.create()

function getToken () {
  return localStorage.getItem('token')
}

req.interceptors.request.use(
  config => {
    config.headers.authorization = getToken()
    return config
  },
  err => {
    return Promise.reject(err)
  }
)

req.interceptors.response.use(
  resp => {
    if (resp.data.code) {
      switch (resp.data.code) {
        case 1006: { // token outdate
          // warn user & reInput info

          // const token = getToken()
          // Refresh(token).then(resp => {
          //   localStorage.setItem('token', resp.data.data)
          // })
          break
        }
      }
    }

    return resp
  },
  error => {
    const code = error.response.data.code
    switch (code) {
      case 1005: { // invalid token
        // warn user & reInput info
        localStorage.removeItem('token')
        break
      }
    }
    return Promise.reject(error.response.status)
  }
)

export default req
