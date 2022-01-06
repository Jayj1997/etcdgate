/*
 * @Author       : jayj
 * @Date         : 2022-01-06 14:40:54
 * @Description  :
 * @LastEditors  : jayj
 * @LastEditTime : 2022-01-06 15:16:57
 */

import req from "./interceptors";
import qs from "qs";

const base = "/v3";

// Auth
// addr: address
// un: username
// pwd: password
function Auth(addr, un, pwd) {
  return req.post(
    base + "/auth",
    qs.stringify({
      address: addr,
      username: un,
      password: pwd,
    }),
    {
      headers: { "Content-Type": "application/x-www-form-urlencoded" },
    }
  );
}

export { Auth };
