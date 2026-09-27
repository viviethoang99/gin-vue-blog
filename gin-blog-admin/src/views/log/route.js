const Layout = () => import('@/layout/index.vue')

export default {
  name: 'Log',
  path: '/log',
  component: Layout,
  redirect: '/log/operation',
  meta: {
    title: 'Operation Logs',
    icon: 'mdi:math-log',
    order: 6,
  },
  children: [
    {
      name: 'OperatingLog',
      path: 'operation',
      component: () => import('./operation/index.vue'),
      meta: {
        title: 'Operation Logs',
        icon: 'mdi:book-open-page-variant-outline',
        keepAlive: true,
      },
    },
    {
      name: 'LoginLog',
      path: 'login',
      component: () => import('./login/index.vue'),
      meta: {
        title: 'Login Logs',
        icon: 'material-symbols:login',
        keepAlive: true,
      },
    },
  ],
}
