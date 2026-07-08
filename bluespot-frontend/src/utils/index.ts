
// toFixed 修复版： 修复 1.335.toFixed(2) => 1.33
export function toFixed (val: number, s = 2) {
  const num = val || 0
  if (/^\d+$/.test(String(num))) { // 整数
    return num.toFixed(s)
  }
  const times = Math.pow(10, s)
  let des = num * times + 0.5
  des = parseInt(String(des), 10) / times
  return des + ''
}