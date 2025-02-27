<script setup lang="ts">
import { Modal, message } from 'ant-design-vue'
import type { Ref } from 'vue'
import type { Site } from '@/api/domain'
import domain from '@/api/domain'
import ChatGPT from '@/components/ChatGPT/ChatGPT.vue'
import { formatDateTime } from '@/lib/helper'
import Deploy from '@/views/domain/components/Deploy.vue'
import { useSettingsStore } from '@/pinia'
import type { ChatComplicationMessage } from '@/api/openai'
import type { CheckedType } from '@/types'

const settings = useSettingsStore()

const configText = inject('configText') as Ref<string>
const enabled = inject('enabled') as Ref<boolean>
const name = inject('name') as Ref<string>
const filepath = inject('filepath') as Ref<string>
const history_chatgpt_record = inject('history_chatgpt_record') as Ref<ChatComplicationMessage[]>
const filename = inject('filename') as Ref<string | number | undefined>
const data = inject('data') as Ref<Site>

const [modal, ContextHolder] = Modal.useModal()

const active_key = ref(['1', '2', '3'])

function enable() {
  domain.enable(name.value).then(() => {
    message.success($gettext('Enabled successfully'))
    enabled.value = true
  }).catch(r => {
    message.error($gettext('Failed to enable %{msg}', { msg: r.message ?? '' }), 10)
  })
}

function disable() {
  domain.disable(name.value).then(() => {
    message.success($gettext('Disabled successfully'))
    enabled.value = false
  }).catch(r => {
    message.error($gettext('Failed to disable %{msg}', { msg: r.message ?? '' }))
  })
}

function on_change_enabled(checked: CheckedType) {
  modal.confirm({
    title: checked ? $gettext('Do you want to enable this site?') : $gettext('Do you want to disable this site?'),
    mask: false,
    centered: true,
    okText: $gettext('OK'),
    cancelText: $gettext('Cancel'),
    async onOk() {
      if (checked)
        enable()
      else
        disable()
    },
  })
}

</script>

<template>
  <ACard
    class="right-settings"
    :bordered="false"
  >
    <ContextHolder />
    <ACollapse
      v-model:active-key="active_key"
      ghost
    >
      <ACollapsePanel
        key="1"
        :header="$gettext('Basic')"
      >
        <AFormItem :label="$gettext('Enabled')">
          <ASwitch
            :checked="enabled"
            @change="on_change_enabled"
          />
        </AFormItem>
        <AFormItem :label="$gettext('Name')">
          <AInput v-model:value="filename" />
        </AFormItem>
        <AFormItem :label="$gettext('Updated at')">
          {{ formatDateTime(data.modified_at) }}
        </AFormItem>
      </ACollapsePanel>
      <ACollapsePanel
        v-if="!settings.is_remote"
        key="2"
        :header="$gettext('Deploy')"
      >
        <Deploy />
      </ACollapsePanel>
      <ACollapsePanel
        key="3"
        header="ChatGPT"
      >
        <ChatGPT
          v-model:history-messages="history_chatgpt_record"
          :content="configText"
          :path="filepath"
        />
      </ACollapsePanel>
    </ACollapse>
  </ACard>
</template>

<style scoped lang="less">
.right-settings {
  position: sticky;
  top: 78px;

  :deep(.ant-card-body) {
    max-height: 100vh;
    overflow-y: scroll;
  }
}

:deep(.ant-collapse-ghost > .ant-collapse-item > .ant-collapse-content > .ant-collapse-content-box) {
  padding: 0;
}

:deep(.ant-collapse > .ant-collapse-item > .ant-collapse-header) {
  padding: 0 0 10px 0;
}
</style>
