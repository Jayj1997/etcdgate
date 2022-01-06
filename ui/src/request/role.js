/*
 * @Author       : jayj
 * @Date         : 2022-01-06 15:06:40
 * @Description  :
 * @LastEditors  : jayj
 * @LastEditTime : 2022-01-06 15:16:16
 */

import req from "./interceptors";
import qs from "qs";

const base = "/v3";

// Roles
// get roles
function Roles() {
  return req.get(base + "/roles");
}

// RolePerm
// get role's permission
function RolePerm(name) {
  return req.get(base + "/role/" + name);
}

// RoleAdd
// create a role
function RoleAdd(name) {
  return req.get(base + "/role_add/" + name);
}

// RoleDelete
// delete a role
function RoleDelete(name) {
  return req.get(base + "/role_delete/" + name);
}

// RoleGrant
// grant keys to a role
// role_name: role name
// key: grant key
// range_end: grant key range end
// type: 0 read 1 write 2 read_write
// f.e. test, aa, ab, 2
// aa->ab mean's aa* til ab. [aa, aaa,aab, ... ,aaz,ab)
function RoleGrant(role_name, key, range_end, type) {
  return req.post(
    base + "role_grant",
    qs.stringify({ role_name, key, range_end, type }),
    {
      headers: { "Content-Type": "application/x-www-form-urlencoded" },
    }
  );
}

// RoleRevoke
// revoke keys to a role
// role_name: role name
// key: grant key
// range_end: grant key range end
// f.e. test, aa, ab
// aa->ab mean's aa* til ab. [aa, aaa,aab, ... ,aaz,ab)
function RoleRevoke(role_name, key, range_end) {
  return req.post(
    base + "role_revoke",
    qs.stringify({ role_name, key, range_end }),
    {
      headers: { "Content-Type": "application/x-www-form-urlencoded" },
    }
  );
}

export { Roles, RolePerm, RoleAdd, RoleDelete, RoleGrant, RoleRevoke };
