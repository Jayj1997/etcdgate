/*
 * @Author       : jayj
 * @Date         : 2021-12-15 16:51:54
 * @Description  :
 * @LastEditors  : jayj
 * @LastEditTime : 2021-12-15 16:58:19
 */
import req from './interceptors'
import qs from 'qs'

const base = '/v3'

// Auth
// addr: address
// un: username
// pwd: password
function Auth (addr, un, pwd) {
  return req.post(base + '/auth',
    qs.stringify({
      address: addr,
      username: un,
      password: pwd
    }),
    {
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' }
    })
}

// Get
// key:
// rev: get reversion of key
function Get (key, rev) {
  return req.post(base + '/get',
    qs.stringify({ key, rev }),
    {
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' }
    })
}

// Put
// key
// val
function Put (key, val) {
  return req.post(base + '/put',
    qs.stringify({ key, val }),
    {
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' }
    })
}

// Delete
// key: target key
// dir: is delete key/*
function Delete (key, dir) {
  return req.post(base + '/del',
    qs.stringify({ key, dir }),
    {
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' }
    })
}

function Directory () {
  return req.get(base + '/directory')
}

export { Auth, Get, Put, Delete, Directory }
