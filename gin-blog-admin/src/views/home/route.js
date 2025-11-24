const Layout = () => import('@/layout/index.vue')

export default {
  name: 'Home',
  path: '/',
  component: Layout,
  redirect: '/home',
  meta: {
    order: 0,
  },
  isCatalogue: true,
  children: [
    {
      name: 'Home2',
      path: 'home2',
      component: () => import('./index.vue'),
      meta: {
        title: 'Dashboard',
        icon: 'ic:sharp-home',
        order: 0,
      },
    },
  ],
}
