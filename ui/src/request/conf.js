/*
 * @Author       : jayj
 * @Date         : 2021-12-15 16:51:54
 * @Description  :
 * @LastEditors  : jayj
 * @LastEditTime : 2022-01-06 15:17:03
 */
import req from "./interceptors";
import qs from "qs";

const base = "/v3";

// ConfGet
// key:
// rev: get reversion of key
function ConfGet(key, rev) {
  return req.post(base + "/get", qs.stringify({ key, rev }), {
    headers: { "Content-Type": "application/x-www-form-urlencoded" },
  });
}

// ConfPut
// key
// val
function ConfPut(key, val) {
  return req.post(base + "/put", qs.stringify({ key, val }), {
    headers: { "Content-Type": "application/x-www-form-urlencoded" },
  });
}

// ConfDelete
// key: target key
// dir: is delete key/*
function ConfDelete(key, dir) {
  return req.post(base + "/del", qs.stringify({ key, dir }), {
    headers: { "Content-Type": "application/x-www-form-urlencoded" },
  });
}

// Directory
function Directory() {
  return req.get(base + "/directory");
}

export { ConfGet, ConfPut, ConfDelete, Directory };
