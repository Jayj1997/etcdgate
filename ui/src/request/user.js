/*
 * @Author       : jayj
 * @Date         : 2022-01-06 14:41:41
 * @Description  :
 * @LastEditors  : jayj
 * @LastEditTime : 2022-01-06 14:50:42
 */
import req from './interceptors'
import qs from 'qs'

const base = '/v3'

// Users
function Users () {
  return req.get(base + '/get')
}

// UserRole get user role
// username
function UserRole(username) {
  return req.get(base + '/user/' + username)
}

// UserAdd create a user
function UserAdd(name, pwd) {
  return req.post(base + '/user_add',
    qs.stringify({name, pwd}),
    {
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' }
    }
  )
}


// UserDelete
function UserDelete(name) {
  return req.get(base+'/user_del/'+name)
}

// UserGrant
// grant user a role
function UserGrant(name, role) {
  return req.post(base + '/user_grant',
  qs.stringify({name, role}),
  {
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' }
  })
}

// UserRevoke
// revoke user a role
function UserRevoke(name, role) {
  return req.post(base + '/user_revoke',
  qs.stringify({name, role}),
  {
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' }
  })
}

export {Users, UserRole, UserAdd,UserDelete,UserGrant,UserRevoke}
