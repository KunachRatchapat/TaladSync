// วันนี้แบบ local (YYYY-MM-DD) — เลี่ยง toISOString ที่เป็น UTC แล้วเพี้ยนข้ามวัน
export function todayLocal(): string {
  const d = new Date()
  const local = new Date(d.getTime() - d.getTimezoneOffset() * 60000)
  return local.toISOString().slice(0, 10)
}

// YYYY-MM-DD -> "ศุกร์ที่ 10 กรกฎาคม 2569" (พ.ศ. อัตโนมัติจาก locale ไทย)
export function formatThaiDate(iso: string): string {
  const d = new Date(`${iso}T00:00:00`)
  return d.toLocaleDateString('th-TH', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  })
}
