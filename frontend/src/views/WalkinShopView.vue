<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import SvgIcon from '@jamescoyle/vue-icon'
import { mdiAccountMultiple, mdiInformation } from '@mdi/js'

import FormInput from '@/components/FormInput.vue'
import FormSelect from '@/components/FormSelect.vue'
import PageHeader from '@/components/PageHeader.vue'
import StateBlock from '@/components/StateBlock.vue'
import { PRODUCT_TYPES, fetchZones, saveShopProfile } from '@/api/shop'
import { useAuthStore } from '@/stores/auth'
import type { ShopProfileForm } from '@/types'

const router = useRouter()
const auth = useAuthStore()

const form = reactive<ShopProfileForm>({
  shopName: '',
  zone: '',
  productType: '',
  ownerName: '',
  phone: '',
  lineId: '',
})
const zones = ref<string[]>([])
const submitting = ref(false)
const errorMessage = ref('')

const phoneError = computed(() =>
  form.phone !== '' && !/^0\d{9}$/.test(form.phone) ? 'เบอร์โทรต้องเป็นตัวเลข 10 หลัก ขึ้นต้นด้วย 0' : '',
)

const canSubmit = computed(
  () =>
    form.shopName.trim() !== '' &&
    form.zone !== '' &&
    form.productType !== '' &&
    form.ownerName.trim() !== '' &&
    form.phone !== '' &&
    phoneError.value === '' &&
    !submitting.value,
)

async function load() {
  try {
    await auth.ensureRole('WALKIN')
  } catch {
    return
  }
  form.ownerName = auth.user?.name ?? ''
  form.phone = auth.user?.phone ?? ''
  try {
    zones.value = await fetchZones()
  } catch {
    errorMessage.value = 'โหลดรายการโซนไม่สำเร็จ ลองใหม่อีกครั้ง'
  }
}

async function onSubmit() {
  submitting.value = true
  errorMessage.value = ''
  try {
    await saveShopProfile({
      ...form,
      shopName: form.shopName.trim(),
      ownerName: form.ownerName.trim(),
      lineId: form.lineId.trim(),
    })
    router.push({ name: 'walkin' })
  } catch {
    errorMessage.value = 'บันทึกไม่สำเร็จ ลองใหม่อีกครั้ง'
  } finally {
    submitting.value = false
  }
}

function onCancel() {
  router.push({ name: 'walkin' })
}

onMounted(load)
</script>

<template>
  <div class="min-h-screen bg-brand-50">
    <PageHeader title="TaladSync · ตลาดเทศบาล" :show-back="false" />

    <main class="mx-auto max-w-app px-4 py-6">
      <StateBlock
        v-if="auth.loading || auth.error"
        :state="auth.loading ? 'loading' : 'error'"
        :message="auth.error"
        @retry="load"
      />

      <form
        v-else
        class="flex flex-col gap-4 rounded-3xl bg-white p-5 shadow-sm"
        @submit.prevent="onSubmit"
      >
        <h2 class="flex items-center gap-2 text-base font-semibold text-slate-600">
          <SvgIcon type="mdi" :path="mdiInformation" :size="24" class="text-emerald-500" />
          เพิ่มข้อมูลร้านค้าของคุณ
        </h2>

        <FormInput v-model="form.shopName" label="ชื่อร้านค้า" placeholder="เช่น ร้านผักป้าแดง" />
        <FormSelect v-model="form.zone" label="โซนจำหน่ายสินค้า" placeholder="เลือกโซน" :options="zones" />
        <FormSelect
          v-model="form.productType"
          label="ประเภทสินค้า"
          placeholder="เลือกประเภทสินค้า"
          :options="PRODUCT_TYPES"
        />

        <h2 class="mt-2 flex items-center gap-2 text-base font-semibold text-slate-600">
          <SvgIcon type="mdi" :path="mdiAccountMultiple" :size="24" class="text-emerald-500" />
          ข้อมูลติดต่อ
        </h2>

        <FormInput v-model="form.ownerName" label="ชื่อผู้ค้า" placeholder="ชื่อ-นามสกุล" />
        <FormInput
          v-model="form.phone"
          label="เบอร์โทรศัพท์"
          placeholder="0812345678"
          inputmode="numeric"
          :maxlength="10"
          :error="phoneError"
        />
        <FormInput v-model="form.lineId" label="LINE ID" placeholder="ไม่บังคับ" :maxlength="50" />

        <p
          v-if="errorMessage"
          class="rounded-xl bg-danger-50 px-4 py-3 text-sm font-medium text-danger-700"
          role="alert"
        >
          {{ errorMessage }}
        </p>

        <div class="mt-2 grid grid-cols-2 gap-4">
          <button
            type="button"
            class="h-11 rounded-2xl bg-slate-300 text-lg font-bold text-slate-500 active:bg-slate-400"
            @click="onCancel"
          >
            ยกเลิก
          </button>
          <button
            type="submit"
            :disabled="!canSubmit"
            class="h-11 rounded-2xl bg-emerald-500 text-lg font-bold text-white active:bg-emerald-600 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {{ submitting ? 'กำลังบันทึก...' : 'บันทึก' }}
          </button>
        </div>
      </form>
    </main>
  </div>
</template>
