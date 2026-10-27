<template>
  <div class="u-page">
    <UserPageHeader :eyebrow="t('hero.redeemEyebrow')" :title="t('home.heroTitle')" :subtitle="t('home.heroSub')" badges />
    <div class="u-main max-w-5xl">
      <section class="public-start card animate-slideInUp">
        <div class="public-start-copy">
          <span class="public-kicker">01</span>
          <h2 class="text-2xl font-bold text-ink">{{ services[0].title }}</h2>
          <p class="mt-2 text-sm text-muted">{{ t('home.startHint') }}</p>
          <router-link :to="services[0].to" class="btn-primary mt-6">{{ services[0].cta }} <span aria-hidden="true">→</span></router-link>
        </div>
        <div class="public-start-steps">
          <div v-for="(line, index) in services[0].points" :key="line" class="public-mini-step">
            <span>{{ index + 1 }}</span>
            <p>{{ line }}</p>
          </div>
        </div>
      </section>

      <div class="public-section-heading">
        <div>
          <span class="public-kicker">02</span>
          <h2 class="text-xl font-bold text-ink">{{ t('home.moreTitle') }}</h2>
        </div>
        <p class="text-sm text-muted">{{ t('home.moreHint') }}</p>
      </div>

      <div class="grid gap-4 md:grid-cols-3 animate-slideInUp">
        <router-link v-for="svc in secondaryServices" :key="svc.to" :to="svc.to" class="public-service card card-hover group">
          <span class="public-service-icon" v-html="svc.icon" />
          <div class="mt-5 flex items-start justify-between gap-3">
            <div>
              <h3 class="text-lg font-bold text-ink">{{ svc.title }}</h3>
              <p class="mt-1 text-xs text-subtle">{{ svc.en }}</p>
            </div>
            <span class="public-arrow" aria-hidden="true">↗</span>
          </div>
          <p class="mt-4 text-sm leading-6 text-muted">{{ svc.points[0] }} · {{ svc.points[1] }}</p>
        </router-link>
      </div>

      <section class="public-flow card">
        <div class="public-section-heading !mb-6">
          <div>
            <span class="public-kicker">03</span>
            <h2 class="text-xl font-bold text-ink">{{ t('home.flowTitle') }}</h2>
          </div>
          <p class="text-sm text-muted">{{ t('home.flowHint') }}</p>
        </div>
        <div class="public-flow-grid">
          <div v-for="(flow, index) in primaryFlow" :key="flow.title" class="public-flow-item">
            <span>{{ index + 1 }}</span>
            <div>
              <h3 class="font-semibold text-ink">{{ flow.title }}</h3>
              <p class="mt-1 text-sm text-muted">{{ flow.steps[0] }}</p>
            </div>
          </div>
        </div>
        <router-link :to="services[0].to" class="public-text-link">{{ t('home.startNow') }} <span aria-hidden="true">→</span></router-link>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import UserPageHeader from '../components/UserPageHeader.vue'

const { t } = useI18n({ useScope: 'global' })

const services = computed(() => [
  {
    to: '/recharge',
    title: t('home.services.recharge.title'),
    en: t('home.services.recharge.en'),
    cta: t('home.services.recharge.cta'),
    icon: '<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><path d="M14 2v6h6"/><path d="m9 14 2 2 4-4"/></svg>',
    points: [t('home.services.recharge.p1'), t('home.services.recharge.p2'), t('home.services.recharge.p3')],
  },
  {
    to: '/batch',
    title: t('home.services.batch.title'),
    en: t('home.services.batch.en'),
    cta: t('home.services.batch.cta'),
    icon: '<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/></svg>',
    points: [t('home.services.batch.p1'), t('home.services.batch.p2'), t('home.services.batch.p3')],
  },
  {
    to: '/history',
    title: t('home.services.lookup.title'),
    en: t('home.services.lookup.en'),
    cta: t('home.services.lookup.cta'),
    icon: '<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"/><path d="m21 21-4.3-4.3"/></svg>',
    points: [t('home.services.lookup.p1'), t('home.services.lookup.p2'), t('home.services.lookup.p3')],
  },
  {
    to: '/billing',
    title: t('home.services.billing.title'),
    en: t('home.services.billing.en'),
    cta: t('home.services.billing.cta'),
    icon: '<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="20" x2="18" y2="10"/><line x1="12" y1="20" x2="12" y2="4"/><line x1="6" y1="20" x2="6" y2="14"/></svg>',
    points: [t('home.services.billing.p1'), t('home.services.billing.p2'), t('home.services.billing.p3')],
  },
])

const secondaryServices = computed(() => services.value.slice(1))
const primaryFlow = computed(() => [
  { title: t('home.flows.submit.title'), steps: [t('home.flows.submit.s1')] },
  { title: t('home.flows.submit.s2').split(' — ')[0], steps: [t('home.flows.submit.s2')] },
  { title: t('home.flows.submit.s3').split(' — ')[0], steps: [t('home.flows.submit.s3')] },
])
</script>
