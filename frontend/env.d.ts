/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_URL: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<Record<string, never>, Record<string, never>, unknown>
  export default component
}


declare module '@jamescoyle/vue-icon' {
  import type { DefineComponent } from 'vue'
  const SvgIcon: DefineComponent<{
    type?: string
    path?: string
    size?: number | string
    viewbox?: string
    flip?: string
    rotate?: number
  }>
  export default SvgIcon
}
